//go:generate stringer -type=PacketType packets.go

// Package protocol defines the packets used in the Minecraft wire protocol.
// Currently supported protocol version is v754, for Minecraft 1.16.5.
package protocol

import (
	"github.com/alexykot/cncraft/pkg/buffer"
)

// Version defines the version of the minecraft wire protocol the current implementation supports.
const Version = 754

// SPacket is server bound packet type.
type SPacket interface {
	ProtocolID() ProtocolPacketID     // Return protocol ID of the packet.
	Type() PacketType                 // Return globally unique type ID of the packet.
	Pull(reader *buffer.Buffer) error // decode the server_data from provided reader into this packet
}

// CPacket is client bound packet type.
type CPacket interface {
	ProtocolID() ProtocolPacketID // Return protocol ID of the packet.
	Type() PacketType             // Return globally unique type ID of the packet.
	Push(writer *buffer.Buffer)   // encode the server_data from this packet into provided writer
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
	protocolSTeleportConfirm            ProtocolPacketID = 0x00
	protocolSQueryBlockNBT              ProtocolPacketID = 0x01
	protocolSSetDifficulty              ProtocolPacketID = 0x02
	protocolSChatMessage                ProtocolPacketID = 0x03
	protocolSClientStatus               ProtocolPacketID = 0x04
	protocolSClientSettings             ProtocolPacketID = 0x05
	protocolSTabComplete                ProtocolPacketID = 0x06
	protocolSWindowConfirmation         ProtocolPacketID = 0x07
	protocolSClickWindowButton          ProtocolPacketID = 0x08
	protocolSClickWindow                ProtocolPacketID = 0x09
	protocolSCloseWindow                ProtocolPacketID = 0x0A
	protocolSPluginMessage              ProtocolPacketID = 0x0B
	protocolSEditBook                   ProtocolPacketID = 0x0C
	protocolSQueryEntityNBT             ProtocolPacketID = 0x0D
	protocolSInteractEntity             ProtocolPacketID = 0x0E
	protocolSGenerateStructure          ProtocolPacketID = 0x0F
	protocolSKeepAlive                  ProtocolPacketID = 0x10
	protocolSLockDifficulty             ProtocolPacketID = 0x11
	protocolSPlayerPosition             ProtocolPacketID = 0x12
	protocolSPlayerPosAndRotation       ProtocolPacketID = 0x13
	protocolSPlayerRotation             ProtocolPacketID = 0x14
	protocolSPlayerMovement             ProtocolPacketID = 0x15
	protocolSVehicleMove                ProtocolPacketID = 0x16
	protocolSSteerBoat                  ProtocolPacketID = 0x17
	protocolSPickItem                   ProtocolPacketID = 0x18
	protocolSCraftRecipeRequest         ProtocolPacketID = 0x19
	protocolSPlayerAbilities            ProtocolPacketID = 0x1A
	protocolSPlayerDigging              ProtocolPacketID = 0x1B
	protocolSEntityAction               ProtocolPacketID = 0x1C
	protocolSSteerVehicle               ProtocolPacketID = 0x1D
	protocolSSetDisplayedRecipe         ProtocolPacketID = 0x1E
	protocolSSetRecipeBookState         ProtocolPacketID = 0x1F
	protocolSNameItem                   ProtocolPacketID = 0x20
	protocolSResourcePackStatus         ProtocolPacketID = 0x21
	protocolSAdvancementTab             ProtocolPacketID = 0x22
	protocolSSelectTrade                ProtocolPacketID = 0x23
	protocolSSetBeaconEffect            ProtocolPacketID = 0x24
	protocolSHeldItemChange             ProtocolPacketID = 0x25
	protocolSUpdateCommandBlock         ProtocolPacketID = 0x26
	protocolSUpdateCommandBlockMinecart ProtocolPacketID = 0x27
	protocolSCreativeInventoryAction    ProtocolPacketID = 0x28
	protocolSUpdateJigsawBlock          ProtocolPacketID = 0x29
	protocolSUpdateStructureBlock       ProtocolPacketID = 0x2A
	protocolSUpdateSign                 ProtocolPacketID = 0x2B
	protocolSAnimation                  ProtocolPacketID = 0x2C
	protocolSSpectate                   ProtocolPacketID = 0x2D
	protocolSPlayerBlockPlacement       ProtocolPacketID = 0x2E
	protocolSUseItem                    ProtocolPacketID = 0x2F
)

