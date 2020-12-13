//go:generate stringer -type=PacketID packets.go

package protocol

// ProtocolPacketID is the official ID of the packet as per the protocol.
type ProtocolPacketID int32

// server bound (incoming) packets, protocol definitions
const (
	// Shake state packets
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
	// Shake state packets
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
	protocolCChatMessage      ProtocolPacketID = 0x0F
	protocolCJoinGame         ProtocolPacketID = 0x26
	protocolCPluginMessage    ProtocolPacketID = 0x19
	protocolCPlayerLocation   ProtocolPacketID = 0x36
	protocolCKeepAlive        ProtocolPacketID = 0x21
	protocolCServerDifficulty ProtocolPacketID = 0x0E
	protocolCPlayerAbilities  ProtocolPacketID = 0x32
	protocolCHeldItemChange   ProtocolPacketID = 0x40
	protocolCDeclareRecipes   ProtocolPacketID = 0x5B
	protocolCChunkData        ProtocolPacketID = 0x22
	protocolCPlayerInfo       ProtocolPacketID = 0x34
	protocolCEntityMetadata   ProtocolPacketID = 0x44
)

type packetDirection int32

const ServerBound = packetDirection(0x1000)
const ClientBound = packetDirection(0xF000)

// PacketID combines direction (1 - server, 2 - client), state ID and the actual protocol ID to make a globally unique PacketID.
// E.g. PacketID 0x1101 means:
//  0x 1 1 01
//     ^ ^ ^^--- the protocol packet ID for server-bound Ping packet;
//     | |------ connection state, 1 for Status state;
//     |-------- server bound packet (1 - server, F - client);
type PacketID int32

const stateShake = 0x0000
const stateStatus = 0x0100
const stateLogin = 0x0200
const statePlay = 0x0300

// server bound (incoming) packets
const (
	// Shake state packets
	SHandshake = PacketID(int32(ServerBound) + stateShake + int32(protocolSHandshake)) // 0x1000

	// Status state packets
	SRequest = PacketID(int32(ServerBound) + stateStatus + int32(protocolSRequest)) // 0x1100
	SPing    = PacketID(int32(ServerBound) + stateStatus + int32(protocolSPing))    // 0x1101

	// Login state packets
	SLoginStart          = PacketID(int32(ServerBound) + stateLogin + int32(protocolSLoginStart))          // 0x1200
	SEncryptionResponse  = PacketID(int32(ServerBound) + stateLogin + int32(protocolSEncryptionResponse))  // 0x1201
	SLoginPluginResponse = PacketID(int32(ServerBound) + stateLogin + int32(protocolSLoginPluginResponse)) // 0x1202

	// Play state packets
	STeleportConfirm = PacketID(int32(ServerBound) + statePlay + int32(protocolSTeleportConfirm)) // 0x1300
	SQueryBlockNBT   = PacketID(int32(ServerBound) + statePlay + int32(protocolSQueryBlockNBT))   // 0x1301
	SSetDifficulty   = PacketID(int32(ServerBound) + statePlay + int32(protocolSSetDifficulty))   // 0x1302
	SChatMessage     = PacketID(int32(ServerBound) + statePlay + int32(protocolSChatMessage))     // 0x1303
	SClientStatus    = PacketID(int32(ServerBound) + statePlay + int32(protocolSClientStatus))    // 0x1304
	SClientSettings  = PacketID(int32(ServerBound) + statePlay + int32(protocolSClientSettings))  // 0x1305
	SPluginMessage   = PacketID(int32(ServerBound) + statePlay + int32(protocolSPluginMessage))   // 0x130B
	SKeepAlive       = PacketID(int32(ServerBound) + statePlay + int32(protocolSKeepAlive))       // 0x130F
	SPlayerPosition  = PacketID(int32(ServerBound) + statePlay + int32(protocolSPlayerPosition))  // 0x1311
	SPlayerLocation  = PacketID(int32(ServerBound) + statePlay + int32(protocolSPlayerLocation))  // 0x1312
	SPlayerRotation  = PacketID(int32(ServerBound) + statePlay + int32(protocolSPlayerRotation))  // 0x1313
	SPlayerAbilities = PacketID(int32(ServerBound) + statePlay + int32(protocolSPlayerAbilities)) // 0x1319
)

