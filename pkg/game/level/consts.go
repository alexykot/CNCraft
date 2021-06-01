package level

const (
	ChunkX = 16  // Width
	ChunkZ = 16  // Length
	ChunkY = 256 // Height

	SectionX = 16 // Width
	SectionZ = 16 // Length
	SectionY = 16 // Height
)

type BlockArr [SectionY][SectionZ][SectionX]Block
