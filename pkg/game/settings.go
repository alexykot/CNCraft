//go:generate stringer -type=Dimension settings.go
//go:generate stringer -type=Difficulty settings.go

package game

import (
	"fmt"
	"time"

	"github.com/alexykot/cncraft/pkg/buffer"
)

type Difficulty byte

type Tick int64

const TickSpeed = time.Millisecond * 50

func (d Tick) AsTime() time.Time {
	return time.Unix(0, int64(d))
}

const (
	Peaceful Difficulty = iota
	Easy
	Normal
	Hard
)

func (d *Difficulty) Pull(reader *buffer.Buffer) error {
	id := reader.PullByte()

	switch int(id) {
	case 0, 1, 2, 3:
		*d = Difficulty(id)
	default:
		return fmt.Errorf("no difficulty for id %d", id)
	}
	return nil
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
	}
	return ""
}
