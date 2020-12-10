package protocol

import (
	"encoding/json"

	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/apis/data"
	"github.com/alexykot/cncraft/apis/data/msgs"
	"github.com/alexykot/cncraft/apis/ents"
	"github.com/alexykot/cncraft/apis/game"
	"github.com/alexykot/cncraft/apis/game/level"
	apisConn "github.com/alexykot/cncraft/impl/conn"
	"github.com/alexykot/cncraft/impl/data/client"
	"github.com/alexykot/cncraft/impl/data/plugin"
	"github.com/alexykot/cncraft/impl/data/status"
)

// HANDSHAKE STATE CLIENT BOUND PACKETS DO NOT EXIST

// STATUS STATE PACKETS
type CPacketResponse struct {
	Status status.Response
}

func (p *CPacketResponse) ID() PacketID { return CResponse }
func (p *CPacketResponse) Push(writer buff.Buffer) {
	if text, err := json.Marshal(p.Status); err != nil {
		panic(err)
	} else {
		writer.PushTxt(string(text))
	}
}

type CPacketPong struct {
	Ping int64
}

func (p *CPacketPong) ID() PacketID { return CPong }
func (p *CPacketPong) Push(writer buff.Buffer) {
	writer.PushI64(p.Ping)
}

// LOGIN STATE PACKETS
type CPacketDisconnect struct {
	Reason msgs.Message
}

func (p *CPacketDisconnect) ID() PacketID { return CDisconnect }
func (p *CPacketDisconnect) Push(writer buff.Buffer) {
	message := p.Reason

	writer.PushTxt(message.AsJson())
}

type CPacketEncryptionRequest struct {
	Server string // unused?
	Public []byte
	Verify []byte
}

func (p *CPacketEncryptionRequest) ID() PacketID { return CEncryptionRequest }
func (p *CPacketEncryptionRequest) Push(writer buff.Buffer) {
	writer.PushTxt(p.Server)
	writer.PushUAS(p.Public, true)
	writer.PushUAS(p.Verify, true)
}

type CPacketLoginSuccess struct {
	PlayerUUID string
	PlayerName string
}

func (p *CPacketLoginSuccess) ID() PacketID { return CLoginSuccess }
func (p *CPacketLoginSuccess) Push(writer buff.Buffer) {
	writer.PushTxt(p.PlayerUUID)
	writer.PushTxt(p.PlayerName)
}

type CPacketSetCompression struct {
	Threshold int32
}

func (p *CPacketSetCompression) ID() PacketID { return CSetCompression }
func (p *CPacketSetCompression) Push(writer buff.Buffer) {
	writer.PushVrI(p.Threshold)
}

type CPacketLoginPluginRequest struct {
	MessageID int32
	Channel   string
	OptData   []byte
}

