package level

const (
	// Chunk height is not constant by design
	ChunkX = 16 // Width
	ChunkZ = 16 // Length

	SectionX = 16 // Width
	SectionZ = 16 // Length
	SectionY = 16 // Height
)

type BlockArr [SectionY][SectionZ][SectionX]Block
