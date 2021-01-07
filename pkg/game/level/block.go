package level

import (
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

type Block interface {
	ID() blocks.BlockID
}

type block struct {
	pos data.PositionI
	id  blocks.BlockID
}

func NewBlock(id blocks.BlockID) Block {
	return &block{id: id}
}

func (b *block) ID() blocks.BlockID {
	return b.id
}
