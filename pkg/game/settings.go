package game

import "fmt"

type Difficulty byte

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
	Nether    = -1
	Overworld = 0
	TheEnd    = 1
)

type Gamemode uint8

const (
	Survival Gamemode = iota
	Creative
	Adventure
	Spectator
)

func (g Gamemode) Encoded(hardcore bool) byte {

	bit := 0
	if hardcore {
		bit = 0x8
	}

	return byte(g) | byte(bit)
}

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

var typeToName = map[WorldType]string{
	WorldDefault:     "default",
	WorldFlat:        "flat",
	WorldLargebiomes: "largeBiomes",
	WorldAmplified:   "amplified",
	WorldCustomized:  "customized",
	WorldBuffet:      "buffet",
	WorldDefault11:   "default_1_1",
}

func (l WorldType) String() string {
	return typeToName[l]
}