// client bound (outgoing) packets
const (
	// Handshake state packets
	// no client bound handshake packets defined in the protocol

	// Status state packets
	protocolCResponse ProtocolPacketID = 0x00
	protocolCPong     ProtocolPacketID = 0x01

	// Login state packets
	protocolCDisconnectLogin    ProtocolPacketID = 0x00
	protocolCEncryptionRequest  ProtocolPacketID = 0x01
	protocolCLoginSuccess       ProtocolPacketID = 0x02
	protocolCSetCompression     ProtocolPacketID = 0x03
	protocolCLoginPluginRequest ProtocolPacketID = 0x04

	// Play state packets
	protocolCSpawnEntity               ProtocolPacketID = 0x00
	protocolCSpawnExperienceOrb        ProtocolPacketID = 0x01
	protocolCSpawnLivingEntity         ProtocolPacketID = 0x02
	protocolCSpawnPainting             ProtocolPacketID = 0x03
	protocolCSpawnPlayer               ProtocolPacketID = 0x04
	protocolCEntityAnimation           ProtocolPacketID = 0x05
	protocolCStatistics                ProtocolPacketID = 0x06
	protocolCAcknowledgePlayerDigging  ProtocolPacketID = 0x07
	protocolCBlockBreakAnimation       ProtocolPacketID = 0x08
	protocolCBlockEntityData           ProtocolPacketID = 0x09
	protocolCBlockAction               ProtocolPacketID = 0x0A
	protocolCBlockChange               ProtocolPacketID = 0x0B
	protocolCBossBar                   ProtocolPacketID = 0x0C
	protocolCServerDifficulty          ProtocolPacketID = 0x0D
	protocolCChatMessage               ProtocolPacketID = 0x0E
	protocolCTabComplete               ProtocolPacketID = 0x0F
	protocolCDeclareCommands           ProtocolPacketID = 0x10
	protocolCWindowConfirmation        ProtocolPacketID = 0x11
	protocolCCloseWindow               ProtocolPacketID = 0x12
	protocolCWindowItems               ProtocolPacketID = 0x13
	protocolCWindowProperty            ProtocolPacketID = 0x14
	protocolCSetSlot                   ProtocolPacketID = 0x15
	protocolCSetCooldown               ProtocolPacketID = 0x16
	protocolCPluginMessage             ProtocolPacketID = 0x17
	protocolCNamedSoundEffect          ProtocolPacketID = 0x18
	protocolCDisconnectPlay            ProtocolPacketID = 0x19
	protocolCEntityStatus              ProtocolPacketID = 0x1A
	protocolCExplosion                 ProtocolPacketID = 0x1B
	protocolCUnloadChunk               ProtocolPacketID = 0x1C
	protocolCChangeGameState           ProtocolPacketID = 0x1D
	protocolCOpenHorseWindow           ProtocolPacketID = 0x1E
	protocolCKeepAlive                 ProtocolPacketID = 0x1F
	protocolCChunkData                 ProtocolPacketID = 0x20
	protocolCEffect                    ProtocolPacketID = 0x21
	protocolCParticle                  ProtocolPacketID = 0x22
	protocolCUpdateLight               ProtocolPacketID = 0x23
	protocolCJoinGame                  ProtocolPacketID = 0x24
	protocolCMapData                   ProtocolPacketID = 0x25
	protocolCTradeList                 ProtocolPacketID = 0x26
	protocolCEntityPosition            ProtocolPacketID = 0x27
	protocolCEntityPositionandRotation ProtocolPacketID = 0x28
	protocolCEntityRotation            ProtocolPacketID = 0x29
	protocolCEntityMovement            ProtocolPacketID = 0x2A
	protocolCVehicleMove               ProtocolPacketID = 0x2B
	protocolCOpenBook                  ProtocolPacketID = 0x2C
	protocolCOpenWindow                ProtocolPacketID = 0x2D
	protocolCOpenSignEditor            ProtocolPacketID = 0x2E
	protocolCCraftRecipeResponse       ProtocolPacketID = 0x2F
	protocolCPlayerAbilities           ProtocolPacketID = 0x30
	protocolCCombatEvent               ProtocolPacketID = 0x31
	protocolCPlayerInfo                ProtocolPacketID = 0x32
	protocolCFacePlayer                ProtocolPacketID = 0x33
	protocolCPlayerPositionAndLook     ProtocolPacketID = 0x34
	protocolCUnlockRecipes             ProtocolPacketID = 0x35
	protocolCDestroyEntities           ProtocolPacketID = 0x36
	protocolCRemoveEntityEffect        ProtocolPacketID = 0x37
	protocolCResourcePackSend          ProtocolPacketID = 0x38
	protocolCRespawn                   ProtocolPacketID = 0x39
	protocolCEntityHeadLook            ProtocolPacketID = 0x3A
	protocolCMultiBlockChange          ProtocolPacketID = 0x3B
	protocolCSelectAdvancementTab      ProtocolPacketID = 0x3C
	protocolCWorldBorder               ProtocolPacketID = 0x3D
	protocolCCamera                    ProtocolPacketID = 0x3E
	protocolCHeldItemChange            ProtocolPacketID = 0x3F
	protocolCUpdateViewPosition        ProtocolPacketID = 0x40
	protocolCUpdateViewDistance        ProtocolPacketID = 0x41
	protocolCSpawnPosition             ProtocolPacketID = 0x42
	protocolCDisplayScoreboard         ProtocolPacketID = 0x43
	protocolCEntityMetadata            ProtocolPacketID = 0x44
	protocolCAttachEntity              ProtocolPacketID = 0x45
	protocolCEntityVelocity            ProtocolPacketID = 0x46
	protocolCEntityEquipment           ProtocolPacketID = 0x47
	protocolCSetExperience             ProtocolPacketID = 0x48
	protocolCUpdateHealth              ProtocolPacketID = 0x49
	protocolCScoreboardObjective       ProtocolPacketID = 0x4A
	protocolCSetPassengers             ProtocolPacketID = 0x4B
	protocolCTeams                     ProtocolPacketID = 0x4C
	protocolCUpdateScore               ProtocolPacketID = 0x4D
	protocolCTimeUpdate                ProtocolPacketID = 0x4E
	protocolCTitle                     ProtocolPacketID = 0x4F
	protocolCEntitySoundEffect         ProtocolPacketID = 0x50
	protocolCSoundEffect               ProtocolPacketID = 0x51
	protocolCStopSound                 ProtocolPacketID = 0x52
	protocolCPlayerListHeaderAndFooter ProtocolPacketID = 0x53
	protocolCNBTQueryResponse          ProtocolPacketID = 0x54
	protocolCCollectItem               ProtocolPacketID = 0x55
	protocolCEntityTeleport            ProtocolPacketID = 0x56
	protocolCAdvancements              ProtocolPacketID = 0x57
	protocolCEntityProperties          ProtocolPacketID = 0x58
	protocolCEntityEffect              ProtocolPacketID = 0x59
	protocolCDeclareRecipes            ProtocolPacketID = 0x5A
	protocolCTags                      ProtocolPacketID = 0x5B
)

