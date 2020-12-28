package entities

import (
	"sync/atomic"

	"github.com/alexykot/cncraft/pkg/game"
)

var entityCounter int32

type Entity interface {
	// TODO add spatial data here?
	ID() int32
	Name() string
}

type Living interface {
	GetHealth() float32
	SetHealth(float32)
}

type PlayerCharacter interface {
	Entity
	Living

	GetGameMode() game.Gamemode
	SetGameMode(game.Gamemode)
}

type entity struct {
	id   int32
	name string
}

func (e *entity) ID() int32    { return e.id }
func (e *entity) Name() string { return e.name }

type living struct {
	maxHealth float32
	health    float32
}

func (l *living) GetHealth() float32 { return l.health }
func (l *living) SetHealth(health float32) {
	l.health = health
	if health > l.maxHealth {
		l.health = l.maxHealth
	}
}

type playerCharacter struct {
	entity
	living

	gamemode game.Gamemode
}

func (pc *playerCharacter) GetGameMode() game.Gamemode { return pc.gamemode }
func (pc *playerCharacter) SetGameMode(mode game.Gamemode) {
	switch mode {
	case game.Survival, game.Creative, game.Adventure, game.Spectator:
		pc.gamemode = mode
	}
}

func NewPC(name string, maxHealth float32) PlayerCharacter {
	// DEBT no idea how entity IDs are to be generated
	atomic.AddInt32(&entityCounter, 1)

	return &playerCharacter{
		entity: entity{
			id:   entityCounter,
			name: name,
		},
		living: living{
			maxHealth: maxHealth,
		},
		gamemode: 0,
	}
}
