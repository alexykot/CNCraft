package envelope

import (
	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

func Empty() *E {
	return &E{}
}

func SPacket(spacket *pb.SPacket) *E {
	return &E{
		Envelope: pb.Envelope{
			Meta:    nil, // not using meta at the moment, reimplement when/if will start using it
			Message: &pb.Envelope_Spacket{Spacket: spacket},
		},
	}
}

func CPacket(cpacket *pb.CPacket) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_Cpacket{Cpacket: cpacket},
		},
	}
}

func JoinedPlayer(joinedPlayer *pb.JoinedPlayer) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_JoinedPlayer{JoinedPlayer: joinedPlayer},
		},
	}
}
