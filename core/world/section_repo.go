package world

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/pkg/game/level"
	pBlocks "github.com/alexykot/cncraft/pkg/protocol/blocks"
)

// SectionRepo loads sections from persistence and handles saving world block updates back into persistence
type SectionRepo struct {
	log *zap.Logger
	db  *sql.DB
}

func newRepo(log *zap.Logger, db *sql.DB) *SectionRepo {
	return &SectionRepo{log, db}
}

// LoadSection - TODO implement this properly when implementing actual world persistence
func (r SectionRepo) LoadSection(x, z int64, index uint8) (level.Section, error) {
	return loadDefaultSection(index), nil
}

// SaveSection - TODO implement this when implementing actual world persistence
func (r SectionRepo) SaveSection(section level.Section) error {
	return nil
}

func loadDefaultSection(index uint8) level.Section {
	var blocks level.BlockArr
	if index == 0 {
		for z := 0; z < level.SectionZ; z++ {
			for x := 0; x < level.SectionX; x++ {
				blocks[0][z][x] = level.NewBlock(pBlocks.Bedrock)
				blocks[1][z][x] = level.NewBlock(pBlocks.Dirt)
				blocks[2][z][x] = level.NewBlock(pBlocks.Dirt)
				blocks[3][z][x] = level.NewBlock(pBlocks.GrassBlock_SnowyTrue)
				blocks[4][z][x] = level.NewBlock(pBlocks.Air)
				blocks[5][z][x] = level.NewBlock(pBlocks.Air)
				blocks[6][z][x] = level.NewBlock(pBlocks.Air)
				blocks[7][z][x] = level.NewBlock(pBlocks.Air)
				blocks[8][z][x] = level.NewBlock(pBlocks.Air)
				blocks[9][z][x] = level.NewBlock(pBlocks.Air)
				blocks[10][z][x] = level.NewBlock(pBlocks.Air)
				blocks[11][z][x] = level.NewBlock(pBlocks.Air)
				blocks[12][z][x] = level.NewBlock(pBlocks.Air)
				blocks[13][z][x] = level.NewBlock(pBlocks.Air)
				blocks[14][z][x] = level.NewBlock(pBlocks.Air)
				blocks[15][z][x] = level.NewBlock(pBlocks.Air)
			}
		}
	} else {
		for z := 0; z < level.SectionZ; z++ {
			for x := 0; x < level.SectionX; x++ {
				blocks[0][z][x] = level.NewBlock(pBlocks.Air)
				blocks[1][z][x] = level.NewBlock(pBlocks.Air)
				blocks[2][z][x] = level.NewBlock(pBlocks.Air)
				blocks[3][z][x] = level.NewBlock(pBlocks.Air)
				blocks[4][z][x] = level.NewBlock(pBlocks.Air)
				blocks[5][z][x] = level.NewBlock(pBlocks.Air)
				blocks[6][z][x] = level.NewBlock(pBlocks.Air)
				blocks[7][z][x] = level.NewBlock(pBlocks.Air)
				blocks[8][z][x] = level.NewBlock(pBlocks.Air)
				blocks[9][z][x] = level.NewBlock(pBlocks.Air)
				blocks[10][z][x] = level.NewBlock(pBlocks.Air)
				blocks[11][z][x] = level.NewBlock(pBlocks.Air)
				blocks[12][z][x] = level.NewBlock(pBlocks.Air)
				blocks[13][z][x] = level.NewBlock(pBlocks.Air)
				blocks[14][z][x] = level.NewBlock(pBlocks.Air)
				blocks[15][z][x] = level.NewBlock(pBlocks.Air)
			}
		}
	}

	return level.NewSection(blocks, index)
}
