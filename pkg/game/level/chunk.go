package level

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

type HeightMap struct {
	MotionBlocking []uint8 `nbt:"motion_blocking"`
	WorldSurface   []uint8 `nbt:"world_surface"` // purpose unknown, left empty
}

type ChunkID string

// 16*16*255 blocks column
type Chunk interface {
	ID() ChunkID

	X() int64 // block coordinates of the lowest X block in the chunk, NOT the chunk coord (divided by 16 rounded down)
	Z() int64 // block coordinates of the lowest Z block in the chunk, NOT the chunk coord (divided by 16 rounded down)

	Sections() []Section
	HeightMap() HeightMap

	// supports values /**/x:[0:15] y:[0:255] z: [0:15]
	// GetBlock(x, y, z int) Block
}

// DEBT no performance considerations applied here yet. Likely will have to be redesigned for RAM/CPU efficiency.

type chunk struct {
	// Shown below is the 0-0 chunk in the top right quarter of the coord grid.
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
	// y coord makes no sense for chunk as the chunk always occupies full height of the world.
	x int64
	z int64

	sections []Section
}

func NewChunk() Chunk {
	return &chunk{}
}

// NewDefaultChunk creates a flatworld hardcoded chunk.
func NewDefaultChunk(x, z int64) Chunk {
	return &chunk{
		x: x,
		z: z,
		sections: []Section{
			NewDefaultSection(0),
		},
	}
}

func (c *chunk) ID() ChunkID { return ChunkID(fmt.Sprintf("chunk-%d-%d", c.x, c.z)) }
func (c *chunk) X() int64    { return c.x }
func (c *chunk) Z() int64    { return c.z }

func (c *chunk) Sections() []Section {
	return c.sections
}

func (c *chunk) HeightMap() HeightMap {
	var sectionIndex int

	// find the topmost non-empty section to start from
	for index, chunkSection := range c.sections {
		if chunkSection != nil {
			sectionIndex = index
		}
	}

	heights := [ChunkX][ChunkZ]uint8{}

	var heightsFound int
	// walk through sections down
	for ; sectionIndex >= 0; sectionIndex-- {
		// walk every column in the section
		for x := 0; x < ChunkX; x++ {
			for z := 0; z < ChunkZ; z++ {
				if heights[x][z] != 0 { // skip if the given column already has a height
					continue
				}

				// scan column top-down and look for non-air blocks
				for y := SectionY; y > 0; y-- {
					sectionBlock := c.sections[sectionIndex].GetBlock(x, y-1, z)
					// DEBT check for solid block rather than non-air, start at https://minecraft.gamepedia.com/Solid_block
					if sectionBlock.ID() != blocks.Air {
						heights[x][z] = uint8(y + sectionIndex*SectionY)
						heightsFound++
						break
					}
				}
			}
		}

		// if all non-air heights are founds - stop scanning
		if heightsFound == ChunkX*ChunkZ {
			break
		}
	}

	heightMap := HeightMap{
		MotionBlocking: make([]uint8, ChunkX*ChunkZ, ChunkX*ChunkZ),
	}
	var i int
	for _, zHeights := range heights {
		for _, height := range zHeights {
			heightMap.MotionBlocking[i] = height
		}
	}

	return heightMap
}