func (p *CPacketLoginPluginRequest) ID() PacketID { return CLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Push(writer buff.Buffer) {
	writer.PushVrI(p.MessageID)
	writer.PushTxt(p.Channel)
	writer.PushUAS(p.OptData, false)
}

// PLAY STATE PACKETS
type CPacketChatMessage struct {
	Message         msgs.Message
	MessagePosition msgs.MessagePosition
}

func (p *CPacketChatMessage) ID() PacketID { return CChatMessage }
func (p *CPacketChatMessage) Push(writer buff.Buffer) {
	message := p.Message

	if p.MessagePosition == msgs.HotBarText {
		message = *msgs.New(message.AsText())
	}

	writer.PushTxt(message.AsJson())
	writer.PushByt(byte(p.MessagePosition))
}

type CPacketJoinGame struct {
	EntityID      int32
	Hardcore      bool
	GameMode      game.GameMode
	Dimension     game.Dimension
	HashedSeed    int64
	MaxPlayers    int
	LevelType     game.LevelType
	ViewDistance  int32
	ReduceDebug   bool
	RespawnScreen bool
}

func (p *CPacketJoinGame) ID() PacketID { return CJoinGame }
func (p *CPacketJoinGame) Push(writer buff.Buffer) {
	writer.PushI32(p.EntityID)
	writer.PushByt(p.GameMode.Encoded(p.Hardcore /* pull this value from somewhere */))
	writer.PushI32(int32(p.Dimension))
	writer.PushI64(p.HashedSeed)
	writer.PushByt(uint8(p.MaxPlayers))
	writer.PushTxt(p.LevelType.String())
	writer.PushVrI(p.ViewDistance)
	writer.PushBit(p.ReduceDebug)
	writer.PushBit(p.RespawnScreen)
}

type CPacketPluginMessage struct {
	Message plugin.Message
}

func (p *CPacketPluginMessage) ID() PacketID { return CPluginMessage }
func (p *CPacketPluginMessage) Push(writer buff.Buffer) {
	writer.PushTxt(p.Message.Chan())
	p.Message.Push(writer)
}

type CPacketPlayerLocation struct {
	Location data.Location
	Relative client.Relativity

	SomeID int32 // no idea what ID is this, the packet type 3/0x36 in the protocol 754 does not have this field
}

func (p *CPacketPlayerLocation) ID() PacketID { return CPlayerLocation }
func (p *CPacketPlayerLocation) Push(writer buff.Buffer) {
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

func (p *CPacketKeepAlive) ID() PacketID { return CKeepAlive }
func (p *CPacketKeepAlive) Push(writer buff.Buffer) {
	writer.PushI64(p.KeepAliveID)
}

type CPacketServerDifficulty struct {
	Difficulty game.Difficulty
	Locked     bool // should probably always be true
}

func (p *CPacketServerDifficulty) ID() PacketID { return CServerDifficulty }
func (p *CPacketServerDifficulty) Push(writer buff.Buffer) {
	writer.PushByt(byte(p.Difficulty))
	writer.PushBit(p.Locked)
}

type CPacketPlayerAbilities struct {
	Abilities   client.PlayerAbilities
	FlyingSpeed float32
	FieldOfView float32
}

func (p *CPacketPlayerAbilities) ID() PacketID { return CPlayerAbilities }
func (p *CPacketPlayerAbilities) Push(writer buff.Buffer) {
	p.Abilities.Push(writer)

	writer.PushF32(p.FlyingSpeed)
	writer.PushF32(p.FieldOfView)
}

type CPacketHeldItemChange struct {
	Slot client.HotBarSlot
}

func (p *CPacketHeldItemChange) ID() PacketID { return CHeldItemChange }
func (p *CPacketHeldItemChange) Push(writer buff.Buffer) {
	writer.PushByt(byte(p.Slot))
}

type CPacketDeclareRecipes struct {
	// Recipes []*Recipe // this doesn't exist yet ;(
	RecipeCount int32
}

func (p *CPacketDeclareRecipes) ID() PacketID { return CDeclareRecipes }
func (p *CPacketDeclareRecipes) Push(writer buff.Buffer) {
	writer.PushVrI(p.RecipeCount)
	// when recipes are implemented, instead of holding a recipe count, simply write the size of the slice, Recipe will implement BufferPush
}

type CPacketChunkData struct {
	Chunk level.Chunk
}

func (p *CPacketChunkData) ID() PacketID { return CChunkData }
func (p *CPacketChunkData) Push(writer buff.Buffer) {
	writer.PushI32(int32(p.Chunk.ChunkX()))
	writer.PushI32(int32(p.Chunk.ChunkZ()))

	// full chunk (for now >:D)
	writer.PushBit(true)

	chunkData := apisConn.NewBuffer()
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
	Action client.PlayerInfoAction
	Values []client.PlayerInfo
}

func (p *CPacketPlayerInfo) ID() PacketID { return CPlayerInfo }
func (p *CPacketPlayerInfo) Push(writer buff.Buffer) {
	writer.PushVrI(int32(p.Action))
	writer.PushVrI(int32(len(p.Values)))

	for _, value := range p.Values {
		value.Push(writer)
	}
}

type CPacketEntityMetadata struct {
	Entity ents.Entity
}

func (p *CPacketEntityMetadata) ID() PacketID { return CEntityMetadata }
func (p *CPacketEntityMetadata) Push(writer buff.Buffer) {
	writer.PushVrI(int32(p.Entity.EntityUUID())) // questionable...

	// only supporting player metadata for now
	_, ok := p.Entity.(ents.Player)
	if ok {

		writer.PushByt(16) // index | displayed skin parts
		writer.PushVrI(0)  // type | byte

		skin := client.SkinParts{
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
