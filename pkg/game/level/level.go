package level

type Level interface {
	Name() string

	Chunks() map[ChunkID]Chunk
	GetChunk(ChunkID) Chunk

	Edges() edges
	//	GetBlock(x, y, z int) Block
}

type edges struct {
	NegativeX int64
	NegativeZ int64
	PositiveX int64
	PositiveZ int64
}

type level struct {
	name   string
	chunks map[ChunkID]Chunk

	boundaries edges
}

var defaultLevel *level

func NewLevel(name string) Level {
	// TODO replace with actual level loading
	return getDefaultLevel(name)
}

func (l *level) Name() string              { return l.name }
func (l *level) Chunks() map[ChunkID]Chunk { return l.chunks }
func (l *level) GetChunk(id ChunkID) Chunk { return l.chunks[id] }
func (l *level) Edges() edges              { return l.boundaries }

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

		// Assuming the level is not going to change dynamically, and is rectangular shape.
		for _, chunk := range defaultLevel.chunks {
			if defaultLevel.boundaries.NegativeX > chunk.X() {
				defaultLevel.boundaries.NegativeX = chunk.X()
			}
			if defaultLevel.boundaries.NegativeZ > chunk.Z() {
				defaultLevel.boundaries.NegativeZ = chunk.Z()
			}
			if defaultLevel.boundaries.PositiveX < chunk.X() {
				defaultLevel.boundaries.PositiveX = chunk.X()
			}
			if defaultLevel.boundaries.PositiveZ < chunk.Z() {
				defaultLevel.boundaries.PositiveZ = chunk.Z()
			}
		}
	}

	return defaultLevel
}
