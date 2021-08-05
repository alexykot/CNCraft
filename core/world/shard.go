package world

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/core/world/events"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type shard struct {
	sync.Mutex

	id       ShardID
	log      *zap.Logger
	ps       nats.PubSub
	dimID    uuid.UUID
	chunkIDs []level.ChunkID
	events   []*envelope.E

	tickHandlers  map[string]events.TickHandler
	eventHandlers map[pb.OneOfEvent]map[string]events.EventHandler
}

func newShard(log *zap.Logger, ps nats.PubSub, id ShardID, dimID uuid.UUID, chunkIDs []level.ChunkID) (*shard, error) {
	if len(chunkIDs) < 1 {
		// not starting a shard if no chunks provided
		return nil, fmt.Errorf("cannot instantiate shard with zero chunks; shard %s, dim %s", id.String(), dimID.String())
	}
	return &shard{
		id:            id,
		log:           log.With(zap.String("shard", string(id))),
		ps:            ps,
		dimID:         dimID,
		chunkIDs:      chunkIDs,
		tickHandlers:  make(map[string]events.TickHandler),
		eventHandlers: make(map[pb.OneOfEvent]map[string]events.EventHandler),
	}, nil
}

func (s *shard) dispatch(ctx context.Context, roster *players.Roster, failSignaller chan startMessage, world *World) error {
	if err := s.initiateHandlers(s.chunkIDs, world, roster); err != nil {
		return fmt.Errorf("failed to instantiate world processors: %w", err)
	}

	if err := s.ps.Subscribe(subj.MkShardEvent(string(s.id)), s.incomingEventHandler); err != nil {
		return fmt.Errorf("failed to register shard events handler: %w", err)
	}

	go s.runEventLoop(ctx, failSignaller)
	return nil
}

func (s *shard) incomingEventHandler(lope *envelope.E) {
	s.Lock()
	defer s.Unlock()

	shardEvent := lope.GetShardEvent()
	if shardEvent == nil {
		return // if not a shard event - silently ignore
	}

	s.events = append(s.events, lope)
}

func (s *shard) runEventLoop(ctx context.Context, failSignaller chan startMessage) {
	var err error
	defer func() {
		s.ps.Unsubscribe(subj.MkShardEvent(string(s.id))) // remove obsolete subscription
		if r := recover(); r != nil {
			err = fmt.Errorf("shard event loop crashed: %v", r)
		}

		// Make sure the shard restart attempted whenever it fails for any reason.
		// If the server is stopping - sharder will ignore this and not restart the shard.
		failSignaller <- startMessage{
			id:          s.id,
			dimensionID: s.dimID,
			chunkIDs:    s.chunkIDs,
			err:         err,
		}
	}()

	ticker := time.NewTicker(game.TickSpeed)

	s.log.Debug("starting event loop")
	for {
		select {
		case <-ctx.Done():
			s.log.Info("stopping shard")
			return // trigger defer and make it return a message, error should be nil
		case tickTime := <-ticker.C:
			tick := game.Tick(tickTime.UnixNano()) // Round to milliseconds maybe?

			if err = s.handleTick(tick, s.cutEvents()); err != nil {
				err = fmt.Errorf("failed to handle tick events: %w", err)
				return // trigger defer and make it return a message with error attached
			}
		}
	}
}

// cutEvents returns a copy of the current outstanding events ready for handling and nullifies the events list itself.
func (s *shard) cutEvents() []*envelope.E {
	s.Lock()
	defer s.Unlock()

	if len(s.events) == 0 {
		return nil
	}

	eventsCopy := s.events
	s.events = nil

	return eventsCopy
}

func mkShardIDFromChunks(dimName string, shardChunks []level.ChunkID) ShardID {
	var leastX, leastZ int64
	for _, chunkID := range shardChunks {
		x, z := level.XZFromChunkID(chunkID)
		if leastX > x {
			leastX = x
		}
		if leastZ > z {
			leastZ = z
		}
	}

	return MkShardIDFromCoords(dimName, leastX, leastZ)
}

func MkShardIDFromCoords(dimName string, x, z int64) ShardID {
	return ShardID(fmt.Sprintf("shard.%s.%d.%d", dimName, x, z))
}

func (s *shard) handleTick(tick game.Tick, tickEvents []*envelope.E) error {
	for name, tickHandler := range s.tickHandlers {
		userOutLopes, err := tickHandler(tick)
		if err != nil {
			return fmt.Errorf("failed to handle tick in handler `%s` of shard `%s`: %w", name, s.id, err)
		}
		for userId, outLopes := range userOutLopes {
			if err := s.ps.Publish(subj.MkConnTransmit(userId), outLopes...); err != nil {
				return fmt.Errorf("failed to publish conn.transmit message for conn `%s`, shard `%s`: %w", userId, s.id, err)
			}
		}
	}

	for _, event := range tickEvents {
		// Don't see a simpler better way to enumerate and find actual message inside a one-off type.
		if playerDigging := event.ShardEvent.GetPlayerDigging(); playerDigging != nil {
			for name, eventHandler := range s.eventHandlers[pb.Event_PlayerDigging] {
				userOutLopes, err := eventHandler(tick, event)
				if err != nil {
					return fmt.Errorf("failed to handle tick event `%s` in handler `%s` of shard `%s`: %w",
						pb.Event_PlayerDigging, name, s.id, err)
				}
				for userId, outLopes := range userOutLopes {
					if err := s.ps.Publish(subj.MkConnTransmit(userId), outLopes...); err != nil {
						return fmt.Errorf("failed to publish conn.transmit message for conn `%s`, shard `%s`: %w", userId, s.id, err)
					}
				}
			}
		}
	}
	return nil
}

func (s *shard) initiateHandlers(chunkIDs []level.ChunkID, world *World, roster *players.Roster) error {
	var err error
	chunks := make([]level.Chunk, len(chunkIDs), len(chunkIDs))
	for index, chunkID := range chunkIDs {
		chunks[index], err = world.getChunk(s.dimID, chunkID)
		if err != nil {
			return fmt.Errorf("failed to retrieve chunk %s, dim %s: %w", chunkID, s.dimID, err)
		}
	}

	for _, handler := range events.NewHandlers(chunks, roster) {
		if tickHandler := handler.GetTickHandler(); tickHandler != nil {
			s.tickHandlers[handler.Name()] = tickHandler
		}

		for eventType, evHandler := range handler.GetEventHandlers() {
			if s.eventHandlers[eventType] == nil {
				s.eventHandlers[eventType] = make(map[string]events.EventHandler)
			}
			s.eventHandlers[eventType][handler.Name()] = evHandler
		}
	}

	return nil
}
