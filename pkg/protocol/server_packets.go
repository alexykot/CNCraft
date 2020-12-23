package protocol

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/player"
)

// HANDSHAKE STATE PACKETS
type SPacketHandshake struct {
	version int32

	host string
	port uint16

	NextState State
}

func (p *SPacketHandshake) ProtocolID() ProtocolPacketID { return protocolSHandshake }
func (p *SPacketHandshake) Type() PacketType             { return SHandshake }
func (p *SPacketHandshake) Pull(reader buffer.B) error {
	var err error

	p.version = reader.PullVrI()
	p.host = reader.PullTxt()
	p.port = reader.PullU16()

	nextState := reader.PullVrI()

	if p.NextState, err = IntToState(int(nextState)); err != nil {
		return fmt.Errorf("failed to parse handshake  next state: %w", err)
	}

	return nil
}

// STATUS STATE PACKETS
type SPacketRequest struct{}

func (p *SPacketRequest) ProtocolID() ProtocolPacketID { return protocolSRequest }
func (p *SPacketRequest) Type() PacketType             { return SRequest }
func (p *SPacketRequest) Pull(reader buffer.B) error {
	// no fields
	return nil
}

type SPacketPing struct {
	Payload int64
}

func (p *SPacketPing) ProtocolID() ProtocolPacketID { return protocolSPing }
func (p *SPacketPing) Type() PacketType             { return SPing }
func (p *SPacketPing) Pull(reader buffer.B) error {
	p.Payload = reader.PullI64()
	return nil // DEBT actually check for errors
}

// LOGIN STATE PACKETS
type SPacketLoginStart struct {
	Username string
}

func (p *SPacketLoginStart) ProtocolID() ProtocolPacketID { return protocolSLoginStart }
func (p *SPacketLoginStart) Type() PacketType             { return SLoginStart }
func (p *SPacketLoginStart) Pull(reader buffer.B) error {
	p.Username = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketEncryptionResponse struct {
	Secret []byte
	Verify []byte
}

func (p *SPacketEncryptionResponse) ProtocolID() ProtocolPacketID { return protocolSEncryptionResponse }
func (p *SPacketEncryptionResponse) Type() PacketType             { return SEncryptionResponse }
func (p *SPacketEncryptionResponse) Pull(reader buffer.B) error {
	p.Secret = reader.PullUAS()
	p.Verify = reader.PullUAS()
	return nil // DEBT actually check for errors
}

type SPacketLoginPluginResponse struct {
	Message int32
	Success bool
	OptData []byte
}

func (p *SPacketLoginPluginResponse) ProtocolID() ProtocolPacketID {
	return protocolSLoginPluginResponse
}
func (p *SPacketLoginPluginResponse) Type() PacketType { return SLoginPluginResponse }
func (p *SPacketLoginPluginResponse) Pull(reader buffer.B) error {
	p.Message = reader.PullVrI()
	p.Success = reader.PullBit()
	p.OptData = reader.UAS()[reader.InI():reader.Len()]
	return nil // DEBT actually check for errors
}

// PLAY STATE PACKETS
type SPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *SPacketKeepAlive) ProtocolID() ProtocolPacketID { return protocolSKeepAlive }
func (p *SPacketKeepAlive) Type() PacketType             { return SKeepAlive }
func (p *SPacketKeepAlive) Pull(reader buffer.B) error {
	p.KeepAliveID = reader.PullI64()
	return nil // DEBT actually check for errors
}

type SPacketChatMessage struct {
	Message string
}

