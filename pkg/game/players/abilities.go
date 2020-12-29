package players

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/mask"
)

// PlayerMaxHealth is max health player can have.
const PlayerMaxHealth = 20.0

type PlayerAbilities struct {
	mask.Masking

	Invulnerable bool
	Flying       bool
	AllowFlight  bool
	InstantBuild bool
}

func (p *PlayerAbilities) Push(writer buffer.B) {
	flags := byte(0)

	p.Set(&flags, 0x01, p.Invulnerable)
	p.Set(&flags, 0x02, p.Flying)
	p.Set(&flags, 0x04, p.AllowFlight)
	p.Set(&flags, 0x08, p.InstantBuild)

	writer.PushByt(flags)
}

func (p *PlayerAbilities) Pull(reader buffer.B) {
	flags := reader.PullByt()

	p.Invulnerable = p.Has(flags, 0x01)
	p.Flying = p.Has(flags, 0x02)
	p.AllowFlight = p.Has(flags, 0x04)
	p.InstantBuild = p.Has(flags, 0x08)
}

type PlayerSettings struct {
	ViewDistance int32
	FlyingSpeed  float32
	FoVModifier  float32 // Field of View Modifier
}
