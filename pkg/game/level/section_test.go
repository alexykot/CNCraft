package level

import (
	"encoding/binary"
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

func TestCompactBlocksBpb4(t *testing.T) {
	t.Run("2_block", func(t *testing.T) {
		palette := []blocks.BlockID{
			blocks.Air,
			blocks.Dirt,
		}
		blockList := [16][16][16]Block{}
		for x, zBlocks := range blockList {
			for z, _ := range zBlocks {
				blockList[x][z][0] = NewBlock(blocks.Dirt)
				blockList[x][z][1] = NewBlock(blocks.Dirt)
				blockList[x][z][2] = NewBlock(blocks.Dirt)
				blockList[x][z][3] = NewBlock(blocks.Dirt)
				blockList[x][z][4] = NewBlock(blocks.Air)
				blockList[x][z][5] = NewBlock(blocks.Air)
				blockList[x][z][6] = NewBlock(blocks.Air)
				blockList[x][z][7] = NewBlock(blocks.Air)
				blockList[x][z][8] = NewBlock(blocks.Air)
				blockList[x][z][9] = NewBlock(blocks.Air)
				blockList[x][z][10] = NewBlock(blocks.Air)
				blockList[x][z][11] = NewBlock(blocks.Air)
				blockList[x][z][12] = NewBlock(blocks.Air)
				blockList[x][z][13] = NewBlock(blocks.Air)
				blockList[x][z][14] = NewBlock(blocks.Air)
				blockList[x][z][15] = NewBlock(blocks.Air)
			}
		}

		compacted := compactBlocksBpb4(palette, blockList)
		expectedLong, _ := binary.Uvarint([]byte{
			0x11, // 0001 0001
			0x11, // 0001 0001
			0x00, // 0000 0000
			0x00, // 0000 0000
			0x00, // 0000 0000
			0x00, // 0000 0000
			0x00, // 0000 0000
			0x00, // 0000 0000
		})
		require.Len(t, compacted, 256)
		assert.Equal(t, expectedLong, compacted[0])
	})
}

func TestCompactBlocksBpb5(t *testing.T) {

}

func TestCompactBlocksBpb6(t *testing.T) {

}

func TestCompactBlocksBpb7(t *testing.T) {

}

func TestCompactBlocksBpb8(t *testing.T) {

}

func TestCompactBlocksBpb14(t *testing.T) {

}
