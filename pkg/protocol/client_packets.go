package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/game/items"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/chat"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/level"
	"github.com/alexykot/cncraft/pkg/game/player"
	"github.com/alexykot/cncraft/pkg/nbt"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
	"github.com/alexykot/cncraft/pkg/protocol/status"
	"github.com/alexykot/cncraft/pkg/protocol/tags"
)

// HANDSHAKE STATE CLIENT BOUND PACKETS DO NOT EXIST

// STATUS STATE PACKETS
type CPacketResponse struct {
	Status status.Response
}

func (p *CPacketResponse) ProtocolID() ProtocolPacketID { return protocolCResponse }
func (p *CPacketResponse) Type() PacketType             { return CResponse }
func (p *CPacketResponse) Push(writer buffer.B) {
	if text, err := json.Marshal(p.Status); err != nil {
		panic(err)
	} else {
		writer.PushString(string(text))
	}
}

type CPacketPong struct {
	Payload int64
}

func (p *CPacketPong) ProtocolID() ProtocolPacketID { return protocolCPong }
func (p *CPacketPong) Type() PacketType             { return CPong }
func (p *CPacketPong) Push(writer buffer.B) {
	writer.PushInt64(p.Payload)
}

// LOGIN STATE PACKETS
type CPacketDisconnectLogin struct {
	Reason *chat.Message
}

func (p *CPacketDisconnectLogin) ProtocolID() ProtocolPacketID { return protocolCDisconnectLogin }
func (p *CPacketDisconnectLogin) Type() PacketType             { return CDisconnectLogin }
func (p *CPacketDisconnectLogin) Push(writer buffer.B) {
	message := p.Reason

	writer.PushString(message.AsJson())
}

func (p *CPacketDisconnectLogin) Pull(reader buffer.B) {
	p.Reason = chat.New(reader.PullString())
}

type CPacketEncryptionRequest struct {
	ServerID    string // Appears to be unused by the Notchian client.
	PublicKey   []byte
	VerifyToken []byte
}

func (p *CPacketEncryptionRequest) ProtocolID() ProtocolPacketID { return protocolCEncryptionRequest }
func (p *CPacketEncryptionRequest) Type() PacketType             { return CEncryptionRequest }
func (p *CPacketEncryptionRequest) Push(writer buffer.B) {
	writer.PushString(p.ServerID)
	writer.PushBytes(p.PublicKey, true)
	writer.PushBytes(p.VerifyToken, true)
}

type CPacketLoginSuccess struct {
	PlayerUUID uuid.UUID
	PlayerName string
}

func (p *CPacketLoginSuccess) ProtocolID() ProtocolPacketID { return protocolCLoginSuccess }
func (p *CPacketLoginSuccess) Type() PacketType             { return CLoginSuccess }
func (p *CPacketLoginSuccess) Push(writer buffer.B) {
	writer.PushUUID(p.PlayerUUID)
	writer.PushString(p.PlayerName)
}

type CPacketSetCompression struct {
	Threshold int32
}

func (p *CPacketSetCompression) ProtocolID() ProtocolPacketID { return protocolCSetCompression }
func (p *CPacketSetCompression) Type() PacketType             { return CSetCompression }
func (p *CPacketSetCompression) Push(writer buffer.B) {
	writer.PushVarInt(p.Threshold)
}

type CPacketLoginPluginRequest struct {
	MessageID int32
	Channel   string
	OptData   []byte
}

