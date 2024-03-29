package world

import (
	"context"
	"errors"
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
	shardControl chan startMessage
	log          *zap.Logger
	ps           nats.PubSub

	roster     players.Roster
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
	err         error
}

func NewSharder(log *zap.Logger, control chan control.Command, conf control.WorldConf, ps nats.PubSub, world *World, roster players.Roster) *Sharder {
	return &Sharder{
		control:      control,
		shardControl: make(chan startMessage),
		log:          log,
		ps:           ps,
		roster:       roster,
		shardSizeX:   int64(conf.ShardSize),
		shardSizeZ:   int64(conf.ShardSize),
		world:        world,
		shards:       make(map[ShardID]*shard),
	}
}

func (sh *Sharder) Start(ctx context.Context) {
	go sh.dispatchSharderLoop(ctx)
	sh.log.Info("sharder started")
}

func (sh *Sharder) FindShardID(dimID uuid.UUID, coords data.PositionI) (ShardID, bool) {
	dim, ok := sh.world.Dimensions[dimID]
	if !ok {
		sh.log.Debug("dimension not found", zap.String("dimID", dimID.String()))
		return "", false
	}

	dimEdges := dim.Edges()
	if dimEdges.PositiveX < coords.X {
		sh.log.Debug("coords outside dim edges", zap.Any("edges", dimEdges), zap.Int64("X", coords.X))
		return "", false
	}
	if dimEdges.PositiveZ < coords.Z {
		sh.log.Debug("coords outside dim edges", zap.Any("edges", dimEdges), zap.Int64("Z", coords.Z))
		return "", false
	}
	if dimEdges.NegativeX > coords.X {
		sh.log.Debug("coords outside dim edges", zap.Any("edges", dimEdges), zap.Int64("X", coords.X))
		return "", false
	}
	if dimEdges.NegativeZ > coords.Z {
		sh.log.Debug("coords outside dim edges", zap.Any("edges", dimEdges), zap.Int64("Z", coords.Z))
		return "", false
	}

	lowestX := int64(math.Floor(float64(coords.X)/float64(sh.shardSizeX*level.ChunkX))) * sh.shardSizeX * level.ChunkX
	lowestZ := int64(math.Floor(float64(coords.Z)/float64(sh.shardSizeZ*level.ChunkZ))) * sh.shardSizeZ * level.ChunkZ

	id := MkShardIDFromCoords(dim.Name(), lowestX, lowestZ)
	if _, ok := sh.shards[id]; !ok {
		sh.log.Error("couldn't find shard for valid coordinates", zap.String("shardID", string(id)), zap.Any("coords", coords))
		return "", false
	}

	return id, true
}

func (sh *Sharder) dispatchSharderLoop(ctx context.Context) {
	sh.signal(control.STARTING, nil)

	defer func() {
		if !sh.isStopping {
			err := errors.New("sharder stopped unexpectedly")
			if r := recover(); r != nil {
				err = fmt.Errorf("sharder crashed: %v", r)
			}
			sh.signal(control.FAILED, err)
		}
	}()

	go sh.bootstrapShards()

	for {
		select {
		case <-ctx.Done():
			sh.Lock()
			sh.isStopping = true

			sh.log.Info("server context closed, sharder shutdown sequence initiated")
			for {
				select { // DEBT maybe have a failsafe timeout for waiting for shards to stop
				// If shutdown initiated - all shards will close on context cancellation, and will report over
				// the shardControl channel. Once all shards have reported and have been deleted - stop the loop,
				// signal to global control and exit the sharder loop.
				case shardStartMsg := <-sh.shardControl:
					delete(sh.shards, shardStartMsg.id)
					if len(sh.shards) == 0 {
						sh.log.Info("no shards left, stopping sharder")
						sh.Unlock()
						sh.signal(control.STOPPED, nil)
						return // all shards are now removed and the sharder loop can stop
					}
				}
			}
		case shardStartMsg := <-sh.shardControl:
			sh.Lock()
			if _, ok := sh.shards[shardStartMsg.id]; !ok {
				var err error
				sh.shards[shardStartMsg.id], err = newShard(sh.log, sh.ps, shardStartMsg.id, shardStartMsg.dimensionID, shardStartMsg.chunkIDs)
				if err != nil {
					sh.log.Error("failed to instantiate shard, signalling shard failure", zap.Error(err))
					sh.signal(control.FAILED, fmt.Errorf("failed to start shard %s: %w", shardStartMsg.id, err))
					sh.Unlock()
					continue
				}
			}

			sh.log.Debug("starting shard", zap.String("id", string(shardStartMsg.id)), zap.Int("chunks", len(shardStartMsg.chunkIDs)))

			if err := sh.shards[shardStartMsg.id].dispatch(ctx, sh.roster, sh.shardControl, sh.world); err != nil {
				sh.log.Error("failed to restart shard, signalling shard failure", zap.Error(err))
				sh.signal(control.FAILED, fmt.Errorf("failed to restart shard %s: %w", shardStartMsg.id, err))
			}

			sh.Unlock()
		}
	}
}

