package world

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/pkg/game/level"
	"github.com/alexykot/cncraft/pkg/protocol/objects"
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
				blocks[0][z][x] = level.NewBlock(objects.BlockBedrock)
				blocks[1][z][x] = level.NewBlock(objects.BlockDirt)
				blocks[2][z][x] = level.NewBlock(objects.BlockDirt)
				blocks[3][z][x] = level.NewBlock(objects.BlockGrassBlock_SnowyTrue)
				blocks[4][z][x] = level.NewBlock(objects.BlockAir)
				blocks[5][z][x] = level.NewBlock(objects.BlockAir)
				blocks[6][z][x] = level.NewBlock(objects.BlockAir)
				blocks[7][z][x] = level.NewBlock(objects.BlockAir)
				blocks[8][z][x] = level.NewBlock(objects.BlockAir)
				blocks[9][z][x] = level.NewBlock(objects.BlockAir)
				blocks[10][z][x] = level.NewBlock(objects.BlockAir)
				blocks[11][z][x] = level.NewBlock(objects.BlockAir)
				blocks[12][z][x] = level.NewBlock(objects.BlockAir)
				blocks[13][z][x] = level.NewBlock(objects.BlockAir)
				blocks[14][z][x] = level.NewBlock(objects.BlockAir)
				blocks[15][z][x] = level.NewBlock(objects.BlockAir)
			}
		}
	} else {
		for z := 0; z < level.SectionZ; z++ {
			for x := 0; x < level.SectionX; x++ {
				blocks[0][z][x] = level.NewBlock(objects.BlockAir)
				blocks[1][z][x] = level.NewBlock(objects.BlockAir)
				blocks[2][z][x] = level.NewBlock(objects.BlockAir)
				blocks[3][z][x] = level.NewBlock(objects.BlockAir)
				blocks[4][z][x] = level.NewBlock(objects.BlockAir)
				blocks[5][z][x] = level.NewBlock(objects.BlockAir)
				blocks[6][z][x] = level.NewBlock(objects.BlockAir)
				blocks[7][z][x] = level.NewBlock(objects.BlockAir)
				blocks[8][z][x] = level.NewBlock(objects.BlockAir)
				blocks[9][z][x] = level.NewBlock(objects.BlockAir)
				blocks[10][z][x] = level.NewBlock(objects.BlockAir)
				blocks[11][z][x] = level.NewBlock(objects.BlockAir)
				blocks[12][z][x] = level.NewBlock(objects.BlockAir)
				blocks[13][z][x] = level.NewBlock(objects.BlockAir)
				blocks[14][z][x] = level.NewBlock(objects.BlockAir)
				blocks[15][z][x] = level.NewBlock(objects.BlockAir)
			}
		}
	}

	return level.NewSection(blocks, index)
}
