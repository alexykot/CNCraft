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
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		palette := s.makePalette()
		require.Len(t, palette, 2)
		assert.Equal(t, blocks.Dirt, palette[0])
		assert.Equal(t, blocks.Air, palette[1])
	})
	t.Run("16_blocks", func(t *testing.T) {
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Air)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Stone)
				s.blocks[x][z][3] = NewBlock(blocks.Grass)
				s.blocks[x][z][4] = NewBlock(blocks.Granite)
				s.blocks[x][z][5] = NewBlock(blocks.Gravel)
				s.blocks[x][z][6] = NewBlock(blocks.Sand)
				s.blocks[x][z][7] = NewBlock(blocks.Sandstone)
				s.blocks[x][z][8] = NewBlock(blocks.Ice)
				s.blocks[x][z][9] = NewBlock(blocks.BlackWool)
				s.blocks[x][z][10] = NewBlock(blocks.WhiteWool)
				s.blocks[x][z][11] = NewBlock(blocks.PinkWool)
				s.blocks[x][z][12] = NewBlock(blocks.GrayWool)
				s.blocks[x][z][13] = NewBlock(blocks.BlueWool)
				s.blocks[x][z][14] = NewBlock(blocks.RedWool)
				s.blocks[x][z][15] = NewBlock(blocks.GreenWool)
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
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		require.Len(t, compacted, 256)
		for i := range compacted {
			assert.Equal(t, uint64(0x1111000000000000), compacted[i])
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
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Air)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Stone)
				s.blocks[x][z][3] = NewBlock(blocks.Grass)
				s.blocks[x][z][4] = NewBlock(blocks.Granite)
				s.blocks[x][z][5] = NewBlock(blocks.Gravel)
				s.blocks[x][z][6] = NewBlock(blocks.Sand)
				s.blocks[x][z][7] = NewBlock(blocks.Sandstone)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Dirt)
				s.blocks[x][z][10] = NewBlock(blocks.Stone)
				s.blocks[x][z][11] = NewBlock(blocks.Grass)
				s.blocks[x][z][12] = NewBlock(blocks.Granite)
				s.blocks[x][z][13] = NewBlock(blocks.Gravel)
				s.blocks[x][z][14] = NewBlock(blocks.Sand)
				s.blocks[x][z][15] = NewBlock(blocks.Sandstone)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		require.Len(t, compacted, 256)
		for i := range compacted {
			assert.Equal(t, uint64(0x0123456701234567), compacted[i])
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

		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Air)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Stone)
				s.blocks[x][z][3] = NewBlock(blocks.Grass)
				s.blocks[x][z][4] = NewBlock(blocks.Granite)
				s.blocks[x][z][5] = NewBlock(blocks.Gravel)
				s.blocks[x][z][6] = NewBlock(blocks.Sand)
				s.blocks[x][z][7] = NewBlock(blocks.Sandstone)
				s.blocks[x][z][8] = NewBlock(blocks.Ice)
				s.blocks[x][z][9] = NewBlock(blocks.BlackWool)
				s.blocks[x][z][10] = NewBlock(blocks.WhiteWool)
				s.blocks[x][z][11] = NewBlock(blocks.PinkWool)
				s.blocks[x][z][12] = NewBlock(blocks.GrayWool)
				s.blocks[x][z][13] = NewBlock(blocks.BlueWool)
				s.blocks[x][z][14] = NewBlock(blocks.RedWool)
				s.blocks[x][z][15] = NewBlock(blocks.GreenWool)
			}
		}

		compacted, err := s.makeBlockData(4, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 256)
		for i := range compacted {
			assert.Equal(t, uint64(0x0123456789ABCDEF), compacted[i])
		}
	})
}

func TestMakeBlockData5(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}

		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(5, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 342)
		var i int
		for ii := range compacted {
			if ii == len(compacted)-1 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii]) // last long is unfinished and only has few blocks
				break
			}

			i++
			if i == 1 {
				assert.Equal(t, uint64(0x0084210000000000), compacted[ii])
			}
			if i == 2 {
				assert.Equal(t, uint64(0x0000000842100000), compacted[ii])
			}
			if i == 3 {
				assert.Equal(t, uint64(0x0000000000008421), compacted[ii])
			}
			if i == 4 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii])
				i = 0
			}
		}
	})
}

