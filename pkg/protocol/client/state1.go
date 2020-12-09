package client

import (
	"encoding/json"
	"github.com/alexykot/cncraft/apis/buff"
	"github.com/alexykot/cncraft/impl/base"
	"github.com/alexykot/cncraft/impl/data/status"
	"github.com/alexykot/cncraft/impl/protocol"
)

// done

type CPacketResponse struct {
	Status status.Response
}

func (p *CPacketResponse) ID() protocol.PacketID { return protocol.CResponse }
func (p *CPacketResponse) Push(writer buff.Buffer, conn base.Connection) {
	if text, err := json.Marshal(p.Status); err != nil {
		panic(err)
	} else {
		writer.PushTxt(string(text))
	}
}

type CPacketPong struct {
	Ping int64
}

func (p *CPacketPong) ID() protocol.PacketID { return protocol.CPong }
func (p *CPacketPong) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI64(p.Ping)
}
