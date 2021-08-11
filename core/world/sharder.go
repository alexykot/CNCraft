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

	ctx context.Context

	control      chan control.Command
	shardControl chan startMessage
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
	err         error
}

func NewSharder(ctx context.Context, conf control.WorldConf, log *zap.Logger, ps nats.PubSub, control chan control.Command, world *World, roster *players.Roster) *Sharder {
	return &Sharder{
		ctx:          ctx,
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

func (sh *Sharder) Start() {
	go sh.dispatchSharderLoop()
	sh.log.Info("sharder started")
}

func (sh *Sharder) FindShardID(dimID uuid.UUID, coords data.PositionI) (ShardID, bool) {
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
		case <-sh.ctx.Done():
			sh.Lock()
			sh.isStopping = true

			sh.log.Info("server context closed, sharder shutdown sequence initiated")
			for {
				select { // DEBT maybe have a failsafe timeout for waiting shards to stop
				// If shutdown initiated - all shards will close on context cancellation, and will report over
				// the shardControl channel. Once all shards have reported and have been deleted - stop the loop,
				// signal to global control and exit sharder.
				case shardStartMsg := <-sh.shardControl:
					delete(sh.shards, shardStartMsg.id)
					if len(sh.shards) == 0 {
						sh.log.Info("no shards left, stopping sharder")
						sh.signal(control.STOPPED, nil)
						sh.Unlock()
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

			if err := sh.shards[shardStartMsg.id].dispatch(sh.ctx, sh.roster, sh.shardControl, sh.world); err != nil {
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
	for lowerX < higherX {
		for lowerZ < higherZ {
			chunkIDs = append(chunkIDs, level.MkChunkID(lowerX, lowerZ))
			lowerZ += level.ChunkZ
		}
		lowerX += level.ChunkX
	}
	return chunkIDs
}

// DEBT This function has a bug and doesn't split correctly for world edges -48,-48,48,48 and shard size 2x2.
func splitDimensionShards(dimID uuid.UUID, dimName string, edges level.Edges, shardX, shardZ int64) map[ShardID]startMessage {
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

func (sh *Sharder) signal(state control.ComponentState, err error) {
	sh.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.SHARDER,
		State:     state,
		Err:       err,
	}
}
