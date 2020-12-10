// Package protocol defines the packets used in the Minecraft wire protocol.
// Currently supported protocol version is v578, for Minecraft 1.15.2.
package protocol

import (
	"fmt"

	"go.uber.org/zap"

	buff "github.com/alexykot/cncraft/pkg/buffers"
)

type Packet interface {
	ID() PacketID
}

type SPacket interface {
	Packet

	// decode the server_data from provided reader into this packet
	Pull(reader buff.Buffer) error
}

type CPacket interface {
	Packet

	// encode the server_data from this packet into provided writer
	Push(writer buff.Buffer)
}

type PacketFactory interface {
	GetSPacket(id PacketID) (SPacket, error)
	GetCPacket(id PacketID) (CPacket, error)
}

type packets struct {
	log      *zap.Logger
	sPackets map[PacketID]func() SPacket
}

func NewPacketFactory(log *zap.Logger) PacketFactory {
	return &packets{
		log:      log,
		sPackets: createSPacketsMap(),
	}
}

func (p *packets) GetCPacket(newPacketID PacketID) (CPacket, error) {
	return nil, fmt.Errorf("paketFactory.GetCPacket is not implemented")
}

func (p *packets) GetSPacket(newPacketID PacketID) (SPacket, error) {
	creator, ok := p.sPackets[newPacketID]
	if !ok {
		return nil, fmt.Errorf("packet ID %d is not recognised", newPacketID)
	}

	return creator(), nil
}

func createSPacketsMap() map[PacketID]func() SPacket {
	return map[PacketID]func() SPacket{
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