type packetDirection int32

const ServerBound = packetDirection(0x1000)

const ClientBound = packetDirection(0xF000)

// PacketType combines direction (1 - server, 2 - client), connection state and the actual protocol packet ID to
// make a globally unique PacketType. E.g. PacketType 0x1101 means:
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
	STeleportConfirm            = PacketType(int32(ServerBound) + statePlay + int32(protocolSTeleportConfirm))            // 0x1300
	SQueryBlockNBT              = PacketType(int32(ServerBound) + statePlay + int32(protocolSQueryBlockNBT))              // 0x1301
	SQueryEntityNBT             = PacketType(int32(ServerBound) + statePlay + int32(protocolSQueryEntityNBT))             // 0x1302
	SSetDifficulty              = PacketType(int32(ServerBound) + statePlay + int32(protocolSSetDifficulty))              // 0x1303
	SChatMessage                = PacketType(int32(ServerBound) + statePlay + int32(protocolSChatMessage))                // 0x1304
	SClientStatus               = PacketType(int32(ServerBound) + statePlay + int32(protocolSClientStatus))               // 0x1305
	SClientSettings             = PacketType(int32(ServerBound) + statePlay + int32(protocolSClientSettings))             // 0x1306
	STabComplete                = PacketType(int32(ServerBound) + statePlay + int32(protocolSTabComplete))                // 0x1307
	SWindowConfirmation         = PacketType(int32(ServerBound) + statePlay + int32(protocolSWindowConfirmation))         // 0x1308
	SClickWindowButton          = PacketType(int32(ServerBound) + statePlay + int32(protocolSClickWindowButton))          // 0x1309
	SClickWindow                = PacketType(int32(ServerBound) + statePlay + int32(protocolSClickWindow))                // 0x130A
	SCloseWindow                = PacketType(int32(ServerBound) + statePlay + int32(protocolSCloseWindow))                // 0x130B
	SPluginMessage              = PacketType(int32(ServerBound) + statePlay + int32(protocolSPluginMessage))              // 0x130C
	SEditBook                   = PacketType(int32(ServerBound) + statePlay + int32(protocolSEditBook))                   // 0x130D
	SInteractEntity             = PacketType(int32(ServerBound) + statePlay + int32(protocolSInteractEntity))             // 0x130E
	SGenerateStructure          = PacketType(int32(ServerBound) + statePlay + int32(protocolSGenerateStructure))          // 0x130F
	SKeepAlive                  = PacketType(int32(ServerBound) + statePlay + int32(protocolSKeepAlive))                  // 0x1310
	SLockDifficulty             = PacketType(int32(ServerBound) + statePlay + int32(protocolSLockDifficulty))             // 0x1311
	SPlayerPosition             = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerPosition))             // 0x1312
	SPlayerPosAndRotation       = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerPosAndRotation))       // 0x1313
	SPlayerRotation             = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerRotation))             // 0x1314
	SPlayerMovement             = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerMovement))             // 0x1315
	SVehicleMove                = PacketType(int32(ServerBound) + statePlay + int32(protocolSVehicleMove))                // 0x1316
	SSteerBoat                  = PacketType(int32(ServerBound) + statePlay + int32(protocolSSteerBoat))                  // 0x1317
	SPickItem                   = PacketType(int32(ServerBound) + statePlay + int32(protocolSPickItem))                   // 0x1318
	SCraftRecipeRequest         = PacketType(int32(ServerBound) + statePlay + int32(protocolSCraftRecipeRequest))         // 0x1319
	SPlayerAbilities            = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerAbilities))            // 0x131A
	SPlayerDigging              = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerDigging))              // 0x131B
	SEntityAction               = PacketType(int32(ServerBound) + statePlay + int32(protocolSEntityAction))               // 0x131C
	SSteerVehicle               = PacketType(int32(ServerBound) + statePlay + int32(protocolSSteerVehicle))               // 0x131D
	SSetDisplayedRecipe         = PacketType(int32(ServerBound) + statePlay + int32(protocolSSetDisplayedRecipe))         // 0x131E
	SSetRecipeBookState         = PacketType(int32(ServerBound) + statePlay + int32(protocolSSetRecipeBookState))         // 0x131F
	SNameItem                   = PacketType(int32(ServerBound) + statePlay + int32(protocolSNameItem))                   // 0x1320
	SResourcePackStatus         = PacketType(int32(ServerBound) + statePlay + int32(protocolSResourcePackStatus))         // 0x1321
	SAdvancementTab             = PacketType(int32(ServerBound) + statePlay + int32(protocolSAdvancementTab))             // 0x1322
	SSelectTrade                = PacketType(int32(ServerBound) + statePlay + int32(protocolSSelectTrade))                // 0x1323
	SSetBeaconEffect            = PacketType(int32(ServerBound) + statePlay + int32(protocolSSetBeaconEffect))            // 0x1324
	SHeldItemChange             = PacketType(int32(ServerBound) + statePlay + int32(protocolSHeldItemChange))             // 0x1325
	SUpdateCommandBlock         = PacketType(int32(ServerBound) + statePlay + int32(protocolSUpdateCommandBlock))         // 0x1326
	SUpdateCommandBlockMinecart = PacketType(int32(ServerBound) + statePlay + int32(protocolSUpdateCommandBlockMinecart)) // 0x1327
	SCreativeInventoryAction    = PacketType(int32(ServerBound) + statePlay + int32(protocolSCreativeInventoryAction))    // 0x1328
	SUpdateJigsawBlock          = PacketType(int32(ServerBound) + statePlay + int32(protocolSUpdateJigsawBlock))          // 0x1329
	SUpdateStructureBlock       = PacketType(int32(ServerBound) + statePlay + int32(protocolSUpdateStructureBlock))       // 0x132A
	SUpdateSign                 = PacketType(int32(ServerBound) + statePlay + int32(protocolSUpdateSign))                 // 0x132B
	SAnimation                  = PacketType(int32(ServerBound) + statePlay + int32(protocolSAnimation))                  // 0x132C
	SSpectate                   = PacketType(int32(ServerBound) + statePlay + int32(protocolSSpectate))                   // 0x132D
	SPlayerBlockPlacement       = PacketType(int32(ServerBound) + statePlay + int32(protocolSPlayerBlockPlacement))       // 0x132E
	SUseItem                    = PacketType(int32(ServerBound) + statePlay + int32(protocolSUseItem))                    // 0x132F
)

