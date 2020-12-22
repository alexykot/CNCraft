package envelope

import (
	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

func NewEmpty() *E {
	return &E{}
}

func NewWithHandshake(hs *pb.Handshake, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_Handshake{Handshake: hs},
		},
	}
}

func NewWithConnState(connState *pb.SetConnState, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_ConnState{ConnState: connState},
		},
	}
}

func NewWithSPacket(spacket *pb.SPacket, meta map[string]string) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    meta,
			Message: &pb.Envelope_Spacket{Spacket: spacket},
		},
	}
}
