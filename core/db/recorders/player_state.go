package recorders

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
)

var ctxGet func() context.Context

// RegisterPlayerStateHandlers registers handlers for envelopes carrying updates of the player state that
// need to be persisted.
func RegisterPlayerStateHandlers(ctxGetter func() context.Context, ps nats.PubSub, log *zap.Logger, db *sql.DB) error {
	ctxGet = ctxGetter

	if err := ps.Subscribe(subj.MkNewPlayerJoined(), handleNewPlayer(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	if err := ps.Subscribe(subj.MkPlayerLeft(), handlePlayerLeft(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	if err := ps.Subscribe(subj.MkPlayerSpatialUpdate(), handlePlayerSpatial(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	if err := ps.Subscribe(subj.MkPlayerInventoryUpdate(), handlePlayerInventory(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}
	return nil
}

func handleNewPlayer(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		var err error
		log := log
		db := db

		newPlayer := inLope.GetNewPlayer()
		if newPlayer == nil {
			log.Error("envelope does not contain new player", zap.Any("envelope", inLope))
			return
		}

		dbPlayer := &orm.Player{
			ConnID:    null.StringFrom(newPlayer.ConnId),
			Username:  newPlayer.Username,
			PositionX: newPlayer.Pos.X,
			PositionY: newPlayer.Pos.Y,
			PositionZ: newPlayer.Pos.Z,
		}
		if dbPlayer.ID, err = uuid.Parse(newPlayer.PlayerId); err != nil {
			log.Error("failed to parse player ID", zap.String("id", newPlayer.PlayerId), zap.Error(err))
			return
		}
		if err := dbPlayer.Insert(ctxGet(), db, boil.Infer()); err != nil {
			log.Error("failed to insert player", zap.String("id", newPlayer.PlayerId), zap.Error(err))
			return
		}
	}
}

func handlePlayerLeft(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		var err error
		log := log
		db := db

		playerLeft := inLope.GetPlayerLeft()
		if playerLeft == nil {
			log.Error("envelope does not contain leaving player", zap.Any("envelope", inLope))
			return
		}

		dbPlayer := &orm.Player{ConnID: null.NewString("", false)}
		if dbPlayer.ID, err = uuid.Parse(playerLeft.PlayerId); err != nil {
			log.Error("failed to parse player ID", zap.String("id", playerLeft.PlayerId), zap.Error(err))
			return
		}
		if _, err = dbPlayer.Update(ctxGet(), db, boil.Whitelist(orm.PlayerColumns.ConnID)); err != nil {
			log.Error("failed to update player", zap.String("id", playerLeft.PlayerId), zap.Error(err))
			return
		}
	}
}

func handlePlayerSpatial(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		log := log
		db := db

		spatial := inLope.GetPlayerSpatial()
		if spatial == nil {
			log.Error("envelope does not contain spatial update", zap.Any("envelope", inLope))
			return
		}

		playerId, err := uuid.Parse(spatial.PlayerId)
		if err != nil {
			log.Error("failed to parse player ID", zap.String("id", spatial.PlayerId), zap.Error(err))
			return
		}

		player, err := orm.FindPlayer(ctxGet(), db, playerId)
		if err == sql.ErrNoRows {
			log.Warn("received posUpdate for nonexistent player, ignoring", zap.String("id", spatial.PlayerId))
			return
		} else if err != nil {
			log.Error("failed to fetch user by UUID", zap.String("id", spatial.PlayerId), zap.Error(err))
			return
		}

		player.PositionX = spatial.Pos.X
		player.PositionY = spatial.Pos.Y
		player.PositionZ = spatial.Pos.Z

		player.Yaw = float64(spatial.Rot.Yaw)
		player.Pitch = float64(spatial.Rot.Pitch)

		player.OnGround = spatial.OnGround

		if _, err = player.Update(ctxGet(), db,
			boil.Whitelist(
				orm.PlayerColumns.PositionX, orm.PlayerColumns.PositionY, orm.PlayerColumns.PositionZ,
				orm.PlayerColumns.Yaw, orm.PlayerColumns.Pitch, orm.PlayerColumns.OnGround,
			)); err != nil {
			log.Error("failed to save user position", zap.String("id", spatial.PlayerId), zap.Error(err))
		}
	}
}

func handlePlayerInventory(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		log := log
		db := db

		inventory := inLope.GetPlayerInventory()
		if inventory == nil {
			log.Error("envelope does not contain inventory update", zap.Any("envelope", inLope))
			return
		}

		playerId, err := uuid.Parse(inventory.PlayerId)
		if err != nil {
			log.Error("failed to parse player ID", zap.String("id", inventory.PlayerId), zap.Error(err))
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Error("failed to start tx", zap.Error(err))
			return
		}

		dbPlayer, err := orm.FindPlayer(ctxGet(), tx, playerId)
		if err == sql.ErrNoRows {
			log.Warn("received posUpdate for nonexistent player, ignoring", zap.String("id", inventory.PlayerId))
			return
		} else if err != nil {
			log.Error("failed to fetch player by UUID", zap.String("id", inventory.PlayerId), zap.Error(err))
			return
		}

		if _, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(playerId)).DeleteAll(ctxGet(), tx); err != nil {
			log.Error("failed to wipe player inventory", zap.String("id", inventory.PlayerId), zap.Error(err))
			_ = tx.Rollback()
			return
		}

		for _, item := range inventory.Inventory {
			dbItem := orm.Inventory{
				PlayerID:   playerId,
				SlotNumber: int16(item.SlotId),
				ItemID:     int16(item.ItemId),
				ItemCount:  int16(item.ItemCount),
			}
			if err := dbItem.Insert(ctxGet(), tx, boil.Infer()); err != nil {
				log.Error("failed to wipe player inventory", zap.String("id", inventory.PlayerId), zap.Error(err))
				_ = tx.Rollback()
				return
			}
		}

		dbPlayer.CurrentHotbar = int16(inventory.CurrentHotbar)
		if _, err = dbPlayer.Update(ctxGet(), tx, boil.Whitelist(orm.PlayerColumns.CurrentHotbar)); err != nil {
			log.Error("failed to save user hotbar active slot", zap.String("id", inventory.PlayerId), zap.Error(err))
		}

		if err := tx.Commit(); err != nil {
			log.Error("failed to commit transaction", zap.String("id", inventory.PlayerId), zap.Error(err))
		}
	}
}
