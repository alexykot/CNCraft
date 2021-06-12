package player

import (
	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/items"
)

type State struct {
	Dimension uuid.UUID
	Location  data.Location
	Inventory *items.Inventory
}
