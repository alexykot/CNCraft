package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/items"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
)

func HandleSPluginMessage(log *zap.Logger, player *players.Player, sPacket protocol.SPacket) error {
	pluginMessage, ok := sPacket.(*protocol.SPacketPluginMessage)
	if !ok {
		return fmt.Errorf("received packet is not a pluginMessage: %v", sPacket)
	}

	if pluginMessage.Message.Chan() != plugin.ChannelBrand {
		log.Warn("cannot handle messages for plugin channel", zap.String("chan", string(pluginMessage.Message.Chan())))
		return nil
	}

	brand, ok := pluginMessage.Message.(*plugin.Brand)
	if !ok {
		return fmt.Errorf("unexpected type of plugin message for channel %s", plugin.ChannelBrand)
	}

	current := player.GetSettings()
	current.ClientBrand = brand.Name
	player.SetSettings(current)

	return nil
}

func HandleSClientSettings(player *players.Player, sPacket protocol.SPacket) error {
	clientSettings, ok := sPacket.(*protocol.SPacketClientSettings)
	if !ok {
		return fmt.Errorf("received packet is not clientSettings: %v", sPacket)
	}

	current := player.GetSettings()
	current.Locale = clientSettings.Locale
	current.ViewDistance = int32(clientSettings.ViewDistance)
	current.Skin = clientSettings.SkinParts
	current.ChatMode = clientSettings.ChatMode
	current.ChatColors = clientSettings.ChatColors
	player.SetSettings(current)

	return nil
}

func HandleSKeepAlive(aliveRecorder func(uuid.UUID, int64), connID uuid.UUID, sPacket protocol.SPacket) error {
	keepAlive, ok := sPacket.(*protocol.SPacketKeepAlive)
	if !ok {
		return fmt.Errorf("received packet is not a keepAlive: %v", sPacket)
	}

	aliveRecorder(connID, keepAlive.KeepAliveID)
	return nil
}

func HandleSPlayerSpatial(locSetter func(uuid.UUID, *data.PositionF, *data.RotationF, *bool), connID uuid.UUID, sPacket protocol.SPacket) error {
	if playerPos, ok := sPacket.(*protocol.SPacketPlayerPosition); ok {
		locSetter(connID, &playerPos.Position, nil, nil)
		return nil
	}

	if playerMove, ok := sPacket.(*protocol.SPacketPlayerMovement); ok {
		locSetter(connID, nil, nil, &playerMove.OnGround)
		return nil
	}

	return fmt.Errorf("received packet is not a spatial update (must be one of playerPosition, playerMovement): %v", sPacket)
}

func HandleSHeldItemChange(heldItemSetter func(connID uuid.UUID, heldItem uint8), connID uuid.UUID, sPacket protocol.SPacket) error {
	heldItem, ok := sPacket.(*protocol.SPacketHeldItemChange)
	if !ok {
		return fmt.Errorf("received packet is not a heldItemChange: %v", sPacket)
	}

	heldItemSetter(connID, heldItem.Slot)
	return nil
}

func HandleSClickWindow(inventory *items.Inventory, log *zap.Logger, sPacket protocol.SPacket) (bool, []protocol.CPacket, error) {
	windowClick, ok := sPacket.(*protocol.SPacketClickWindow)
	if !ok {
		return false, nil, fmt.Errorf("received packet is not a clickWindow: %v", sPacket)
	}

	cPacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CWindowConfirmation) // Predefined packet is expected to always exist.
	windowConfirm := cPacket.(*protocol.CPacketWindowConfirmation)                      // And always be of the correct type.
	windowConfirm.WindowID = windowClick.WindowID
	windowConfirm.ActionID = windowClick.ActionID
	windowConfirm.Accepted = true

	var cPackets []protocol.CPacket
	var isInventoryUpdated bool

	switch windowClick.WindowID {
	case items.InventoryWindow:
		var err error
		var droppedItem *items.Slot
		droppedItem, isInventoryUpdated, err = inventory.HandleClick(
			windowClick.ActionID, windowClick.SlotID, windowClick.Mode, windowClick.Button, windowClick.ClickedItem)
		if err != nil {
			log.Warn("invalid window click received", zap.Error(err))
			windowConfirm.Accepted = false

			cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CWindowItems)
			windowItems := cpacket.(*protocol.CPacketWindowItems)
			inventorySlots := inventory.ToArray()
			windowItems.SlotCount = int16(len(inventorySlots))
			windowItems.Slots = inventorySlots

			cpacket, _ = protocol.GetPacketFactory().MakeCPacket(protocol.CSetSlot)
			setSlot := cpacket.(*protocol.CPacketSetSlot)
			setSlot.WindowID = items.CursorWindow
			setSlot.SlotID = items.CursorSlot
			setSlot.Slot = inventory.GetCursor()
			log.Debug("resetting cursor", zap.Any("cursor", setSlot.Slot))

			cPackets = append(cPackets, windowConfirm, windowItems, setSlot)
			break
		}

		if droppedItem != nil {
			// TODO handle dropped item
		}
		cPackets = append(cPackets, windowConfirm)
	default:
		return false, nil, fmt.Errorf("window ID %d is not implemented", windowClick.WindowID)
	}

	return isInventoryUpdated, cPackets, nil
}

