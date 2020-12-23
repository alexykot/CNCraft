package players

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game/entities"
)

type Tally struct {
	log *zap.Logger
	ps  nats.PubSub

	players map[uuid.UUID]Player // ID of the user always matches ID of the corresponding connection.
}

type Player struct {
	PC       entities.PlayerCharacter
	Username string
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
	playerJoinedHandler := func(lope *envelope.E) {
		log := u.log

		joined := lope.GetJoinedPlayer()
		if joined == nil {
			log.Error("failed to parse envelope - there is no player inside", zap.Any("envelope", lope))
			return
		}

		userId, err := uuid.Parse(joined.Id)
		if err != nil {
			log.Error("failed to parse user ID as UUID", zap.Any("id", joined.Id))
			return
		}

		u.players[userId] = Player{
			Username: joined.Username,
		}

		log.Debug("added new player", zap.String("username", joined.Username))
	}

	if err := u.ps.Subscribe(subj.MkJoinedPlayers(), playerJoinedHandler); err != nil {
		return fmt.Errorf("failed to subscribe for added players: %w", err)
	}

	u.log.Info("global player handlers registered")
	return nil
}

func (u *Tally) AddUnauthenticated(id uuid.UUID, username string) {
	u.players[id] = Player{
		Username: username,
	}
}
