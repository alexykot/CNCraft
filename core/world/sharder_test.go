package world

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexykot/cncraft/pkg/game/level"
)

func TestSplitDimensionShards(t *testing.T) {
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
			name: "4_NE_only",
			edges: level.Edges{
				NegativeX: 0,
				NegativeZ: 0,
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
			name: "2_NE_NW_only",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: 0,
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
				MkShardIDFromCoords(dimName, -48, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -48, 48): {
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
				MkShardIDFromCoords(dimName, -64, 48): {
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
				PositiveX: 63,
				PositiveZ: 63,
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
				PositiveX: 63,
				PositiveZ: 63,
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
						ch(-64, -16), ch(-48, -16), ch(-32, -16), ch(-16, -16),
						ch(-64, -32), ch(-48, -32), ch(-32, -32), ch(-16, -32),
						ch(-64, -48), ch(-48, -48), ch(-32, -48), ch(-16, -48),
						ch(-64, -64), ch(-48, -64), ch(-32, -64), ch(-16, -64),
					},
				},
			},
		},
		{
			name: "4_NE_NW_SE_SW_underfill",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 63,
				PositiveZ: 63,
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
				MkShardIDFromCoords(dimName, 0, -64): {
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
						ch(-64, -16), ch(-48, -16), ch(-32, -16), ch(-16, -16),
						ch(-64, -32), ch(-48, -32), ch(-32, -32), ch(-16, -32),
						ch(-64, -48), ch(-48, -48), ch(-32, -48), ch(-16, -48),
						ch(-64, -64), ch(-48, -64), ch(-32, -64), ch(-16, -64),
					},
				},
			},
		},
		{
			name: "16_NE_NW_SE_SW",
			edges: level.Edges{
				NegativeX: -64,
				NegativeZ: -64,
				PositiveX: 63,
				PositiveZ: 63,
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

				// SE
				MkShardIDFromCoords(dimName, 0, -64): {
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
				MkShardIDFromCoords(dimName, 0, -32): {
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
						ch(-64, -48), ch(-48, -48),
						ch(-64, -64), ch(-48, -64),
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
						ch(-64, -16), ch(-48, -16),
						ch(-64, -32), ch(-48, -32),
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
				PositiveX: 63,
				PositiveZ: 63,
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
				MkShardIDFromCoords(dimName, -48, 0): {
					chunkIDs: []level.ChunkID{
						ch(-48, 32), ch(-32, 32), ch(-16, 32),
						ch(-48, 16), ch(-32, 16), ch(-16, 16),
						ch(-48, 00), ch(-32, 00), ch(-16, 00),
					},
				},
				MkShardIDFromCoords(dimName, -48, 48): {
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
				MkShardIDFromCoords(dimName, -64, 48): {
					chunkIDs: []level.ChunkID{
						ch(-64, 48),
					},
				},

				// SE
				MkShardIDFromCoords(dimName, 0, -48): {
					chunkIDs: []level.ChunkID{
						ch(00, -16), ch(16, -16), ch(32, -16),
						ch(00, -32), ch(16, -32), ch(32, -32),
						ch(00, -48), ch(16, -48), ch(32, -48),
					},
				},
				MkShardIDFromCoords(dimName, 48, -48): {
					chunkIDs: []level.ChunkID{
						ch(48, -16),
						ch(48, -32),
						ch(48, -48),
					},
				},
				MkShardIDFromCoords(dimName, 0, -64): {
					chunkIDs: []level.ChunkID{
						ch(00, -64), ch(16, -64), ch(32, -64),
					},
				},
				MkShardIDFromCoords(dimName, 48, -64): {
					chunkIDs: []level.ChunkID{
						ch(48, -64),
					},
				},

				// SW
				MkShardIDFromCoords(dimName, -48, -48): {
					chunkIDs: []level.ChunkID{
						ch(-48, -16), ch(-32, -16), ch(-16, -16),
						ch(-48, -32), ch(-32, -32), ch(-16, -32),
						ch(-48, -48), ch(-32, -48), ch(-16, -48),
					},
				},
				MkShardIDFromCoords(dimName, -64, -48): {
					chunkIDs: []level.ChunkID{
						ch(-64, -16),
						ch(-64, -32),
						ch(-64, -48),
					},
				},
				MkShardIDFromCoords(dimName, -48, -64): {
					chunkIDs: []level.ChunkID{
						ch(-48, -64), ch(-32, -64), ch(-16, -64),
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
			shardStarts := splitDimensionShards(dimID, dimName, test.edges, test.shardSize, test.shardSize)
			require.Equal(t, len(test.expectShardStarts), len(shardStarts), fmt.Sprintf("unexpected shardStarts lendth %d", len(shardStarts)))

			for shardID, expectStartMess := range test.expectShardStarts {
				startMess, ok := shardStarts[shardID]
				require.True(t, ok, fmt.Sprintf("shardStart not found for %s", shardID))
				assert.Equal(t, dimID, startMess.dimensionID)
				assert.Equal(t, shardID, startMess.id)

				require.Equal(t, len(expectStartMess.chunkIDs), len(startMess.chunkIDs),
					fmt.Sprintf("unexpected number of chunks - %d in shard %s", len(startMess.chunkIDs), startMess.id))
				for _, expectChunkID := range expectStartMess.chunkIDs {
					var found bool
					for _, chunkID := range startMess.chunkIDs {
						if chunkID == expectChunkID {
							found = true
							break
						}
					}
					assert.True(t, found, fmt.Sprintf("chunk %s not found in the shard %s", expectChunkID, shardID))
				}
			}
		})
	}
}

func TestSplitAreaChunks(t *testing.T) {
	type testCase struct {
		name             string
		lowerX, lowerZ   int64
		higherX, higherZ int64
		expect           []level.ChunkID
	}

	cases := []testCase{
		{
			name:    "ok_one_positive",
			lowerX:  0,
			lowerZ:  0,
			higherX: 16,
			higherZ: 16,
			expect:  []level.ChunkID{level.MkChunkID(0, 0)},
		},
		{
			name:    "ok_two_positive",
			lowerX:  0,
			lowerZ:  0,
			higherX: 16,
			higherZ: 32,
			expect: []level.ChunkID{
				level.MkChunkID(0, 0),
				level.MkChunkID(0, 16),
			},
		},
		{
			name:    "ok_four_positive",
			lowerX:  0,
			lowerZ:  0,
			higherX: 32,
			higherZ: 32,
			expect: []level.ChunkID{
				level.MkChunkID(0, 0),
				level.MkChunkID(0, 16),
				level.MkChunkID(16, 0),
				level.MkChunkID(16, 16),
			},
		},
		{
			name:    "ok_one_negative",
			lowerX:  -16,
			lowerZ:  -16,
			higherX: 0,
			higherZ: 0,
			expect:  []level.ChunkID{level.MkChunkID(-16, -16)},
		},
		{
			name:    "ok_two_negative",
			lowerX:  -16,
			lowerZ:  -32,
			higherX: 0,
			higherZ: 0,
			expect: []level.ChunkID{
				level.MkChunkID(-16, -32),
				level.MkChunkID(-16, -16),
			},
		},
		{
			name:    "ok_four_negative",
			lowerX:  -32,
			lowerZ:  -32,
			higherX: 0,
			higherZ: 0,
			expect: []level.ChunkID{
				level.MkChunkID(-32, -32),
				level.MkChunkID(-32, -16),
				level.MkChunkID(-16, -32),
				level.MkChunkID(-16, -16),
			},
		},
		{
			name:    "ok_four_negative",
			lowerX:  -32,
			lowerZ:  -32,
			higherX: 0,
			higherZ: 0,
			expect: []level.ChunkID{
				level.MkChunkID(-32, -32),
				level.MkChunkID(-32, -16),
				level.MkChunkID(-16, -32),
				level.MkChunkID(-16, -16),
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			res := splitAreaChunks(test.lowerX, test.lowerZ, test.higherX, test.higherZ)
			require.Equal(t, len(test.expect), len(res))
			for i, expectID := range test.expect {
				assert.Equal(t, expectID, res[i])
			}
		})
	}
}

// func TestFindShardID(t *testing.T) {
//     world := getTestWorld()
//     dimID := world.StartDimension
//     dimName := world.Dimensions[world.StartDimension].Name()
//     sharder := NewSharder(log.MustGetTestLogger(), nil, control.WorldConf{ShardSize: 3}, nil, world, nil)
//     go func() {
//         for {
//             select {
//             case <-sharder.shardControl:
//                 // drain the control loop
//             }
//         }
//     }()
//
//     type testCase struct {
//         name          string
//         dimID         uuid.UUID
//         coords        data.PositionI
//         expectShardID ShardID
//         expectFound   bool
//     }
//
//     cases := []testCase{
//         {
//             name:          "yes_wtf",
//             dimID:         dimID,
//             coords:        data.PositionI{X: -1, Y: 3, Z: 0},
//             expectShardID: MkShardIDFromCoords(dimName, -48, 0),
//             expectFound:   true,
//         },
//         {
//             name:  "no_dim_not_found",
//             dimID: uuid.New(),
//         },
//         {
//             name:   "no_beyond_egde_positive_x",
//             dimID:  dimID,
//             coords: data.PositionI{X: 49, Y: 3, Z: 0},
//         },
//         {
//             name:   "no_beyond_egde_negative_x",
//             dimID:  dimID,
//             coords: data.PositionI{X: -49, Y: 3, Z: 0},
//         },
//         {
//             name:   "no_beyond_egde_positive_z",
//             dimID:  dimID,
//             coords: data.PositionI{X: 0, Y: 3, Z: 49},
//         },
//         {
//             name:   "no_beyond_egde_negative_z",
//             dimID:  dimID,
//             coords: data.PositionI{X: 0, Y: 3, Z: -49},
//         },
//     }
//
//     for _, test := range cases {
//         t.Run(test.name, func(t *testing.T) {
//             shardID, found := sharder.FindShardID(test.dimID, test.coords)
//             assert.Equal(t, test.expectFound, found)
//             assert.Equal(t, test.expectShardID, shardID)
//         })
//     }
// }

func ch(x, z int64) level.ChunkID {
	return level.MkChunkID(x, z)
}
