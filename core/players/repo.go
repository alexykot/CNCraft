// Package players contains implementation for players list (Roster), the player repo that allows to load players
// from persistence and the Player details itself.
package players

import "C"
import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/db"
	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/items"
	"github.com/alexykot/cncraft/pkg/game/player"
	"github.com/alexykot/cncraft/pkg/protocol/objects"
)

// repo is a player repository, it implements handling the persistent storage of player data.
type repo struct {
	windowLog *zap.Logger
	db        *sql.DB
}

func newRepo(log *zap.Logger, db *sql.DB) *repo {
	return &repo{log, db}
}

func (r *repo) InitPlayer(username string, connID, dimensionID uuid.UUID) (p *Player, isNew bool, err error) {
	tx, err := r.db.BeginTx(db.Ctx(), nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbPlayer, err := orm.Players(orm.PlayerWhere.Username.EQ(username)).One(db.Ctx(), tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, fmt.Errorf("failed to load player data: %w", err)
	}

	if err == sql.ErrNoRows {
		p = r.createNewPlayer(username, connID, dimensionID)
		isNew = true
	} else {
		if p, err = r.loadPlayer(tx, dbPlayer, username, connID, dimensionID); err != nil {
			return nil, false, fmt.Errorf("failed to load player: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit tx: %w", err)
	}

	return p, isNew, nil
}

func (r *repo) createNewPlayer(username string, connID, dimensionID uuid.UUID) *Player {
	inventory := items.NewInventory(r.windowLog)

	return &Player{
		ID:       uuid.New(),
		ConnID:   connID,
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: &player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		Abilities: &player.Abilities{},
		State: &player.State{
			Dimension: dimensionID,
			Inventory: inventory,
			Location: data.Location{ // DEBT this needs to be replaced with proper spawn point and starting conditions.
				PositionF: data.PositionF{
					X: 0,
					Y: 10,
					Z: 0,
				},
			},
		},
	}
}

func (r *repo) loadPlayer(tx *sql.Tx, dbPlayer *orm.Player, username string, connID, dimensionID uuid.UUID) (*Player, error) {
	dbPlayer.ConnID = null.StringFrom(connID.String())
	_, err := dbPlayer.Update(db.Ctx(), tx, boil.Whitelist(orm.PlayerColumns.ConnID))

	dbInventories, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(dbPlayer.ID)).All(db.Ctx(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to query player inventory: %w", err)
	}

	inventory := items.NewInventory(r.windowLog)
	inventory.CurrentHotbarSlot = uint8(dbPlayer.CurrentHotbar)

	for _, dbItem := range dbInventories {
		inventory.SetSlot(dbItem.SlotNumber, items.Slot{ItemID: objects.ItemID(dbItem.ItemID), ItemCount: dbItem.ItemCount})
	}

	return &Player{
		ID:       dbPlayer.ID,
		ConnID:   connID,
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: &player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		Abilities: &player.Abilities{},
		State: &player.State{
			// not using the previously saved dimension for the player here because player may
			// join a dimension different from what they left previously.
			Dimension: dimensionID,
			Inventory: inventory,
			Location: data.Location{
				// DEBT this does not account for the fact that the player may join a different dimension from
				//  what they left previously, and last recorded location in that dimension will be irrelevant.
				PositionF: data.PositionF{
					X: dbPlayer.PositionX,
					Y: dbPlayer.PositionY,
					Z: dbPlayer.PositionZ,
				},
			},
		},
	}, nil
}
