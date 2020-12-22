package handlers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// HandleSHandshake handles the Handshake packet.
func HandleSHandshake(ps nats.PubSub, connID uuid.UUID, spacket protocol.SPacket) error {
	packet, ok := spacket.(*protocol.SPacketHandshake)
	if !ok {
		return fmt.Errorf("received packet is not a handshake: %v", spacket)
	}

	var nextState pb.ConnState
	switch packet.NextState {
	case protocol.Shake:
		nextState = pb.ConnState_HANDSHAKE
	case protocol.Status:
		nextState = pb.ConnState_STATUS
	case protocol.Login:
		nextState = pb.ConnState_LOGIN
	default:
		return fmt.Errorf("unexpected next state received: %d", packet.NextState)
	}

	lope := envelope.ConnState(&pb.SetConnState{State: nextState}, nil)
	if err := ps.Publish(subj.MkConnStateChange(connID), lope); err != nil {
		return fmt.Errorf("failed to publish connstate change: %w", err)
	}

	return nil
}
