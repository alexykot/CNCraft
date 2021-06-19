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
				assert.Equal(t, uint8(2), heights[x][z])
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
			// if i < len(compacted)-1 {
			// 	// 000000100 000000100 000000100 000000100 000000100 000000100 000000100 0
			// 	// 00000010 00000001 00000000 10000000 01000000 00100000 00010000 00001000
			// 	// 0x0201008040201008
			// 	expect = int64(0x0201008040201008)
			// } else {
			// 	// 000000000 000000000 000000000 000000100 000000100 000000100 000000100 0
			// 	// 00000000 00000000 00000000 00000000 01000000 00100000 00010000 00001000
			// 	// 0x0000000040201008
			// 	expect = int64(0x0000000040201008)
			// }

			// with Notchian -2 adjustment
			if i < len(compacted)-1 {
				expect = int64(0x0100804020100804)
			} else {
				expect = int64(0x0000000020100804)
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

func TestGetLocalXZ(t *testing.T) {
	type testCase struct {
		provideGlobal int64
		expectLocal   int64
	}

	cases := []testCase{
		{provideGlobal: 0, expectLocal: 0},
		{provideGlobal: 15, expectLocal: 15},
		{provideGlobal: 16, expectLocal: 0},
		{provideGlobal: 32, expectLocal: 0},
		{provideGlobal: 34, expectLocal: 2},
		{provideGlobal: 37, expectLocal: 5},
		{provideGlobal: 159, expectLocal: 15},
		{provideGlobal: 160, expectLocal: 0},
		{provideGlobal: -15, expectLocal: 1},
		{provideGlobal: -16, expectLocal: 0},
		{provideGlobal: -17, expectLocal: 15},
		{provideGlobal: -31, expectLocal: 1},
		{provideGlobal: -32, expectLocal: 0},
		{provideGlobal: -34, expectLocal: 14},
		{provideGlobal: -37, expectLocal: 11},
		{provideGlobal: -160, expectLocal: 0},
		{provideGlobal: -161, expectLocal: 15},
	}

	for _, test := range cases {
		assert.Equal(t, test.expectLocal, getLocalXZ(test.provideGlobal))
	}
}

func TestGetChunkXZ(t *testing.T) {
	type testCase struct {
		provideBlock int64
		expectChunk  int64
	}

	cases := []testCase{
		{provideBlock: 0, expectChunk: 0},
		{provideBlock: 15, expectChunk: 0},
		{provideBlock: 16, expectChunk: 16},
		{provideBlock: 32, expectChunk: 32},
		{provideBlock: 34, expectChunk: 32},
		{provideBlock: 37, expectChunk: 32},
		{provideBlock: 159, expectChunk: 144},
		{provideBlock: 160, expectChunk: 160},
		{provideBlock: -1, expectChunk: -16},
		{provideBlock: -15, expectChunk: -16},
		{provideBlock: -16, expectChunk: -16},
		{provideBlock: -17, expectChunk: -32},
		{provideBlock: -31, expectChunk: -32},
		{provideBlock: -32, expectChunk: -32},
		{provideBlock: -34, expectChunk: -48},
		{provideBlock: -37, expectChunk: -48},
		{provideBlock: -160, expectChunk: -160},
		{provideBlock: -161, expectChunk: -176},
	}

	for _, test := range cases {
		assert.Equal(t, test.expectChunk, getChunkXZ(test.provideBlock))
	}
}
