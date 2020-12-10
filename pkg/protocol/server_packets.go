package protocol

import (
	"fmt"

	buff "github.com/alexykot/cncraft/pkg/buffers"
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

func (p *SPacketHandshake) ID() PacketID { return SHandshake }
func (p *SPacketHandshake) Pull(reader buff.Buffer) error {
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
type SPacketRequest struct {
}

func (p *SPacketRequest) ID() PacketID { return SRequest }
func (p *SPacketRequest) Pull(reader buff.Buffer) error {
	// no fields
	return nil
}

type SPacketPing struct {
	Ping int64
}

func (p *SPacketPing) ID() PacketID { return SPing }
func (p *SPacketPing) Pull(reader buff.Buffer) error {
	p.Ping = reader.PullI64()
	return nil // DEBT actually check for errors
}

// LOGIN STATE PACKETS
type SPacketLoginStart struct {
	PlayerName string
}

func (p *SPacketLoginStart) ID() PacketID { return SLoginStart }
func (p *SPacketLoginStart) Pull(reader buff.Buffer) error {
	p.PlayerName = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketEncryptionResponse struct {
	Secret []byte
	Verify []byte
}

func (p *SPacketEncryptionResponse) ID() PacketID { return SEncryptionResponse }
func (p *SPacketEncryptionResponse) Pull(reader buff.Buffer) error {
	p.Secret = reader.PullUAS()
	p.Verify = reader.PullUAS()
	return nil // DEBT actually check for errors
}

type SPacketLoginPluginResponse struct {
	Message int32
	Success bool
	OptData []byte
}

func (p *SPacketLoginPluginResponse) ID() PacketID { return SLoginPluginResponse }
func (p *SPacketLoginPluginResponse) Pull(reader buff.Buffer) error {
	p.Message = reader.PullVrI()
	p.Success = reader.PullBit()
	p.OptData = reader.UAS()[reader.InI():reader.Len()]
	return nil // DEBT actually check for errors
}

// PLAY STATE PACKETS
type SPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *SPacketKeepAlive) ID() PacketID { return SKeepAlive }
func (p *SPacketKeepAlive) Pull(reader buff.Buffer) error {
	p.KeepAliveID = reader.PullI64()
	return nil // DEBT actually check for errors
}

type SPacketChatMessage struct {
	Message string
}

func (p *SPacketChatMessage) ID() PacketID { return SChatMessage }
func (p *SPacketChatMessage) Pull(reader buff.Buffer) error {
	p.Message = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketTeleportConfirm struct {
	TeleportID int32
}

func (p *SPacketTeleportConfirm) ID() PacketID { return STeleportConfirm }
func (p *SPacketTeleportConfirm) Pull(reader buff.Buffer) error {
	p.TeleportID = reader.PullVrI()
	return nil // DEBT actually check for errors
}

type SPacketQueryBlockNBT struct {
	TransactionID int32
	Position      data.PositionI
}

func (p *SPacketQueryBlockNBT) ID() PacketID { return SQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Pull(reader buff.Buffer) error {
	p.TransactionID = reader.PullVrI()
	p.Position = reader.PullPos()
	return nil // DEBT actually check for errors
}

type SPacketSetDifficulty struct {
	Difficult game.Difficulty
}

func (p *SPacketSetDifficulty) ID() PacketID { return SSetDifficulty }
func (p *SPacketSetDifficulty) Pull(reader buff.Buffer) error {
	p.Difficult = game.DifficultyValueOf(reader.PullByt())
	return nil // DEBT actually check for errors
}

// TODO plugins are not supported
//type SPacketPluginMessage struct {
//	Message plugin.Message
//}
//
//func (p *SPacketPluginMessage) ID() PacketID { return SPluginMessage }
//func (p *SPacketPluginMessage) Pull(reader buff.Buffer) error {
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

func (p *SPacketClientStatus) ID() PacketID { return SClientStatus }
func (p *SPacketClientStatus) Pull(reader buff.Buffer) error {
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

func (p *SPacketClientSettings) ID() PacketID { return SClientSettings }
func (p *SPacketClientSettings) Pull(reader buff.Buffer) error {
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

func (p *SPacketPlayerAbilities) ID() PacketID { return SPlayerAbilities }
func (p *SPacketPlayerAbilities) Pull(reader buff.Buffer) error {
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

func (p *SPacketPlayerPosition) ID() PacketID { return SPlayerPosition }
func (p *SPacketPlayerPosition) Pull(reader buff.Buffer) error {
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

func (p *SPacketPlayerLocation) ID() PacketID { return SPlayerLocation }
func (p *SPacketPlayerLocation) Pull(reader buff.Buffer) error {
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

func (p *SPacketPlayerRotation) ID() PacketID { return SPlayerRotation }
func (p *SPacketPlayerRotation) Pull(reader buff.Buffer) error {
	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.OnGround = reader.PullBit()
	return nil // DEBT actually check for errors
}
