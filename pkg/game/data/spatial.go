package data

import (
	"fmt"
	"math"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/mask"
)

type PositionI struct {
	X int64
	Y int64
	Z int64
}

type PositionF struct {
	X float64
	Y float64
	Z float64
}

type RotationF struct {
	Yaw   float32
	Pitch float32
}

type Location struct {
	PositionF
	RotationF
	OnGround bool
}

type Relativity struct {
	mask.Masking

	X bool
	Y bool
	Z bool

	Yaw   bool
	Pitch bool
}

func (p PositionF) ToInt() PositionI {
	return PositionI{
		X: int64(math.Round(p.X)),
		Y: int64(math.Round(p.Y)),
		Z: int64(math.Round(p.Z)),
	}
}

func (p PositionF) String() string {
	return fmt.Sprintf("%f:%f:%f", p.X, p.Y, p.Z)
}

func (p PositionI) ToFloat() PositionF {
	return PositionF{
		X: float64(p.X),
		Y: float64(p.Y),
		Z: float64(p.Z),
	}
}

func (p PositionI) String() string {
	return fmt.Sprintf("%d:%d:%d", p.X, p.Y, p.Z)
}

func (p PositionI) Pull(reader *buffer.Buffer) {
	val := reader.PullUint64()
	p.X = int64(val) >> 38
	p.Y = int64(val) & 0xFFF
	p.Z = int64(val) << 26 >> 38
}

func (p PositionI) Push(writer *buffer.Buffer) {
	writer.PushInt64(((p.X & 0x3FFFFFF) << 38) | ((p.Z & 0x3FFFFFF) << 12) | (p.Y & 0xFFF))
}

func (r *Relativity) Push(writer *buffer.Buffer) {
	flags := byte(0)

	r.Set(&flags, 0x01, r.X)
	r.Set(&flags, 0x02, r.Y)
	r.Set(&flags, 0x04, r.Z)

	// the fact that these are flipped deeply bothers me.
	r.Set(&flags, 0x08, r.Pitch)
	r.Set(&flags, 0x10, r.Yaw)

	writer.PushByte(flags)
}
