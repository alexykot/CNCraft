package protocol

import (
	"fmt"
)

type PacketFactory interface {
	MakeSPacket(id PacketType) (SPacket, error)
	MakeCPacket(id PacketType) (CPacket, error)
}

type packetFactory struct {
	sPackets map[PacketType]func() SPacket
	cPackets map[PacketType]func() CPacket
}

var factorySingleton packetFactory

func init() {
	factorySingleton = packetFactory{
		sPackets: createSPacketsMap(),
		cPackets: createCPacketsMap(),
	}
}

// GetPacketFactory returns the singleton of the packet factory object.
func GetPacketFactory() PacketFactory {
	return &factorySingleton
}

func (p *packetFactory) MakeCPacket(packetType PacketType) (CPacket, error) {
	creator, ok := p.cPackets[packetType]
	if !ok {
		return nil, fmt.Errorf("packet Type %d is not recognised", packetType)
	}

	return creator(), nil
}

func (p *packetFactory) MakeSPacket(packetType PacketType) (SPacket, error) {
	creator, ok := p.sPackets[packetType]
	if !ok {
		return nil, fmt.Errorf("packet Type %d is not recognised", packetType)
	}

	return creator(), nil
}

func createSPacketsMap() map[PacketType]func() SPacket {
	return map[PacketType]func() SPacket{
		// Handshake state
		SHandshake: func() SPacket {
			return &SPacketHandshake{}
		},

		// Status state
		SRequest: func() SPacket {
			return &SPacketRequest{}
		},
		SPing: func() SPacket {
			return &SPacketPing{}
		},

		// Login state
		SLoginStart: func() SPacket {
			return &SPacketLoginStart{}
		},
		SEncryptionResponse: func() SPacket {
			return &SPacketEncryptionResponse{}
		},
		SLoginPluginResponse: func() SPacket {
			return &SPacketLoginPluginResponse{}
		},

		// Play state
		STeleportConfirm: func() SPacket {
			return &SPacketTeleportConfirm{}
		},
		SQueryBlockNBT: func() SPacket {
			return &SPacketQueryBlockNBT{}
		},
		SSetDifficulty: func() SPacket {
			return &SPacketSetDifficulty{}
		},
		SChatMessage: func() SPacket {
			return &SPacketChatMessage{}
		},
		SClientStatus: func() SPacket {
			return &SPacketClientStatus{}
		},
		SClientSettings: func() SPacket {
			return &SPacketClientSettings{}
		},
		// TODO plugins are not supported atm
		//SPluginMessage: func() SPacket {
		//	return &SPacketPluginMessage{}
		//},
		SKeepAlive: func() SPacket {
			return &SPacketKeepAlive{}
		},
		SPlayerPosition: func() SPacket {
			return &SPacketPlayerPosition{}
		},
		SPlayerLocation: func() SPacket {
			return &SPacketPlayerLocation{}
		},
		SPlayerRotation: func() SPacket {
			return &SPacketPlayerRotation{}
		},
		SPlayerAbilities: func() SPacket {
			return &SPacketPlayerAbilities{}
		},
	}
}

func createCPacketsMap() map[PacketType]func() CPacket {
	return map[PacketType]func() CPacket{
		// Status state packets
		CResponse: func() CPacket {
			return &CPacketResponse{}
		},
		CPong: func() CPacket {
			return &CPacketPong{}
		},

		// Login state packets
		CDisconnect: func() CPacket {
			return &CPacketDisconnect{}
		},
		CEncryptionRequest: func() CPacket {
			return &CPacketEncryptionRequest{}
		},
		CLoginSuccess: func() CPacket {
			return &CPacketLoginSuccess{}
		},
		CSetCompression: func() CPacket {
			return &CPacketSetCompression{}
		},
		CLoginPluginRequest: func() CPacket {
			return &CPacketLoginPluginRequest{}
		},

		// Play state packets
		CChatMessage: func() CPacket {
			return &CPacketChatMessage{}
		},
		CJoinGame: func() CPacket {
			return &CPacketJoinGame{}
		},
		// TODO plugins are not supported atm
		//protocol.CPluginMessage: func() protocol.CPacket{
		//	return &protocol.CPacketPluginMessage{}
		//},
		CPlayerLocation: func() CPacket {
			return &CPacketPlayerLocation{}
		},
		CKeepAlive: func() CPacket {
			return &CPacketKeepAlive{}
		},
		CServerDifficulty: func() CPacket {
			return &CPacketServerDifficulty{}
		},
		CPlayerAbilities: func() CPacket {
			return &CPacketPlayerAbilities{}
		},
		CHeldItemChange: func() CPacket {
			return &CPacketHeldItemChange{}
		},
		CDeclareRecipes: func() CPacket {
			return &CPacketDeclareRecipes{}
		},
		CChunkData: func() CPacket {
			return &CPacketChunkData{}
		},
		CPlayerInfo: func() CPacket {
			return &CPacketPlayerInfo{}
		},
		CEntityMetadata: func() CPacket {
			return &CPacketEntityMetadata{}
		},
	}
}
