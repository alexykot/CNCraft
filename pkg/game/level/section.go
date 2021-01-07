package level

import (
	"encoding/binary"

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
}

type section struct {
	// DEBT will need to store compacted paletted block map and unpack on request to save RAM
	blocks [16][16][16]Block // x,z,y block coordinates
	index  uint8
}

func (s *section) Index() int { return int(s.index) }

func (s *section) GetBlock(x, y, z int) Block {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return nil
	}

	return s.blocks[x][z][y]
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

	compactData := s.makeBlockData(bpb, palette)
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

func (s *section) makeBlockData(bpb uint8, palette []blocks.BlockID) []uint64 {
	switch bpb {
	case 4:
		return compactBlocksBpb4(palette, s.blocks)
	case 5:
		return compactBlocksBpb5(palette, s.blocks)
	case 6:
		return compactBlocksBpb6(palette, s.blocks)
	case 7:
		return compactBlocksBpb7(palette, s.blocks)
	case 8:
		return compactBlocksBpb8(palette, s.blocks)
	case 14:
		return compactBlocksBpb14(s.blocks)
	}

	return nil
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

// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000          Long, 8 bytes
// 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000  bpb=4, 16 blocks
// 0000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000 00000     bpb=5, 12 blocks
// 0000 000000 000000 000000 000000 000000 000000 000000 000000 000000 000000       bpb=6, 10 blocks
// 0 0000000 0000000 0000000 0000000 0000000 0000000 0000000 0000000 0000000        bpb=7, 9 blocks
// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000          bpb=8, 8 blocks
// 00000000 00000000000000 00000000000000 00000000000000 00000000000000             bpb=14, 4 blocks

func compactBlocksBpb4(palette []blocks.BlockID, blockList [16][16][16]Block) []uint64 {
	var compactData []uint64
	blockTuple := [16]Block{}
	paletteIndices := [16]uint8{}
	var i uint8
	for _, zBlocks := range blockList {
		for _, yBlocks := range zBlocks {
			for _, block := range yBlocks {
				blockTuple[i] = block
				for paletteIndex, blockID := range palette {
					if blockID == block.ID() {
						paletteIndices[i] = uint8(paletteIndex)
						break // DEBT this does not handle case when blockID is not found in the palette
					}
				}

				if i == 15 {
					dataTuple := make([]byte, 8, 8)
					dataTuple[0] = (paletteIndices[0] << 4) | paletteIndices[1]
					dataTuple[1] = (paletteIndices[2] << 4) | paletteIndices[3]
					dataTuple[2] = (paletteIndices[3] << 4) | paletteIndices[5]
					dataTuple[3] = (paletteIndices[4] << 4) | paletteIndices[7]
					dataTuple[4] = (paletteIndices[6] << 4) | paletteIndices[9]
					dataTuple[5] = (paletteIndices[8] << 4) | paletteIndices[11]
					dataTuple[6] = (paletteIndices[10] << 4) | paletteIndices[13]
					dataTuple[7] = (paletteIndices[12] << 4) | paletteIndices[15]

					long, _ := binary.Uvarint(dataTuple)
					compactData = append(compactData, long)

					i = 0
				} else {
					i++
				}
			}
		}
	}
	return compactData
}

func compactBlocksBpb5(palette []blocks.BlockID, blocks [16][16][16]Block) []uint64 {
	return nil
}

func compactBlocksBpb6(palette []blocks.BlockID, blocks [16][16][16]Block) []uint64 {
	return nil
}

func compactBlocksBpb7(palette []blocks.BlockID, blocks [16][16][16]Block) []uint64 {
	return nil
}

func compactBlocksBpb8(palette []blocks.BlockID, blocks [16][16][16]Block) []uint64 {
	return nil
}

func compactBlocksBpb14(blocks [16][16][16]Block) []uint64 {
	return nil
}
