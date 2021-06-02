package world

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type Sharder struct {
	control      chan control.Command
	shardStarter chan startShard
	log          *zap.Logger
	ps           nats.PubSub

	shardSizeX int64
	shardSizeZ int64
	world      *World
	shards     map[ShardID]*Shard
	isStopping bool
}

type ShardID string

type startShard struct {
	id     ShardID
	chunks []level.ChunkID
}

func NewSharder(conf control.WorldConf, log *zap.Logger, ps nats.PubSub, control chan control.Command, world *World) *Sharder {
	return &Sharder{
		control:      control,
		shardStarter: make(chan startShard),
		log:          log,
		ps:           ps,
		shardSizeX:   int64(conf.ShardSize),
		shardSizeZ:   int64(conf.ShardSize),
		world:        world,
		shards:       make(map[ShardID]*Shard),
	}
}

func (sh Sharder) Start() {
	go sh.dispatchSharderLoop()
	sh.log.Info("sharder started")
}

func (sh Sharder) dispatchSharderLoop() {
	defer func() {
		message := "sharder stopped unexpectedly"
		if r := recover(); r != nil {
			message = fmt.Sprintf("sharder panicked: %v", r)
		}
		// stop the server if Sharder exits for any reason
		sh.control <- control.Command{Signal: control.SERVER_FAIL, Message: message}
	}()

	go sh.bootstrapShards()

	for {
		select {
		case command := <-sh.control:
			switch command.Signal {
			case control.SHARD_FAIL:
				sh.control <- control.Command{ // signalling critical failure
					Signal:  control.SERVER_FAIL,
					Message: fmt.Sprintf("failed to start a shard: %s", command.Message),
				}
				sh.log.Info("shard failed, signalling server failure", zap.String("message", command.Message))
			case control.SERVER_STOP:
				sh.isStopping = true
				if len(sh.shards) == 0 { // stop the loop if all shards are already stopped
					sh.log.Info("stopping sharder")
					return
				} else { // otherwise command shards to stop
					sh.control <- control.Command{ // signalling shards to stop
						Signal:  control.SHARD_STOP,
						Message: "stop all shards",
					}
					sh.log.Info("signalling shards to stop")
				}
			}
		case shardStarter := <-sh.shardStarter:
			_, ok := sh.shards[shardStarter.id]
			if !ok {
				sh.shards[shardStarter.id] = newShard(sh.log, shardStarter.id, shardStarter.chunks)
			}
			if sh.isStopping {
				delete(sh.shards, shardStarter.id)
				if len(sh.shards) == 0 {
					sh.log.Info("stopping sharder")
					return // all shards are now removed and the loop can stop
				}
			} else {
				sh.log.Debug("starting shard", zap.String("id", string(shardStarter.id)), zap.Int("chunks", len(shardStarter.chunks)))
				go sh.shards[shardStarter.id].dispatch(sh.ps, sh.control, sh.shardStarter)
			}
		}
	}
}

func (sh Sharder) bootstrapShards() {
	defer func() {
		if r := recover(); r != nil {
			sh.control <- control.Command{
				Signal:  control.SERVER_FAIL,
				Message: fmt.Sprintf("sharder bootstrap panicked: %v", r),
			}
		}
	}()

	var shardCount int
	for _, worldLevel := range sh.world.Levels {
		levelEdges := worldLevel.Edges()
		sh.log.Debug("bootstrapping shards", zap.String("level", worldLevel.Name()), zap.Any("edges", levelEdges))

		// starting from 0.0 coords - cover all four quadrants of the map.
		// This assumes 0.0 coords is actually within the boundaries of the map. It can be on the edge though.
		var shardEdgeStartX, shardEdgeStartZ int64
		for shardEdgeStartX < levelEdges.PositiveX {
			for shardEdgeStartZ < levelEdges.PositiveZ {
				shardChunks := splitAreaChunks(
					shardEdgeStartX,
					shardEdgeStartZ,
					shardEdgeStartX+sh.shardSizeX*level.ChunkX,
					shardEdgeStartZ+sh.shardSizeZ*level.ChunkZ)
				sh.shardStarter <- startShard{
					id:     mkShardIDFromChunks(worldLevel.Name(), shardChunks),
					chunks: shardChunks,
				}
				shardEdgeStartZ += sh.shardSizeZ * level.ChunkZ
				shardCount++
			}
			shardEdgeStartX += sh.shardSizeX * level.ChunkX
		}

		shardEdgeStartX = 0
		shardEdgeStartZ = 0
		for shardEdgeStartX < levelEdges.PositiveX {
			for shardEdgeStartZ > levelEdges.NegativeZ {
				shardChunks := splitAreaChunks(
					shardEdgeStartX,
					shardEdgeStartZ-sh.shardSizeZ*level.ChunkZ,
					shardEdgeStartX+sh.shardSizeX*level.ChunkX,
					shardEdgeStartZ)
				sh.shardStarter <- startShard{
					id:     mkShardIDFromChunks(worldLevel.Name(), shardChunks),
					chunks: shardChunks,
				}
				shardEdgeStartZ -= sh.shardSizeZ * level.ChunkZ
				shardCount++
			}
			shardEdgeStartX += sh.shardSizeX * level.ChunkX
		}

		shardEdgeStartX = 0
		shardEdgeStartZ = 0
		for shardEdgeStartX > levelEdges.NegativeX {
			for shardEdgeStartZ < levelEdges.PositiveZ {
				shardChunks := splitAreaChunks(
					shardEdgeStartX-sh.shardSizeX*level.ChunkX,
					shardEdgeStartZ,
					shardEdgeStartX,
					shardEdgeStartZ+sh.shardSizeZ*level.ChunkZ)
				sh.shardStarter <- startShard{
					id:     mkShardIDFromChunks(worldLevel.Name(), shardChunks),
					chunks: shardChunks,
				}
				shardEdgeStartZ += sh.shardSizeZ * level.ChunkZ
				shardCount++
			}
			shardEdgeStartX -= sh.shardSizeX * level.ChunkX
		}

		shardEdgeStartX = 0
		shardEdgeStartZ = 0
		for shardEdgeStartX > levelEdges.NegativeX {
			for shardEdgeStartZ > levelEdges.NegativeZ {
				shardChunks := splitAreaChunks(
					shardEdgeStartX-sh.shardSizeX*level.ChunkX,
					shardEdgeStartZ-sh.shardSizeZ*level.ChunkZ,
					shardEdgeStartZ,
					shardEdgeStartX)
				sh.shardStarter <- startShard{
					id:     mkShardIDFromChunks(worldLevel.Name(), shardChunks),
					chunks: shardChunks,
				}
				shardEdgeStartZ -= sh.shardSizeZ * level.ChunkZ
				shardCount++
			}
			shardEdgeStartX -= sh.shardSizeX * level.ChunkX
		}
	}
	sh.log.Info(fmt.Sprintf("%d shards bootstrapped", shardCount))
}

// Split an area defined by given bottom left and top right points into chunks and return list of chunk IDs.
func splitAreaChunks(lowerX, lowerZ, higherX, higherZ int64) []level.ChunkID {
	var chunkIDs []level.ChunkID
	for lowerX < higherX {
		for lowerZ < higherZ {
			chunkIDs = append(chunkIDs, level.MkChunkID(lowerX, lowerZ))
			lowerZ += level.ChunkZ
		}
		lowerX += level.ChunkX
	}
	return chunkIDs
}
