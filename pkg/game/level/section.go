package level

import (
	"fmt"

	buff "github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

// 16*16*16 blocks cubic section, part of the chunk
type Section interface {
	buff.BPush

	// position in the chunk, 0 to 15
	Index() int

	// supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int) Block

	// supports values x:[0:15] y:[0:15] z: [0:15]
	SetBlock(x, y, z int, block Block) error
}

type section struct {
	// DEBT will need to store compacted paletted block map and unpack on request to save RAM
	blocks [16][16][16]Block // x,z,y block coordinates
	index  uint8
}

func NewSection(index uint8) Section {
	return &section{
		blocks: [16][16][16]Block{},
		index:  index,
	}
}

func (s *section) Index() int { return int(s.index) }

func (s *section) GetBlock(x, y, z int) Block {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return nil
	}

	return s.blocks[x][z][y]
}

func (s *section) SetBlock(x, y, z int, b Block) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("block coords x,y,z: %d,%d,%d out of range", x, y, z)
	}
	s.blocks[x][z][y] = b
	return nil
}

func (s *section) Push(writer buff.B) {
	// push count of non-air blocks
	writer.PushInt16(4096) // DEBT this does not consider non-air blocks yet

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

func (s *section) makePalette() []blocks.BlockID {
	paletteMap := make(map[blocks.BlockID]struct{})

	for _, zBlocks := range s.blocks {
		for _, yBlocks := range zBlocks {
			for _, block := range yBlocks {
				paletteMap[block.ID()] = struct{}{}
			}
		}
	}

	palette := make([]blocks.BlockID, len(paletteMap), len(paletteMap))
	var i int
	for blockId := range paletteMap {
		palette[i] = blockId
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
func (s *section) makeBlockData(bpb uint8, palette []blocks.BlockID) ([]uint64, error) {
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
	for _, zBlocks := range s.blocks {
		for _, yBlocks := range zBlocks {
			for _, sectionBlock := range yBlocks {
				blocksTuple[i] = sectionBlock

				if useGlobalPalette {
					paletteIndices[i] = uint32(sectionBlock.ID())
				} else {
					var found bool
					for paletteIndex, blockID := range palette {
						if blockID == sectionBlock.ID() {
							paletteIndices[i] = uint32(paletteIndex)
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("block ID %d not found in palette", sectionBlock.ID())
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
