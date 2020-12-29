//go:generate stringer -type=PacketType packets.go

// Package protocol defines the packets used in the Minecraft wire protocol.
// Currently supported protocol version is v578, for Minecraft 1.15.2.
// DEBT protocol should be interfaced and swappable implementation for simple plugging of different protocol versions.
//  This is assuming that the meaning and naming of the packets does not change between protocol versions and only
//  new packets are added over time. Assumption to be verified.
package protocol

import (
	"github.com/alexykot/cncraft/pkg/buffer"
)

// Version defines the version of the minecraft wire protocol the current implementation supports.
const Version = 578

// SPacket is server bound packet type.
type SPacket interface {
	ProtocolID() ProtocolPacketID // Return protocol ID of the packet.
	Type() PacketType             // Return globally unique type ID of the packet.
	Pull(reader buffer.B) error   // decode the server_data from provided reader into this packet
}

// CPacket is client bound packet type.
type CPacket interface {
	ProtocolID() ProtocolPacketID // Return protocol ID of the packet.
	Type() PacketType             // Return globally unique type ID of the packet.
	Push(writer buffer.B)         // encode the server_data from this packet into provided writer
}

// ProtocolPacketID is the official Type of the packet as per the protocol.
type ProtocolPacketID int32

// server bound (incoming) packets, protocol definitions
const (
	// Handshake state packets
	protocolSHandshake ProtocolPacketID = 0x00

	// Status state packets
	protocolSRequest ProtocolPacketID = 0x00
	protocolSPing    ProtocolPacketID = 0x01

	// Login state packets
	protocolSLoginStart          ProtocolPacketID = 0x00
	protocolSEncryptionResponse  ProtocolPacketID = 0x01
	protocolSLoginPluginResponse ProtocolPacketID = 0x02

	// Play state packets
	protocolSTeleportConfirm ProtocolPacketID = 0x00
	protocolSQueryBlockNBT   ProtocolPacketID = 0x01
	protocolSSetDifficulty   ProtocolPacketID = 0x02
	protocolSChatMessage     ProtocolPacketID = 0x03
	protocolSClientStatus    ProtocolPacketID = 0x04
	protocolSClientSettings  ProtocolPacketID = 0x05
	protocolSPluginMessage   ProtocolPacketID = 0x0B
	protocolSKeepAlive       ProtocolPacketID = 0x0F
	protocolSPlayerPosition  ProtocolPacketID = 0x11
	protocolSPlayerLocation  ProtocolPacketID = 0x12
	protocolSPlayerRotation  ProtocolPacketID = 0x13
	protocolSPlayerAbilities ProtocolPacketID = 0x19
)

// client bound (outgoing) packets
const (
	// Handshake state packets
	// no client bound handshake packets defined in the protocol

	// Status state packets
	protocolCResponse ProtocolPacketID = 0x00
	protocolCPong     ProtocolPacketID = 0x01

	// Login state packets
	protocolCDisconnect         ProtocolPacketID = 0x00
	protocolCEncryptionRequest  ProtocolPacketID = 0x01
	protocolCLoginSuccess       ProtocolPacketID = 0x02
	protocolCSetCompression     ProtocolPacketID = 0x03
	protocolCLoginPluginRequest ProtocolPacketID = 0x04

	// Play state packets
	protocolCChatMessage           ProtocolPacketID = 0x0F
	protocolCJoinGame              ProtocolPacketID = 0x26
	protocolCPluginMessage         ProtocolPacketID = 0x19
	protocolCPlayerPositionAndLook ProtocolPacketID = 0x36
	protocolCKeepAlive             ProtocolPacketID = 0x21
	protocolCServerDifficulty      ProtocolPacketID = 0x0E
	protocolCPlayerAbilities       ProtocolPacketID = 0x32
	protocolCHeldItemChange        ProtocolPacketID = 0x40
	protocolCDeclareRecipes        ProtocolPacketID = 0x5B
	protocolCChunkData             ProtocolPacketID = 0x22
	protocolCPlayerInfo            ProtocolPacketID = 0x34
	protocolCEntityMetadata        ProtocolPacketID = 0x44
)

type packetDirection int32

const ServerBound = packetDirection(0x1000)

const ClientBound = packetDirection(0xF000)

// PacketType combines direction (1 - server, 2 - client), state Type and the actual protocol Type to make a globally unique PacketType.
// E.g. PacketType 0x1101 means:
//  0x 1 1 01
//     ^ ^ ^^--- the protocol packet Type for server-bound Ping packet;
//     | |------ connection state, 1 for Status state;
//     |-------- server bound packet (1 - server, F - client);
type PacketType int32

func (i PacketType) Value() int32 { return int32(i) }
func (i PacketType) ProtocolID() ProtocolPacketID {
	return ProtocolPacketID(0x00FF & int32(i))
}

const stateShake = 0x0000
const stateStatus = 0x0100
const stateLogin = 0x0200
const statePlay = 0x0300

const TypeUnspecified = PacketType(-0x0001) // packet type unspecified

