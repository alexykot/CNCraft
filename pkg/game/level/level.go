package level

type Level interface {
	Name() string

	Chunks() map[ChunkID]Chunk
	GetChunk(ChunkID) Chunk

	//	GetBlock(x, y, z int) Block
}

type level struct {
	name   string
	chunks map[ChunkID]Chunk
}

var defaultLevel *level

func NewLevel(name string) Level {
	// TODO replace with actual level loading
	return getDefaultLevel(name)
}

func (l *level) Name() string              { return l.name }
func (l *level) Chunks() map[ChunkID]Chunk { return l.chunks }
func (l *level) GetChunk(id ChunkID) Chunk { return l.chunks[id] }

func getDefaultLevel(name string) Level {
	chunks := map[ChunkID]Chunk{}
	for x := -48; x <= 48; x = x + 16 {
		for z := -48; z <= 48; z = z + 16 {
			chunk := NewChunk(int64(x), int64(z))
			chunks[chunk.ID()] = chunk
		}
	}

	if defaultLevel == nil {
		defaultLevel = &level{
			name:   name,
			chunks: chunks,
		}
	}
	return defaultLevel
}
