package server

import (
	"fmt"
	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/protocol"
)

type SPacketHandshake struct {
	version int32

	host string
	port uint16

	NextState protocol.State
}

func (p *SPacketHandshake) ID() protocol.PacketID { return protocol.SHandshake }
func (p *SPacketHandshake) Pull(reader buff.Buffer, conn base.Connection) error {
	var err error

	p.version = reader.PullVrI()
	p.host = reader.PullTxt()
	p.port = reader.PullU16()

	nextState := reader.PullVrI()

	if p.NextState, err = protocol.IntToState(int(nextState)); err != nil {
		return fmt.Errorf("failed to parse handshake  next state: %w", err)
	}

	return nil
}
