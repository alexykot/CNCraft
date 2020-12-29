package protocol

import (
	"encoding/json"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/chat"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/entities"
	"github.com/alexykot/cncraft/pkg/game/level"
	"github.com/alexykot/cncraft/pkg/game/players"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
	"github.com/alexykot/cncraft/pkg/protocol/status"
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
		writer.PushTxt(string(text))
	}
}

type CPacketPong struct {
	Payload int64
}

func (p *CPacketPong) ProtocolID() ProtocolPacketID { return protocolCPong }
func (p *CPacketPong) Type() PacketType             { return CPong }
func (p *CPacketPong) Push(writer buffer.B) {
	writer.PushI64(p.Payload)
}

// LOGIN STATE PACKETS
type CPacketDisconnect struct {
	Reason chat.Message
}

func (p *CPacketDisconnect) ProtocolID() ProtocolPacketID { return protocolCDisconnect }
func (p *CPacketDisconnect) Type() PacketType             { return CDisconnect }
func (p *CPacketDisconnect) Push(writer buffer.B) {
	message := p.Reason

	writer.PushTxt(message.AsJson())
}

type CPacketEncryptionRequest struct {
	ServerID    string // Appears to be unused by the Notchian client.
	PublicKey   []byte
	VerifyToken []byte
}

func (p *CPacketEncryptionRequest) ProtocolID() ProtocolPacketID { return protocolCEncryptionRequest }
func (p *CPacketEncryptionRequest) Type() PacketType             { return CEncryptionRequest }
func (p *CPacketEncryptionRequest) Push(writer buffer.B) {
	writer.PushTxt(p.ServerID)
	writer.PushUAS(p.PublicKey, true)
	writer.PushUAS(p.VerifyToken, true)
}

type CPacketLoginSuccess struct {
	PlayerUUID string
	PlayerName string
}

func (p *CPacketLoginSuccess) ProtocolID() ProtocolPacketID { return protocolCLoginSuccess }
func (p *CPacketLoginSuccess) Type() PacketType             { return CLoginSuccess }
func (p *CPacketLoginSuccess) Push(writer buffer.B) {
	writer.PushTxt(p.PlayerUUID)
	writer.PushTxt(p.PlayerName)
}

type CPacketSetCompression struct {
	Threshold int32
}

func (p *CPacketSetCompression) ProtocolID() ProtocolPacketID { return protocolCSetCompression }
func (p *CPacketSetCompression) Type() PacketType             { return CSetCompression }
func (p *CPacketSetCompression) Push(writer buffer.B) {
	writer.PushVrI(p.Threshold)
}

type CPacketLoginPluginRequest struct {
	MessageID int32
	Channel   string
	OptData   []byte
}

