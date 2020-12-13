package level

import "github.com/google/uuid"

type Level interface {
	Name() string
	ID() uuid.UUID

	Chunks() []Chunk

	GetChunk(x, z int) Chunk

	GetChunkIfLoaded(x, z int) Chunk

	GetBlock(x, y, z int) Block
}
