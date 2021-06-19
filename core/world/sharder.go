package world

import (
	"fmt"
	"math"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type Sharder struct {
	sync.Mutex

	control      chan control.Command
	shardStarter chan startMessage
	log          *zap.Logger
	ps           nats.PubSub

	roster     *players.Roster
	world      *World
	shardSizeX int64
	shardSizeZ int64
	shards     map[ShardID]*shard
	isStopping bool
}

// ShardID is formatted as `shard.<levelName>.<lowestX>.<lowestZ>`, e.g. `shard.Overworld.0.-160`.
type ShardID string

func (sid ShardID) String() string {
	return string(sid)
}

type startMessage struct {
	id          ShardID
	dimensionID uuid.UUID
	chunkIDs    []level.ChunkID
}

func NewSharder(conf control.WorldConf, log *zap.Logger, ps nats.PubSub, control chan control.Command, world *World, roster *players.Roster) *Sharder {
	return &Sharder{
		control:      control,
		shardStarter: make(chan startMessage),
		log:          log,
		ps:           ps,
		roster:       roster,
		shardSizeX:   int64(conf.ShardSize),
		shardSizeZ:   int64(conf.ShardSize),
		world:        world,
		shards:       make(map[ShardID]*shard),
	}
}

func (sh *Sharder) Start() {
	go sh.dispatchSharderLoop()
	sh.log.Info("sharder started")
}

func (sh *Sharder) FindShardID(dimID uuid.UUID, coords *data.PositionI) (ShardID, bool) {
	dim, ok := sh.world.Dimensions[dimID]
	if !ok {
		return "", false
	}

	dimEdges := dim.Edges()
	if dimEdges.PositiveX < coords.X {
		return "", false
	}
	if dimEdges.PositiveZ < coords.Z {
		return "", false
	}
	if dimEdges.NegativeX < coords.X {
		return "", false
	}
	if dimEdges.NegativeZ < coords.Z {
		return "", false
	}

	lowestX := int64(math.Floor(float64(coords.X) / float64(sh.shardSizeX)))
	lowestZ := int64(math.Floor(float64(coords.Z) / float64(sh.shardSizeZ)))

	id := MkShardIDFromCoords(dim.Name(), lowestX, lowestZ)
	if _, ok := sh.shards[id]; !ok {
		sh.log.Error("couldn't find shard for valid coordinates", zap.String("shardID", string(id)), zap.Any("coords", coords))
		return "", false
	}

	return id, true
}

func (sh *Sharder) dispatchSharderLoop() {
	defer func() {
		if !sh.isStopping {
			message := "sharder stopped unexpectedly"
			if r := recover(); r != nil {
				message = fmt.Sprintf("sharder panicked: %v", r)
			}
			// stop the server if Sharder exits for any reason
			sh.control <- control.Command{Signal: control.SERVER_FAIL, Message: message}
		}
	}()

	go sh.bootstrapShards()

	for {
		select {
		case command := <-sh.control:
			switch command.Signal {
			case control.SHARD_FAIL:
				sh.log.Error("shard failed, signalling server failure", zap.String("message", command.Message))
				sh.control <- control.Command{ // signalling critical failure
					Signal:  control.SERVER_FAIL,
					Message: fmt.Sprintf("failed to start a shard: %s", command.Message),
				}
			case control.SERVER_STOP, control.SERVER_FAIL:
				sh.Lock()
				sh.isStopping = true
				if len(sh.shards) == 0 { // stop the loop if all shards are already stopped
					sh.log.Info("no shards left, stopping sharder")
					sh.Unlock()
					return
				} else { // otherwise command shards to stop
					sh.control <- control.Command{ // signalling shards to stop
						Signal:  control.SHARD_STOP,
						Message: "stop all shards",
					}
					sh.log.Info("signalling shards to stop")
				}
				sh.Unlock()
			}
		case shardStarter := <-sh.shardStarter:
			sh.Lock()
			_, ok := sh.shards[shardStarter.id]
			if !ok {
				var err error
				sh.shards[shardStarter.id], err = newShard(sh.log, sh.ps, shardStarter.id, shardStarter.dimensionID, shardStarter.chunkIDs)
				if err != nil {
					sh.log.Error("failed to instantiate shard, signalling shard failure", zap.Error(err))
					sh.control <- control.Command{ // signalling critical failure
						Signal:  control.SHARD_FAIL,
						Message: fmt.Sprintf("failed to start a shard: %s", err.Error()),
					}
					sh.Unlock()
					continue
				}
			}
			if sh.isStopping {
				delete(sh.shards, shardStarter.id)
				if len(sh.shards) == 0 {
					sh.log.Info("no shards left, stopping sharder")
					sh.Unlock()
					return // all shards are now removed and the sharder loop can stop
				}
			} else {
				sh.log.Debug("starting shard", zap.String("id", string(shardStarter.id)), zap.Int("chunks", len(shardStarter.chunkIDs)))

				go sh.shards[shardStarter.id].dispatch(sh.roster, sh.control, sh.shardStarter, sh.world)
			}
			sh.Unlock()
		}
	}
}

func (sh *Sharder) bootstrapShards() {
	defer func() {
		if r := recover(); r != nil {
			sh.control <- control.Command{
				Signal:  control.SERVER_FAIL,
				Message: fmt.Sprintf("sharder bootstrap panicked: %v", r),
			}
		}
	}()

	var shardCount int
	for id, dimension := range sh.world.Dimensions {
		sh.log.Debug("bootstrapping shards", zap.String("level", dimension.Name()), zap.Any("edges", dimension.Edges()))
		shardStarts := splitDimShards(id, dimension.Name(), dimension.Edges(), sh.shardSizeX, sh.shardSizeZ)
		for _, start := range shardStarts {
			sh.shardStarter <- start
			shardCount++
		}
	}
	sh.log.Info(fmt.Sprintf("%d shards started in %d dimensions", shardCount, len(sh.world.Dimensions)))
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

func splitDimShards(dimID uuid.UUID, dimName string, edges level.Edges, shardX, shardZ int64) map[ShardID]startMessage {
	var shardStarts = make(map[ShardID]startMessage)

	// starting from 0.0 coords - cover all four quadrants of the map.
	// This assumes 0.0 coords is actually within the boundaries of the map. It can be on the edge though.
	var shardEdgeStartX, shardEdgeStartZ int64
	for shardEdgeStartX < edges.PositiveX {
		for shardEdgeStartZ < edges.PositiveZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX,
				shardEdgeStartZ,
				shardEdgeStartX+shardX*level.ChunkX,
				shardEdgeStartZ+shardZ*level.ChunkZ)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ += shardZ * level.ChunkZ
		}
		shardEdgeStartX += shardX * level.ChunkX
	}

	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX < edges.PositiveX {
		for shardEdgeStartZ > edges.NegativeZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX,
				shardEdgeStartZ-shardZ*level.ChunkZ,
				shardEdgeStartX+shardX*level.ChunkX,
				shardEdgeStartZ)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ -= shardZ * level.ChunkZ
		}
		shardEdgeStartX += shardX * level.ChunkX
	}

	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX > edges.NegativeX {
		for shardEdgeStartZ < edges.PositiveZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX-shardX*level.ChunkX,
				shardEdgeStartZ,
				shardEdgeStartX,
				shardEdgeStartZ+shardZ*level.ChunkZ)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ += shardZ * level.ChunkZ
		}
		shardEdgeStartX -= shardX * level.ChunkX
	}

	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX > edges.NegativeX {
		for shardEdgeStartZ > edges.NegativeZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX-shardX*level.ChunkX,
				shardEdgeStartZ-shardZ*level.ChunkZ,
				shardEdgeStartZ,
				shardEdgeStartX)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ -= shardZ * level.ChunkZ
		}
		shardEdgeStartX -= shardX * level.ChunkX
	}

	return shardStarts
}
