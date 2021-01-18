package level

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

type HeightMap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"` // purpose unknown, left empty
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
	heights := c.findHeights()
	heightMap := HeightMap{
		MotionBlocking: c.compactHeights(heights),
		// WorldSurface purpose is unknown, but in the Notchian packet it's contents is same as of MotionBlocking.
		WorldSurface: c.compactHeights(heights),
	}
	return heightMap
}

func (c *chunk) findHeights() [ChunkX][ChunkZ]uint8 {
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

		// if all non-air heights are found - stop scanning
		if heightsFound == ChunkX*ChunkZ {
			break
		}
	}
	return heights
}

func (c *chunk) compactHeights(heights [ChunkX][ChunkZ]uint8) []int64 {
	const tupleSize = 7
	const bitsPerHeight = 9

	var i int
	var resI int
	var long uint64
	uRes := make([]uint64, 37, 37)
	for x, zHeignts := range heights {
		for z := range zHeignts {
			long = long << bitsPerHeight
			long = long | uint64(heights[x][z])
			i++
			if i == tupleSize {
				uRes[resI] = long << 1
				long = 0
				resI++
				i = 0
			}
		}
	}

	if i != 0 {
		uRes[resI] = long << 1
	}

	res := make([]int64, 37, 37)
	for i, _ := range uRes {
		res[i] = int64(uRes[i])
	}
	return res
}
