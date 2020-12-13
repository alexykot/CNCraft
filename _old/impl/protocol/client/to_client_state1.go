package client

import (
	"encoding/json"
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/status"
	"github.com/golangmc/minecraft-server/impl/protocol"
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
