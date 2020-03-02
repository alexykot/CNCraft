package server

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/protocol"
)

// done

type SPacketRequest struct {
}

func (p *SPacketRequest) ID() protocol.PacketID { return protocol.SRequest }
func (p *SPacketRequest) Pull(reader buff.Buffer, conn base.Connection) error {
	// no fields
	return nil
}

type SPacketPing struct {
	Ping int64
}

func (p *SPacketPing) ID() protocol.PacketID { return protocol.SPing }
func (p *SPacketPing) Pull(reader buff.Buffer, conn base.Connection) error {
	p.Ping = reader.PullI64()
	return nil // DEBT actually check for errors
}
