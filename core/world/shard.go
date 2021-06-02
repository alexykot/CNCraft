package world

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game/level"
)

const tickSpeed = time.Millisecond * 200

type Shard struct {
	id     ShardID
	log    *zap.Logger
	chunks []level.ChunkID
}

// TODO next:
//  - make shard event queue
//  - decide on event generic structure to allow for multiple event types
//  - make shard router in sharder, by event coords and level reference
//  - handle sample incoming packet (e.g. StartMining) and put it into event queue for correct shard
//  - process sample incoming event (e.g. StartMining) and dispatch async response to correct client

func newShard(log *zap.Logger, id ShardID, chunks []level.ChunkID) *Shard {
	return &Shard{
		id:     id,
		log:    log.With(zap.String("shard", string(id))),
		chunks: chunks,
	}
}

func (s Shard) dispatch(ps nats.PubSub, controller chan control.Command, restarter chan startShard) {
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

		restarter <- startShard{ // make sure the shard is restarted whenever it fails for any reason
			id:     s.id,
			chunks: s.chunks,
		}
	}()

	s.runEventLoop(controller)
}

func (s Shard) handleIncomingEvent() func(lope *envelope.E) {
	return func(lope *envelope.E) {

	}
}

func (s Shard) runEventLoop(controller chan control.Command) {
	ticker := time.NewTicker(tickSpeed)

	s.log.Info("starting event loop")
	for {
		select {
		case command := <-controller:
			switch command.Signal {
			case control.SHARD_STOP:
				s.log.Info("stopping shard")
				return
			}

		case <-ticker.C:
			// TODO event loop goes here
		}
	}
}

func mkShardIDFromChunks(levelName string, shardChunks []level.ChunkID) ShardID {
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

	return MkShardIDFromCoords(levelName, leastX, leastZ)
}

func MkShardIDFromCoords(levelName string, x, z int64) ShardID {
	return ShardID(fmt.Sprintf("shard.%s.%d.%d", levelName, x, z))
}
