package level

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/protocol/objects"
)

// Section - 16*16*16 blocks cubic section, part of the chunk
type Section interface {
	buffer.BPush

	// Index - position in the chunk, 0 to 15
	Index() int

	// GetBlock - supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int64) Block

	// SetBlock - supports values x:[0:15] y:[0:15] z: [0:15]
	SetBlock(x, y, z int64, block Block) error
}

func NewSection(blocks BlockArr, index uint8) Section {
	return &section{
		blocks: blocks,
		index:  index,
	}
}

type section struct {
	// DEBT will need to store compacted paletted block map and unpack on request to save RAM
	blocks BlockArr // x,z,y block coordinates
	index  uint8
}

func (s *section) Index() int { return int(s.index) }

func (s *section) GetBlock(x, y, z int64) Block {
	if x < 0 || x > SectionX-1 || y < 0 || y > SectionX-1 || z < 0 || z > SectionZ-1 {
		return nil
	}

	return s.blocks[x][z][y]
}

func (s *section) SetBlock(x, y, z int64, b Block) error {
	if x < 0 || x > SectionX-1 || y < 0 || y > SectionY-1 || z < 0 || z > SectionZ-1 {
		return fmt.Errorf("block coords x,y,z: %d,%d,%d out of range", x, y, z)
	}
	s.blocks[x][z][y] = b
	return nil
}

func (s *section) Push(writer *buffer.Buffer) {
	// push count of non-air blocks
	writer.PushInt16(SectionY * SectionZ * SectionX) // DEBT this does not consider non-air blocks yet

	palette := s.makePalette()
	bpb := bitsPerBlock(len(palette))

	// push bits-per-block value
	writer.PushByte(bpb)

	// push palette only if 8 or less bits per block, otherwise use global palette directly
	if bpb < 9 {
		writer.PushVarInt(int32(len(palette)))
		for _, blockID := range palette {
			writer.PushVarInt(int32(blockID))
		}
	}

	compactData, err := s.makeBlockData(bpb, palette)
	if err != nil {
		// DEBT update buffer interface to support errors.
	}
	writer.PushVarInt(int32(len(compactData)))
	for _, long := range compactData {
		writer.PushUint64(long)
	}
}

func (s *section) makePalette() []objects.BlockID {
	paletteMap := make(map[objects.BlockID]struct{})

	for y := 0; y < SectionY; y++ {
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				paletteMap[s.blocks[y][z][x].ID()] = struct{}{}
			}
		}
	}

	palette := make([]objects.BlockID, len(paletteMap), len(paletteMap))
	var i int
	for blockId := range paletteMap {
		palette[i] = blockId
		i++
	}

	return palette
}

func bitsPerBlock(len int) uint8 {
	var bpb uint8
	palleteSize := 1

	for palleteSize < len {
		palleteSize = palleteSize << 1
		bpb++
	}

	if bpb < 4 {
		bpb = 4
	} else if bpb > 8 {
		bpb = 14 // DEBT reconsider this if we ever support plugins
	}

	return bpb
}

// makeBlockData implements the palette-based compaction algorithm for section blocks array.
// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000          Long, 8 bytes
// 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000  bpb=4, 16 blocks
// 0000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000     bpb=5, 12 blocks
// 0000 000000 000000 000000 000000 000000 000000 000000 000000 000000 000000       bpb=6, 10 blocks
// 0 0000000 0000000 0000000 0000000 0000000 0000000 0000000 0000000 0000000        bpb=7, 9 blocks
// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000          bpb=8, 8 blocks
// 00000000 00000000000000 00000000000000 00000000000000 00000000000000             bpb=14, 4 blocks
func (s *section) makeBlockData(bpb uint8, palette []objects.BlockID) ([]uint64, error) {
	var useGlobalPalette bool

	var tupleSize uint8
	switch bpb {
	case 4:
		tupleSize = 16
	case 5:
		tupleSize = 12
	case 6:
		tupleSize = 10
	case 7:
		tupleSize = 9
	case 8:
		tupleSize = 8
	case 14:
		tupleSize = 4
		useGlobalPalette = true
	default:
		return nil, fmt.Errorf("bpb value %d not supported", bpb)
	}

	var compactData []uint64
	blocksTuple := make([]Block, tupleSize, tupleSize)
	paletteIndices := make([]uint32, tupleSize, tupleSize)
	var i uint8
	for y := 0; y < SectionY; y++ {
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				blocksTuple[i] = s.blocks[y][z][x]

				if useGlobalPalette {
					paletteIndices[i] = uint32(s.blocks[y][z][x].ID())
				} else {
					var found bool
					for paletteIndex, blockID := range palette {
						if blockID == s.blocks[y][z][x].ID() {
							paletteIndices[i] = uint32(paletteIndex)
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("block ID %d not found in palette", s.blocks[y][z][x].ID())
					}
				}

				if i == tupleSize-1 {
					var long uint64
					for _, paletteIndex := range paletteIndices {
						long = long << bpb
						long = long | uint64(paletteIndex)
					}
					compactData = append(compactData, long)

					i = 0
				} else {
					i++
				}
			}
		}
	}

	if i > 0 { // append any remaining blocks that did not fill a full long
		var long uint64
		for _, paletteIndex := range paletteIndices {
			long = long << bpb
			long = long | uint64(paletteIndex)
		}
		compactData = append(compactData, long)
	}

	return compactData, nil
}
