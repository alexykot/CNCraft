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

var getCtx func() context.Context

// RegisterPlayerStateHandlers registers handlers for envelopes carrying updates of the player state that
// need to be persisted.
func RegisterPlayerStateHandlers(ctxGetter func() context.Context, ps nats.PubSub, log *zap.Logger, db *sql.DB) error {
	getCtx = ctxGetter

	if err := ps.Subscribe(subj.MkPlayerJoined(), handlePlayerJoined(log, db)); err != nil {
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

func handlePlayerJoined(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		var err error
		log := log
		db := db

		joinedPlayer := inLope.GetPlayerJoined()
		if joinedPlayer == nil {
			log.Error("envelope does not contain a joining player", zap.Any("envelope", inLope))
			return
		}

		dbPlayer := &orm.Player{
			ConnID:    null.StringFrom(joinedPlayer.ConnId),
			Username:  joinedPlayer.Username,
			PositionX: joinedPlayer.Pos.X,
			PositionY: joinedPlayer.Pos.Y,
			PositionZ: joinedPlayer.Pos.Z,
		}
		if dbPlayer.ID, err = uuid.Parse(joinedPlayer.PlayerId); err != nil {
			log.Error("failed to parse player ID", zap.String("id", joinedPlayer.PlayerId), zap.Error(err))
			return
		}
		if dbPlayer.DimensionID, err = uuid.Parse(joinedPlayer.DimensionId); err != nil {
			log.Error("failed to parse dimension ID", zap.String("id", joinedPlayer.DimensionId), zap.Error(err))
			return
		}
		exists, err := orm.PlayerExists(getCtx(), db, dbPlayer.ID)
		if err != nil {
			log.Error("failed to check if player exists", zap.Error(err))
			return
		}

		if exists {
			if rows, err := dbPlayer.Update(getCtx(), db, boil.Whitelist(
				orm.PlayerColumns.ConnID,
				orm.PlayerColumns.DimensionID,
				orm.PlayerColumns.PositionX,
				orm.PlayerColumns.PositionY,
				orm.PlayerColumns.PositionZ,
			)); err != nil {
				log.Error("failed to update player", zap.String("id", joinedPlayer.PlayerId), zap.Error(err))
			} else if rows == 0 {
				log.Error("failed to update player: player not found", zap.String("id", joinedPlayer.PlayerId))
			} else if rows > 1 { // this should be impossible with update by primary key
				log.Error("more than one player updated by id", zap.String("id", joinedPlayer.PlayerId))
			}
		} else {
			if err := dbPlayer.Insert(getCtx(), db, boil.Infer()); err != nil {
				log.Error("failed to insert new player", zap.String("id", joinedPlayer.PlayerId), zap.Error(err))
			}
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
		if _, err = dbPlayer.Update(getCtx(), db, boil.Whitelist(orm.PlayerColumns.ConnID)); err != nil {
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

		player, err := orm.FindPlayer(getCtx(), db, playerId)
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

		if _, err = player.Update(getCtx(), db,
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

		dbPlayer, err := orm.FindPlayer(getCtx(), tx, playerId)
		if err == sql.ErrNoRows {
			log.Warn("received posUpdate for nonexistent player, ignoring", zap.String("id", inventory.PlayerId))
			return
		} else if err != nil {
			log.Error("failed to fetch player by UUID", zap.String("id", inventory.PlayerId), zap.Error(err))
			return
		}

		if _, err := orm.Inventories(orm.InventoryWhere.PlayerID.EQ(playerId)).DeleteAll(getCtx(), tx); err != nil {
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
			if err := dbItem.Insert(getCtx(), tx, boil.Infer()); err != nil {
				log.Error("failed to wipe player inventory", zap.String("id", inventory.PlayerId), zap.Error(err))
				_ = tx.Rollback()
				return
			}
		}

		dbPlayer.CurrentHotbar = int16(inventory.CurrentHotbar)
		if _, err = dbPlayer.Update(getCtx(), tx, boil.Whitelist(orm.PlayerColumns.CurrentHotbar)); err != nil {
			log.Error("failed to save user hotbar active slot", zap.String("id", inventory.PlayerId), zap.Error(err))
		}

		if err := tx.Commit(); err != nil {
			log.Error("failed to commit transaction", zap.String("id", inventory.PlayerId), zap.Error(err))
		}
	}
}