// client bound (outgoing) packets
const (
	// Handshake state packets
	// no client bound handshake packets defined in the protocol

	// Status state packets
	CResponse = PacketType(int32(ClientBound) + stateStatus + int32(protocolCResponse)) // 0xF100
	CPong     = PacketType(int32(ClientBound) + stateStatus + int32(protocolCPong))     // 0xF101

	// Login state packets
	CDisconnectLogin    = PacketType(int32(ClientBound) + stateLogin + int32(protocolCDisconnectLogin))    // 0xF200
	CEncryptionRequest  = PacketType(int32(ClientBound) + stateLogin + int32(protocolCEncryptionRequest))  // 0xF201
	CLoginSuccess       = PacketType(int32(ClientBound) + stateLogin + int32(protocolCLoginSuccess))       // 0xF202
	CSetCompression     = PacketType(int32(ClientBound) + stateLogin + int32(protocolCSetCompression))     // 0xF203
	CLoginPluginRequest = PacketType(int32(ClientBound) + stateLogin + int32(protocolCLoginPluginRequest)) // 0xF204

	// Play state packets

	CSpawnEntity               = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnEntity))               // 0xF300
	CSpawnExperienceOrb        = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnExperienceOrb))        // 0xF301
	CSpawnLivingEntity         = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnLivingEntity))         // 0xF302
	CSpawnPainting             = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnPainting))             // 0xF303
	CSpawnPlayer               = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnPlayer))               // 0xF304
	CEntityAnimation           = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityAnimation))           // 0xF305
	CStatistics                = PacketType(int32(ClientBound) + statePlay + int32(protocolCStatistics))                // 0xF306
	CAcknowledgePlayerDigging  = PacketType(int32(ClientBound) + statePlay + int32(protocolCAcknowledgePlayerDigging))  // 0xF307
	CBlockBreakAnimation       = PacketType(int32(ClientBound) + statePlay + int32(protocolCBlockBreakAnimation))       // 0xF308
	CBlockEntityData           = PacketType(int32(ClientBound) + statePlay + int32(protocolCBlockEntityData))           // 0xF309
	CBlockAction               = PacketType(int32(ClientBound) + statePlay + int32(protocolCBlockAction))               // 0xF30A
	CBlockChange               = PacketType(int32(ClientBound) + statePlay + int32(protocolCBlockChange))               // 0xF30B
	CBossBar                   = PacketType(int32(ClientBound) + statePlay + int32(protocolCBossBar))                   // 0xF30C
	CServerDifficulty          = PacketType(int32(ClientBound) + statePlay + int32(protocolCServerDifficulty))          // 0xF30D
	CChatMessage               = PacketType(int32(ClientBound) + statePlay + int32(protocolCChatMessage))               // 0xF30E
	CTabComplete               = PacketType(int32(ClientBound) + statePlay + int32(protocolCTabComplete))               // 0xF30F
	CDeclareCommands           = PacketType(int32(ClientBound) + statePlay + int32(protocolCDeclareCommands))           // 0xF310
	CWindowConfirmation        = PacketType(int32(ClientBound) + statePlay + int32(protocolCWindowConfirmation))        // 0xF311
	CCloseWindow               = PacketType(int32(ClientBound) + statePlay + int32(protocolCCloseWindow))               // 0xF312
	CWindowItems               = PacketType(int32(ClientBound) + statePlay + int32(protocolCWindowItems))               // 0xF313
	CWindowProperty            = PacketType(int32(ClientBound) + statePlay + int32(protocolCWindowProperty))            // 0xF314
	CSetSlot                   = PacketType(int32(ClientBound) + statePlay + int32(protocolCSetSlot))                   // 0xF315
	CSetCooldown               = PacketType(int32(ClientBound) + statePlay + int32(protocolCSetCooldown))               // 0xF316
	CPluginMessage             = PacketType(int32(ClientBound) + statePlay + int32(protocolCPluginMessage))             // 0xF317
	CNamedSoundEffect          = PacketType(int32(ClientBound) + statePlay + int32(protocolCNamedSoundEffect))          // 0xF318
	CDisconnectPlay            = PacketType(int32(ClientBound) + statePlay + int32(protocolCDisconnectPlay))            // 0xF319
	CEntityStatus              = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityStatus))              // 0xF31A
	CExplosion                 = PacketType(int32(ClientBound) + statePlay + int32(protocolCExplosion))                 // 0xF31B
	CUnloadChunk               = PacketType(int32(ClientBound) + statePlay + int32(protocolCUnloadChunk))               // 0xF31C
	CChangeGameState           = PacketType(int32(ClientBound) + statePlay + int32(protocolCChangeGameState))           // 0xF31D
	COpenHorseWindow           = PacketType(int32(ClientBound) + statePlay + int32(protocolCOpenHorseWindow))           // 0xF31E
	CKeepAlive                 = PacketType(int32(ClientBound) + statePlay + int32(protocolCKeepAlive))                 // 0xF31F
	CChunkData                 = PacketType(int32(ClientBound) + statePlay + int32(protocolCChunkData))                 // 0xF320
	CEffect                    = PacketType(int32(ClientBound) + statePlay + int32(protocolCEffect))                    // 0xF321
	CParticle                  = PacketType(int32(ClientBound) + statePlay + int32(protocolCParticle))                  // 0xF322
	CUpdateLight               = PacketType(int32(ClientBound) + statePlay + int32(protocolCUpdateLight))               // 0xF323
	CJoinGame                  = PacketType(int32(ClientBound) + statePlay + int32(protocolCJoinGame))                  // 0xF324
	CMapData                   = PacketType(int32(ClientBound) + statePlay + int32(protocolCMapData))                   // 0xF325
	CTradeList                 = PacketType(int32(ClientBound) + statePlay + int32(protocolCTradeList))                 // 0xF326
	CEntityPosition            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityPosition))            // 0xF327
	CEntityPositionandRotation = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityPositionandRotation)) // 0xF328
	CEntityRotation            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityRotation))            // 0xF329
	CEntityMovement            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityMovement))            // 0xF32A
	CVehicleMove               = PacketType(int32(ClientBound) + statePlay + int32(protocolCVehicleMove))               // 0xF32B
	COpenBook                  = PacketType(int32(ClientBound) + statePlay + int32(protocolCOpenBook))                  // 0xF32C
	COpenWindow                = PacketType(int32(ClientBound) + statePlay + int32(protocolCOpenWindow))                // 0xF32D
	COpenSignEditor            = PacketType(int32(ClientBound) + statePlay + int32(protocolCOpenSignEditor))            // 0xF32E
	CCraftRecipeResponse       = PacketType(int32(ClientBound) + statePlay + int32(protocolCCraftRecipeResponse))       // 0xF32F
	CPlayerAbilities           = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerAbilities))           // 0xF330
	CCombatEvent               = PacketType(int32(ClientBound) + statePlay + int32(protocolCCombatEvent))               // 0xF331
	CPlayerInfo                = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerInfo))                // 0xF332
	CFacePlayer                = PacketType(int32(ClientBound) + statePlay + int32(protocolCFacePlayer))                // 0xF333
	CPlayerPositionAndLook     = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerPositionAndLook))     // 0xF334
	CUnlockRecipes             = PacketType(int32(ClientBound) + statePlay + int32(protocolCUnlockRecipes))             // 0xF335
	CDestroyEntities           = PacketType(int32(ClientBound) + statePlay + int32(protocolCDestroyEntities))           // 0xF336
	CRemoveEntityEffect        = PacketType(int32(ClientBound) + statePlay + int32(protocolCRemoveEntityEffect))        // 0xF337
	CResourcePackSend          = PacketType(int32(ClientBound) + statePlay + int32(protocolCResourcePackSend))          // 0xF338
	CRespawn                   = PacketType(int32(ClientBound) + statePlay + int32(protocolCRespawn))                   // 0xF339
	CEntityHeadLook            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityHeadLook))            // 0xF33A
	CMultiBlockChange          = PacketType(int32(ClientBound) + statePlay + int32(protocolCMultiBlockChange))          // 0xF33B
	CSelectAdvancementTab      = PacketType(int32(ClientBound) + statePlay + int32(protocolCSelectAdvancementTab))      // 0xF33C
	CWorldBorder               = PacketType(int32(ClientBound) + statePlay + int32(protocolCWorldBorder))               // 0xF33D
	CCamera                    = PacketType(int32(ClientBound) + statePlay + int32(protocolCCamera))                    // 0xF33E
	CHeldItemChange            = PacketType(int32(ClientBound) + statePlay + int32(protocolCHeldItemChange))            // 0xF33F
	CUpdateViewPosition        = PacketType(int32(ClientBound) + statePlay + int32(protocolCUpdateViewPosition))        // 0xF340
	CUpdateViewDistance        = PacketType(int32(ClientBound) + statePlay + int32(protocolCUpdateViewDistance))        // 0xF341
	CSpawnPosition             = PacketType(int32(ClientBound) + statePlay + int32(protocolCSpawnPosition))             // 0xF342
	CDisplayScoreboard         = PacketType(int32(ClientBound) + statePlay + int32(protocolCDisplayScoreboard))         // 0xF343
	CEntityMetadata            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityMetadata))            // 0xF344
	CAttachEntity              = PacketType(int32(ClientBound) + statePlay + int32(protocolCAttachEntity))              // 0xF345
	CEntityVelocity            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityVelocity))            // 0xF346
	CEntityEquipment           = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityEquipment))           // 0xF347
	CSetExperience             = PacketType(int32(ClientBound) + statePlay + int32(protocolCSetExperience))             // 0xF348
	CUpdateHealth              = PacketType(int32(ClientBound) + statePlay + int32(protocolCUpdateHealth))              // 0xF349
	CScoreboardObjective       = PacketType(int32(ClientBound) + statePlay + int32(protocolCScoreboardObjective))       // 0xF34A
	CSetPassengers             = PacketType(int32(ClientBound) + statePlay + int32(protocolCSetPassengers))             // 0xF34B
	CTeams                     = PacketType(int32(ClientBound) + statePlay + int32(protocolCTeams))                     // 0xF34C
	CUpdateScore               = PacketType(int32(ClientBound) + statePlay + int32(protocolCUpdateScore))               // 0xF34D
	CTimeUpdate                = PacketType(int32(ClientBound) + statePlay + int32(protocolCTimeUpdate))                // 0xF34E
	CTitle                     = PacketType(int32(ClientBound) + statePlay + int32(protocolCTitle))                     // 0xF34F
	CEntitySoundEffect         = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntitySoundEffect))         // 0xF350
	CSoundEffect               = PacketType(int32(ClientBound) + statePlay + int32(protocolCSoundEffect))               // 0xF351
	CStopSound                 = PacketType(int32(ClientBound) + statePlay + int32(protocolCStopSound))                 // 0xF352
	CPlayerListHeaderAndFooter = PacketType(int32(ClientBound) + statePlay + int32(protocolCPlayerListHeaderAndFooter)) // 0xF353
	CNBTQueryResponse          = PacketType(int32(ClientBound) + statePlay + int32(protocolCNBTQueryResponse))          // 0xF354
	CCollectItem               = PacketType(int32(ClientBound) + statePlay + int32(protocolCCollectItem))               // 0xF355
	CEntityTeleport            = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityTeleport))            // 0xF356
	CAdvancements              = PacketType(int32(ClientBound) + statePlay + int32(protocolCAdvancements))              // 0xF357
	CEntityProperties          = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityProperties))          // 0xF358
	CEntityEffect              = PacketType(int32(ClientBound) + statePlay + int32(protocolCEntityEffect))              // 0xF359
	CDeclareRecipes            = PacketType(int32(ClientBound) + statePlay + int32(protocolCDeclareRecipes))            // 0xF35A
	CTags                      = PacketType(int32(ClientBound) + statePlay + int32(protocolCTags))                      // 0xF35B
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