// client bound (outgoing) packets
const (
	// Shake state packets
	// no client bound handshake packets defined in the protocol

	// Status state packets
	CResponse = PacketID(int32(ClientBound) + stateStatus + int32(protocolCResponse)) // 0xF100
	CPong     = PacketID(int32(ClientBound) + stateStatus + int32(protocolCPong))     // 0xF101

	// Login state packets
	CDisconnect         = PacketID(int32(ClientBound) + stateLogin + int32(protocolCDisconnect))         // 0xF200
	CEncryptionRequest  = PacketID(int32(ClientBound) + stateLogin + int32(protocolCEncryptionRequest))  // 0xF201
	CLoginSuccess       = PacketID(int32(ClientBound) + stateLogin + int32(protocolCLoginSuccess))       // 0xF202
	CSetCompression     = PacketID(int32(ClientBound) + stateLogin + int32(protocolCSetCompression))     // 0xF203
	CLoginPluginRequest = PacketID(int32(ClientBound) + stateLogin + int32(protocolCLoginPluginRequest)) // 0xF204

	// Play state packets
	CChatMessage      = PacketID(int32(ClientBound) + statePlay + int32(protocolCChatMessage))      // 0xF30F
	CJoinGame         = PacketID(int32(ClientBound) + statePlay + int32(protocolCJoinGame))         // 0xF326
	CPluginMessage    = PacketID(int32(ClientBound) + statePlay + int32(protocolCPluginMessage))    // 0xF319
	CPlayerLocation   = PacketID(int32(ClientBound) + statePlay + int32(protocolCPlayerLocation))   // 0xF336
	CKeepAlive        = PacketID(int32(ClientBound) + statePlay + int32(protocolCKeepAlive))        // 0xF321
	CServerDifficulty = PacketID(int32(ClientBound) + statePlay + int32(protocolCServerDifficulty)) // 0xF30E
	CPlayerAbilities  = PacketID(int32(ClientBound) + statePlay + int32(protocolCPlayerAbilities))  // 0xF332
	CHeldItemChange   = PacketID(int32(ClientBound) + statePlay + int32(protocolCHeldItemChange))   // 0xF340
	CDeclareRecipes   = PacketID(int32(ClientBound) + statePlay + int32(protocolCDeclareRecipes))   // 0xF35B
	CChunkData        = PacketID(int32(ClientBound) + statePlay + int32(protocolCChunkData))        // 0xF322
	CPlayerInfo       = PacketID(int32(ClientBound) + statePlay + int32(protocolCPlayerInfo))       // 0xF334
	CEntityMetadata   = PacketID(int32(ClientBound) + statePlay + int32(protocolCEntityMetadata))   // 0xF344
)

var serverBound []PacketID

func init() {
	serverBound = []PacketID{
		SHandshake,
		SRequest,
		SPing,
		SLoginStart,
		SEncryptionResponse,
		SLoginPluginResponse,
		STeleportConfirm,
		SQueryBlockNBT,
		SSetDifficulty,
		SChatMessage,
		SClientStatus,
		SClientSettings,
		SPluginMessage,
		SKeepAlive,
		SPlayerPosition,
		SPlayerLocation,
		SPlayerRotation,
		SPlayerAbilities,
	}
}

func MakeID(direction packetDirection, state State, pID ProtocolPacketID) PacketID {
	return PacketID(int32(direction) + int32(state) + int32(pID))
}

func MakePacketTopic(id PacketID) string {
	return "packet." + id.String()
}
