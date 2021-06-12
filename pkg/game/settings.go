//go:generate stringer -type=Dimension settings.go

package game

import (
	"fmt"
	"time"
)

type Difficulty byte

type Tick int64

const TickSpeed = time.Millisecond * 50

const (
	Peaceful Difficulty = iota
	Easy
	Normal
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Peaceful:
		return "Peaceful"
	case Easy:
		return "Easy"
	case Normal:
		return "Normal"
	case Hard:
		return "Hard"
	default:
		panic(fmt.Errorf("no difficulty for id %d", byte(d)))
	}
}

func ValueOfDifficulty(d Difficulty) byte {
	return byte(d)
}

func DifficultyValueOf(id byte) Difficulty {
	switch id {
	case 0:
		return Peaceful
	case 1:
		return Easy
	case 2:
		return Normal
	case 3:
		return Hard
	default:
		panic(fmt.Errorf("no difficulty for id %d", id))
	}
}

type Dimension int

const (
	Nether    Dimension = -1
	Overworld Dimension = 0
	TheEnd    Dimension = 1
)

type Gamemode uint8

const (
	Survival Gamemode = iota
	Creative
	Adventure
	Spectator
)

type Coreness bool

const (
	Hardcore Coreness = true
	Softcore Coreness = false
)

type WorldType int

const (
	WorldDefault WorldType = iota
	WorldFlat
	WorldLargebiomes
	WorldAmplified
	WorldCustomized
	WorldBuffet
	WorldDefault11
)

func (l WorldType) String() string {
	switch l {
	case WorldDefault:
		return "default"
	case WorldFlat:
		return "flat"
	case WorldLargebiomes:
		return "largeBiomes"
	case WorldAmplified:
		return "amplified"
	case WorldCustomized:
		return "customized"
	case WorldBuffet:
		return "buffet"
	case WorldDefault11:
		return "default_1_1"
	}
	return ""
}