// server bound (incoming) packets
const (
	// Handshake state packets
	SHandshake = PacketType(int32(ServerBound) + stateShake + int32(protocolSHandshake)) // 0x1000

	// Status state packets
	SRequest = PacketType(int32(ServerBound) + stateStatus + int32(protocolSRequest)) // 0x1100
	SPing    = PacketType(int32(ServerBound) + stateStatus + int32(protocolSPing))    // 0x1101

	// Login state packets
	SLoginStart          = PacketType(int32(ServerBound) + stateLogin + int32(protocolSLoginStart))          // 0x1200
	SEncryptionResponse  = PacketType(int32(ServerBound) + stateLogin + int32(protocolSEncryptionResponse))  // 0x1201
	SLoginPluginResponse = PacketType(int32(ServerBound) + stateLogin + int32(protocolSLoginPluginResponse)) // 0x1202

	// Play state packets
	STeleportConfirm = PacketType(int32(ServerBound) + statePlay + int32(protocolSTeleportConfirm)) // 0x1300
	SQueryBlockNBT   = PacketType(int32(ServerBound) + statePlay + int32(protocolSQueryBlockNBT))   // 0x1301
	SSetDifficulty   = PacketType(int32(ServerBound) + statePlay + int32(protocolSSetDifficulty))   // 0x1302
	SChatMessage     = PacketType(int32(ServerBound) + statePlay + int32(protocolSChatMessage))     // 0x1303
	SClientStatus    = PacketType(int32(ServerBound) + statePlay + int32(protocolSClientStatus))    // 0x1304
	SClientSettings  = PacketType(int32(ServerBound) + statePlay + int32(protocolSClientSettings))  // 0x1305
	SPluginMessage   = PacketType(int32(ServerBound) + statePlay + int32(protocolSPluginMessage))   // 0x130B
	SKeepAlive       = PacketType(int32(ServerBound) + statePlay + int32(protocolSKeepAlive))       // 0x130F
	SPlayerPosition  = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerPosition))  // 0x1311
	SPlayerLocation  = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerLocation))  // 0x1312
	SPlayerRotation  = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerRotation))  // 0x1313
	SPlayerAbilities = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerAbilities)) // 0x1319
)

// client bound (outgoing) packets
const (
	// Handshake state packets
	// no client bound handshake packets defined in the protocol

	// Status state packets
	CResponse = PacketType(int32(ClientBound) + stateStatus + int32(protocolCResponse)) // 0xF100
	CPong     = PacketType(int32(ClientBound) + stateStatus + int32(protocolCPong))     // 0xF101

	// Login state packets
	CDisconnect         = PacketType(int32(ClientBound) + stateLogin + int32(protocolCDisconnect))         // 0xF200
	CEncryptionRequest  = PacketType(int32(ClientBound) + stateLogin + int32(protocolCEncryptionRequest))  // 0xF201
	CLoginSuccess       = PacketType(int32(ClientBound) + stateLogin + int32(protocolCLoginSuccess))       // 0xF202
	CSetCompression     = PacketType(int32(ClientBound) + stateLogin + int32(protocolCSetCompression))     // 0xF203
	CLoginPluginRequest = PacketType(int32(ClientBound) + stateLogin + int32(protocolCLoginPluginRequest)) // 0xF204

	// Play state packets
	CChatMessage           = PacketType(int32(ClientBound) + statePlay + int32(protocolCChatMessage))           // 0xF30F
	CJoinGame              = PacketType(int32(ClientBound) + statePlay + int32(protocolCJoinGame))              // 0xF326
	CPluginMessage         = PacketType(int32(ClientBound) + statePlay + int32(protocolCPluginMessage))         // 0xF319
	CPlayerPositionAndLook = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerPositionAndLook)) // 0xF336
	CKeepAlive             = PacketType(int32(ClientBound) + statePlay + int32(protocolCKeepAlive))             // 0xF321
	CServerDifficulty      = PacketType(int32(ClientBound) + statePlay + int32(protocolCServerDifficulty))      // 0xF30E
	CPlayerAbilities       = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerAbilities))       // 0xF332
	CHeldItemChange        = PacketType(int32(ClientBound) + statePlay + int32(protocolCHeldItemChange))        // 0xF340
	CDeclareRecipes        = PacketType(int32(ClientBound) + statePlay + int32(protocolCDeclareRecipes))        // 0xF35B
	CChunkData             = PacketType(int32(ClientBound) + statePlay + int32(protocolCChunkData))             // 0xF322
	CPlayerInfo            = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerInfo))            // 0xF334
	CEntityMetadata        = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityMetadata))        // 0xF344
)

func makeType(direction packetDirection, state State, pID ProtocolPacketID) PacketType {
	stateInt := int32(state) * 0x100
	return PacketType(int32(direction) + stateInt + int32(pID))
}

// MakeSType creates type ID for server bound packets
func MakeSType(state State, pID ProtocolPacketID) PacketType {
	return makeType(ServerBound, state, pID)
}

// MakeCType creates type ID for client bound packets
func MakeCType(state State, pID ProtocolPacketID) PacketType {
	return makeType(ClientBound, state, pID)
}
