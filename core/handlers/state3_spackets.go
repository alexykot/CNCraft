package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

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

// HandleSCloseWindow
// TODO nothing to do here yet, to implement later
func HandleSCloseWindow(sPacket protocol.SPacket) error {
	if _, ok := sPacket.(*protocol.SPacketCloseWindow); !ok {
		return fmt.Errorf("received packet is not a closeWindow: %v", sPacket)
	}

	return nil
}

// HandleSEntityAction
// TODO nothing to do here yet, to implement later
func HandleSEntityAction(sPacket protocol.SPacket) error {
	if _, ok := sPacket.(*protocol.SPacketEntityAction); !ok {
		return fmt.Errorf("received packet is not a closeWindow: %v", sPacket)
	}

	return nil
}
