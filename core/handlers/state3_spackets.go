package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/users"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
)

func HandleSPluginMessage(log *zap.Logger, tally *users.Roster, connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
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

	current := tally.GetPlayerSettings(userID)
	current.ClientBrand = brand.Name
	tally.SetPlayerSettings(userID, current)

	return nil, nil
}

func HandleSClientSettings(tally *users.Roster, connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	clientSettings, ok := sPacket.(*protocol.SPacketClientSettings)
	if !ok {
		return nil, fmt.Errorf("received packet is not a pluginMessage: %v", sPacket)
	}

	userID := connID // By design connection ID is also the auth user ID and then the player ID.

	current := tally.GetPlayerSettings(userID)
	current.Locale = clientSettings.Locale
	current.ViewDistance = int32(clientSettings.ViewDistance)
	current.Skin = clientSettings.SkinParts
	current.ChatMode = clientSettings.ChatMode
	current.ChatColors = clientSettings.ChatColors
	tally.SetPlayerSettings(userID, current)

	return nil, nil
}
