package handlers

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
)

func mkCpacketEnvelope(cpacket protocol.CPacket) *envelope.E {
	bufOut := buffer.New()
	cpacket.Push(bufOut)
	return envelope.CPacket(&pb.CPacket{Bytes: bufOut.UAS(), PacketType: cpacket.Type().Value()})
}
