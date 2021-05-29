package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/pkg/game/items"

	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/game/data"
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

func HandleSClickWindow(connID uuid.UUID, inventory *items.Inventory,
	inventoryUpdater func(uuid.UUID), log *zap.Logger, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	windowClick, ok := sPacket.(*protocol.SPacketClickWindow)
	if !ok {
		return nil, fmt.Errorf("received packet is not a clickWindow: %v", sPacket)
	}

	cPacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CWindowConfirmation) // Predefined packet is expected to always exist.
	windowConfirm := cPacket.(*protocol.CPacketWindowConfirmation)                      // And always be of the correct type.

	windowConfirm.WindowID = windowClick.WindowID
	windowConfirm.Accepted = true
	windowConfirm.ActionID = windowClick.ActionID

	switch windowClick.WindowID {
	case items.InventoryWindow:
		droppedItem, isInventoryUpdated, err := inventory.HandleClick(
			windowClick.ActionID, windowClick.SlotID, windowClick.Mode, windowClick.Button, windowClick.ClickedItem)
		if err != nil {
			log.Warn("invalid window click received", zap.Error(err), zap.String("conn", connID.String()))
			windowConfirm.Accepted = false
			break
		}

		if isInventoryUpdated {
			inventoryUpdater(connID)
		}

		if droppedItem != nil {
			// TODO handle dropped item
		}
	default:
		return nil, fmt.Errorf("window ID %d is not implemented", windowClick.WindowID)
	}

	return []protocol.CPacket{windowConfirm}, nil
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
