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

func PlayerLoading(loadingPlayer *pb.PlayerLoading) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_PlayerLoading{PlayerLoading: loadingPlayer},
		},
	}
}

func PlayerJoined(joinedPlayer *pb.PlayerJoined) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_PlayerJoined{PlayerJoined: joinedPlayer},
		},
	}
}

func PlayerLeft(leftPlayer *pb.PlayerLeft) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_PlayerLeft{PlayerLeft: leftPlayer},
		},
	}
}

func CloseConn(closeConn *pb.CloseConn) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_CloseConn{CloseConn: closeConn},
		},
	}
}

func PlayerSpatialUpdate(spatialUpdate *pb.PlayerSpatialUpdate) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_PlayerSpatial{PlayerSpatial: spatialUpdate},
		},
	}
}

func PlayerInventoryUpdate(inventoryUpdate *pb.PlayerInventoryUpdate) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_PlayerInventory{PlayerInventory: inventoryUpdate},
		},
	}
}

func NewPlayerJoined(newPlayer *pb.NewPlayerJoined) *E {
	return &E{
		Envelope: pb.Envelope{
			Message: &pb.Envelope_NewPlayer{NewPlayer: newPlayer},
		},
	}
}

func PlayerDigging(digging *pb.PlayerDigging) *E {
	return &E{
		Envelope: pb.Envelope{
			ShardEvent: &pb.ShardEvent{
				Event: &pb.ShardEvent_PlayerDigging{PlayerDigging: digging},
			},
		},
	}
}
