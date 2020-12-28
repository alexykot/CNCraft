package handlers

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// RegisterHandlersState3 registers handlers for envelopes broadcast in the Play connection state.
func RegisterHandlersState3(ps nats.PubSub, log *zap.Logger, tally *players.Tally) error {
	if err := ps.Subscribe(subj.MkPlayerLoading(), handlePlayerLoading(ps, log, tally)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	return nil
}

func handlePlayerLoading(ps nats.PubSub, log *zap.Logger, tally *players.Tally) func(lope *envelope.E) {

	hander := func(lope *envelope.E) {
		//ps := ps
		log := log
		tally := tally

		loading := lope.GetPlayerLoading()
		if loading == nil {
			log.Error("failed to parse envelope - no LoadingPlayer inside", zap.Any("envelope", lope))
			return
		}

		userId, err := uuid.Parse(loading.Id)
		if err != nil {
			log.Error("failed to parse user ID as UUID", zap.String("id", loading.Id))
			return
		}
		p := tally.AddPlayer(userId, loading.Username)
		currentWorld := world.GetWorld()

		cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CJoinGame) // Predefined packet is expected to always exist.
		joinGame := cpacket.(*protocol.CPacketJoinGame)                           // And always be of the correct type.

		joinGame.EntityID = p.PC.ID()
		joinGame.GameMode = currentWorld.Gamemode
		joinGame.Dimension = game.Overworld
		joinGame.IsHardcore = currentWorld.Coreness
		joinGame.LevelType = currentWorld.Type
		joinGame.HashedSeed = int64(binary.LittleEndian.Uint64(currentWorld.SeedHash[:]))
		joinGame.ViewDistance = p.Settings.ViewDistance
		joinGame.RespawnScreen = currentConf.EnableRespawnScreen
	}

	return hander
}
