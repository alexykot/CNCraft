package level

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexykot/cncraft/pkg/protocol/objects"
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
				s.blocks[0][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[3][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[4][z][x] = NewBlock(objects.BlockAir)
				s.blocks[5][z][x] = NewBlock(objects.BlockAir)
				s.blocks[6][z][x] = NewBlock(objects.BlockAir)
				s.blocks[7][z][x] = NewBlock(objects.BlockAir)
				s.blocks[8][z][x] = NewBlock(objects.BlockAir)
				s.blocks[9][z][x] = NewBlock(objects.BlockAir)
				s.blocks[10][z][x] = NewBlock(objects.BlockAir)
				s.blocks[11][z][x] = NewBlock(objects.BlockAir)
				s.blocks[12][z][x] = NewBlock(objects.BlockAir)
				s.blocks[13][z][x] = NewBlock(objects.BlockAir)
				s.blocks[14][z][x] = NewBlock(objects.BlockAir)
				s.blocks[15][z][x] = NewBlock(objects.BlockAir)
			}
		}

		palette := s.makePalette()
		require.Len(t, palette, 2)
		assert.Equal(t, objects.BlockDirt, palette[0])
		assert.Equal(t, objects.BlockAir, palette[1])
	})
	t.Run("16_blocks", func(t *testing.T) {
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(objects.BlockAir)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockStone)
				s.blocks[3][z][x] = NewBlock(objects.BlockGrass)
				s.blocks[4][z][x] = NewBlock(objects.BlockGranite)
				s.blocks[5][z][x] = NewBlock(objects.BlockGravel)
				s.blocks[6][z][x] = NewBlock(objects.BlockSand)
				s.blocks[7][z][x] = NewBlock(objects.BlockSandstone)
				s.blocks[8][z][x] = NewBlock(objects.BlockIce)
				s.blocks[9][z][x] = NewBlock(objects.BlockBlackWool)
				s.blocks[10][z][x] = NewBlock(objects.BlockWhiteWool)
				s.blocks[11][z][x] = NewBlock(objects.BlockPinkWool)
				s.blocks[12][z][x] = NewBlock(objects.BlockGrayWool)
				s.blocks[13][z][x] = NewBlock(objects.BlockBlueWool)
				s.blocks[14][z][x] = NewBlock(objects.BlockRedWool)
				s.blocks[15][z][x] = NewBlock(objects.BlockGreenWool)
			}
		}

		palette := s.makePalette()
		require.Len(t, palette, 16)
		expectedBlocks := []objects.BlockID{
			objects.BlockAir,
			objects.BlockDirt,
			objects.BlockStone,
			objects.BlockGrass,
			objects.BlockGranite,
			objects.BlockGravel,
			objects.BlockSand,
			objects.BlockSandstone,
			objects.BlockIce,
			objects.BlockBlackWool,
			objects.BlockWhiteWool,
			objects.BlockPinkWool,
			objects.BlockGrayWool,
			objects.BlockBlueWool,
			objects.BlockRedWool,
			objects.BlockGreenWool,
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
		palette := []objects.BlockID{
			objects.BlockAir,
			objects.BlockDirt,
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[3][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[4][z][x] = NewBlock(objects.BlockAir)
				s.blocks[5][z][x] = NewBlock(objects.BlockAir)
				s.blocks[6][z][x] = NewBlock(objects.BlockAir)
				s.blocks[7][z][x] = NewBlock(objects.BlockAir)
				s.blocks[8][z][x] = NewBlock(objects.BlockAir)
				s.blocks[9][z][x] = NewBlock(objects.BlockAir)
				s.blocks[10][z][x] = NewBlock(objects.BlockAir)
				s.blocks[11][z][x] = NewBlock(objects.BlockAir)
				s.blocks[12][z][x] = NewBlock(objects.BlockAir)
				s.blocks[13][z][x] = NewBlock(objects.BlockAir)
				s.blocks[14][z][x] = NewBlock(objects.BlockAir)
				s.blocks[15][z][x] = NewBlock(objects.BlockAir)
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
		palette := []objects.BlockID{
			objects.BlockAir,       // 0
			objects.BlockDirt,      // 1
			objects.BlockStone,     // 2
			objects.BlockGrass,     // 3
			objects.BlockGranite,   // 4
			objects.BlockGravel,    // 5
			objects.BlockSand,      // 6
			objects.BlockSandstone, // 7
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(objects.BlockAir)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockStone)
				s.blocks[3][z][x] = NewBlock(objects.BlockGrass)
				s.blocks[4][z][x] = NewBlock(objects.BlockGranite)
				s.blocks[5][z][x] = NewBlock(objects.BlockGravel)
				s.blocks[6][z][x] = NewBlock(objects.BlockSand)
				s.blocks[7][z][x] = NewBlock(objects.BlockSandstone)
				s.blocks[8][z][x] = NewBlock(objects.BlockAir)
				s.blocks[9][z][x] = NewBlock(objects.BlockAir)
				s.blocks[10][z][x] = NewBlock(objects.BlockAir)
				s.blocks[11][z][x] = NewBlock(objects.BlockAir)
				s.blocks[12][z][x] = NewBlock(objects.BlockAir)
				s.blocks[13][z][x] = NewBlock(objects.BlockAir)
				s.blocks[14][z][x] = NewBlock(objects.BlockAir)
				s.blocks[15][z][x] = NewBlock(objects.BlockAir)
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
		palette := []objects.BlockID{
			objects.BlockAir,
			objects.BlockDirt,
			objects.BlockStone,
			objects.BlockGrass,
			objects.BlockGranite,
			objects.BlockGravel,
			objects.BlockSand,
			objects.BlockSandstone,
			objects.BlockIce,
			objects.BlockBlackWool,
			objects.BlockWhiteWool,
			objects.BlockPinkWool,
			objects.BlockGrayWool,
			objects.BlockBlueWool,
			objects.BlockRedWool,
			objects.BlockGreenWool,
		}

		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(objects.BlockAir)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockStone)
				s.blocks[3][z][x] = NewBlock(objects.BlockGrass)
				s.blocks[4][z][x] = NewBlock(objects.BlockGranite)
				s.blocks[5][z][x] = NewBlock(objects.BlockGravel)
				s.blocks[6][z][x] = NewBlock(objects.BlockSand)
				s.blocks[7][z][x] = NewBlock(objects.BlockSandstone)
				s.blocks[8][z][x] = NewBlock(objects.BlockIce)
				s.blocks[9][z][x] = NewBlock(objects.BlockBlackWool)
				s.blocks[10][z][x] = NewBlock(objects.BlockWhiteWool)
				s.blocks[11][z][x] = NewBlock(objects.BlockPinkWool)
				s.blocks[12][z][x] = NewBlock(objects.BlockGrayWool)
				s.blocks[13][z][x] = NewBlock(objects.BlockBlueWool)
				s.blocks[14][z][x] = NewBlock(objects.BlockRedWool)
				s.blocks[15][z][x] = NewBlock(objects.BlockGreenWool)
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
		palette := []objects.BlockID{
			objects.BlockAir,  //  0, 00000000000000
			objects.BlockDirt, // 10, 00000000001010
		}
		s := &section{index: 0, blocks: [SectionY][SectionZ][SectionX]Block{}}
		for z := 0; z < SectionZ; z++ {
			for x := 0; x < SectionX; x++ {
				s.blocks[0][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[1][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[2][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[3][z][x] = NewBlock(objects.BlockDirt)
				s.blocks[4][z][x] = NewBlock(objects.BlockAir)
				s.blocks[5][z][x] = NewBlock(objects.BlockAir)
				s.blocks[6][z][x] = NewBlock(objects.BlockAir)
				s.blocks[7][z][x] = NewBlock(objects.BlockAir)
				s.blocks[8][z][x] = NewBlock(objects.BlockAir)
				s.blocks[9][z][x] = NewBlock(objects.BlockAir)
				s.blocks[10][z][x] = NewBlock(objects.BlockAir)
				s.blocks[11][z][x] = NewBlock(objects.BlockAir)
				s.blocks[12][z][x] = NewBlock(objects.BlockAir)
				s.blocks[13][z][x] = NewBlock(objects.BlockAir)
				s.blocks[14][z][x] = NewBlock(objects.BlockAir)
				s.blocks[15][z][x] = NewBlock(objects.BlockAir)
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
