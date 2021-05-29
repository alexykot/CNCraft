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

func HandleSClickWindow(connID uuid.UUID, inventory *items.Inventory, inventoryUpdater func(uuid.UUID), sPacket protocol.SPacket) error {
	click, ok := sPacket.(*protocol.SPacketClickWindow)
	if !ok {
		return fmt.Errorf("received packet is not a clickWindow: %v", sPacket)
	}

	switch click.WindowID {
	case uint8(items.InventoryWindow):
		droppedItem, isInventoryUpdated, err := inventory.HandleClick(click.ActionID, click.SlotID, click.Mode, click.Button, click.ClickedItem)
		if err != nil {
			// TODO this is not a packet handling error, this is an illegitimate inventory action from the client.
			//  This needs to trigger negative WindowConfirmation response packet and put inventory into "upset" state.
			return fmt.Errorf("error handling inventory click: %w", err)
		}

		if isInventoryUpdated {
			inventoryUpdater(connID)
		}

		if droppedItem != nil {
			// TODO handle dropped item
		}
	default:
		return fmt.Errorf("window ID %d is not implemented", click.WindowID)
	}

	return nil
}

func HandleSCloseWindow(player *players.Player, sPacket protocol.SPacket) error {
	closeWindow, ok := sPacket.(*protocol.SPacketCloseWindow)
	if !ok {
		return fmt.Errorf("received packet is not a closeWindow: %v", sPacket)
	}

	if closeWindow.WindowID == items.InventoryWindow {
		droppedItem := player.State.Inventory.CloseWindow()
		if droppedItem.IsPresent {
			// TODO handle dropped item
		}
	} else {
		return fmt.Errorf("cannot handle closeWindow for WindowID: %d", closeWindow.WindowID)
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
