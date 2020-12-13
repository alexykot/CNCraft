package player

import (
	"github.com/alexykot/cncraft/pkg/buffers"
	"github.com/alexykot/cncraft/pkg/mask"
)

type PlayerAbilities struct {
	mask.Masking

	Invulnerable bool
	Flying       bool
	AllowFlight  bool
	InstantBuild bool // creative??
}

func (p *PlayerAbilities) Push(writer buffers.Buffer) {
	flags := byte(0)

	p.Set(&flags, 0x01, p.Invulnerable)
	p.Set(&flags, 0x02, p.Flying)
	p.Set(&flags, 0x04, p.AllowFlight)
	p.Set(&flags, 0x08, p.InstantBuild)

	writer.PushByt(flags)
}

func (p *PlayerAbilities) Pull(reader buffers.Buffer) {
	flags := reader.PullByt()

	p.Invulnerable = p.Has(flags, 0x01)
	p.Flying = p.Has(flags, 0x02)
	p.AllowFlight = p.Has(flags, 0x04)
	p.InstantBuild = p.Has(flags, 0x08)
}
