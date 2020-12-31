package data

import (
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
}

type Relativity struct {
	mask.Masking

	X bool
	Y bool
	Z bool

	Yaw   bool
	Pitch bool
}

func (r *PositionI) Pull(reader buffer.B) {
	val := reader.PullUint64()
	r.X = int64(val) >> 38
	r.Y = int64(val) & 0xFFF
	r.Z = int64(val) << 26 >> 38
}

func (r *PositionI) Push(writer buffer.B) {
	writer.PushInt64(((r.X & 0x3FFFFFF) << 38) | ((r.Z & 0x3FFFFFF) << 12) | (r.Y & 0xFFF))
}

func (r *Relativity) Push(writer buffer.B) {
	flags := byte(0)

	r.Set(&flags, 0x01, r.X)
	r.Set(&flags, 0x02, r.Y)
	r.Set(&flags, 0x04, r.Z)

	// the fact that these are flipped deeply bothers me.
	r.Set(&flags, 0x08, r.Pitch)
	r.Set(&flags, 0x10, r.Yaw)

	writer.PushByte(flags)
}
