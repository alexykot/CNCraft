package world

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
)

// DEBT this is wholly hardcoded to single world per server. Will likely have to redesign for multitenancy if
//  things like lobby will be needed. Not a today's problem, will only be needed if project takes off.
//  Or maybe lobby can be supported as a separate dimension (level) within the same world? Need more data.

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
		}
		binary.LittleEndian.PutUint32(defaultWorld.Seed, rand.Uint32())
		defaultWorld.SeedHash = sha256.Sum256(defaultWorld.Seed)

		defaultWorld.Levels = make(map[string]level.Level)
		defaultWorld.Levels[game.Overworld.String()] = level.GetDefaultLevel()
	}

	return defaultWorld
}
