package players

import (
	"sync"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/player"
)

type Player struct {
	ID        uuid.UUID
	ConnID    uuid.UUID
	PC        entities.PlayerCharacter
	Username  string
	Settings  *player.Settings
	Abilities *player.Abilities
	State     *player.State

	mu sync.Mutex
}

func (p *Player) GetState() *player.State {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.State
}

func (p *Player) SetState(state *player.State) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State = state
}

func (p *Player) GetSettings() *player.Settings {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Settings
}

func (p *Player) SetSettings(settings *player.Settings) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Settings = settings
}

func (p *Player) GetLocation() data.Location {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.State.Location
}

func (p *Player) SetPosition(position data.PositionF) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State.Location.PositionF = position
}

func (p *Player) SetRotation(rotation data.RotationF) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State.Location.RotationF = rotation
}

func (p *Player) SetOnGround(onGround bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State.Location.OnGround = onGround
}
