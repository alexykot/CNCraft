package client

import (
	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/apis/data/msgs"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/protocol"
)

// done

type CPacketDisconnect struct {
	Reason msgs.Message
}

func (p *CPacketDisconnect) ID() protocol.PacketID { return protocol.CDisconnect }
func (p *CPacketDisconnect) Push(writer buff.Buffer, conn base.Connection) {
	message := p.Reason

	writer.PushTxt(message.AsJson())
}

type CPacketEncryptionRequest struct {
	Server string // unused?
	Public []byte
	Verify []byte
}

func (p *CPacketEncryptionRequest) ID() protocol.PacketID { return protocol.CEncryptionRequest }
func (p *CPacketEncryptionRequest) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.Server)
	writer.PushUAS(p.Public, true)
	writer.PushUAS(p.Verify, true)
}

type CPacketLoginSuccess struct {
	PlayerUUID string
	PlayerName string
}

func (p *CPacketLoginSuccess) ID() protocol.PacketID { return protocol.CLoginSuccess }
func (p *CPacketLoginSuccess) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushTxt(p.PlayerUUID)
	writer.PushTxt(p.PlayerName)
}

type CPacketSetCompression struct {
	Threshold int32
}

func (p *CPacketSetCompression) ID() protocol.PacketID { return protocol.CSetCompression }
func (p *CPacketSetCompression) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.Threshold)
}

type CPacketLoginPluginRequest struct {
	MessageID int32
	Channel   string
	OptData   []byte
}

func (p *CPacketLoginPluginRequest) ID() protocol.PacketID { return protocol.CLoginPluginRequest }
func (p *CPacketLoginPluginRequest) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushVrI(p.MessageID)
	writer.PushTxt(p.Channel)
	writer.PushUAS(p.OptData, false)
}
