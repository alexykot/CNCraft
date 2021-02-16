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

func HandleSPluginMessage(log *zap.Logger, roster *players.Roster, connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	pluginMessage, ok := sPacket.(*protocol.SPacketPluginMessage)
	if !ok {
		return nil, fmt.Errorf("received packet is not a pluginMessage: %v", sPacket)
	}

	userID := connID // By design connection ID is also the auth user ID and then the player ID.

	if pluginMessage.Message.Chan() != plugin.ChannelBrand {
		log.Warn("cannot handle messages for plugin channel", zap.String("chan", string(pluginMessage.Message.Chan())))
		return nil, nil
	}

	brand, ok := pluginMessage.Message.(*plugin.Brand)
	if !ok {
		return nil, fmt.Errorf("unexpected type of plugin message for channel %s", plugin.ChannelBrand)
	}

	player, ok := roster.GetPlayer(userID)
	if !ok {
		return nil, fmt.Errorf("player no found fo conn %s", userID)
	}

	current := player.GetSettings()
	current.ClientBrand = brand.Name
	player.SetSettings(current)

	return nil, nil
}

func HandleSClientSettings(roster *players.Roster, connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	clientSettings, ok := sPacket.(*protocol.SPacketClientSettings)
	if !ok {
		return nil, fmt.Errorf("received packet is not clientSettings: %v", sPacket)
	}

	userID := connID // By design connection ID is also the auth user ID and then the player ID.
	player, ok := roster.GetPlayer(userID)
	if !ok {
		return nil, fmt.Errorf("player no found fo conn %s", userID)
	}

	current := player.GetSettings()
	current.Locale = clientSettings.Locale
	current.ViewDistance = int32(clientSettings.ViewDistance)
	current.Skin = clientSettings.SkinParts
	current.ChatMode = clientSettings.ChatMode
	current.ChatColors = clientSettings.ChatColors
	player.SetSettings(current)

	return nil, nil
}

func HandleSKeepAlive(aliveRecorder func(uuid.UUID, int64), connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	keepAlive, ok := sPacket.(*protocol.SPacketKeepAlive)
	if !ok {
		return nil, fmt.Errorf("received packet is not a keepAlive: %v", sPacket)
	}

	aliveRecorder(connID, keepAlive.KeepAliveID)
	return nil, nil
}

func HandleSPlayerPosition(posSetter func(data.PositionF), sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	playerPos, ok := sPacket.(*protocol.SPacketPlayerPosition)
	if !ok {
		return nil, fmt.Errorf("received packet is not a keepAlive: %v", sPacket)
	}

	posSetter(playerPos.Position)
	return nil, nil
}