func HandleSCloseWindow(player *players.Player, sPacket protocol.SPacket) error {
	closeWindow, ok := sPacket.(*protocol.SPacketCloseWindow)
	if !ok {
		return fmt.Errorf("received packet is not a closeWindow: %v", sPacket)
	}

	switch closeWindow.WindowID {
	case items.InventoryWindow:
		droppedItem := player.State.Inventory.CloseWindow()
		if droppedItem.IsPresent {
			// TODO handle dropped item
		}
	default:
		return fmt.Errorf("window ID %d is not implemented", closeWindow.WindowID)
	}

	return nil
}

func HandleSWindowConfirmation(inventory *items.Inventory, sPacket protocol.SPacket) error {
	windowConfirm, ok := sPacket.(*protocol.SPacketWindowConfirmation)
	if !ok {
		return fmt.Errorf("received packet is not a windowConfirmation: %v", sPacket)
	}

	switch windowConfirm.WindowID {
	case items.InventoryWindow:
		inventory.Apologise(windowConfirm.ActionID)
	default:
		return fmt.Errorf("window ID %d is not implemented", windowConfirm.WindowID)
	}

	return nil
}

// HandleSEntityAction - TODO nothing to do here yet, to implement later
func HandleSEntityAction(sPacket protocol.SPacket) error {
	if _, ok := sPacket.(*protocol.SPacketEntityAction); !ok {
		return fmt.Errorf("received packet is not a closeWindow: %v", sPacket)
	}

	return nil
}

// HandleSAnimation - TODO nothing to do here yet, to implement later
func HandleSAnimation(sPacket protocol.SPacket) error {
	if _, ok := sPacket.(*protocol.SPacketAnimation); !ok {
		return fmt.Errorf("received packet is not Animation: %v", sPacket)
	}

	return nil
}

func HandleSPlayerDigging(ps nats.PubSub, sharder *world.Sharder, player *players.Player, sPacket protocol.SPacket) error {
	dig, ok := sPacket.(*protocol.SPacketPlayerDigging)
	if !ok {
		return fmt.Errorf("received packet is not a heldItemChange: %v", sPacket)
	}

	// DEBT this should check player position and ensure dig position is legal
	shardID, ok := sharder.FindShardID(player.State.Dimension, dig.Position)
	if !ok {
		return fmt.Errorf("could not find shard for coords provided: x.%d z.%d", dig.Position.X, dig.Position.Z)
	}

	lope := envelope.PlayerDigging(&pb.PlayerDigging{
		PlayerId: player.ConnID.String(),
		Action:   pb.PlayerDigging_Action(dig.Status),
		Pos: &pb.Position{
			X: float64(dig.Position.X),
			Y: float64(dig.Position.Y),
			Z: float64(dig.Position.Z),
		},
		BlockFace: pb.BlockFace(dig.Face),
	})

	if err := ps.Publish(subj.MkShardEvent(string(shardID)), lope); err != nil {
		return fmt.Errorf("failed to publish shard PlayerDigging event: %w", err)
	}

	return nil
}
