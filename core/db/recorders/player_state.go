package recorders

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
)

// RegisterPlayerStateHandlers registers handlers for envelopes carrying updates of the player state that
// need to be persisted.
func RegisterPlayerStateHandlers(ps nats.PubSub, log *zap.Logger, db *sql.DB) error {
	if err := ps.Subscribe(subj.MkPlayerPosUpdate(), handlePlayerPos(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	if err := ps.Subscribe(subj.MkNewPlayerJoined(), handleNewPlayer(log, db)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	return nil
}

func handleNewPlayer(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		log := log
		db := db

		newPlayer := inLope.GetNewPlayer()
		if newPlayer == nil {
			log.Error("envelope does not contain new player", zap.Any("envelope", inLope))
			return
		}

		dbPlayer := &orm.Player{
			ID:        newPlayer.Id,
			Username:  newPlayer.Username,
			PositionX: newPlayer.Pos.X,
			PositionY: newPlayer.Pos.Y,
			PositionZ: newPlayer.Pos.Z,
		}
		if err := dbPlayer.Insert(context.TODO(), db, boil.Infer()); err != nil {
			log.Error("failed to insert player", zap.String("id", newPlayer.Id), zap.Error(err))
			return
		}
	}
}

func handlePlayerPos(log *zap.Logger, db *sql.DB) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		log := log
		db := db

		posUpdate := inLope.GetPlayerPos()
		if posUpdate == nil {
			log.Error("envelope does not contain position update", zap.Any("envelope", inLope))
			return
		}
		player, err := orm.FindPlayer(context.TODO(), db, posUpdate.Id)
		if err == sql.ErrNoRows {
			log.Warn("received posUpdate for nonexistent player, ignoring", zap.String("id", posUpdate.Id))
			return
		} else if err != nil {
			log.Error("failed to fetch user by UUID", zap.String("id", posUpdate.Id), zap.Error(err))
			return
		}

		player.PositionX = posUpdate.Pos.X
		player.PositionY = posUpdate.Pos.Y
		player.PositionZ = posUpdate.Pos.Z

		if _, err = player.Update(context.TODO(), db,
			boil.Whitelist(orm.PlayerColumns.PositionX, orm.PlayerColumns.PositionY, orm.PlayerColumns.PositionZ)); err != nil {
			log.Error("failed to save user position", zap.String("id", posUpdate.Id), zap.Error(err))
			return
		}
	}
}
