package level

import (
	"github.com/alexykot/cncraft/pkg/buffer"
)

type ChunkID string

type Chunk interface {
	buffer.BPush

	ID() ChunkID

	ChunkX() int
	ChunkZ() int

	Slices() []Slice

	// supports values y:[0:15]
	GetSlice(y int) Slice

	// supports values x:[0:15] y:[0:255] z: [0:15]
	GetBlock(x, y, z int) Block

	//HeightMapNbtCompound() *tags.NbtCompound
}

// DEBT no performance considerations applied here yet. Likely will have to be redesigned for RAM/CPU efficiency.

type chunk struct {
	blocks map[int]map[int]map[int]Block // x,y,z coords
}
