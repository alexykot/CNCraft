package handlers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"

	"github.com/alexykot/cncraft/core/nats"
)

// HandleSHandshake handles the Handshake packet.
func HandleSHandshake(ps nats.PubSub, connID uuid.UUID, spacket protocol.SPacket) error {
	packet, ok := spacket.(*protocol.SPacketHandshake)
	if !ok {
		return fmt.Errorf("received packet is not a handshake: %v", spacket)
	}

	nextState := pb.ConnState_LOGIN
	if packet.NextState == protocol.Status {
		nextState = pb.ConnState_STATUS
	}

	lope := envelope.NewWithConnState(&pb.SetConnState{State: nextState}, nil)
	if err := ps.Publish(subj.MkConnStateChange(connID), lope); err != nil {
		return fmt.Errorf("failed to publish connstate change: %w", err)
	}

	return nil
}
