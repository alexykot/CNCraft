package level

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/buffer"
)

type ChunkID string

// 16*16*255 blocks column
type Chunk interface {
	buffer.BPush

	ID() ChunkID

	ChunkX() int
	ChunkZ() int

	Sections() []Section

	// supports values /**/x:[0:15] y:[0:255] z: [0:15]
	GetBlock(x, y, z int) Block

	//HeightMapNbtCompound() *tags.NbtCompound
}

// DEBT no performance considerations applied here yet. Likely will have to be redesigned for RAM/CPU efficiency.

type chunk struct {
	// Shown below is the 0-0 chunk on the top right of the coord grid from zero.
	//                    ^
	//                  +z|
	//                    |
	//                    |--+
	//   -x              0|  |             +x
	//   ----------------------------------->
	//                    |
	//                    |
	//                    |
	//                  -z|
	//
	//
	// y coord makes no sense for chunk as the chunk always occupies whole height of the world.
	x int64
	z int64

	sections []Section
}

func (c *chunk) ID() ChunkID {
	return ChunkID(fmt.Sprintf("chunk-%d-%d", c.x, c.z))
}