func (p *SPacketChatMessage) ProtocolID() ProtocolPacketID { return protocolSChatMessage }
func (p *SPacketChatMessage) Type() PacketType             { return SChatMessage }
func (p *SPacketChatMessage) Pull(reader buffer.B) error {
	p.Message = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketTeleportConfirm struct {
	TeleportID int32
}

func (p *SPacketTeleportConfirm) ProtocolID() ProtocolPacketID { return protocolSTeleportConfirm }
func (p *SPacketTeleportConfirm) Type() PacketType             { return STeleportConfirm }
func (p *SPacketTeleportConfirm) Pull(reader buffer.B) error {
	p.TeleportID = reader.PullVrI()
	return nil // DEBT actually check for errors
}

type SPacketQueryBlockNBT struct {
	TransactionID int32
	Position      data.PositionI
}

func (p *SPacketQueryBlockNBT) ProtocolID() ProtocolPacketID { return protocolSQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Type() PacketType             { return SQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Pull(reader buffer.B) error {
	p.TransactionID = reader.PullVrI()
	p.Position = reader.PullPos()
	return nil // DEBT actually check for errors
}

type SPacketSetDifficulty struct {
	Difficult game.Difficulty
}

func (p *SPacketSetDifficulty) ProtocolID() ProtocolPacketID { return protocolSSetDifficulty }
func (p *SPacketSetDifficulty) Type() PacketType             { return SSetDifficulty }
func (p *SPacketSetDifficulty) Pull(reader buffer.B) error {
	p.Difficult = game.DifficultyValueOf(reader.PullByt())
	return nil // DEBT actually check for errors
}

// TODO plugins are not supported
//type SPacketPluginMessage struct {
//	Message plugin.Message
//}
//
//func (p *SPacketPluginMessage) Type() PacketType { return SPluginMessage }
//func (p *SPacketPluginMessage) Pull(reader buffers.Buffer) error {
//	channel := reader.PullTxt()
//	message := plugin.GetMessageForChannel(channel)
//
//	if message == nil {
//		return fmt.Errorf("channel `%s` not found ", channel)
//	}
//
//	message.Pull(reader)
//
//	p.Message = message
//
//	return nil // DEBT actually check for errors
//}

type SPacketClientStatus struct {
	Action player.StatusAction
}

func (p *SPacketClientStatus) ProtocolID() ProtocolPacketID { return protocolSClientStatus }
func (p *SPacketClientStatus) Type() PacketType             { return SClientStatus }
func (p *SPacketClientStatus) Pull(reader buffer.B) error {
	p.Action = player.StatusAction(reader.PullVrI())
	return nil // DEBT actually check for errors
}

type SPacketClientSettings struct {
	Locale       string
	ViewDistance byte
	ChatMode     player.ChatMode
	ChatColors   bool // if false, strip messages of colors before sending
	SkinParts    player.SkinParts
	MainHand     player.MainHand
}

func (p *SPacketClientSettings) ProtocolID() ProtocolPacketID { return protocolSClientSettings }
func (p *SPacketClientSettings) Type() PacketType             { return SClientSettings }
func (p *SPacketClientSettings) Pull(reader buffer.B) error {
	p.Locale = reader.PullTxt()
	p.ViewDistance = reader.PullByt()
	p.ChatMode = player.ChatMode(reader.PullVrI())
	p.ChatColors = reader.PullBit()

	parts := player.SkinParts{}
	parts.Pull(reader)

	p.SkinParts = parts
	p.MainHand = player.MainHand(reader.PullVrI())
	return nil // DEBT actually check for errors
}

type SPacketPlayerAbilities struct {
	Abilities   player.PlayerAbilities
	FlightSpeed float32
	GroundSpeed float32
}

func (p *SPacketPlayerAbilities) ProtocolID() ProtocolPacketID { return protocolSPlayerAbilities }
func (p *SPacketPlayerAbilities) Type() PacketType             { return SPlayerAbilities }
func (p *SPacketPlayerAbilities) Pull(reader buffer.B) error {
	abilities := player.PlayerAbilities{}
	abilities.Pull(reader)

	p.Abilities = abilities

	p.FlightSpeed = reader.PullF32()
	p.GroundSpeed = reader.PullF32()
	return nil // DEBT actually check for errors
}

type SPacketPlayerPosition struct {
	Position data.PositionF
	OnGround bool
}

func (p *SPacketPlayerPosition) ProtocolID() ProtocolPacketID { return protocolSPlayerPosition }
func (p *SPacketPlayerPosition) Type() PacketType             { return SPlayerPosition }
func (p *SPacketPlayerPosition) Pull(reader buffer.B) error {
	p.Position = data.PositionF{
		X: reader.PullF64(),
		Y: reader.PullF64(),
		Z: reader.PullF64(),
	}

	p.OnGround = reader.PullBit()
	return nil // DEBT actually check for errors
}

type SPacketPlayerLocation struct {
	Location data.Location
	OnGround bool
}

func (p *SPacketPlayerLocation) ProtocolID() ProtocolPacketID { return protocolSPlayerLocation }
func (p *SPacketPlayerLocation) Type() PacketType             { return SPlayerLocation }
func (p *SPacketPlayerLocation) Pull(reader buffer.B) error {
	p.Location = data.Location{
		PositionF: data.PositionF{
			X: reader.PullF64(),
			Y: reader.PullF64(),
			Z: reader.PullF64(),
		},
		RotationF: data.RotationF{
			AxisX: reader.PullF32(),
			AxisY: reader.PullF32(),
		},
	}

	p.OnGround = reader.PullBit()
	return nil // DEBT actually check for errors
}

type SPacketPlayerRotation struct {
	Rotation data.RotationF
	OnGround bool
}

func (p *SPacketPlayerRotation) ProtocolID() ProtocolPacketID { return protocolSPlayerRotation }
func (p *SPacketPlayerRotation) Type() PacketType             { return SPlayerRotation }
func (p *SPacketPlayerRotation) Pull(reader buffer.B) error {
	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.OnGround = reader.PullBit()
	return nil // DEBT actually check for errors
}