func (p *CPacketLoginPluginRequest) ProtocolID() ProtocolPacketID { return protocolCLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Type() PacketType             { return CLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Push(writer buffer.B) {
	writer.PushVrI(p.MessageID)
	writer.PushTxt(p.Channel)
	writer.PushUAS(p.OptData, false)
}

// PLAY STATE PACKETS
type CPacketChatMessage struct {
	Message         chat.Message
	MessagePosition chat.MessagePosition
}

func (p *CPacketChatMessage) ProtocolID() ProtocolPacketID { return protocolCChatMessage }
func (p *CPacketChatMessage) Type() PacketType             { return CChatMessage }
func (p *CPacketChatMessage) Push(writer buffer.B) {
	message := p.Message

	if p.MessagePosition == chat.HotBarText {
		message = *chat.New(message.AsText())
	}

	writer.PushTxt(message.AsJson())
	writer.PushByt(byte(p.MessagePosition))
}

type CPacketJoinGame struct {
	EntityID      int32
	IsHardcore    game.Coreness
	GameMode      game.Gamemode
	Dimension     game.Dimension
	HashedSeed    int64
	LevelType     game.WorldType
	ViewDistance  int32
	ReduceDebug   bool
	RespawnScreen bool
}

func (p *CPacketJoinGame) ProtocolID() ProtocolPacketID { return protocolCJoinGame }
func (p *CPacketJoinGame) Type() PacketType             { return CJoinGame }
func (p *CPacketJoinGame) Push(writer buffer.B) {
	writer.PushI32(p.EntityID)
	writer.PushByt(p.GameMode.Encoded(bool(p.IsHardcore)))
	writer.PushI32(int32(p.Dimension))
	writer.PushI64(p.HashedSeed)
	writer.PushByt(uint8(0)) // is ignored by the Notchian client
	writer.PushTxt(p.LevelType.String())
	writer.PushVrI(p.ViewDistance)
	writer.PushBit(p.ReduceDebug)
	writer.PushBit(p.RespawnScreen)
}

type CPacketPluginMessage struct {
	Message plugin.Message
}

func (p *CPacketPluginMessage) ProtocolID() ProtocolPacketID { return protocolCPluginMessage }
func (p *CPacketPluginMessage) Type() PacketType             { return CPluginMessage }
func (p *CPacketPluginMessage) Push(writer buffer.B) {
	writer.PushTxt(p.Message.Chan())
	p.Message.Push(writer)
}

type CPacketPlayerLocation struct {
	Location data.Location
	Relative players.Relativity

	SomeID int32 // no idea what ID is this, the packet type 3/0x36 in the protocol 754 does not have this field
}

func (p *CPacketPlayerLocation) ProtocolID() ProtocolPacketID { return protocolCPlayerLocation }
func (p *CPacketPlayerLocation) Type() PacketType             { return CPlayerLocation }
func (p *CPacketPlayerLocation) Push(writer buffer.B) {
	writer.PushF64(p.Location.X)
	writer.PushF64(p.Location.Y)
	writer.PushF64(p.Location.Z)

	writer.PushF32(p.Location.AxisX)
	writer.PushF32(p.Location.AxisY)

	p.Relative.Push(writer)

	writer.PushVrI(p.SomeID)
}

type CPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *CPacketKeepAlive) ProtocolID() ProtocolPacketID { return protocolCKeepAlive }
func (p *CPacketKeepAlive) Type() PacketType             { return CKeepAlive }
func (p *CPacketKeepAlive) Push(writer buffer.B) {
	writer.PushI64(p.KeepAliveID)
}

type CPacketServerDifficulty struct {
	Difficulty game.Difficulty
	Locked     bool // should probably always be true
}

func (p *CPacketServerDifficulty) ProtocolID() ProtocolPacketID { return protocolCServerDifficulty }
func (p *CPacketServerDifficulty) Type() PacketType             { return CServerDifficulty }
func (p *CPacketServerDifficulty) Push(writer buffer.B) {
	writer.PushByt(byte(p.Difficulty))
	writer.PushBit(p.Locked)
}

type CPacketPlayerAbilities struct {
	Abilities   players.PlayerAbilities
	FlyingSpeed float32
	FieldOfView float32
}

func (p *CPacketPlayerAbilities) ProtocolID() ProtocolPacketID { return protocolCPlayerAbilities }
func (p *CPacketPlayerAbilities) Type() PacketType             { return CPlayerAbilities }
func (p *CPacketPlayerAbilities) Push(writer buffer.B) {
	p.Abilities.Push(writer)

	writer.PushF32(p.FlyingSpeed)
	writer.PushF32(p.FieldOfView)
}

type CPacketHeldItemChange struct {
	Slot players.HotBarSlot
}

func (p *CPacketHeldItemChange) ProtocolID() ProtocolPacketID { return protocolCHeldItemChange }
func (p *CPacketHeldItemChange) Type() PacketType             { return CHeldItemChange }
func (p *CPacketHeldItemChange) Push(writer buffer.B) {
	writer.PushByt(byte(p.Slot))
}

type CPacketDeclareRecipes struct {
	// Recipes []*Recipe // this doesn't exist yet ;(
	RecipeCount int32
}

func (p *CPacketDeclareRecipes) ProtocolID() ProtocolPacketID { return protocolCDeclareRecipes }
func (p *CPacketDeclareRecipes) Type() PacketType             { return CDeclareRecipes }
func (p *CPacketDeclareRecipes) Push(writer buffer.B) {
	writer.PushVrI(p.RecipeCount)
	// when recipes are implemented, instead of holding a recipe count, simply write the size of the slice, Recipe will implement BufferPush
}

type CPacketChunkData struct {
	Chunk level.Chunk
}

func (p *CPacketChunkData) ProtocolID() ProtocolPacketID { return protocolCChunkData }
func (p *CPacketChunkData) Type() PacketType             { return CChunkData }
func (p *CPacketChunkData) Push(writer buffer.B) {
	writer.PushI32(int32(p.Chunk.ChunkX()))
	writer.PushI32(int32(p.Chunk.ChunkZ()))

	// full chunk (for now >:D)
	writer.PushBit(true)

	chunkData := buffer.New()
	p.Chunk.Push(chunkData) // write chunk data and primary bit mask

	// pull primary bit mask and push to writer
	writer.PushVrI(chunkData.PullVrI())

	// write height-maps
	writer.PushNbt(p.Chunk.HeightMapNbtCompound())

	biomes := make([]int32, 1024, 1024)
	for i := range biomes {
		biomes[i] = 0 // void biome
	}

	for _, biome := range biomes {
		writer.PushI32(biome)
	}

	// data, prefixed with len
	writer.PushUAS(chunkData.UAS(), true)

	// write block entities
	writer.PushVrI(0)
}

type CPacketPlayerInfo struct {
	Action players.PlayerInfoAction
	Values []players.PlayerInfo
}

func (p *CPacketPlayerInfo) ProtocolID() ProtocolPacketID { return protocolCPlayerInfo }
func (p *CPacketPlayerInfo) Type() PacketType             { return CPlayerInfo }
func (p *CPacketPlayerInfo) Push(writer buffer.B) {
	writer.PushVrI(int32(p.Action))
	writer.PushVrI(int32(len(p.Values)))

	for _, value := range p.Values {
		value.Push(writer)
	}
}

type CPacketEntityMetadata struct {
	Entity entities.Entity
}

func (p *CPacketEntityMetadata) ProtocolID() ProtocolPacketID { return protocolCEntityMetadata }
func (p *CPacketEntityMetadata) Type() PacketType             { return CEntityMetadata }
func (p *CPacketEntityMetadata) Push(writer buffer.B) {
	writer.PushVrI(int32(p.Entity.ID())) // questionable...

	// only supporting player metadata for now
	_, ok := p.Entity.(entities.PlayerCharacter)
	if ok {

		writer.PushByt(16) // index | displayed skin parts
		writer.PushVrI(0)  // type | byte

		skin := players.SkinParts{
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

	writer.PushByt(0xFF)
}
