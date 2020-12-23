package handlers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/mojang"
)

// HandleSRequest handles the StatusRequest packet.
func HandleSLoginStart(newbieAdder func(uuid.UUID, string), crypter mojang.RSACrypter, pacFac protocol.PacketFactory,
	connID uuid.UUID, spacket protocol.SPacket) (protocol.CPacket, error) {
	loginStart, ok := spacket.(*protocol.SPacketLoginStart)
	if !ok {
		return nil, fmt.Errorf("received packet is not a loginStart: %v", spacket)
	}

	// By design connection ID is also the user ID and then the player ID.
	newbieAdder(connID, loginStart.Username)

	if currentConf.IsCracked {
		cpacket, _ := pacFac.MakeCPacket(protocol.CLoginSuccess)             // Predefined packet is expected to always exist.
		cpacket.(*protocol.CPacketLoginSuccess).PlayerUUID = connID.String() // And always be of the correct type.
		cpacket.(*protocol.CPacketLoginSuccess).PlayerName = loginStart.Username
		return cpacket, nil
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CEncryptionRequest) // Predefined packet is expected to always exist.
	encRequest := cpacket.(*protocol.CPacketEncryptionRequest)    // And always be of the correct type.

	encRequest.ServerID = currentConf.ServerID
	encRequest.PublicKey = crypter.GetPubKey()

	return encRequest, nil
}
