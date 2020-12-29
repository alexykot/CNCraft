package users

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/player"
)

// Roster holds the map of all players logged into this server.
type Roster struct {
	log *zap.Logger
	ps  nats.PubSub
	mu  sync.Mutex

	users map[uuid.UUID]User // ID of the user always matches ID of the corresponding connection.
}

type User struct {
	PC        entities.PlayerCharacter
	Username  string
	Settings  player.Settings
	Abilities player.Abilities
	State     player.State
}

func NewRoster(log *zap.Logger, ps nats.PubSub) *Roster {
	return &Roster{
		log: log,
		ps:  ps,

		users: make(map[uuid.UUID]User),
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

func (r *Roster) AddPlayer(userID uuid.UUID, username string) *User {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.users[userID]; ok {
		r.log.Error("player already exists", zap.String("id", userID.String()),
			zap.String("existing", existing.Username), zap.String("new", username))
		return nil
	}

	// TODO whole user state hardcoded until persistence is properly implemented
	p := User{
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		State: player.State{CurrentHotbarSlot: player.Slot0},
	}

	r.users[userID] = p
	return &p
}

func (r *Roster) playerJoinedHandler(lope *envelope.E) {
	//joined := lope.GetPlayerJoined()
	//if joined == nil {
	//	r.log.Error("failed to parse envelope - no JoinedPlayer inside", zap.Any("envelope", lope))
	//	return
	//}
	//
	//userId, err := uuid.Parse(joined.Id)
	//if err != nil {
	//	r.log.Error("failed to parse user ID as UUID", zap.Any("id", joined.Id))
	//	return
	//}

	// TODO implement this
}

func (r *Roster) GetPlayer(userID uuid.UUID) (User, bool) {
	p, ok := r.users[userID]
	return p, ok
}

func (r *Roster) GetPlayerState(userID uuid.UUID) player.State {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[userID]; !ok {
		return player.State{}
	}

	return r.users[userID].State
}

func (r *Roster) SetPlayerState(userID uuid.UUID, state player.State) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.users[userID]
	if !ok {
		return
	}
	p.State = state
	r.users[userID] = p
}

func (r *Roster) GetPlayerSettings(userID uuid.UUID) player.Settings {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[userID]; !ok {
		return player.Settings{}
	}

	return r.users[userID].Settings
}

func (r *Roster) SetPlayerSettings(userID uuid.UUID, settings player.Settings) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.users[userID]
	if !ok {
		return
	}
	p.Settings = settings
	r.users[userID] = p
}
