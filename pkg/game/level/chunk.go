package level

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

// SectionRepo - persistence-aware interface for loading and saving sections.
type SectionRepo interface {
	// LoadSection - loads section from persistence into memory
	LoadSection(x, z int64, index uint8) (Section, error)

	// SaveSection - saves section state
	// DEBT This does not allow for differential updates, will be ineffective to save whole section every time
	//  a block in the section is updated. Will need to be optimised for diff updates only, eventually.
	SaveSection(Section) error
}

type heightMap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
	WorldSurface   []int64 `nbt:"WORLD_SURFACE"` // purpose unknown, left empty
}

type ChunkID string

func MkChunkID(x, z int64) ChunkID {
	return ChunkID(fmt.Sprintf("chunk.%d.%d", x, z))
}

func XZFromChunkID(ID ChunkID) (x, z int64) {
	pieces := strings.Split(string(ID), ".")
	if len(pieces) != 3 {
		panic(fmt.Sprintf("invalid chunkID `%s`", ID))
	}
	xInt, err := strconv.Atoi(pieces[1])
	if err != nil {
		panic(fmt.Sprintf("invalid chunkID `%s`", ID))
	}
	zInt, err := strconv.Atoi(pieces[2])
	if err != nil {
		panic(fmt.Sprintf("invalid chunkID `%s`", ID))
	}
	return int64(xInt), int64(zInt)
}

// Chunk - 16*16*255 blocks column
type Chunk interface {
	ID() ChunkID

	X() int64 // block coordinates of the lowest X block in the chunk, NOT the chunk coord (divided by 16 rounded down)
	Z() int64 // block coordinates of the lowest Z block in the chunk, NOT the chunk coord (divided by 16 rounded down)

	Load(repo SectionRepo) error
	Unload()

	Sections() []Section
	HeightMap() heightMap

	// supports values x:[0:15] y:[0:255] z: [0:15]
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

// NewChunk creates new chunk (not loaded yet)
func NewChunk(x, z int64) Chunk {
	return &chunk{
		x: x,
		z: z,
	}
}

func (c *chunk) ID() ChunkID { return MkChunkID(c.x, c.z) }
func (c *chunk) X() int64    { return c.x }
func (c *chunk) Z() int64    { return c.z }

func (c *chunk) Sections() []Section {
	return c.sections
}

func (c *chunk) Load(repo SectionRepo) error {
	c.Unload()
	var err error

	// DEBT make this configurable for supporting taller worlds
	c.sections = make([]Section, 8, 8)

	if c.sections[0], err = repo.LoadSection(c.x, c.z, 0); err != nil {
		return fmt.Errorf("failed to load section %d: %w", 0, err)
	}

	// TODO when all sections are set, even with just Air blocks - something fails on the client.
	// if c.sections[1], err = repo.LoadSection(c.x, c.z, 1); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 1, err)
	// }
	// if c.sections[2], err = repo.LoadSection(c.x, c.z, 2); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 2, err)
	// }
	// if c.sections[3], err = repo.LoadSection(c.x, c.z, 3); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 3, err)
	// }
	// if c.sections[4], err = repo.LoadSection(c.x, c.z, 4); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 4, err)
	// }
	// if c.sections[5], err = repo.LoadSection(c.x, c.z, 5); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 5, err)
	// }
	// if c.sections[6], err = repo.LoadSection(c.x, c.z, 6); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 6, err)
	// }
	// if c.sections[7], err = repo.LoadSection(c.x, c.z, 7); err != nil {
	// 	return fmt.Errorf("failed to load section %d: %w", 7, err)
	// }

	return nil
}

func (c *chunk) Unload() {
	c.sections = nil // DEBT is this enough to unload section data from memory 🤔
}

func (c *chunk) HeightMap() heightMap {
	heights := c.findHeights()
	heightMap := heightMap{
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
						// DEBT for unknown reason Notchian server supplies height of "2" in the flatworld the
						//  solid block height is "4". Adjusting until figure out why.
						heights[x][z] = uint8(y + sectionIndex*SectionY - 2)
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
