package users

import (
	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/game/entities"
)

type Users struct {
	users map[uuid.UUID]User
}

type User struct {
	PC       entities.PlayerCharacter
	Username string
}
