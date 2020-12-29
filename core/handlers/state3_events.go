package handlers

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/users"
	"github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
)

// RegisterEventHandlersState3 registers handlers for envelopes broadcast in the Play connection state.
//  Play state handlers are entirely asynchronous, so NATS subscriptions need to be created at boot time.
func RegisterEventHandlersState3(ps nats.PubSub, log *zap.Logger, tally *users.Roster) error {
	if err := ps.Subscribe(subj.MkPlayerLoading(), handlePlayerLoading(ps, log, tally)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	return nil
}

func handlePlayerLoading(ps nats.PubSub, log *zap.Logger, tally *users.Roster) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		ps := ps
		log := log
		tally := tally

		loading := inLope.GetPlayerLoading()
		if loading == nil {
			log.Error("failed to parse envelope - no PlayerLoading inside", zap.Any("envelope", inLope))
			return
		}

		userId, err := uuid.Parse(loading.Id)
		if err != nil {
			log.Error("failed to parse user ID as UUID", zap.String("id", loading.Id))
			return
		}
		log.Debug("handling player loading", zap.String("user", userId.String()))

		p := tally.AddPlayer(userId, loading.Username)
		currentWorld := world.GetWorld()
		var outLopes []*envelope.E

		cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CJoinGame) // Predefined packet is expected to always exist.
		joinGame := cpacket.(*protocol.CPacketJoinGame)                           // And always be of the correct type.

		joinGame.EntityID = p.PC.ID()
		joinGame.GameMode = currentWorld.Gamemode
		joinGame.Dimension = game.Overworld
		joinGame.IsHardcore = currentWorld.Coreness
		joinGame.LevelType = currentWorld.Type
		joinGame.HashedSeed = int64(binary.LittleEndian.Uint64(currentWorld.SeedHash[:]))
		joinGame.ViewDistance = p.Settings.ViewDistance
		joinGame.RespawnScreen = control.GetCurrentConfig().EnableRespawnScreen
		outLopes = append(outLopes, mkCpacketEnvelope(joinGame))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPluginMessage)
		pluginMessage := cpacket.(*protocol.CPacketPluginMessage)
		pluginMessage.Message = &plugin.Brand{Name: control.GetCurrentConfig().Brand}
		outLopes = append(outLopes, mkCpacketEnvelope(pluginMessage))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CServerDifficulty)
		difficulty := cpacket.(*protocol.CPacketServerDifficulty)
		difficulty.Difficulty = currentWorld.Difficulty
		difficulty.Locked = currentWorld.DifficultyIsLocked
		outLopes = append(outLopes, mkCpacketEnvelope(difficulty))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPlayerAbilities)
		abilities := cpacket.(*protocol.CPacketPlayerAbilities)
		abilities.Abilities = p.Abilities
		abilities.FlyingSpeed = p.Settings.FlyingSpeed
		abilities.FieldOfView = p.Settings.FoVModifier
		outLopes = append(outLopes, mkCpacketEnvelope(abilities))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CHeldItemChange)
		heldItemChange := cpacket.(*protocol.CPacketHeldItemChange)
		heldItemChange.Slot = p.State.CurrentHotbarSlot
		outLopes = append(outLopes, mkCpacketEnvelope(heldItemChange))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CDeclareRecipes)
		declareRecipes := cpacket.(*protocol.CPacketDeclareRecipes)
		declareRecipes.RecipeCount = 0 // TODO probably will be a static list of recipes defined for current server version
		outLopes = append(outLopes, mkCpacketEnvelope(declareRecipes))

		// TODO CTags packet is not defined
		// TODO CEntityStatus packet is not defined
		// TODO CDeclareCommands packet is not defined
		// TODO CUnlockRecipes packet is not defined

		// Player Position And Look
		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPlayerPositionAndLook)
		posAndLook := cpacket.(*protocol.CPacketPlayerPositionAndLook)
		posAndLook.Location = p.State.CurrentLocation // Relative is always False here.
		outLopes = append(outLopes, mkCpacketEnvelope(posAndLook))

		if err := ps.Publish(subj.MkConnTransmit(userId), outLopes...); err != nil {
			log.Error("failed to publish conn.transmit message", zap.Error(err), zap.Any("conn", userId))
			return
		}
	}
}
