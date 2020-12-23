package handlers

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/status"
)

// HandleSPing handles the Ping packet.
func HandleSPing(pacFac protocol.PacketFactory, spacket protocol.SPacket) (protocol.CPacket, error) {
	ping, ok := spacket.(*protocol.SPacketPing)
	if !ok {
		return nil, fmt.Errorf("received packet is not a ping: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CPong)       // Predefined packet is expected to always exist.
	cpacket.(*protocol.CPacketPong).Payload = ping.Payload // And always be of the correct type.
	return cpacket, nil
}

// HandleSRequest handles the StatusRequest packet.
func HandleSRequest(pacFac protocol.PacketFactory, spacket protocol.SPacket) (protocol.CPacket, error) {
	_, ok := spacket.(*protocol.SPacketRequest)
	if !ok {
		return nil, fmt.Errorf("received packet is not a status request: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CResponse)                                  // Predefined packet is expected to always exist.
	cpacket.(*protocol.CPacketResponse).Status = status.DefaultResponse(protocol.Version) // And always be of the correct type.
	return cpacket, nil
}
