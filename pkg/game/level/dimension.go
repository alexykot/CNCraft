package level

type Dimension interface {
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

type dimension struct {
	name   string
	chunks map[ChunkID]Chunk

	boundaries edges
}

var defaultDim *dimension

func NewDimension(name string) Dimension {
	// TODO replace with actual dimension loading
	return getDefaultDimension(name)
}

func (d *dimension) Name() string              { return d.name }
func (d *dimension) Chunks() map[ChunkID]Chunk { return d.chunks }
func (d *dimension) GetChunk(id ChunkID) Chunk { return d.chunks[id] }
func (d *dimension) Edges() edges              { return d.boundaries }

func getDefaultDimension(name string) Dimension {
	chunks := map[ChunkID]Chunk{}
	for x := -48; x <= 48; x = x + 16 {
		for z := -48; z <= 48; z = z + 16 {
			chunk := NewChunk(int64(x), int64(z))
			chunks[chunk.ID()] = chunk
		}
	}

	if defaultDim == nil {
		defaultDim = &dimension{
			name:   name,
			chunks: chunks,
		}

		// Assuming the dimension is not going to change dynamically, and is rectangular shape.
		for _, chunk := range defaultDim.chunks {
			if defaultDim.boundaries.NegativeX > chunk.X() {
				defaultDim.boundaries.NegativeX = chunk.X()
			}
			if defaultDim.boundaries.NegativeZ > chunk.Z() {
				defaultDim.boundaries.NegativeZ = chunk.Z()
			}
			if defaultDim.boundaries.PositiveX < chunk.X() {
				defaultDim.boundaries.PositiveX = chunk.X()
			}
			if defaultDim.boundaries.PositiveZ < chunk.Z() {
				defaultDim.boundaries.PositiveZ = chunk.Z()
			}
		}
	}

	return defaultDim
}
