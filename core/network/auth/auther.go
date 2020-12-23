package auth

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// Auther handles user authentication before joining the server.
type Auther struct {
	log *zap.Logger
	ps  nats.PubSub

	pacFac protocol.PacketFactory

	users map[uuid.UUID]user
}

type userState int

const (
	newbie userState = iota // LoginStart just received
	joined                  // user fully authenticated and ready to join the game
)

type user struct {
	username string
	state    userState
}

func NewAuther(log *zap.Logger, ps nats.PubSub) *Auther {
	return &Auther{
		log: log,
		ps:  ps,

		users: make(map[uuid.UUID]user),
	}
}

func (a *Auther) AddNewbie(id uuid.UUID, username string) {
	a.users[id] = user{
		username: username,
		state:    newbie,
	}
}

func (a *Auther) RunMojangSessionAuth(sharedSecret []byte, username string) {

}
