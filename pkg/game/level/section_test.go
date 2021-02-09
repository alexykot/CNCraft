package level

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

func TestBitsPerBlock(t *testing.T) {
	assert.Equal(t, uint8(4), bitsPerBlock(1))
	assert.Equal(t, uint8(4), bitsPerBlock(16))

	assert.Equal(t, uint8(5), bitsPerBlock(17))
	assert.Equal(t, uint8(5), bitsPerBlock(32))

	assert.Equal(t, uint8(6), bitsPerBlock(33))
	assert.Equal(t, uint8(6), bitsPerBlock(64))

	assert.Equal(t, uint8(7), bitsPerBlock(65))
	assert.Equal(t, uint8(7), bitsPerBlock(128))

	assert.Equal(t, uint8(8), bitsPerBlock(129))
	assert.Equal(t, uint8(8), bitsPerBlock(256))

	assert.Equal(t, uint8(14), bitsPerBlock(257))
	assert.Equal(t, uint8(14), bitsPerBlock(1024))
	assert.Equal(t, uint8(14), bitsPerBlock(4096))
}

func TestMakePalette(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Dirt)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Dirt)
				s.blocks[3][z][x] = NewBlock(blocks.Dirt)
				s.blocks[4][z][x] = NewBlock(blocks.Air)
				s.blocks[5][z][x] = NewBlock(blocks.Air)
				s.blocks[6][z][x] = NewBlock(blocks.Air)
				s.blocks[7][z][x] = NewBlock(blocks.Air)
				s.blocks[8][z][x] = NewBlock(blocks.Air)
				s.blocks[9][z][x] = NewBlock(blocks.Air)
				s.blocks[10][z][x] = NewBlock(blocks.Air)
				s.blocks[11][z][x] = NewBlock(blocks.Air)
				s.blocks[12][z][x] = NewBlock(blocks.Air)
				s.blocks[13][z][x] = NewBlock(blocks.Air)
				s.blocks[14][z][x] = NewBlock(blocks.Air)
				s.blocks[15][z][x] = NewBlock(blocks.Air)
			}
		}

		palette := s.makePalette()
		require.Len(t, palette, 2)
		assert.Equal(t, blocks.Dirt, palette[0])
		assert.Equal(t, blocks.Air, palette[1])
	})
	t.Run("16_blocks", func(t *testing.T) {
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Air)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Stone)
				s.blocks[3][z][x] = NewBlock(blocks.Grass)
				s.blocks[4][z][x] = NewBlock(blocks.Granite)
				s.blocks[5][z][x] = NewBlock(blocks.Gravel)
				s.blocks[6][z][x] = NewBlock(blocks.Sand)
				s.blocks[7][z][x] = NewBlock(blocks.Sandstone)
				s.blocks[8][z][x] = NewBlock(blocks.Ice)
				s.blocks[9][z][x] = NewBlock(blocks.BlackWool)
				s.blocks[10][z][x] = NewBlock(blocks.WhiteWool)
				s.blocks[11][z][x] = NewBlock(blocks.PinkWool)
				s.blocks[12][z][x] = NewBlock(blocks.GrayWool)
				s.blocks[13][z][x] = NewBlock(blocks.BlueWool)
				s.blocks[14][z][x] = NewBlock(blocks.RedWool)
				s.blocks[15][z][x] = NewBlock(blocks.GreenWool)
			}
		}

		palette := s.makePalette()
		require.Len(t, palette, 16)
		expectedBlocks := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
			blocks.Stone,
			blocks.Grass,
			blocks.Granite,
			blocks.Gravel,
			blocks.Sand,
			blocks.Sandstone,
			blocks.Ice,
			blocks.BlackWool,
			blocks.WhiteWool,
			blocks.PinkWool,
			blocks.GrayWool,
			blocks.BlueWool,
			blocks.RedWool,
			blocks.GreenWool,
		}

		for _, expected := range expectedBlocks {
			var found bool
			for _, actual := range palette {
				if actual == expected {
					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("block ID %d not found in palette", expected))
		}
	})
}

