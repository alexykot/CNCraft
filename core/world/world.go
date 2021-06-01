package world

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"go.uber.org/zap"

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

	// TODO not clear where this should be saved and come from.
	//  Maybe it should be hardcoded as server defaults and not saved with the world at all.
	DimensionCodec tags.DimensionCodec
	Dimension      tags.Dimension

	Levels map[string]level.Level

	repo *SectionRepo
	log  *zap.Logger
}

// NewWorld - creates world from persisted settigns. Does NOT load world data.
func NewWorld(id uuid.UUID, log *zap.Logger, db *sql.DB) (*World, error) {
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
			DimensionCodec:     vanillaDimentionsCodec,
			Dimension:          vanillaDimentionsCodec.Dimensions.RegistryEntries[0].Element,
		}
		binary.LittleEndian.PutUint32(defaultWorld.Seed, rand.Uint32())
		defaultWorld.SeedHash = sha256.Sum256(defaultWorld.Seed)

		defaultWorld.Levels = make(map[string]level.Level)
		defaultWorld.Levels[game.Overworld.String()] = level.NewLevel(game.Overworld.String())
	}

	return defaultWorld
}

func (w *World) Load() error {
	if w.repo == nil {
		return fmt.Errorf("world section repo not initialised")
	}

	for name, worldLevel := range w.Levels {
		w.log.Debug(fmt.Sprintf("loading level %s", name))

		chunks := worldLevel.Chunks()
		for _, chunk := range chunks {
			if err := chunk.Load(w.repo); err != nil {
				return fmt.Errorf("failed to load chunk %s: %w", chunk.ID(), err)
			}
		}
	}

	w.log.Info(fmt.Sprintf("world `%s` loaded", w.Name))

	return nil
}
