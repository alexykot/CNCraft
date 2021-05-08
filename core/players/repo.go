package players

import "C"
import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/items"
	"github.com/alexykot/cncraft/pkg/game/player"
)

// repo is a player repository, it implements handling the persistent storage of player data.
type repo struct {
	db *sql.DB
}

func newRepo(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) InitPlayer(username string, connID uuid.UUID) (p *Player, isNew bool, err error) {
	tx, err := r.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbPlayer, err := orm.Players(orm.PlayerWhere.Username.EQ(username)).One(context.TODO(), tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, fmt.Errorf("failed to load player data: %w", err)
	}

	if err == sql.ErrNoRows {
		p = r.createNewPlayer(username, connID)
		isNew = true
	} else {
		if p, err = r.loadPlayer(tx, dbPlayer, username, connID); err != nil {
			return nil, false, fmt.Errorf("failed to load player: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit tx: %w", err)
	}

	return p, isNew, nil
}

// TODO this needs to be replaced with proper spawn point and starting conditions.
func (r *repo) createNewPlayer(username string, connID uuid.UUID) *Player {
	return &Player{
		ID:       uuid.New(),
		ConnID:   connID,
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		State: player.State{
			Location: data.Location{
				PositionF: data.PositionF{
					X: 0,
					Y: 10,
					Z: 0,
				},
			},
		},
	}
}

func (r *repo) loadPlayer(tx *sql.Tx, dbPlayer *orm.Player, username string, connID uuid.UUID) (*Player, error) {
	dbPlayer.ConnID = null.StringFrom(connID.String())
	_, err := dbPlayer.Update(context.TODO(), tx, boil.Whitelist(orm.PlayerColumns.ConnID))

	dbInventories, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(dbPlayer.ID)).All(context.TODO(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to query player inventory: %w", err)
	}

	inventory := items.Inventory{
		CurrentHotbarSlot: items.HotBarSlot(dbPlayer.CurrentHotbar),
	}
	for _, dbItem := range dbInventories {
		inventory.AssignSlot(int(dbItem.SlotNumber), int(dbItem.ItemID), int(dbItem.ItemCount))
	}

	return &Player{
		ID:       dbPlayer.ID,
		ConnID:   connID,
		PC:       entities.NewPC(username, player.MaxHealth),
		Username: username,
		Settings: player.Settings{
			ViewDistance: 7,
			FlyingSpeed:  0.05,
			FoVModifier:  0.1,
		},
		State: player.State{
			Inventory: inventory,
			Location: data.Location{
				PositionF: data.PositionF{
					X: dbPlayer.PositionX,
					Y: dbPlayer.PositionY,
					Z: dbPlayer.PositionZ,
				},
			},
		},
	}, nil
}