func TestMakeBlockData6(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(6, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 410)
		var i int
		for ii := range compacted {
			if ii == len(compacted)-1 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii]) // last long is unfinished and only has few blocks
				break
			}

			i++
			switch i {
			case 1:
				assert.Equal(t, uint64(0x0041041000000000), compacted[ii])
			case 2:
				assert.Equal(t, uint64(0x0000000000041041), compacted[ii])
			case 3:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii])
			case 4:
				assert.Equal(t, uint64(0x0000041041000000), compacted[ii])
			case 5:
				assert.Equal(t, uint64(0x0000000000000041), compacted[ii])
			case 6:
				assert.Equal(t, uint64(0x0041000000000000), compacted[ii])
			case 7:
				assert.Equal(t, uint64(0x0000000041041000), compacted[ii])
			case 8:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii])
				i = 0
			}
		}
	})
}

func TestMakeBlockData7(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(7, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 456)

		var i int
		for ii := range compacted {
			if ii == len(compacted)-1 {
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii]) // last long is unfinished and only has few blocks
				break
			}

			i++
			switch i {
			case 1:
				assert.Equal(t, uint64(0x0102040800000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 2:
				assert.Equal(t, uint64(0x0000000000000081), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 3:
				assert.Equal(t, uint64(0x0102000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 4:
				assert.Equal(t, uint64(0x0000000000204081), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 5:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 6:
				assert.Equal(t, uint64(0x0000000810204000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 7:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 8:
				assert.Equal(t, uint64(0x0002040810000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 9:
				assert.Equal(t, uint64(0x0000000000000001), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 10:
				assert.Equal(t, uint64(0x0102040000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 11:
				assert.Equal(t, uint64(0x0000000000004081), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 12:
				assert.Equal(t, uint64(0x0100000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 13:
				assert.Equal(t, uint64(0x0000000010204080), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 14:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 15:
				assert.Equal(t, uint64(0x0000040810200000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 16:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
				i = 0
			}

		}
	})
}

func TestMakeBlockData8(t *testing.T) {
	t.Run("2_blocks", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(8, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 512)

		var i int
		for ii := range compacted {
			i++
			switch i {
			case 1:
				assert.Equal(t, uint64(0x0101010100000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 2:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
				i = 0
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
		s := &section{index: 0, blocks: [16][16][16]Block{}}
		for x, zBlocks := range s.blocks {
			for z := range zBlocks {
				s.blocks[x][z][0] = NewBlock(blocks.Dirt)
				s.blocks[x][z][1] = NewBlock(blocks.Dirt)
				s.blocks[x][z][2] = NewBlock(blocks.Dirt)
				s.blocks[x][z][3] = NewBlock(blocks.Dirt)
				s.blocks[x][z][4] = NewBlock(blocks.Air)
				s.blocks[x][z][5] = NewBlock(blocks.Air)
				s.blocks[x][z][6] = NewBlock(blocks.Air)
				s.blocks[x][z][7] = NewBlock(blocks.Air)
				s.blocks[x][z][8] = NewBlock(blocks.Air)
				s.blocks[x][z][9] = NewBlock(blocks.Air)
				s.blocks[x][z][10] = NewBlock(blocks.Air)
				s.blocks[x][z][11] = NewBlock(blocks.Air)
				s.blocks[x][z][12] = NewBlock(blocks.Air)
				s.blocks[x][z][13] = NewBlock(blocks.Air)
				s.blocks[x][z][14] = NewBlock(blocks.Air)
				s.blocks[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted, err := s.makeBlockData(14, palette)
		require.NoError(t, err)
		assert.Len(t, compacted, 1024)

		var i int
		for ii := range compacted {
			i++
			switch i {
			case 1:
				assert.Equal(t, uint64(0x00002800A002800A), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 2:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 3:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
			case 4:
				assert.Equal(t, uint64(0x0000000000000000), compacted[ii], fmt.Sprintf("line %d failed, long #%d", i, ii))
				i = 0
			}
		}
	})
}
