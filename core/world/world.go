package world

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"github.com/alexykot/cncraft/pkg/protocol/tags"

	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
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

	// TODO not clear where this should be saved and come from. Maybe it should be hardcoded as server defaults and
	//  not saved with the world at all.
	DimensionCodec tags.DimensionCodec
	Dimension      tags.Dimension

	Levels map[string]level.Level
}

var defaultWorld *World // TODO replace this with actual world loading from persistence

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
		defaultWorld.Levels[game.Overworld.String()] = level.GetDefaultLevel()
	}

	return defaultWorld
}
