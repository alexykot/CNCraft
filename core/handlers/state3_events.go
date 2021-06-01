package handlers

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
)

// RegisterEventHandlersState3 registers handlers for envelopes broadcast in the Play connection state.
//  Play state handlers are entirely asynchronous, so NATS subscriptions need to be created at boot time.
func RegisterEventHandlersState3(log *zap.Logger, ps nats.PubSub, roster *players.Roster, world *world.World) error {
	if err := ps.Subscribe(subj.MkPlayerLoading(), handlePlayerLoading(ps, log, roster, world)); err != nil {
		return fmt.Errorf("failed to register PlayerLoading handler: %w", err)
	}

	log.Info("Play state event handlers registered")

	return nil
}

func handlePlayerLoading(ps nats.PubSub, log *zap.Logger, roster *players.Roster, world *world.World) func(lope *envelope.E) {
	return func(inLope *envelope.E) {
		ps := ps
		log := log

		loading := inLope.GetPlayerLoading()
		if loading == nil {
			log.Error("failed to parse envelope - no PlayerLoading inside", zap.Any("envelope", inLope))
			return
		}

		userId, err := uuid.Parse(loading.ConnId)
		if err != nil {
			log.Error("failed to parse user ID as UUID", zap.String("id", loading.ConnId), zap.Error(err))
			return
		}
		log.Debug("handling player loading", zap.String("user", userId.String()))

		p, err := roster.AddPlayer(userId, loading.Username)
		if err != nil {
			log.Error("failed add player", zap.Error(err))
			return
		}
		var outLopes []*envelope.E

		cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CJoinGame) // Predefined packet is expected to always exist.
		joinGame := cpacket.(*protocol.CPacketJoinGame)                           // And always be of the correct type.

		joinGame.EntityID = p.PC.ID()
		joinGame.GameMode = world.Gamemode
		joinGame.DimensionCodec = world.DimensionCodec
		joinGame.Dimension = world.Dimension
		joinGame.IsHardcore = world.Coreness
		joinGame.HashedSeed = int64(binary.LittleEndian.Uint64(world.SeedHash[:]))
		joinGame.ViewDistance = p.Settings.ViewDistance
		joinGame.EnableRespawnScreen = control.GetCurrentConfig().EnableRespawnScreen
		outLopes = append(outLopes, mkCpacketEnvelope(joinGame))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPluginMessage)
		pluginMessage := cpacket.(*protocol.CPacketPluginMessage)
		pluginMessage.Message = &plugin.Brand{Name: control.GetCurrentConfig().Brand}
		outLopes = append(outLopes, mkCpacketEnvelope(pluginMessage))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CServerDifficulty)
		difficulty := cpacket.(*protocol.CPacketServerDifficulty)
		difficulty.Difficulty = world.Difficulty
		difficulty.Locked = world.DifficultyIsLocked
		outLopes = append(outLopes, mkCpacketEnvelope(difficulty))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPlayerAbilities)
		abilities := cpacket.(*protocol.CPacketPlayerAbilities)
		abilities.Abilities = *p.Abilities
		abilities.FlyingSpeed = p.Settings.FlyingSpeed
		abilities.FieldOfView = p.Settings.FoVModifier
		outLopes = append(outLopes, mkCpacketEnvelope(abilities))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CDeclareRecipes)
		declareRecipes := cpacket.(*protocol.CPacketDeclareRecipes)
		declareRecipes.RecipeCount = 0 // TODO probably will be a static list of recipes defined for current server version
		outLopes = append(outLopes, mkCpacketEnvelope(declareRecipes))

		// TODO CTags packet is not defined
		// TODO CEntityStatus packet is not defined
		// TODO CDeclareCommands packet is not defined
		// TODO CUnlockRecipes packet is not defined

		// TODO move this to a separate world loader
		chunksToLoad := world.Levels[game.Overworld.String()].Chunks()
		for _, chunk := range chunksToLoad {
			cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CChunkData)
			chunkData := cpacket.(*protocol.CPacketChunkData)
			chunkData.Chunk = chunk
			outLopes = append(outLopes, mkCpacketEnvelope(chunkData))
		}

		// Player Position And Look
		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CPlayerPositionAndLook)
		posAndLook := cpacket.(*protocol.CPacketPlayerPositionAndLook)
		posAndLook.Location = p.State.Location // Relative is always False here.
		outLopes = append(outLopes, mkCpacketEnvelope(posAndLook))

		// Player inventory init
		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CWindowItems)
		winItems := cpacket.(*protocol.CPacketWindowItems)
		inventorySlots := p.State.Inventory.ToArray()
		winItems.SlotCount = int16(len(inventorySlots))
		winItems.Slots = inventorySlots
		outLopes = append(outLopes, mkCpacketEnvelope(winItems))

		cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CHeldItemChange)
		heldItemChange := cpacket.(*protocol.CPacketHeldItemChange)
		heldItemChange.Slot = p.State.Inventory.CurrentHotbarSlot
		outLopes = append(outLopes, mkCpacketEnvelope(heldItemChange))

		if err := ps.Publish(subj.MkConnTransmit(userId), outLopes...); err != nil {
			log.Error("failed to publish conn.transmit message", zap.Error(err), zap.Any("conn", userId))
			return
		}
	}
}
