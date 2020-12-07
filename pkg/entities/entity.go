package entities

import (
	"github.com/golangmc/minecraft-server/pkg/bus"
)

type Sender interface {
	Name() string

	SendMessage(message ...bus.Envelope)
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

type Player interface {
	EntityLiving

	Online(bool)
	IsOnline() bool
}
