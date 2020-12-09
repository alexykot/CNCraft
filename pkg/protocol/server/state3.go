package server

import (
	"fmt"
	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/apis/data"
	"github.com/alexykot/cncraft/apis/game"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/data/client"
	"github.com/alexykot/cncraft/impl/data/plugin"
	"github.com/alexykot/cncraft/impl/protocol"
)

type SPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *SPacketKeepAlive) ID() protocol.PacketID { return protocol.SKeepAlive }
func (p *SPacketKeepAlive) Pull(reader buff.Buffer, conn base.Connection) error {
	p.KeepAliveID = reader.PullI64()
	return nil // DEBT actually check for errors
}

type SPacketChatMessage struct {
	Message string
}

func (p *SPacketChatMessage) ID() protocol.PacketID { return protocol.SChatMessage }
func (p *SPacketChatMessage) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Message = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketTeleportConfirm struct {
	TeleportID int32
}

func (p *SPacketTeleportConfirm) ID() protocol.PacketID { return protocol.STeleportConfirm }
func (p *SPacketTeleportConfirm) Pull(reader buff.Buffer, conn base.Connection) error {
	p.TeleportID = reader.PullVrI()
	return nil // DEBT actually check for errors
}

type SPacketQueryBlockNBT struct {
	TransactionID int32
	Position      data.PositionI
}

func (p *SPacketQueryBlockNBT) ID() protocol.PacketID { return protocol.SQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Pull(reader buff.Buffer, conn base.Connection) error {
	p.TransactionID = reader.PullVrI()
	p.Position = reader.PullPos()
	return nil // DEBT actually check for errors
}

type SPacketSetDifficulty struct {
	Difficult game.Difficulty
}

func (p *SPacketSetDifficulty) ID() protocol.PacketID { return protocol.SSetDifficulty }
func (p *SPacketSetDifficulty) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Difficult = game.DifficultyValueOf(reader.PullByt())
	return nil // DEBT actually check for errors
}

type SPacketPluginMessage struct {
	Message plugin.Message
}

func (p *SPacketPluginMessage) ID() protocol.PacketID { return protocol.SPluginMessage }
func (p *SPacketPluginMessage) Pull(reader buff.Buffer, conn base.Connection) error {
	channel := reader.PullTxt()
	message := plugin.GetMessageForChannel(channel)

	if message == nil {
		return fmt.Errorf("channel `%s` not found ", channel)
	}

	message.Pull(reader)

	p.Message = message

	return nil // DEBT actually check for errors
}

type SPacketClientStatus struct {
	Action client.StatusAction
}

func (p *SPacketClientStatus) ID() protocol.PacketID { return protocol.SClientStatus }
func (p *SPacketClientStatus) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Action = client.StatusAction(reader.PullVrI())
	return nil // DEBT actually check for errors
}

type SPacketClientSettings struct {
	Locale       string
	ViewDistance byte
	ChatMode     client.ChatMode
	ChatColors   bool // if false, strip messages of colors before sending
	SkinParts    client.SkinParts
	MainHand     client.MainHand
}

func (p *SPacketClientSettings) ID() protocol.PacketID { return protocol.SClientSettings }
func (p *SPacketClientSettings) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Locale = reader.PullTxt()
	p.ViewDistance = reader.PullByt()
	p.ChatMode = client.ChatMode(reader.PullVrI())
	p.ChatColors = reader.PullBit()

	parts := client.SkinParts{}
	parts.Pull(reader)

	p.SkinParts = parts
	p.MainHand = client.MainHand(reader.PullVrI())
	return nil // DEBT actually check for errors
}

type SPacketPlayerAbilities struct {
	Abilities   client.PlayerAbilities
	FlightSpeed float32
	GroundSpeed float32
}

func (p *SPacketPlayerAbilities) ID() protocol.PacketID { return protocol.SPlayerAbilities }
func (p *SPacketPlayerAbilities) Pull(reader buff.Buffer, conn base.Connection) error {
	abilities := client.PlayerAbilities{}
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

func (p *SPacketPlayerPosition) ID() protocol.PacketID { return protocol.SPlayerPosition }
func (p *SPacketPlayerPosition) Pull(reader buff.Buffer, conn base.Connection) error {
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

func (p *SPacketPlayerLocation) ID() protocol.PacketID { return protocol.SPlayerLocation }
func (p *SPacketPlayerLocation) Pull(reader buff.Buffer, conn base.Connection) error {
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

func (p *SPacketPlayerRotation) ID() protocol.PacketID { return protocol.SPlayerRotation }
func (p *SPacketPlayerRotation) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.OnGround = reader.PullBit()
	return nil // DEBT actually check for errors
}
