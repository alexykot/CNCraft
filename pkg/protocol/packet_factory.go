// Package protocol defines the packets used in the Minecraft wire protocol.
// Currently supported protocol version is v578, for Minecraft 1.15.2.
package protocol

import (
	"fmt"

	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/apis/logs"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/protocol/server"
)

type Packet interface {
	ID() PacketID
}

type SPacket interface {
	Packet

	// decode the server_data from provided reader into this packet
	Pull(reader buff.Buffer, conn base.Connection) error
}

type CPacket interface {
	Packet

	// encode the server_data from this packet into provided writer
	Push(writer buff.Buffer, conn base.Connection)
}

type PacketFactory interface {
	GetSPacket(id PacketID) (SPacket, error)
	GetCPacket(id PacketID) (CPacket, error)
}

type packets struct {
	logger   *logs.Logging
	sPackets map[PacketID]func() SPacket

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection
}

func NewPacketFactory() PacketFactory {
	return &packets{
		logger:   logs.NewLogging("protocol", logs.EveryLevel...),
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
			return &server.SPacketHandshake{}
		},

		// Status state
		SRequest: func() SPacket {
			return &server.SPacketRequest{}
		},
		SPing: func() SPacket {
			return &server.SPacketPing{}
		},

		// Login state
		SLoginStart: func() SPacket {
			return &server.SPacketLoginStart{}
		},
		SEncryptionResponse: func() SPacket {
			return &server.SPacketEncryptionResponse{}
		},
		SLoginPluginResponse: func() SPacket {
			return &server.SPacketLoginPluginResponse{}
		},

		// Play state
		STeleportConfirm: func() SPacket {
			return &server.SPacketTeleportConfirm{}
		},
		SQueryBlockNBT: func() SPacket {
			return &server.SPacketQueryBlockNBT{}
		},
		SSetDifficulty: func() SPacket {
			return &server.SPacketSetDifficulty{}
		},
		SChatMessage: func() SPacket {
			return &server.SPacketChatMessage{}
		},
		SClientStatus: func() SPacket {
			return &server.SPacketClientStatus{}
		},
		SClientSettings: func() SPacket {
			return &server.SPacketClientSettings{}
		},
		SPluginMessage: func() SPacket {
			return &server.SPacketPluginMessage{}
		},
		SKeepAlive: func() SPacket {
			return &server.SPacketKeepAlive{}
		},
		SPlayerPosition: func() SPacket {
			return &server.SPacketPlayerPosition{}
		},
		SPlayerLocation: func() SPacket {
			return &server.SPacketPlayerLocation{}
		},
		SPlayerRotation: func() SPacket {
			return &server.SPacketPlayerRotation{}
		},
		SPlayerAbilities: func() SPacket {
			return &server.SPacketPlayerAbilities{}
		},
	}
}
