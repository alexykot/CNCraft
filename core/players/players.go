package players

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/player"
)

// Roster holds the map of all players logged into this server.
type Roster struct {
	log *zap.Logger
	ps  nats.PubSub
	mu  sync.Mutex

	players map[uuid.UUID]Player // ID of the user always matches ID of the corresponding connection.
}

type Player struct {
	PC        entities.PlayerCharacter
	Username  string
	Settings  player.Settings
	Abilities player.Abilities
	State     player.State

	mu sync.Mutex
}

func NewRoster(log *zap.Logger, ps nats.PubSub) *Roster {
	return &Roster{
		log: log,
		ps:  ps,

		players: make(map[uuid.UUID]Player),
	}
}

// RegisterHandlers creates subscriptions to all relevant global subjects.
func (r *Roster) RegisterHandlers() error {
	if err := r.ps.Subscribe(subj.MkPlayerJoined(), r.playerJoinedHandler); err != nil {
		return fmt.Errorf("failed to subscribe for joined users: %w", err)
	}

	r.log.Info("global player handlers registered")
	return nil
}

func (r *Roster) AddPlayer(userID uuid.UUID, username string) *Player {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.players[userID]; ok {
		r.log.Error("player already exists", zap.String("id", userID.String()),
			zap.String("existing", existing.Username), zap.String("new", username))
		return nil
	}

	// TODO whole user state hardcoded until persistence is properly implemented
	p := Player{
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		State: player.State{
			CurrentHotbarSlot: player.Slot0,
			CurrentLocation: data.Location{
				PositionF: data.PositionF{
					X: 15,
					Y: 32,
					Z: 15,
				},
			},
		},
	}

	r.players[userID] = p
	return &p
}

func (r *Roster) playerJoinedHandler(lope *envelope.E) {
	// joined := lope.GetPlayerJoined()
	// if joined == nil {
	//	r.log.Error("failed to parse envelope - no JoinedPlayer inside", zap.Any("envelope", lope))
	//	return
	// }
	//
	// userId, err := uuid.Parse(joined.Id)
	// if err != nil {
	//	r.log.Error("failed to parse user ID as UUID", zap.Any("id", joined.Id))
	//	return
	// }

	// TODO implement this
}

func (r *Roster) GetPlayer(userID uuid.UUID) (Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.players[userID]
	return p, ok
}

func (p *Player) GetState() player.State {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.State
}

func (p *Player) SetState(state player.State) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State = state
}

func (p *Player) GetSettings() player.Settings {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Settings
}

func (p *Player) SetSettings(settings player.Settings) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Settings = settings
}

func (p *Player) GetPosition() data.PositionF {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.State.CurrentLocation.PositionF
}

func (p *Player) SetPosition(position data.PositionF) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.State.CurrentLocation.PositionF = position
}
