package handlers

import (
	"fmt"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/pkg/protocol"
)

var currentConf control.ServerConf

func RegisterConf(serverConfig control.ServerConf) {
	currentConf = serverConfig
}

// HandleSHandshake handles the Handshake packet.
func HandleSHandshake(stateSetter func(state protocol.State), spacket protocol.SPacket) error {
	packet, ok := spacket.(*protocol.SPacketHandshake)
	if !ok {
		return fmt.Errorf("received packet is not a handshake: %v", spacket)
	}

	switch packet.NextState {
	case protocol.Handshake, protocol.Status, protocol.Login:
		stateSetter(packet.NextState)
		return nil
	}
	return fmt.Errorf("unexpected next state received: %d", packet.NextState)
}
