package level

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/game/data/tags"
)

type Chunk interface {
	buffer.BufferPush

	ChunkX() int
	ChunkZ() int

	Slices() []Slice

	Level() Level

	// supports values y:[0:15]
	GetSlice(y int) Slice

	// supports values x:[0:15] y:[0:255] z: [0:15]
	GetBlock(x, y, z int) Block

	HeightMapNbtCompound() *tags.NbtCompound
}
