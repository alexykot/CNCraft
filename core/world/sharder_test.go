package world

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexykot/cncraft/pkg/game/level"
)

func TestSplitDimShards(t *testing.T) {
	type testCase struct {
		name              string
		edges             level.Edges
		shardSize         int64
		expectShardStarts map[ShardID]startMessage
	}

	dimID := uuid.New()
	dimName := "test_dim"

	cases := []testCase{
		{
			name: "4_NE",
			edges: level.Edges{
				NegativeX: 0,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 2,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 16), ch(16, 16),
						ch(00, 00), ch(16, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 32): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48),
						ch(00, 32), ch(16, 32),
					},
				},
				MkShardIDFromCoords(dimName, 32, 0): {
					chunkIDs: []level.ChunkID{
						ch(32, 16), ch(48, 16),
						ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 32, 32): {
					chunkIDs: []level.ChunkID{
						ch(32, 48), ch(48, 48),
						ch(32, 32), ch(48, 32),
					},
				},
			},
		},
		{
			name: "1_NE",
			edges: level.Edges{
				NegativeX: 0,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 4,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},
			},
		},
		{
			name: "1_NE_underfill",
			edges: level.Edges{
				NegativeX: 0,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 5,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},
			},
		},
		{
			name: "4_NE_underfill",
			edges: level.Edges{
				NegativeX: 0,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 3,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 32), ch(16, 32), ch(32, 32),
						ch(00, 16), ch(16, 16), ch(32, 16),
						ch(00, 00), ch(16, 00), ch(32, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 48): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48),
					},
				},
				MkShardIDFromCoords(dimName, 48, 0): {
					chunkIDs: []level.ChunkID{
						ch(48, 32),
						ch(48, 16),
						ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 48, 48): {
					chunkIDs: []level.ChunkID{
						ch(48, 48),
					},
				},
			},
		},
		{
			name: "2_NE_NW",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 4,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48), ch(-32, 48), ch(-16, 48),
						ch(-64, 32), ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-64, 16), ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-64, 00), ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
			},
		},
		{
			name: "8_NE_NW",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 2,
			expectShardStarts: map[ShardID]startMessage{
				// NE
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 16), ch(16, 16),
						ch(00, 00), ch(16, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 32): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48),
						ch(00, 32), ch(16, 32),
					},
				},
				MkShardIDFromCoords(dimName, 32, 0): {
					chunkIDs: []level.ChunkID{
						ch(32, 16), ch(48, 16),
						ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 32, 32): {
					chunkIDs: []level.ChunkID{
						ch(32, 48), ch(48, 48),
						ch(32, 32), ch(48, 32),
					},
				},
				// NW
				MkShardIDFromCoords(dimName, -32, 0): {
					chunkIDs: []level.ChunkID{
						ch(-32, 16), ch(-16, 16),
						ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 16), ch(-48, 16),
						ch(-64, 00), ch(-48, 00),
					},
				},
				MkShardIDFromCoords(dimName, -32, 32): {
					chunkIDs: []level.ChunkID{
						ch(-32, 48), ch(-16, 48),
						ch(-32, 32), ch(-16, 32),
					},
				},
				MkShardIDFromCoords(dimName, -64, 32): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48),
						ch(-64, 32), ch(-48, 32),
					},
				},
			},
		},
		{
			name: "8_NE_NW_underfill",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 3,
			expectShardStarts: map[ShardID]startMessage{
				// NE
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 32), ch(16, 32), ch(32, 32),
						ch(00, 16), ch(16, 16), ch(32, 16),
						ch(00, 00), ch(16, 00), ch(32, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 48): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48),
					},
				},
				MkShardIDFromCoords(dimName, 48, 0): {
					chunkIDs: []level.ChunkID{
						ch(48, 32),
						ch(48, 16),
						ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 48, 48): {
					chunkIDs: []level.ChunkID{
						ch(48, 48),
					},
				},
				// NW
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 48), ch(-32, 48), ch(-16, 48),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 32),
						ch(-64, 16),
						ch(-64, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48),
					},
				},
			},
		},
		{
			name: "2_NE_NW_underfill",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: 0,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 5,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48), ch(-32, 48), ch(-16, 48),
						ch(-64, 32), ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-64, 16), ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-64, 00), ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
			},
		},
		{
			name: "4_NE_NW_SE_SW",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 4,
			expectShardStarts: map[ShardID]startMessage{
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48), ch(-32, 48), ch(-16, 48),
						ch(-64, 32), ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-64, 16), ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-64, 00), ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -16), ch(16, -16), ch(32, -16), ch(48, -16),
						ch(00, -32), ch(16, -32), ch(32, -32), ch(48, -32),
						ch(00, -48), ch(16, -48), ch(32, -48), ch(48, -48),
						ch(00, -64), ch(16, -64), ch(32, -64), ch(48, -64),
					},
				},
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(-64, -16), ch(48, -16), ch(-32, -16), ch(-16, -16),
						ch(-64, -32), ch(48, -32), ch(-32, -32), ch(-16, -32),
						ch(-64, -48), ch(48, -48), ch(-32, -48), ch(-16, -48),
						ch(-64, -64), ch(48, -64), ch(-32, -64), ch(-16, -64),
					},
				},
			},
		},
		{
			name: "4_NE_NW_SE_SW_underfill",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 5,
			expectShardStarts: map[ShardID]startMessage{
				// NE
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48), ch(48, 48),
						ch(00, 32), ch(16, 32), ch(32, 32), ch(48, 32),
						ch(00, 16), ch(16, 16), ch(32, 16), ch(48, 16),
						ch(00, 00), ch(16, 00), ch(32, 00), ch(48, 00),
					},
				},

				// NW
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48), ch(-32, 48), ch(-16, 48),
						ch(-64, 32), ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-64, 16), ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-64, 00), ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},

				// SE
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -16), ch(16, -16), ch(32, -16), ch(48, -16),
						ch(00, -32), ch(16, -32), ch(32, -32), ch(48, -32),
						ch(00, -48), ch(16, -48), ch(32, -48), ch(48, -48),
						ch(00, -64), ch(16, -64), ch(32, -64), ch(48, -64),
					},
				},

				// SW
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(-64, -16), ch(48, -16), ch(-32, -16), ch(-16, -16),
						ch(-64, -32), ch(48, -32), ch(-32, -32), ch(-16, -32),
						ch(-64, -48), ch(48, -48), ch(-32, -48), ch(-16, -48),
						ch(-64, -64), ch(48, -64), ch(-32, -64), ch(-16, -64),
					},
				},
			},
		},
		{
			name: "16_NE_NW_SE_SW",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 5,
			expectShardStarts: map[ShardID]startMessage{
				// NE
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 16), ch(16, 16),
						ch(00, 00), ch(16, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 32): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48),
						ch(00, 32), ch(16, 32),
					},
				},
				MkShardIDFromCoords(dimName, 32, 0): {
					chunkIDs: []level.ChunkID{
						ch(32, 16), ch(48, 16),
						ch(32, 00), ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 32, 32): {
					chunkIDs: []level.ChunkID{
						ch(32, 48), ch(48, 48),
						ch(32, 32), ch(48, 32),
					},
				},

				// NW
				MkShardIDFromCoords(dimName, -32, 0): {
					chunkIDs: []level.ChunkID{
						ch(-32, 16), ch(-16, 16),
						ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 16), ch(-48, 16),
						ch(-64, 00), ch(-48, 00),
					},
				},
				MkShardIDFromCoords(dimName, -32, 32): {
					chunkIDs: []level.ChunkID{
						ch(-32, 48), ch(-16, 48),
						ch(-32, 32), ch(-16, 32),
					},
				},
				MkShardIDFromCoords(dimName, -64, 32): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48), ch(-48, 48),
						ch(-64, 32), ch(-48, 32),
					},
				},

				// SE
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -48), ch(16, -48),
						ch(00, -64), ch(16, -64),
					},
				},
				MkShardIDFromCoords(dimName, 32, -64): {
					chunkIDs: []level.ChunkID{
						ch(32, -48), ch(48, -48),
						ch(32, -64), ch(48, -64),
					},
				},
				MkShardIDFromCoords(dimName, 00, -32): {
					chunkIDs: []level.ChunkID{
						ch(00, -16), ch(16, -16),
						ch(00, -32), ch(16, -32),
					},
				},
				MkShardIDFromCoords(dimName, 32, -32): {
					chunkIDs: []level.ChunkID{
						ch(32, -16), ch(48, -16),
						ch(32, -32), ch(48, -32),
					},
				},

				// SW
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(-64, -48), ch(48, -48),
						ch(-64, -64), ch(48, -64),
					},
				},
				MkShardIDFromCoords(dimName, -32, -64): {
					chunkIDs: []level.ChunkID{
						ch(-32, -48), ch(-16, -48),
						ch(-32, -64), ch(-16, -64),
					},
				},
				MkShardIDFromCoords(dimName, -64, -32): {
					chunkIDs: []level.ChunkID{
						ch(-64, -16), ch(48, -16),
						ch(-64, -32), ch(48, -32),
					},
				},
				MkShardIDFromCoords(dimName, -32, -32): {
					chunkIDs: []level.ChunkID{
						ch(-32, -16), ch(-16, -16),
						ch(-32, -32), ch(-16, -32),
					},
				},
			},
		},
		{
			name: "16_NE_NW_SE_SW_underfill",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 48,
				PositiveZ: 48,
			},
			shardSize: 3,
			expectShardStarts: map[ShardID]startMessage{
				// NE
				MkShardIDFromCoords(dimName, 0, 0): {
					chunkIDs: []level.ChunkID{
						ch(00, 32), ch(16, 32), ch(32, 32),
						ch(00, 16), ch(16, 16), ch(32, 16),
						ch(00, 00), ch(16, 00), ch(32, 00),
					},
				},
				MkShardIDFromCoords(dimName, 0, 48): {
					chunkIDs: []level.ChunkID{
						ch(00, 48), ch(16, 48), ch(32, 48),
					},
				},
				MkShardIDFromCoords(dimName, 48, 0): {
					chunkIDs: []level.ChunkID{
						ch(48, 32),
						ch(48, 16),
						ch(48, 00),
					},
				},
				MkShardIDFromCoords(dimName, 48, 48): {
					chunkIDs: []level.ChunkID{
						ch(48, 48),
					},
				},

				// NW
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 48), ch(-32, 48), ch(-16, 48),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 32),
						ch(-64, 16),
						ch(-64, 00),
					},
				},
				MkShardIDFromCoords(dimName, -64, 0): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48),
					},
				},

				// SE
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -16), ch(16, -16), ch(32, -16),
						ch(00, -32), ch(16, -32), ch(32, -32),
						ch(00, -48), ch(16, -48), ch(32, -48),
					},
				},
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(48, -16),
						ch(48, -32),
						ch(48, -48),
					},
				},
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -64), ch(16, -64), ch(32, -64),
					},
				},
				MkShardIDFromCoords(dimName, 00, -64): {
					chunkIDs: []level.ChunkID{
						ch(48, -64),
					},
				},

				// SW
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(48, -16), ch(-32, -16), ch(-16, -16),
						ch(48, -32), ch(-32, -32), ch(-16, -32),
						ch(48, -48), ch(-32, -48), ch(-16, -48),
					},
				},
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(-64, -16),
						ch(-64, -32),
						ch(-64, -48),
					},
				},
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(48, -64), ch(-32, -64), ch(-16, -64),
					},
				},
				MkShardIDFromCoords(dimName, -64, -64): {
					chunkIDs: []level.ChunkID{
						ch(-64, -64),
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			shardStarts := splitDimShards(dimID, dimName, test.edges, test.shardSize, test.shardSize)
			assert.Equal(t, len(test.expectShardStarts), len(shardStarts))
			for shardID, expectStartMess := range test.expectShardStarts {
				startMess, ok := shardStarts[shardID]
				require.True(t, ok)
				assert.Equal(t, dimID, startMess.dimensionID)
				assert.Equal(t, shardID, startMess.id)

				require.Equal(t, len(expectStartMess.chunkIDs), len(startMess.chunkIDs))
				for index, expectChunkID := range expectStartMess.chunkIDs {
					// TODO order is not guaranteed, rework the test
					assert.Equal(t, expectChunkID, startMess.chunkIDs[index])
				}
			}
		})
	}
}

func ch(x, z int64) level.ChunkID {
	return level.MkChunkID(x, z)
}
