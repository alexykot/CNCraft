package level

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexykot/cncraft/pkg/protocol/blocks"
)

func TestFindHeights(t *testing.T) {
	t.Run("flatworld", func(t *testing.T) {
		c := getDefaultChunk()
		heights := c.findHeights()
		for x, zHeignts := range heights {
			for z := range zHeignts {
				assert.Equal(t, uint8(4), heights[x][z])
			}
		}
	})
}

func TestCompactHeights(t *testing.T) {
	t.Run("flatworld", func(t *testing.T) {
		c := getDefaultChunk()
		compacted := c.compactHeights(c.findHeights())
		assert.Len(t, compacted, 37)
		for i := range compacted {

			var expect int64
			if i < len(compacted)-1 {
				// 000000100 000000100 000000100 000000100 000000100 000000100 000000100 0
				// 00000010 00000001 00000000 10000000 01000000 00100000 00010000 00001000
				// 0x0201008040201008
				expect = int64(0x0201008040201008)
			} else {
				// 000000000 000000000 000000000 000000100 000000100 000000100 000000100 0
				// 00000000 00000000 00000000 00000000 01000000 00100000 00010000 00001000
				// 0x0000000040201008
				expect = int64(0x0000000040201008)
			}
			assert.Equal(t, expect, compacted[i], fmt.Sprintf("error at compact index %d", i))
		}
	})
}

func getDefaultChunk() *chunk {
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

	return &chunk{sections: []Section{s}}
}
