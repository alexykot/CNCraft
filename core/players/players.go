package players

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/players"
)

// Tally holds the map of all players logged into this server.
type Tally struct {
	log *zap.Logger
	ps  nats.PubSub
	mu  sync.Mutex

	players map[uuid.UUID]Player // ID of the user always matches ID of the corresponding connection.
}

type Player struct {
	PC       entities.PlayerCharacter
	Username string
	Settings settings
}

// TODO implement proper player settings
type settings struct {
	ViewDistance int32
}

func NewTally(log *zap.Logger, ps nats.PubSub) *Tally {
	return &Tally{
		log: log,
		ps:  ps,

		players: make(map[uuid.UUID]Player),
	}
}

// RegisterHandlers creates subscriptions to all relevant global subjects.
func (u *Tally) RegisterHandlers() error {
	if err := u.ps.Subscribe(subj.MkPlayerJoined(), u.playerJoinedHandler); err != nil {
		return fmt.Errorf("failed to subscribe for joined players: %w", err)
	}

	u.log.Info("global player handlers registered")
	return nil
}

func (u *Tally) AddPlayer(userID uuid.UUID, username string) *Player {
	u.mu.Lock()
	defer u.mu.Lock()

	if existing, ok := u.players[userID]; ok {
		u.log.Error("player already exists", zap.String("id", userID.String()),
			zap.String("existing", existing.Username), zap.String("new", username))
		return nil
	}
	p := Player{
		PC:       entities.NewPC(username, players.PlayerMaxHealth),
		Username: username,
		Settings: settings{ // TODO all settings are hardcoded until properly implemented
			ViewDistance: 7,
		},
	}

	u.players[userID] = p
	return &p
}

func (u *Tally) playerJoinedHandler(lope *envelope.E) {
	//joined := lope.GetPlayerJoined()
	//if joined == nil {
	//	u.log.Error("failed to parse envelope - no JoinedPlayer inside", zap.Any("envelope", lope))
	//	return
	//}
	//
	//userId, err := uuid.Parse(joined.Id)
	//if err != nil {
	//	u.log.Error("failed to parse user ID as UUID", zap.Any("id", joined.Id))
	//	return
	//}

	// TODO implement this
}

func (u *Tally) GetPlayer(userID uuid.UUID) (Player, bool) {
	player, ok := u.players[userID]
	return player, ok
}