func (sh *Sharder) bootstrapShards() {
	defer func() {
		if r := recover(); r != nil {
			sh.log.Error("failed to bootstrap sharder", zap.Any("panic", r))
			sh.signal(control.FAILED, fmt.Errorf("sharder bootstrap panicked: %v", r))
		}
	}()

	var shardCount int
	for id, dimension := range sh.world.Dimensions {
		sh.log.Debug("bootstrapping shards", zap.String("level", dimension.Name()), zap.Any("edges", dimension.Edges()))
		shardStarts := splitDimensionShards(id, dimension.Name(), dimension.Edges(), sh.shardSizeX, sh.shardSizeZ)
		for _, start := range shardStarts {
			sh.shardControl <- start
			shardCount++
		}
	}
	sh.signal(control.READY, nil)
	sh.log.Info(fmt.Sprintf("%d shards started in %d dimensions", shardCount, len(sh.world.Dimensions)))
}

// Split an area defined by given bottom left and top right points into chunks and return list of chunk IDs.
func splitAreaChunks(lowerX, lowerZ, higherX, higherZ int64) []level.ChunkID {
	var chunkIDs []level.ChunkID
	resetLowerZ := lowerZ

	for lowerX < higherX {
		lowerZ = resetLowerZ
		for lowerZ < higherZ {
			chunkIDs = append(chunkIDs, level.MkChunkID(lowerX, lowerZ))
			lowerZ += level.ChunkZ
		}
		lowerX += level.ChunkX
	}
	return chunkIDs
}

func splitDimensionShards(dimID uuid.UUID, dimName string, edges level.Edges, shardX, shardZ int64) map[ShardID]startMessage {
	var shardStarts = make(map[ShardID]startMessage)

	// Start from 0.0 point and cover all four quadrants of the map.
	// This expects that 0.0 point is actually within the boundaries of the map. It can be on the edge though.

	// North-East quadrant
	var shardEdgeStartX, shardEdgeStartZ int64
	for shardEdgeStartX < edges.PositiveX {
		for shardEdgeStartZ < edges.PositiveZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX,
				shardEdgeStartZ,
				min(shardEdgeStartX+shardX*level.ChunkX, edges.PositiveX),
				min(shardEdgeStartZ+shardZ*level.ChunkZ, edges.PositiveZ),
			)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			if _, ok := shardStarts[shardID]; ok {
				println("duplicate shard", shardID)
			}
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ += shardZ * level.ChunkZ
		}
		shardEdgeStartZ = 0
		shardEdgeStartX += shardX * level.ChunkX
	}

	// South-East quadrant
	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX < edges.PositiveX {
		for shardEdgeStartZ > edges.NegativeZ {
			shardChunks := splitAreaChunks(
				shardEdgeStartX,
				max(shardEdgeStartZ-shardZ*level.ChunkZ, edges.NegativeZ),
				min(shardEdgeStartX+shardX*level.ChunkX, edges.PositiveX),
				shardEdgeStartZ)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ -= shardZ * level.ChunkZ
		}
		shardEdgeStartZ = 0
		shardEdgeStartX += shardX * level.ChunkX
	}

	// North-West quadrant
	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX > edges.NegativeX {
		for shardEdgeStartZ < edges.PositiveZ {
			shardChunks := splitAreaChunks(
				max(shardEdgeStartX-shardX*level.ChunkX, edges.NegativeX),
				shardEdgeStartZ,
				shardEdgeStartX,
				min(shardEdgeStartZ+shardZ*level.ChunkZ, edges.PositiveZ),
			)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ += shardZ * level.ChunkZ
		}
		shardEdgeStartZ = 0
		shardEdgeStartX -= shardX * level.ChunkX
	}

	// South-West quadrant
	shardEdgeStartX = 0
	shardEdgeStartZ = 0
	for shardEdgeStartX > edges.NegativeX {
		for shardEdgeStartZ > edges.NegativeZ {
			shardChunks := splitAreaChunks(
				max(shardEdgeStartX-shardX*level.ChunkX, edges.NegativeX),
				max(shardEdgeStartZ-shardZ*level.ChunkZ, edges.NegativeZ),
				shardEdgeStartX,
				shardEdgeStartZ,
			)
			shardID := mkShardIDFromChunks(dimName, shardChunks)
			shardStarts[shardID] = startMessage{
				id:          shardID,
				dimensionID: dimID,
				chunkIDs:    shardChunks,
			}
			shardEdgeStartZ -= shardZ * level.ChunkZ
		}
		shardEdgeStartZ = 0
		shardEdgeStartX -= shardX * level.ChunkX
	}

	return shardStarts
}

func (sh *Sharder) signal(state control.ComponentState, err error) {
	sh.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.SHARDER,
		State:     state,
		Err:       err,
	}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
