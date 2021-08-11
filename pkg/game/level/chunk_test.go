package level

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/protocol/objects"
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
			s.blocks[x][z][0] = NewBlock(objects.BlockDirt)
			s.blocks[x][z][1] = NewBlock(objects.BlockDirt)
			s.blocks[x][z][2] = NewBlock(objects.BlockDirt)
			s.blocks[x][z][3] = NewBlock(objects.BlockDirt)
			s.blocks[x][z][4] = NewBlock(objects.BlockAir)
			s.blocks[x][z][5] = NewBlock(objects.BlockAir)
			s.blocks[x][z][6] = NewBlock(objects.BlockAir)
			s.blocks[x][z][7] = NewBlock(objects.BlockAir)
			s.blocks[x][z][8] = NewBlock(objects.BlockAir)
			s.blocks[x][z][9] = NewBlock(objects.BlockAir)
			s.blocks[x][z][10] = NewBlock(objects.BlockAir)
			s.blocks[x][z][11] = NewBlock(objects.BlockAir)
			s.blocks[x][z][12] = NewBlock(objects.BlockAir)
			s.blocks[x][z][13] = NewBlock(objects.BlockAir)
			s.blocks[x][z][14] = NewBlock(objects.BlockAir)
			s.blocks[x][z][15] = NewBlock(objects.BlockAir)
		}
	}

	return &chunk{sections: []Section{s}}
}

func TestGetLocalPosition(t *testing.T) {
	type testCase struct {
		provideGlobal data.PositionI
		expectLocal   data.PositionI
	}

	cases := []testCase{
		{provideGlobal: data.PositionI{0, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{15, 0, 0}, expectLocal: data.PositionI{15, 0, 0}},
		{provideGlobal: data.PositionI{16, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{32, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{34, 0, 0}, expectLocal: data.PositionI{2, 0, 0}},
		{provideGlobal: data.PositionI{37, 0, 0}, expectLocal: data.PositionI{5, 0, 0}},
		{provideGlobal: data.PositionI{159, 0, 0}, expectLocal: data.PositionI{15, 0, 0}},
		{provideGlobal: data.PositionI{160, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{-15, 0, 0}, expectLocal: data.PositionI{1, 0, 0}},
		{provideGlobal: data.PositionI{-16, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{-17, 0, 0}, expectLocal: data.PositionI{15, 0, 0}},
		{provideGlobal: data.PositionI{-31, 0, 0}, expectLocal: data.PositionI{1, 0, 0}},
		{provideGlobal: data.PositionI{-32, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{-34, 0, 0}, expectLocal: data.PositionI{14, 0, 0}},
		{provideGlobal: data.PositionI{-37, 0, 0}, expectLocal: data.PositionI{11, 0, 0}},
		{provideGlobal: data.PositionI{-160, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{-161, 0, 0}, expectLocal: data.PositionI{15, 0, 0}},

		{provideGlobal: data.PositionI{0, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, 15}, expectLocal: data.PositionI{0, 0, 15}},
		{provideGlobal: data.PositionI{0, 0, 16}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, 32}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, 34}, expectLocal: data.PositionI{0, 0, 2}},
		{provideGlobal: data.PositionI{0, 0, 37}, expectLocal: data.PositionI{0, 0, 5}},
		{provideGlobal: data.PositionI{0, 0, 159}, expectLocal: data.PositionI{0, 0, 15}},
		{provideGlobal: data.PositionI{0, 0, 160}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, -15}, expectLocal: data.PositionI{0, 0, 1}},
		{provideGlobal: data.PositionI{0, 0, -16}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, -17}, expectLocal: data.PositionI{0, 0, 15}},
		{provideGlobal: data.PositionI{0, 0, -31}, expectLocal: data.PositionI{0, 0, 1}},
		{provideGlobal: data.PositionI{0, 0, -32}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, -34}, expectLocal: data.PositionI{0, 0, 14}},
		{provideGlobal: data.PositionI{0, 0, -37}, expectLocal: data.PositionI{0, 0, 11}},
		{provideGlobal: data.PositionI{0, 0, -160}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 0, -161}, expectLocal: data.PositionI{0, 0, 15}},

		{provideGlobal: data.PositionI{0, 0, 0}, expectLocal: data.PositionI{0, 0, 0}},
		{provideGlobal: data.PositionI{0, 15, 0}, expectLocal: data.PositionI{0, 15, 0}},
		{provideGlobal: data.PositionI{0, 16, 0}, expectLocal: data.PositionI{0, 16, 0}},
		{provideGlobal: data.PositionI{0, 32, 0}, expectLocal: data.PositionI{0, 32, 0}},
		{provideGlobal: data.PositionI{0, 34, 0}, expectLocal: data.PositionI{0, 34, 0}},
		{provideGlobal: data.PositionI{0, 37, 0}, expectLocal: data.PositionI{0, 37, 0}},
		{provideGlobal: data.PositionI{0, 159, 0}, expectLocal: data.PositionI{0, 159, 0}},
		{provideGlobal: data.PositionI{0, 160, 0}, expectLocal: data.PositionI{0, 160, 0}},
		{provideGlobal: data.PositionI{0, 254, 0}, expectLocal: data.PositionI{0, 254, 0}},
	}

	for _, test := range cases {
		assert.Equal(t, test.expectLocal, getLocalPosition(test.provideGlobal))
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

	t.Run("assert_chunk_square", func(t *testing.T) {
		assert.True(t, ChunkX == ChunkZ, "only square chunks supported")
	})
}