func (p *CPacketLoginPluginRequest) ProtocolID() ProtocolPacketID { return protocolCLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Type() PacketType             { return CLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Push(writer buffer.B) {
	writer.PushVarInt(p.MessageID)
	writer.PushString(p.Channel)
	writer.PushBytes(p.OptData, false)
}

// PLAY STATE PACKETS
type CPacketSpawnEntity struct{}

func (p *CPacketSpawnEntity) ProtocolID() ProtocolPacketID { return protocolCSpawnEntity }
func (p *CPacketSpawnEntity) Type() PacketType             { return CDisconnectPlay }
func (p *CPacketSpawnEntity) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSpawnExperienceOrb struct{}

func (p *CPacketSpawnExperienceOrb) ProtocolID() ProtocolPacketID { return protocolCSpawnExperienceOrb }
func (p *CPacketSpawnExperienceOrb) Type() PacketType             { return CSpawnExperienceOrb }
func (p *CPacketSpawnExperienceOrb) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSpawnLivingEntity struct{}

func (p *CPacketSpawnLivingEntity) ProtocolID() ProtocolPacketID { return protocolCSpawnLivingEntity }
func (p *CPacketSpawnLivingEntity) Type() PacketType             { return CSpawnLivingEntity }
func (p *CPacketSpawnLivingEntity) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSpawnPainting struct{}

func (p *CPacketSpawnPainting) ProtocolID() ProtocolPacketID { return protocolCSpawnPainting }
func (p *CPacketSpawnPainting) Type() PacketType             { return CSpawnPainting }
func (p *CPacketSpawnPainting) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSpawnPlayer struct{}

func (p *CPacketSpawnPlayer) ProtocolID() ProtocolPacketID { return protocolCSpawnPlayer }
func (p *CPacketSpawnPlayer) Type() PacketType             { return CSpawnPlayer }
func (p *CPacketSpawnPlayer) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityAnimation struct{}

func (p *CPacketEntityAnimation) ProtocolID() ProtocolPacketID { return protocolCEntityAnimation }
func (p *CPacketEntityAnimation) Type() PacketType             { return CEntityAnimation }
func (p *CPacketEntityAnimation) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketStatistics struct{}

func (p *CPacketStatistics) ProtocolID() ProtocolPacketID { return protocolCStatistics }
func (p *CPacketStatistics) Type() PacketType             { return CStatistics }
func (p *CPacketStatistics) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketAcknowledgePlayerDigging struct{}

func (p *CPacketAcknowledgePlayerDigging) ProtocolID() ProtocolPacketID {
	return protocolCAcknowledgePlayerDigging
}
func (p *CPacketAcknowledgePlayerDigging) Type() PacketType     { return CAcknowledgePlayerDigging }
func (p *CPacketAcknowledgePlayerDigging) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketBlockBreakAnimation struct{}

func (p *CPacketBlockBreakAnimation) ProtocolID() ProtocolPacketID {
	return protocolCBlockBreakAnimation
}
func (p *CPacketBlockBreakAnimation) Type() PacketType     { return CBlockBreakAnimation }
func (p *CPacketBlockBreakAnimation) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketBlockEntityData struct{}

func (p *CPacketBlockEntityData) ProtocolID() ProtocolPacketID { return protocolCBlockEntityData }
func (p *CPacketBlockEntityData) Type() PacketType             { return CBlockEntityData }
func (p *CPacketBlockEntityData) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketBlockAction struct{}

func (p *CPacketBlockAction) ProtocolID() ProtocolPacketID { return protocolCBlockAction }
func (p *CPacketBlockAction) Type() PacketType             { return CBlockAction }
func (p *CPacketBlockAction) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketBlockChange struct{}

func (p *CPacketBlockChange) ProtocolID() ProtocolPacketID { return protocolCBlockChange }
func (p *CPacketBlockChange) Type() PacketType             { return CBlockChange }
func (p *CPacketBlockChange) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketBossBar struct{}

func (p *CPacketBossBar) ProtocolID() ProtocolPacketID { return protocolCBossBar }
func (p *CPacketBossBar) Type() PacketType             { return CBossBar }
func (p *CPacketBossBar) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketServerDifficulty struct {
	Difficulty game.Difficulty
	Locked     bool // should probably always be true
}

func (p *CPacketServerDifficulty) ProtocolID() ProtocolPacketID { return protocolCServerDifficulty }
func (p *CPacketServerDifficulty) Type() PacketType             { return CServerDifficulty }
func (p *CPacketServerDifficulty) Push(writer buffer.B) {
	writer.PushByte(byte(p.Difficulty))
	writer.PushBool(p.Locked)
}

type CPacketChatMessage struct {
	Message         chat.Message
	MessagePosition chat.MessagePosition
	Sender          uuid.UUID
}

func (p *CPacketChatMessage) ProtocolID() ProtocolPacketID { return protocolCChatMessage }
func (p *CPacketChatMessage) Type() PacketType             { return CChatMessage }
func (p *CPacketChatMessage) Push(writer buffer.B) {
	message := p.Message

	if p.MessagePosition == chat.HotBarText {
		message = *chat.New(message.AsText())
	}

	writer.PushString(message.AsJson())
	writer.PushByte(byte(p.MessagePosition))
	writer.PushUUID(p.Sender)
}

type CPacketTabComplete struct{}

func (p *CPacketTabComplete) ProtocolID() ProtocolPacketID { return protocolCTabComplete }
func (p *CPacketTabComplete) Type() PacketType             { return CTabComplete }
func (p *CPacketTabComplete) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketDeclareCommands struct{}

func (p *CPacketDeclareCommands) ProtocolID() ProtocolPacketID { return protocolCDeclareCommands }
func (p *CPacketDeclareCommands) Type() PacketType             { return CDeclareCommands }
func (p *CPacketDeclareCommands) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketWindowConfirmation struct {
	WindowID items.WindowID
	ActionID int16
	Accepted bool
}

func (p *CPacketWindowConfirmation) ProtocolID() ProtocolPacketID { return protocolCWindowConfirmation }
func (p *CPacketWindowConfirmation) Type() PacketType             { return CWindowConfirmation }
func (p *CPacketWindowConfirmation) Push(writer buffer.B) {
	writer.PushByte(byte(p.WindowID))
	writer.PushInt16(p.ActionID)
	writer.PushBool(p.Accepted)
}

func (p *CPacketWindowConfirmation) Pull(reader buffer.B) {
	p.WindowID = items.WindowID(reader.PullByte())
	p.ActionID = reader.PullInt16()
	p.Accepted = reader.PullBool()
}

type CPacketCloseWindow struct{}

func (p *CPacketCloseWindow) ProtocolID() ProtocolPacketID { return protocolCCloseWindow }
func (p *CPacketCloseWindow) Type() PacketType             { return CCloseWindow }
func (p *CPacketCloseWindow) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketWindowItems struct {
	WindowID  int8
	SlotCount int16
	Slots     []items.Slot
}

func (p *CPacketWindowItems) ProtocolID() ProtocolPacketID { return protocolCWindowItems }
func (p *CPacketWindowItems) Type() PacketType             { return CWindowItems }
func (p *CPacketWindowItems) Push(writer buffer.B) {
	writer.PushByte(byte(p.WindowID))
	writer.PushInt16(p.SlotCount)

	for _, slotItem := range p.Slots {
		writer.PushBool(slotItem.IsPresent)
		if slotItem.IsPresent {
			writer.PushVarInt(int32(slotItem.ItemID))
			writer.PushByte(byte(slotItem.ItemCount))

			writer.PushByte(0x00) // TODO item NBT data not implemented
		}
	}
}

type CPacketWindowProperty struct{}

func (p *CPacketWindowProperty) ProtocolID() ProtocolPacketID { return protocolCWindowProperty }
func (p *CPacketWindowProperty) Type() PacketType             { return CWindowProperty }
func (p *CPacketWindowProperty) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSetSlot struct{}

func (p *CPacketSetSlot) ProtocolID() ProtocolPacketID { return protocolCSetSlot }
func (p *CPacketSetSlot) Type() PacketType             { return CSetSlot }
func (p *CPacketSetSlot) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSetCooldown struct{}

func (p *CPacketSetCooldown) ProtocolID() ProtocolPacketID { return protocolCSetCooldown }
func (p *CPacketSetCooldown) Type() PacketType             { return CSetCooldown }
func (p *CPacketSetCooldown) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketPluginMessage struct {
	Message plugin.Message
}

func (p *CPacketPluginMessage) ProtocolID() ProtocolPacketID { return protocolCPluginMessage }
func (p *CPacketPluginMessage) Type() PacketType             { return CPluginMessage }
func (p *CPacketPluginMessage) Push(writer buffer.B) {
	writer.PushString(string(p.Message.Chan()))
	p.Message.Push(writer)
}

type CPacketNamedSoundEffect struct{}

func (p *CPacketNamedSoundEffect) ProtocolID() ProtocolPacketID { return protocolCNamedSoundEffect }
func (p *CPacketNamedSoundEffect) Type() PacketType             { return CNamedSoundEffect }
func (p *CPacketNamedSoundEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketDisconnectPlay struct {
	Reason *chat.Message
}

func (p *CPacketDisconnectPlay) ProtocolID() ProtocolPacketID { return protocolCDisconnectPlay }
func (p *CPacketDisconnectPlay) Type() PacketType             { return CDisconnectPlay }
func (p *CPacketDisconnectPlay) Push(writer buffer.B) {
	message := p.Reason

	writer.PushString(message.AsJson())
}

func (p *CPacketDisconnectPlay) Pull(reader buffer.B) {
	p.Reason = chat.New(reader.PullString())
}

type CPacketEntityStatus struct{}

func (p *CPacketEntityStatus) ProtocolID() ProtocolPacketID { return protocolCEntityStatus }
func (p *CPacketEntityStatus) Type() PacketType             { return CEntityStatus }
func (p *CPacketEntityStatus) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketExplosion struct{}

func (p *CPacketExplosion) ProtocolID() ProtocolPacketID { return protocolCExplosion }
func (p *CPacketExplosion) Type() PacketType             { return CExplosion }
func (p *CPacketExplosion) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketUnloadChunk struct{}

func (p *CPacketUnloadChunk) ProtocolID() ProtocolPacketID { return protocolCUnloadChunk }
func (p *CPacketUnloadChunk) Type() PacketType             { return CUnloadChunk }
func (p *CPacketUnloadChunk) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketChangeGameState struct{}

func (p *CPacketChangeGameState) ProtocolID() ProtocolPacketID { return protocolCChangeGameState }
func (p *CPacketChangeGameState) Type() PacketType             { return CChangeGameState }
func (p *CPacketChangeGameState) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketOpenHorseWindow struct{}

func (p *CPacketOpenHorseWindow) ProtocolID() ProtocolPacketID { return protocolCOpenHorseWindow }
func (p *CPacketOpenHorseWindow) Type() PacketType             { return COpenHorseWindow }
func (p *CPacketOpenHorseWindow) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *CPacketKeepAlive) ProtocolID() ProtocolPacketID { return protocolCKeepAlive }
func (p *CPacketKeepAlive) Type() PacketType             { return CKeepAlive }
func (p *CPacketKeepAlive) Push(writer buffer.B) {
	writer.PushInt64(p.KeepAliveID)
}

type CPacketChunkData struct {
	Chunk level.Chunk
}

func (p *CPacketChunkData) ProtocolID() ProtocolPacketID { return protocolCChunkData }
func (p *CPacketChunkData) Type() PacketType             { return CChunkData }
func (p *CPacketChunkData) Push(writer buffer.B) {
	writer.PushInt32(int32(p.Chunk.X() / 16)) // convert block coord into chunk coord
	writer.PushInt32(int32(p.Chunk.Z() / 16)) // convert block coord into chunk coord

	// TODO make this conditional and only use this packet to send large scale updates.
	writer.PushBool(true) // IsFullChunk, see https://wiki.vg/Chunk_Format#Full_chunk

	var bitMask uint8
	for i, chunkSection := range p.Chunk.Sections() {
		if chunkSection != nil {
			bitMask = bitMask | (1 << i)
		}
	}
	// The documentation claims it's a VarInt, but the test chunk data from Notchian is inconclusive. Try a byte and see.
	writer.PushByte(bitMask)

	heightMapBuff := buffer.New()
	if err := nbt.Marshal(heightMapBuff, p.Chunk.HeightMap()); err != nil {
		panic(fmt.Errorf("failed to marshal NBT: %w", err))
	}
	writer.PushBytes(heightMapBuff.Bytes(), false)

	biomesBuff := buffer.New()
	biomesBuff.PushVarInt(1024)
	biomes := [1024]int32{}
	for range biomes {
		biomesBuff.PushVarInt(1) // TODO biomes are hardcoded to Plains for now.
	}
	writer.PushBytes(biomesBuff.Bytes(), false)

	sectionsBuff := buffer.New()
	for _, chunkSection := range p.Chunk.Sections() {
		if chunkSection != nil {
			sectionBuff := buffer.New()
			chunkSection.Push(sectionBuff)
			sectionsBuff.PushBytes(sectionBuff.Bytes(), false)
		}
	}
	writer.PushBytes(sectionsBuff.Bytes(), true)

	blockEntitiesBuff := buffer.New()
	blockEntitiesBuff.PushVarInt(0) // TODO block entities not implemented yet
	// if err := nbt.Marshal(blockEntitiesBuff, struct{}{}); err != nil {
	// 	panic(fmt.Errorf("failed to marshal NBT: %w", err))
	// }
	// println(fmt.Sprintf("blockEntities bytes %X", blockEntitiesBuff.Bytes()))
	writer.PushBytes(blockEntitiesBuff.Bytes(), false)
}

type CPacketEffect struct{}

func (p *CPacketEffect) ProtocolID() ProtocolPacketID { return protocolCEffect }
func (p *CPacketEffect) Type() PacketType             { return CEffect }
func (p *CPacketEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketParticle struct{}

func (p *CPacketParticle) ProtocolID() ProtocolPacketID { return protocolCParticle }
func (p *CPacketParticle) Type() PacketType             { return CParticle }
func (p *CPacketParticle) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketUpdateLight struct{}

func (p *CPacketUpdateLight) ProtocolID() ProtocolPacketID { return protocolCUpdateLight }
func (p *CPacketUpdateLight) Type() PacketType             { return CUpdateLight }
func (p *CPacketUpdateLight) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketJoinGame struct {
	EntityID int32

	WorldNames []string

	GameMode   game.Gamemode
	IsHardcore game.Coreness

	DimensionCodec tags.DimensionCodec
	Dimension      tags.Dimension

	WorldName string
	IsDebug   bool
	IsFlat    bool

	HashedSeed int64

	ViewDistance        int32
	IsDebugReduced      bool
	EnableRespawnScreen bool
}

func (p *CPacketJoinGame) ProtocolID() ProtocolPacketID { return protocolCJoinGame }
func (p *CPacketJoinGame) Type() PacketType             { return CJoinGame }
func (p *CPacketJoinGame) Push(writer buffer.B) {
	writer.PushInt32(p.EntityID)
	writer.PushBool(bool(p.IsHardcore))
	writer.PushByte(byte(p.GameMode))
	writer.PushByte(0xFF) // "Previous Gamemode" field, hardcoded to "-1" which means "none" and ignored.

	writer.PushVarInt(int32(len(p.WorldNames)))
	for _, worldName := range p.WorldNames {
		writer.PushString(worldName)
	}

	// DEBT push packet interface should handle and return marshalling errors
	if err := nbt.Marshal(writer, p.DimensionCodec); err != nil {
		panic(fmt.Errorf("failed to marshal NBT: %w", err))
	}
	if err := nbt.Marshal(writer, p.Dimension); err != nil {
		panic(fmt.Errorf("failed to marshal NBT: %w", err))
	}

	writer.PushString(p.WorldName)
	writer.PushInt64(p.HashedSeed)
	writer.PushVarInt(0) // "Max Players" is ignored by the Notchian client
	writer.PushVarInt(p.ViewDistance)
	writer.PushBool(p.IsDebugReduced)
	writer.PushBool(p.EnableRespawnScreen)
	writer.PushBool(p.IsDebug)
	writer.PushBool(p.IsFlat)
}

type CPacketMapData struct{}

func (p *CPacketMapData) ProtocolID() ProtocolPacketID { return protocolCMapData }
func (p *CPacketMapData) Type() PacketType             { return CMapData }
func (p *CPacketMapData) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketTradeList struct{}

func (p *CPacketTradeList) ProtocolID() ProtocolPacketID { return protocolCTradeList }
func (p *CPacketTradeList) Type() PacketType             { return CTradeList }
func (p *CPacketTradeList) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityPosition struct{}

func (p *CPacketEntityPosition) ProtocolID() ProtocolPacketID { return protocolCEntityPosition }
func (p *CPacketEntityPosition) Type() PacketType             { return CEntityPosition }
func (p *CPacketEntityPosition) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityPositionandRotation struct{}

func (p *CPacketEntityPositionandRotation) ProtocolID() ProtocolPacketID {
	return protocolCEntityPositionandRotation
}
func (p *CPacketEntityPositionandRotation) Type() PacketType     { return CEntityPositionandRotation }
func (p *CPacketEntityPositionandRotation) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketEntityRotation struct{}

func (p *CPacketEntityRotation) ProtocolID() ProtocolPacketID { return protocolCEntityRotation }
func (p *CPacketEntityRotation) Type() PacketType             { return CEntityRotation }
func (p *CPacketEntityRotation) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityMovement struct{}

func (p *CPacketEntityMovement) ProtocolID() ProtocolPacketID { return protocolCEntityMovement }
func (p *CPacketEntityMovement) Type() PacketType             { return CEntityMovement }
func (p *CPacketEntityMovement) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketVehicleMove struct{}

func (p *CPacketVehicleMove) ProtocolID() ProtocolPacketID { return protocolCVehicleMove }
func (p *CPacketVehicleMove) Type() PacketType             { return CVehicleMove }
func (p *CPacketVehicleMove) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketOpenBook struct{}

func (p *CPacketOpenBook) ProtocolID() ProtocolPacketID { return protocolCOpenBook }
func (p *CPacketOpenBook) Type() PacketType             { return COpenBook }
func (p *CPacketOpenBook) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketOpenWindow struct{}

func (p *CPacketOpenWindow) ProtocolID() ProtocolPacketID { return protocolCOpenWindow }
func (p *CPacketOpenWindow) Type() PacketType             { return COpenWindow }
func (p *CPacketOpenWindow) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketOpenSignEditor struct{}

func (p *CPacketOpenSignEditor) ProtocolID() ProtocolPacketID { return protocolCOpenSignEditor }
func (p *CPacketOpenSignEditor) Type() PacketType             { return COpenSignEditor }
func (p *CPacketOpenSignEditor) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketCraftRecipeResponse struct{}

func (p *CPacketCraftRecipeResponse) ProtocolID() ProtocolPacketID {
	return protocolCCraftRecipeResponse
}
func (p *CPacketCraftRecipeResponse) Type() PacketType     { return CCraftRecipeResponse }
func (p *CPacketCraftRecipeResponse) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketPlayerAbilities struct {
	Abilities   player.Abilities
	FlyingSpeed float32
	FieldOfView float32
}

func (p *CPacketPlayerAbilities) ProtocolID() ProtocolPacketID { return protocolCPlayerAbilities }
func (p *CPacketPlayerAbilities) Type() PacketType             { return CPlayerAbilities }
func (p *CPacketPlayerAbilities) Push(writer buffer.B) {
	p.Abilities.Push(writer)

	writer.PushFloat32(p.FlyingSpeed)
	writer.PushFloat32(p.FieldOfView)
}

type CPacketCombatEvent struct{}

func (p *CPacketCombatEvent) ProtocolID() ProtocolPacketID { return protocolCCombatEvent }
func (p *CPacketCombatEvent) Type() PacketType             { return CCombatEvent }
func (p *CPacketCombatEvent) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketPlayerInfo struct {
	Action player.PlayerInfoAction
	Values []player.PlayerInfo
}

func (p *CPacketPlayerInfo) ProtocolID() ProtocolPacketID { return protocolCPlayerInfo }
func (p *CPacketPlayerInfo) Type() PacketType             { return CPlayerInfo }
func (p *CPacketPlayerInfo) Push(writer buffer.B) {
	panic("player.PlayerInfo structure may have changed in 1.16.4, need to recheck")

	writer.PushVarInt(int32(p.Action))
	writer.PushVarInt(int32(len(p.Values)))

	for _, value := range p.Values {
		value.Push(writer)
	}
}

type CPacketFacePlayer struct{}

func (p *CPacketFacePlayer) ProtocolID() ProtocolPacketID { return protocolCFacePlayer }
func (p *CPacketFacePlayer) Type() PacketType             { return CFacePlayer }
func (p *CPacketFacePlayer) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketPlayerPositionAndLook struct {
	Location data.Location
	Relative data.Relativity

	TeleportID int32 // no idea what ID is this, the packet type 3/0x36 in the protocol 754 does not have this field
}

func (p *CPacketPlayerPositionAndLook) ProtocolID() ProtocolPacketID {
	return protocolCPlayerPositionAndLook
}
func (p *CPacketPlayerPositionAndLook) Type() PacketType { return CPlayerPositionAndLook }
func (p *CPacketPlayerPositionAndLook) Push(writer buffer.B) {
	writer.PushFloat64(p.Location.X)
	writer.PushFloat64(p.Location.Y)
	writer.PushFloat64(p.Location.Z)

	writer.PushFloat32(p.Location.Yaw)
	writer.PushFloat32(p.Location.Pitch)

	p.Relative.Push(writer)

	writer.PushVarInt(p.TeleportID)
}

type CPacketUnlockRecipes struct{}

func (p *CPacketUnlockRecipes) ProtocolID() ProtocolPacketID { return protocolCUnlockRecipes }
func (p *CPacketUnlockRecipes) Type() PacketType             { return CUnlockRecipes }
func (p *CPacketUnlockRecipes) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketDestroyEntities struct{}

func (p *CPacketDestroyEntities) ProtocolID() ProtocolPacketID { return protocolCDestroyEntities }
func (p *CPacketDestroyEntities) Type() PacketType             { return CDestroyEntities }
func (p *CPacketDestroyEntities) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketRemoveEntityEffect struct{}

func (p *CPacketRemoveEntityEffect) ProtocolID() ProtocolPacketID { return protocolCRemoveEntityEffect }
func (p *CPacketRemoveEntityEffect) Type() PacketType             { return CRemoveEntityEffect }
func (p *CPacketRemoveEntityEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketResourcePackSend struct{}

func (p *CPacketResourcePackSend) ProtocolID() ProtocolPacketID { return protocolCResourcePackSend }
func (p *CPacketResourcePackSend) Type() PacketType             { return CResourcePackSend }
func (p *CPacketResourcePackSend) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketRespawn struct{}

func (p *CPacketRespawn) ProtocolID() ProtocolPacketID { return protocolCRespawn }
func (p *CPacketRespawn) Type() PacketType             { return CRespawn }
func (p *CPacketRespawn) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityHeadLook struct{}

func (p *CPacketEntityHeadLook) ProtocolID() ProtocolPacketID { return protocolCEntityHeadLook }
func (p *CPacketEntityHeadLook) Type() PacketType             { return CEntityHeadLook }
func (p *CPacketEntityHeadLook) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketMultiBlockChange struct{}

func (p *CPacketMultiBlockChange) ProtocolID() ProtocolPacketID { return protocolCMultiBlockChange }
func (p *CPacketMultiBlockChange) Type() PacketType             { return CMultiBlockChange }
func (p *CPacketMultiBlockChange) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSelectAdvancementTab struct{}

func (p *CPacketSelectAdvancementTab) ProtocolID() ProtocolPacketID {
	return protocolCSelectAdvancementTab
}
func (p *CPacketSelectAdvancementTab) Type() PacketType     { return CSelectAdvancementTab }
func (p *CPacketSelectAdvancementTab) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketWorldBorder struct{}

func (p *CPacketWorldBorder) ProtocolID() ProtocolPacketID { return protocolCWorldBorder }
func (p *CPacketWorldBorder) Type() PacketType             { return CWorldBorder }
func (p *CPacketWorldBorder) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketCamera struct{}

func (p *CPacketCamera) ProtocolID() ProtocolPacketID { return protocolCCamera }
func (p *CPacketCamera) Type() PacketType             { return CCamera }
func (p *CPacketCamera) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketHeldItemChange struct {
	Slot uint8
}

func (p *CPacketHeldItemChange) ProtocolID() ProtocolPacketID { return protocolCHeldItemChange }
func (p *CPacketHeldItemChange) Type() PacketType             { return CHeldItemChange }
func (p *CPacketHeldItemChange) Push(writer buffer.B) {
	writer.PushByte(p.Slot)
}

type CPacketUpdateViewPosition struct{}

func (p *CPacketUpdateViewPosition) ProtocolID() ProtocolPacketID { return protocolCUpdateViewPosition }
func (p *CPacketUpdateViewPosition) Type() PacketType             { return CUpdateViewPosition }
func (p *CPacketUpdateViewPosition) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketUpdateViewDistance struct{}

func (p *CPacketUpdateViewDistance) ProtocolID() ProtocolPacketID { return protocolCUpdateViewDistance }
func (p *CPacketUpdateViewDistance) Type() PacketType             { return CUpdateViewDistance }
func (p *CPacketUpdateViewDistance) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSpawnPosition struct{}

func (p *CPacketSpawnPosition) ProtocolID() ProtocolPacketID { return protocolCSpawnPosition }
func (p *CPacketSpawnPosition) Type() PacketType             { return CSpawnPosition }
func (p *CPacketSpawnPosition) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketDisplayScoreboard struct{}

func (p *CPacketDisplayScoreboard) ProtocolID() ProtocolPacketID { return protocolCDisplayScoreboard }
func (p *CPacketDisplayScoreboard) Type() PacketType             { return CDisplayScoreboard }
func (p *CPacketDisplayScoreboard) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityMetadata struct {
	Entity entities.Entity
}

func (p *CPacketEntityMetadata) ProtocolID() ProtocolPacketID { return protocolCEntityMetadata }
func (p *CPacketEntityMetadata) Type() PacketType             { return CEntityMetadata }
func (p *CPacketEntityMetadata) Push(writer buffer.B) {
	writer.PushVarInt(p.Entity.ID())

	// only supporting player metadata for now
	_, ok := p.Entity.(entities.PlayerCharacter)
	if ok {

		writer.PushByte(16)  // index | displayed skin parts
		writer.PushVarInt(0) // type | byte

		skin := player.SkinParts{
			Cape: true,
			Head: true,
			Body: true,
			ArmL: true,
			ArmR: true,
			LegL: true,
			LegR: true,
		}

		skin.Push(writer)
	}

	writer.PushByte(0xFF)
}

type CPacketAttachEntity struct{}

func (p *CPacketAttachEntity) ProtocolID() ProtocolPacketID { return protocolCAttachEntity }
func (p *CPacketAttachEntity) Type() PacketType             { return CAttachEntity }
func (p *CPacketAttachEntity) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityVelocity struct{}

func (p *CPacketEntityVelocity) ProtocolID() ProtocolPacketID { return protocolCEntityVelocity }
func (p *CPacketEntityVelocity) Type() PacketType             { return CEntityVelocity }
func (p *CPacketEntityVelocity) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityEquipment struct{}

func (p *CPacketEntityEquipment) ProtocolID() ProtocolPacketID { return protocolCEntityEquipment }
func (p *CPacketEntityEquipment) Type() PacketType             { return CEntityEquipment }
func (p *CPacketEntityEquipment) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSetExperience struct{}

func (p *CPacketSetExperience) ProtocolID() ProtocolPacketID { return protocolCSetExperience }
func (p *CPacketSetExperience) Type() PacketType             { return CSetExperience }
func (p *CPacketSetExperience) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketUpdateHealth struct{}

func (p *CPacketUpdateHealth) ProtocolID() ProtocolPacketID { return protocolCUpdateHealth }
func (p *CPacketUpdateHealth) Type() PacketType             { return CUpdateHealth }
func (p *CPacketUpdateHealth) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketScoreboardObjective struct{}

func (p *CPacketScoreboardObjective) ProtocolID() ProtocolPacketID {
	return protocolCScoreboardObjective
}
func (p *CPacketScoreboardObjective) Type() PacketType     { return CScoreboardObjective }
func (p *CPacketScoreboardObjective) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketSetPassengers struct{}

func (p *CPacketSetPassengers) ProtocolID() ProtocolPacketID { return protocolCSetPassengers }
func (p *CPacketSetPassengers) Type() PacketType             { return CSetPassengers }
func (p *CPacketSetPassengers) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketTeams struct{}

func (p *CPacketTeams) ProtocolID() ProtocolPacketID { return protocolCTeams }
func (p *CPacketTeams) Type() PacketType             { return CTeams }
func (p *CPacketTeams) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketUpdateScore struct{}

func (p *CPacketUpdateScore) ProtocolID() ProtocolPacketID { return protocolCUpdateScore }
func (p *CPacketUpdateScore) Type() PacketType             { return CUpdateScore }
func (p *CPacketUpdateScore) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketTimeUpdate struct{}

func (p *CPacketTimeUpdate) ProtocolID() ProtocolPacketID { return protocolCTimeUpdate }
func (p *CPacketTimeUpdate) Type() PacketType             { return CTimeUpdate }
func (p *CPacketTimeUpdate) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketTitle struct{}

func (p *CPacketTitle) ProtocolID() ProtocolPacketID { return protocolCTitle }
func (p *CPacketTitle) Type() PacketType             { return CTitle }
func (p *CPacketTitle) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntitySoundEffect struct{}

func (p *CPacketEntitySoundEffect) ProtocolID() ProtocolPacketID { return protocolCEntitySoundEffect }
func (p *CPacketEntitySoundEffect) Type() PacketType             { return CEntitySoundEffect }
func (p *CPacketEntitySoundEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketSoundEffect struct{}

func (p *CPacketSoundEffect) ProtocolID() ProtocolPacketID { return protocolCSoundEffect }
func (p *CPacketSoundEffect) Type() PacketType             { return CSoundEffect }
func (p *CPacketSoundEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketStopSound struct{}

func (p *CPacketStopSound) ProtocolID() ProtocolPacketID { return protocolCStopSound }
func (p *CPacketStopSound) Type() PacketType             { return CStopSound }
func (p *CPacketStopSound) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketPlayerListHeaderAndFooter struct{}

func (p *CPacketPlayerListHeaderAndFooter) ProtocolID() ProtocolPacketID {
	return protocolCPlayerListHeaderAndFooter
}
func (p *CPacketPlayerListHeaderAndFooter) Type() PacketType     { return CPlayerListHeaderAndFooter }
func (p *CPacketPlayerListHeaderAndFooter) Push(writer buffer.B) { panic("packet not implemented") }

type CPacketNBTQueryResponse struct{}

func (p *CPacketNBTQueryResponse) ProtocolID() ProtocolPacketID { return protocolCNBTQueryResponse }
func (p *CPacketNBTQueryResponse) Type() PacketType             { return CNBTQueryResponse }
func (p *CPacketNBTQueryResponse) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketCollectItem struct{}

func (p *CPacketCollectItem) ProtocolID() ProtocolPacketID { return protocolCCollectItem }
func (p *CPacketCollectItem) Type() PacketType             { return CCollectItem }
func (p *CPacketCollectItem) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityTeleport struct{}

func (p *CPacketEntityTeleport) ProtocolID() ProtocolPacketID { return protocolCEntityTeleport }
func (p *CPacketEntityTeleport) Type() PacketType             { return CEntityTeleport }
func (p *CPacketEntityTeleport) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketAdvancements struct{}

func (p *CPacketAdvancements) ProtocolID() ProtocolPacketID { return protocolCAdvancements }
func (p *CPacketAdvancements) Type() PacketType             { return CAdvancements }
func (p *CPacketAdvancements) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityProperties struct{}

func (p *CPacketEntityProperties) ProtocolID() ProtocolPacketID { return protocolCEntityProperties }
func (p *CPacketEntityProperties) Type() PacketType             { return CEntityProperties }
func (p *CPacketEntityProperties) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketEntityEffect struct{}

func (p *CPacketEntityEffect) ProtocolID() ProtocolPacketID { return protocolCEntityEffect }
func (p *CPacketEntityEffect) Type() PacketType             { return CEntityEffect }
func (p *CPacketEntityEffect) Push(writer buffer.B)         { panic("packet not implemented") }

type CPacketDeclareRecipes struct {
	// Recipes []*Recipe // this doesn't exist yet ;(
	RecipeCount int32
}

func (p *CPacketDeclareRecipes) ProtocolID() ProtocolPacketID { return protocolCDeclareRecipes }
func (p *CPacketDeclareRecipes) Type() PacketType             { return CDeclareRecipes }
func (p *CPacketDeclareRecipes) Push(writer buffer.B) {
	writer.PushVarInt(p.RecipeCount)
	// when recipes are implemented, instead of holding a recipe count, simply write the size of the slice, Recipe will implement BPush
}

type CPacketTags struct{}

func (p *CPacketTags) ProtocolID() ProtocolPacketID { return protocolCTags }
func (p *CPacketTags) Type() PacketType             { return CTags }
func (p *CPacketTags) Push(writer buffer.B)         { panic("packet not implemented") }