func TestMakeBlockData4(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Dirt)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Dirt)
				s.blocks[3][z][x] = NewBlock(blocks.Dirt)
				s.blocks[4][z][x] = NewBlock(blocks.Air)
				s.blocks[5][z][x] = NewBlock(blocks.Air)
				s.blocks[6][z][x] = NewBlock(blocks.Air)
				s.blocks[7][z][x] = NewBlock(blocks.Air)
				s.blocks[8][z][x] = NewBlock(blocks.Air)
				s.blocks[9][z][x] = NewBlock(blocks.Air)
				s.blocks[10][z][x] = NewBlock(blocks.Air)
				s.blocks[11][z][x] = NewBlock(blocks.Air)
				s.blocks[12][z][x] = NewBlock(blocks.Air)
				s.blocks[13][z][x] = NewBlock(blocks.Air)
				s.blocks[14][z][x] = NewBlock(blocks.Air)
				s.blocks[15][z][x] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		require.Len(t, compacted, 256)
		for i := range compacted {
			if i < 64 {
				assert.Equal(t, uint64(0x1111111111111111), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else {
				assert.Equal(t, uint64(0x0000000000000000), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			}
		}
	})
	t.Run("8_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,       // 0
			blocks.Dirt,      // 1
			blocks.Stone,     // 2
			blocks.Grass,     // 3
			blocks.Granite,   // 4
			blocks.Gravel,    // 5
			blocks.Sand,      // 6
			blocks.Sandstone, // 7
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Air)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Stone)
				s.blocks[3][z][x] = NewBlock(blocks.Grass)
				s.blocks[4][z][x] = NewBlock(blocks.Granite)
				s.blocks[5][z][x] = NewBlock(blocks.Gravel)
				s.blocks[6][z][x] = NewBlock(blocks.Sand)
				s.blocks[7][z][x] = NewBlock(blocks.Sandstone)
				s.blocks[8][z][x] = NewBlock(blocks.Air)
				s.blocks[9][z][x] = NewBlock(blocks.Air)
				s.blocks[10][z][x] = NewBlock(blocks.Air)
				s.blocks[11][z][x] = NewBlock(blocks.Air)
				s.blocks[12][z][x] = NewBlock(blocks.Air)
				s.blocks[13][z][x] = NewBlock(blocks.Air)
				s.blocks[14][z][x] = NewBlock(blocks.Air)
				s.blocks[15][z][x] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		require.Len(t, compacted, 256)
		for i := range compacted {
			if i < 16 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 32 {
				assert.Equal(t, uint64(0x1111111111111111), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 48 {
				assert.Equal(t, uint64(0x2222222222222222), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 64 {
				assert.Equal(t, uint64(0x3333333333333333), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 80 {
				assert.Equal(t, uint64(0x4444444444444444), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 96 {
				assert.Equal(t, uint64(0x5555555555555555), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 112 {
				assert.Equal(t, uint64(0x6666666666666666), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 128 {
				assert.Equal(t, uint64(0x7777777777777777), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else {
				assert.Equal(t, uint64(0x0000000000000000), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			}
		}
	})
	t.Run("16_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
			blocks.Stone,
			blocks.Grass,
			blocks.Granite,
			blocks.Gravel,
			blocks.Sand,
			blocks.Sandstone,
			blocks.Ice,
			blocks.BlackWool,
			blocks.WhiteWool,
			blocks.PinkWool,
			blocks.GrayWool,
			blocks.BlueWool,
			blocks.RedWool,
			blocks.GreenWool,
		}

		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Air)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Stone)
				s.blocks[3][z][x] = NewBlock(blocks.Grass)
				s.blocks[4][z][x] = NewBlock(blocks.Granite)
				s.blocks[5][z][x] = NewBlock(blocks.Gravel)
				s.blocks[6][z][x] = NewBlock(blocks.Sand)
				s.blocks[7][z][x] = NewBlock(blocks.Sandstone)
				s.blocks[8][z][x] = NewBlock(blocks.Ice)
				s.blocks[9][z][x] = NewBlock(blocks.BlackWool)
				s.blocks[10][z][x] = NewBlock(blocks.WhiteWool)
				s.blocks[11][z][x] = NewBlock(blocks.PinkWool)
				s.blocks[12][z][x] = NewBlock(blocks.GrayWool)
				s.blocks[13][z][x] = NewBlock(blocks.BlueWool)
				s.blocks[14][z][x] = NewBlock(blocks.RedWool)
				s.blocks[15][z][x] = NewBlock(blocks.GreenWool)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 256)
		for i := range compacted {
			if i < 16 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 32 {
				assert.Equal(t, uint64(0x1111111111111111), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 48 {
				assert.Equal(t, uint64(0x2222222222222222), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 64 {
				assert.Equal(t, uint64(0x3333333333333333), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 80 {
				assert.Equal(t, uint64(0x4444444444444444), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 96 {
				assert.Equal(t, uint64(0x5555555555555555), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 112 {
				assert.Equal(t, uint64(0x6666666666666666), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 128 {
				assert.Equal(t, uint64(0x7777777777777777), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 144 {
				assert.Equal(t, uint64(0x8888888888888888), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 160 {
				assert.Equal(t, uint64(0x9999999999999999), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 176 {
				assert.Equal(t, uint64(0xAAAAAAAAAAAAAAAA), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 192 {
				assert.Equal(t, uint64(0xBBBBBBBBBBBBBBBB), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 208 {
				assert.Equal(t, uint64(0xCCCCCCCCCCCCCCCC), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 224 {
				assert.Equal(t, uint64(0xDDDDDDDDDDDDDDDD), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 240 {
				assert.Equal(t, uint64(0xEEEEEEEEEEEEEEEE), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else if i < 256 {
				assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			}
		}
	})
}

func TestMakeBlockData14(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,  //  0, 00000000000000
			blocks.Dirt, // 10, 00000000001010
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(blocks.Dirt)
				s.blocks[1][z][x] = NewBlock(blocks.Dirt)
				s.blocks[2][z][x] = NewBlock(blocks.Dirt)
				s.blocks[3][z][x] = NewBlock(blocks.Dirt)
				s.blocks[4][z][x] = NewBlock(blocks.Air)
				s.blocks[5][z][x] = NewBlock(blocks.Air)
				s.blocks[6][z][x] = NewBlock(blocks.Air)
				s.blocks[7][z][x] = NewBlock(blocks.Air)
				s.blocks[8][z][x] = NewBlock(blocks.Air)
				s.blocks[9][z][x] = NewBlock(blocks.Air)
				s.blocks[10][z][x] = NewBlock(blocks.Air)
				s.blocks[11][z][x] = NewBlock(blocks.Air)
				s.blocks[12][z][x] = NewBlock(blocks.Air)
				s.blocks[13][z][x] = NewBlock(blocks.Air)
				s.blocks[14][z][x] = NewBlock(blocks.Air)
				s.blocks[15][z][x] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(14, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 1024)

		for i := range compacted {
			if i < 256 {
				assert.Equal(t, uint64(0x00002800a002800a), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			} else {
				assert.Equal(t, uint64(0x0000000000000000), compacted[i], fmt.Sprintf("compacted long %d invalid", i))
			}
		}
	})
}
