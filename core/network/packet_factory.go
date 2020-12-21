package network

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol"
)

type PacketFactory interface {
	MakeSPacket(id protocol.PacketID) (protocol.SPacket, error)
	MakeCPacket(id protocol.PacketID) (protocol.CPacket, error)
}

type packetFactory struct {
	sPackets map[protocol.PacketID]func() protocol.SPacket
}

func NewPacketFactory() PacketFactory {
	return &packetFactory{
		sPackets: createSPacketsMap(),
	}
}

func (p *packetFactory) MakeCPacket(newPacketID protocol.PacketID) (protocol.CPacket, error) {
	return nil, fmt.Errorf("packetFactory.MakeCPacket is not implemented")
}

func (p *packetFactory) MakeSPacket(newPacketID protocol.PacketID) (protocol.SPacket, error) {
	creator, ok := p.sPackets[newPacketID]
	if !ok {
		return nil, fmt.Errorf("packet ID %d is not recognised", newPacketID)
	}

	return creator(), nil
}

func createSPacketsMap() map[protocol.PacketID]func() protocol.SPacket {
	return map[protocol.PacketID]func() protocol.SPacket{
		// Handshake state
		protocol.SHandshake: func() protocol.SPacket {
			return &protocol.SPacketHandshake{}
		},

		// Status state
		protocol.SRequest: func() protocol.SPacket {
			return &protocol.SPacketRequest{}
		},
		protocol.SPing: func() protocol.SPacket {
			return &protocol.SPacketPing{}
		},

		// Login state
		protocol.SLoginStart: func() protocol.SPacket {
			return &protocol.SPacketLoginStart{}
		},
		protocol.SEncryptionResponse: func() protocol.SPacket {
			return &protocol.SPacketEncryptionResponse{}
		},
		protocol.SLoginPluginResponse: func() protocol.SPacket {
			return &protocol.SPacketLoginPluginResponse{}
		},

		// Play state
		protocol.STeleportConfirm: func() protocol.SPacket {
			return &protocol.SPacketTeleportConfirm{}
		},
		protocol.SQueryBlockNBT: func() protocol.SPacket {
			return &protocol.SPacketQueryBlockNBT{}
		},
		protocol.SSetDifficulty: func() protocol.SPacket {
			return &protocol.SPacketSetDifficulty{}
		},
		protocol.SChatMessage: func() protocol.SPacket {
			return &protocol.SPacketChatMessage{}
		},
		protocol.SClientStatus: func() protocol.SPacket {
			return &protocol.SPacketClientStatus{}
		},
		protocol.SClientSettings: func() protocol.SPacket {
			return &protocol.SPacketClientSettings{}
		},
		// TODO plugins are not supported atm
		//SPluginMessage: func() SPacket {
		//	return &SPacketPluginMessage{}
		//},
		protocol.SKeepAlive: func() protocol.SPacket {
			return &protocol.SPacketKeepAlive{}
		},
		protocol.SPlayerPosition: func() protocol.SPacket {
			return &protocol.SPacketPlayerPosition{}
		},
		protocol.SPlayerLocation: func() protocol.SPacket {
			return &protocol.SPacketPlayerLocation{}
		},
		protocol.SPlayerRotation: func() protocol.SPacket {
			return &protocol.SPacketPlayerRotation{}
		},
		protocol.SPlayerAbilities: func() protocol.SPacket {
			return &protocol.SPacketPlayerAbilities{}
		},
	}
}