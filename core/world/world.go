package world

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"github.com/alexykot/cncraft/pkg/game"
)

// DEBT this is wholly hardcoded to single world per server. Will likely have to redesign for multitenancy if
//  things like lobby will be needed. Not a today's problem, will only be needed if project takes off.

// World holds details of the current world.
type World struct {
	Coreness           game.Coreness
	Gamemode           game.Gamemode
	Type               game.WorldType
	Difficulty         game.Difficulty
	DifficultyIsLocked bool
	Name               string
	Seed               []byte
	SeedHash           [32]byte
}

func GetWorld() *World {
	return getDefaultWorld() // replace this with actual loading from persistence
}

func getDefaultWorld() *World {
	defaultWorld := &World{
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

	return defaultWorld
}
