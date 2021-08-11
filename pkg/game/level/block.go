package level

import (
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/protocol/objects"
)

type Block interface {
	ID() objects.BlockID
}

type block struct {
	pos data.PositionI
	id  objects.BlockID
}

func NewBlock(id objects.BlockID) Block {
	return &block{id: id}
}

func (b *block) ID() objects.BlockID {
	return b.id
}
