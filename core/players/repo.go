package players

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
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
			CurrentHotbarSlot: player.Slot0,
			CurrentLocation: data.Location{
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
	// DEBT this hack is needed to replace the old ID of the user that was on server before with the new ID for
	//  the newly created connection for this user.
	slice := orm.PlayerSlice{dbPlayer}
	updates := orm.M{orm.PlayerColumns.ID: replacementID.String()}
	if count, err := slice.UpdateAll(context.TODO(), tx, updates); err != nil {
		return nil, fmt.Errorf("failed to update player ID: %w", err)
	} else if count == 0 {
		return nil, fmt.Errorf("failed to update player ID: no rows updated")
	} else if count > 1 {
		return nil, fmt.Errorf("failed to update player ID: more than one row updated")
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
			CurrentHotbarSlot: player.Slot0,
			CurrentLocation: data.Location{
				PositionF: data.PositionF{
					X: dbPlayer.PositionX,
					Y: dbPlayer.PositionY,
					Z: dbPlayer.PositionZ,
				},
			},
		},
	}, nil
}
