package world

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/world/handlers"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type Shard struct {
	id     ShardID
	log    *zap.Logger
	chunks []level.ChunkID

	mu     *sync.Mutex
	events []*envelope.E

	tickHandlers  []handlers.TickHandler
	eventHandlers map[pb.OneOfEvent][]handlers.EventHandler
}

// TODO next:
//  - handle sample incoming packet (e.g. StartMining) and put it into event queue for correct shard
//  - process sample incoming event (e.g. StartMining) and dispatch async response to correct client
//  - create and register individual instances of world event and tick handlers for every instantiated shard
//  - consider separating shard dispatch from the shard instantiation

func newShard(log *zap.Logger, id ShardID, chunks []level.ChunkID) *Shard {
	shard := &Shard{
		id:            id,
		log:           log.With(zap.String("shard", string(id))),
		chunks:        chunks,
		mu:            &sync.Mutex{},
		eventHandlers: make(map[pb.OneOfEvent][]handlers.EventHandler),
	}

	worldProcessors := handlers.Get()
	for _, proc := range worldProcessors {
		handlerMap := proc.GetEventHandlers()
		for eventType, evHandlers := range handlerMap {
			shard.eventHandlers[eventType] = append(shard.eventHandlers[eventType], evHandlers...)
		}

		tickHandlers := proc.GetTickHandlers()
		shard.tickHandlers = append(shard.tickHandlers, tickHandlers...)
	}

	return shard
}

func (s *Shard) dispatch(ps nats.PubSub, controller chan control.Command, restarter chan startMessage) {
	if len(s.chunks) < 1 {
		return // not starting a shard without any chunks
	}

	if err := ps.Subscribe(subj.MkShardEvent(string(s.id)), s.handleIncomingEvent()); err != nil {
		controller <- control.Command{
			Signal:  control.SHARD_FAIL,
			Message: fmt.Errorf("failed to register PlayerLoading handler: %w", err).Error(),
		}
		return
	}

	defer func() {
		ps.Unsubscribe(subj.MkShardEvent(string(s.id))) // remove old subscription
		if r := recover(); r != nil {
			s.log.Error("shard event loop crashed", zap.Any("panic", r))
		}

		restarter <- startMessage{ // make sure the shard is restarted whenever it fails for any reason
			id:     s.id,
			chunks: s.chunks,
		}
	}()

	s.runEventLoop(controller)
}

func (s *Shard) handleIncomingEvent() func(lope *envelope.E) {
	return func(lope *envelope.E) {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.events = append(s.events, lope)
	}
}

func (s *Shard) runEventLoop(controller chan control.Command) {
	ticker := time.NewTicker(game.TickSpeed)

	s.log.Info("starting event loop")
	for {
		select {
		case command := <-controller:
			switch command.Signal {
			case control.SHARD_STOP:
				s.log.Info("stopping shard")
				return
			}

		case tickTime := <-ticker.C:
			// Round to milliseconds maybe?
			tick := game.Tick(tickTime.UnixNano())

			if err := s.handleTick(tick, s.copyEvents()); err != nil {
				s.log.Error("failed to handle tick events", zap.String("shard", string(s.id)), zap.Error(err))
			}
		}
	}
}

func (s *Shard) copyEvents() []*envelope.E {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) == 0 {
		return nil
	}

	eventsCopy := make([]*envelope.E, len(s.events), len(s.events))
	for i, event := range s.events {
		eventsCopy[i] = event
	}
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

func (s *Shard) handleTick(tick game.Tick, tickEvents []*envelope.E) error {
	for range tickEvents {

	}
	return nil
}
