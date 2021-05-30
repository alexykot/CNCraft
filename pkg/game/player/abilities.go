package player

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/mask"
)

// MaxHealth is max health player can have.
const MaxHealth = 20.0

type Abilities struct {
	mask.Masking

	Invulnerable bool
	Flying       bool
	AllowFlight  bool
	InstantBuild bool
}

func (p *Abilities) Push(writer *buffer.Buffer) {
	flags := byte(0)

	p.Set(&flags, 0x01, p.Invulnerable)
	p.Set(&flags, 0x02, p.Flying)
	p.Set(&flags, 0x04, p.AllowFlight)
	p.Set(&flags, 0x08, p.InstantBuild)

	writer.PushByte(flags)
}

func (p *Abilities) Pull(reader *buffer.Buffer) {
	flags := reader.PullByte()

	p.Invulnerable = p.Has(flags, 0x01)
	p.Flying = p.Has(flags, 0x02)
	p.AllowFlight = p.Has(flags, 0x04)
	p.InstantBuild = p.Has(flags, 0x08)
}
