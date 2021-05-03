package players

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

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

func (r *repo) InitPlayer(userID uuid.UUID, username string) (p *Player, isNew bool, err error) {
	tx, err := r.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbPlayer, err := orm.Players(orm.PlayerWhere.Username.EQ(username)).One(context.TODO(), tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, fmt.Errorf("failed to load player data: %w", err)
	}

	if err == sql.ErrNoRows {
		p = r.createNewPlayer(userID, username)
		isNew = true
	} else {
		if p, err = r.loadPlayer(tx, dbPlayer, userID, username); err != nil {
			return nil, false, fmt.Errorf("failed to load player: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit tx: %w", err)
	}

	return p, isNew, nil
}

// TODO this needs to be replaced with proper spawn point and starting conditions.
func (r *repo) createNewPlayer(userID uuid.UUID, username string) *Player {
	return &Player{
		ID:       userID,
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

func (r *repo) loadPlayer(tx *sql.Tx, dbPlayer *orm.Player, replacementID uuid.UUID, username string) (*Player, error) {

	{
		oldId := dbPlayer.ID
		// DEBT this hack is needed to replace the old ID of the user that was on server before with the new ID for
		//  the newly created connection for this user.
		//  This will need to be done the other way round, i.e. update the ConnID with the one stored before,
		//  instead of updating stored player with the new ConnID.

		// dbPlayer.ID = replacementID.String()
		// err := dbPlayer.Insert(context.TODO(), tx, boil.Infer())

		if _, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(oldId)).
			UpdateAll(context.TODO(), tx, orm.M{orm.InventoryColumns.PlayerID: replacementID.String()}); err != nil {
			return nil, fmt.Errorf("failed to update player ID: %w", err)
		}

		if count, err := orm.Players(orm.PlayerWhere.ID.EQ(oldId)).
			UpdateAll(context.TODO(), tx, orm.M{orm.PlayerColumns.ID: replacementID.String()}); err != nil {
			return nil, fmt.Errorf("failed to update player ID: %w", err)
		} else if count == 0 {
			return nil, fmt.Errorf("failed to update player ID: no rows updated")
		} else if count > 1 {
			return nil, fmt.Errorf("failed to update player ID: more than one row updated")
		}

		dbPlayer.ID = replacementID.String()
		if err := dbPlayer.Reload(context.TODO(), tx); err != nil {
			return nil, fmt.Errorf("failed to reload player: %w", err)
		}
	}

	dbInventory, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(dbPlayer.ID)).All(context.TODO(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to query player inventory: %w", err)
	}

	inventory := items.Inventory{
		CurrentHotbarSlot: items.HotBarSlot(dbPlayer.CurrentHotbar),
	}
	for _, dbItem := range dbInventory {
		inventory.AssignSlot(int(dbItem.SlotNumber), int(dbItem.ItemID), int(dbItem.ItemCount))
	}

	return &Player{
		ID:       replacementID,
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
