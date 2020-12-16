package entities

import (
	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/user"
)

type Sender interface {
	Name() string

	SendMessage(message ...nats.Envelope)
}

type Entity interface {
	Sender

	ID() int64
}

type EntityLiving interface {
	Entity

	GetHealth() float64
	SetHealth(health float64)
}

type PlayerCharacter interface {
	EntityLiving

	GetGameMode() game.GameMode
	SetGameMode(mode game.GameMode)

	Online(bool)
	IsOnline() bool

	GetProfile() *user.Profile
}
