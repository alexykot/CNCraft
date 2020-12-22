package handlers

import (
	"fmt"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol/status"
	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// HandleSPing handles the Ping packet.
func HandleSPing(ps nats.PubSub, pacFac protocol.PacketFactory, connID uuid.UUID, spacket protocol.SPacket) error {
	ping, ok := spacket.(*protocol.SPacketPing)
	if !ok {
		return fmt.Errorf("received packet is not a ping: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CPong)
	pong := cpacket.(*protocol.CPacketPong)
	pong.Payload = ping.Payload

	buff0 := buffer.New()
	pong.Push(buff0)

	lope := envelope.CPacket(&pb.CPacket{
		Bytes: buff0.UAS(),
	}, nil)
	if err := ps.Publish(subj.MkConnSend(connID), lope); err != nil {
		return fmt.Errorf("failed to publish CPong packet: %w", err)
	}

	return nil
}

// HandleSPing handles the Ping packet.
func HandleSRequest(ps nats.PubSub, pacFac protocol.PacketFactory, connID uuid.UUID, spacket protocol.SPacket) error {
	_, ok := spacket.(*protocol.SPacketRequest)
	if !ok {
		return fmt.Errorf("received packet is not a status request: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CResponse)
	statusResponse := cpacket.(*protocol.CPacketResponse)
	statusResponse.Status = status.DefaultResponse(578)

	buff0 := buffer.New()
	statusResponse.Push(buff0)

	lope := envelope.CPacket(&pb.CPacket{Bytes: buff0.UAS()}, nil)
	if err := ps.Publish(subj.MkConnSend(connID), lope); err != nil {
		return fmt.Errorf("failed to publish CPong packet: %w", err)
	}

	return nil
}

//func RegisterHandlersState1(ps nats.PubSub, logger *zap.Logger) {
//	{ // client bound packets
//		ps.Subscribe(protocol.MakePacketTopic(protocol.CResponse), func(envelopeIn envelope.E) {
//			connID, ok := envelopeIn.GetMetaKey(nats.MetaConnID)
//			if !ok {
//				// DEBT figure out logging here
//				return
//			}
//			packet, ok := envelopeIn.GetMessage().(protocol.CPacketResponse)
//			if !ok {
//				// DEBT figure out logging here
//				return
//			}
//			buf := buffer.New()
//			packet.Push(buf)
//			ps.Publish(network.MkConnSubjSend(connID), nats.NewEnvelope(buf, nil))
//		})
//
//		ps.Subscribe(protocol.MakePacketTopic(protocol.CPong), func(envelopeIn envelope.E) {
//			connID, ok := envelopeIn.GetMetaKey(nats.MetaConnID)
//			if !ok {
//				// DEBT figure out logging here
//				return
//			}
//			packet, ok := envelopeIn.GetMessage().(protocol.CPacketPong)
//			if !ok {
//				// DEBT figure out logging here
//				return
//			}
//			buf := buffer.New()
//			packet.Push(buf)
//			ps.Publish(network.MkConnSubjSend(connID), nats.NewEnvelope(buf, nil))
//		})
//	}
//
//	{ // server bound packets
//		ps.Subscribe(protocol.MakePacketTopic(protocol.SRequest), func(envelopeIn envelope.E) {
//			//packet, ok := envelopeIn.GetMessage().(protocol.SPacketRequest)
//			//if !ok {
//			//	// DEBT figure out logging here
//			//	return
//			//}
//
//			ps.Publish(protocol.MakePacketTopic(protocol.CResponse),
//				nats.NewEnvelope(protocol.CPacketResponse{Status: status.DefaultResponse(0)}, envelopeIn.GetMetaMap()))
//		})
//
//		ps.Subscribe(protocol.MakePacketTopic(protocol.SPing), func(envelopeIn envelope.E) {
//			packet, ok := envelopeIn.GetMessage().(protocol.SPacketPing)
//			if !ok {
//				// DEBT figure out logging here
//				return
//			}
//
//			ps.Publish(protocol.MakePacketTopic(protocol.CPong),
//				nats.NewEnvelope(protocol.CPacketPong{Payload: packet.Payload}, envelopeIn.GetMetaMap()))
//		})
//	}
//}
