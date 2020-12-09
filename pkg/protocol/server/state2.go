package server

import (
	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/protocol"
)

// done

type SPacketLoginStart struct {
	PlayerName string
}

func (p *SPacketLoginStart) ID() protocol.PacketID { return protocol.SLoginStart }
func (p *SPacketLoginStart) Pull(reader buff.Buffer, conn base.Connection) error {
	p.PlayerName = reader.PullTxt()
	return nil // DEBT actually check for errors
}

type SPacketEncryptionResponse struct {
	Secret []byte
	Verify []byte
}

func (p *SPacketEncryptionResponse) ID() protocol.PacketID { return protocol.SEncryptionResponse }
func (p *SPacketEncryptionResponse) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Secret = reader.PullUAS()
	p.Verify = reader.PullUAS()
	return nil // DEBT actually check for errors
}

type SPacketLoginPluginResponse struct {
	Message int32
	Success bool
	OptData []byte
}

func (p *SPacketLoginPluginResponse) ID() protocol.PacketID { return protocol.SLoginPluginResponse }
func (p *SPacketLoginPluginResponse) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Message = reader.PullVrI()
	p.Success = reader.PullBit()
	p.OptData = reader.UAS()[reader.InI():reader.Len()]
	return nil // DEBT actually check for errors
}
