package world

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
	"github.com/alexykot/cncraft/pkg/protocol/tags"
)

// DEBT the system is wholly hardcoded to single world per server. May want to redesign for multitenancy later.
//  Not a today's problem, will only be worth it if project takes off.
//  Lobby and similar ancillary places can be supported as a separate dimensions within single world.

// World holds details of the current world.
type World struct {
	Name     string
	Seed     []byte
	SeedHash [32]byte

	Coreness           game.Coreness
	Gamemode           game.Gamemode
	Type               game.WorldType
	Difficulty         game.Difficulty
	DifficultyIsLocked bool

	// TODO not clear where this should be saved and come from. And what it does.
	//  Maybe it should be hardcoded as server defaults and not saved with the world at all.
	NBTDimensionCodec tags.DimensionCodec
	NBTDimension      tags.Dimension

	StartDimension uuid.UUID
	Dimensions     map[uuid.UUID]level.Dimension

	repo *SectionRepo
	log  *zap.Logger
}

// NewWorld - creates world from persisted settigns. Does NOT load world data.
func NewWorld(log *zap.Logger, _ control.WorldConf, db *sql.DB) (*World, error) {
	world := GetDefaultWorld() // TODO load world starting settings from persistence.

	world.log = log
	world.repo = newRepo(log, db)

	return world, nil
}

// TODO replace this with actual world loading from persistence
var defaultWorld *World

func GetDefaultWorld() *World {
	if defaultWorld == nil {
		defaultWorld = &World{
			Name:               "Default World",
			Coreness:           game.Softcore,
			Type:               game.WorldFlat,
			Gamemode:           game.Survival,
			Difficulty:         game.Peaceful,
			DifficultyIsLocked: true,
			Seed:               make([]byte, 4, 4),
			StartDimension:     uuid.NewSHA1(uuid.UUID{}, []byte(game.Overworld.String())),
			Dimensions:         make(map[uuid.UUID]level.Dimension),
			NBTDimensionCodec:  vanillaDimentionsCodec,
			NBTDimension:       vanillaDimentionsCodec.Dimensions.RegistryEntries[0].Element,
		}
		binary.LittleEndian.PutUint32(defaultWorld.Seed, rand.Uint32())
		defaultWorld.SeedHash = sha256.Sum256(defaultWorld.Seed)

		defaultWorld.Dimensions[defaultWorld.StartDimension] = level.NewDimension(game.Overworld.String())
	}

	return defaultWorld
}

func (w *World) Load(_ context.Context, ctrlChan chan control.Command) {
	for name, worldDim := range w.Dimensions {
		w.log.Debug(fmt.Sprintf("loading level %s", name))

		chunks := worldDim.Chunks()
		for _, chunk := range chunks {
			if err := chunk.Load(w.repo); err != nil {
				// World does not have any async loops, so does not need to signal readiness, it's ready as soon
				// as it's loaded, and has no internal components that would need to be stopped.
				// But it can fail while loading and that needs to be signalled.
				ctrlChan <- control.Command{
					Signal:    control.COMPONENT,
					Component: control.WORLD,
					State:     control.FAILED,
					Err:       fmt.Errorf("failed to load world: failed to load chunk %s: %w", chunk.ID(), err),
				}
				return
			}
		}
	}

	w.log.Info(fmt.Sprintf("world `%s` loaded", w.Name))
}

func (w *World) getChunk(dimensionID uuid.UUID, chunkID level.ChunkID) (level.Chunk, error) {
	dim, ok := w.Dimensions[dimensionID]
	if !ok {
		return nil, fmt.Errorf("dimension %s not found", dimensionID.String())
	}

	chunk, ok := dim.GetChunk(chunkID)
	if !ok {
		return nil, fmt.Errorf("chunk %s not found", chunkID.String())
	}
	return chunk, nil
}
