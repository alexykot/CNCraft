package players

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

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

func (r *repo) InitPlayer(userID uuid.UUID, username string) (*Player, error) {
	tx, err := r.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbPlayer, err := orm.Players(orm.PlayerWhere.Username.EQ(username)).One(context.TODO(), tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to load player data: %w", err)
	}

	var p *Player
	if err == sql.ErrNoRows {
		if p, err = r.createPlayer(tx, userID, username); err != nil {
			return nil, fmt.Errorf("failed to create player: %w", err)
		}
	} else {
		if p, err = r.loadPlayer(tx, dbPlayer, userID, username); err != nil {
			return nil, fmt.Errorf("failed to load player: %w", err)
		}
	}

	return p, nil
}

func (r *repo) createPlayer(tx *sql.Tx, userID uuid.UUID, username string) (*Player, error) {
	dbPlayer := &orm.Player{
		ID:        userID.String(),
		Username:  username,
		PositionX: 0,
		PositionY: 0,
		PositionZ: 0,
	}
	if err := dbPlayer.Insert(context.TODO(), tx, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert player: %w", err)
	}

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
					X: dbPlayer.PositionX,
					Y: dbPlayer.PositionY,
					Z: dbPlayer.PositionZ,
				},
			},
		},
	}, nil
}

func (r *repo) loadPlayer(tx *sql.Tx, dbPlayer *orm.Player, replacementID uuid.UUID, username string) (*Player, error) {
	dbPlayer.ID = replacementID.String()
	if _, err := dbPlayer.Update(context.TODO(), tx, boil.Whitelist(orm.PlayerColumns.ID)); err != nil {
		return nil, fmt.Errorf("failed to update player ID: %w", err)
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
