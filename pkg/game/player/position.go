package player

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/mask"
)

type Relativity struct {
	mask.Masking

	X bool
	Y bool
	Z bool

	AxisX bool
	AxisY bool
}

func (r *Relativity) Push(writer buffer.B) {
	flags := byte(0)

	r.Set(&flags, 0x01, r.X)
	r.Set(&flags, 0x02, r.Y)
	r.Set(&flags, 0x04, r.Z)

	// the fact that these are flipped deeply bothers me.
	r.Set(&flags, 0x08, r.AxisY)
	r.Set(&flags, 0x10, r.AxisX)

	writer.PushByt(flags)
}
