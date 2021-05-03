package player

import (
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/items"
)

type State struct {
	Location  data.Location
	Inventory items.Inventory
}
