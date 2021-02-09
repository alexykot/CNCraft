package level

import "github.com/alexykot/cncraft/pkg/game"

type Level interface {
	Name() string

	Chunks() map[ChunkID]Chunk
	GetChunk(ChunkID) Chunk

	//	GetChunkIfLoaded(x, z int) GetChunk
	//
	//	GetBlock(x, y, z int) Block
}

type level struct {
	name   string
	chunks map[ChunkID]Chunk
}

var defaultLevel *level

func NewLevel(name string) Level {
	return &level{
		name: name,
	}
}

func (l *level) Name() string              { return l.name }
func (l *level) Chunks() map[ChunkID]Chunk { return l.chunks }
func (l *level) GetChunk(id ChunkID) Chunk { return l.chunks[id] }

func GetDefaultLevel() Level {
	chunks := map[ChunkID]Chunk{}
	for x := -48; x <= 48; x = x + 16 {
		for z := -48; z <= 48; z = z + 16 {
			chunk := NewDefaultChunk(int64(x), int64(z))
			chunks[chunk.ID()] = chunk
		}
	}

	if defaultLevel == nil {
		defaultLevel = &level{
			name:   game.Overworld.String(),
			chunks: chunks,
		}
	}
	return defaultLevel
}
