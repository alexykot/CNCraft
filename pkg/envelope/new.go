package envelope

import (
	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

func Empty() *E {
	return &E{}
}

func NewConn(conn *pb.NewConnection, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_NewConn{NewConn: conn},
		},
	}
}

func ConnState(connState *pb.SetConnState, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_ConnState{ConnState: connState},
		},
	}
}

func SPacket(spacket *pb.SPacket, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_Spacket{Spacket: spacket},
		},
	}
}

func CPacket(cpacket *pb.CPacket, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_Cpacket{Cpacket: cpacket},
		},
	}
}
