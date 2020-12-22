package handlers

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/status"
)

// HandleSPing handles the Ping packet.
func HandleSPing(transmitter func(protocol.CPacket), pacFac protocol.PacketFactory, spacket protocol.SPacket) error {
	ping, ok := spacket.(*protocol.SPacketPing)
	if !ok {
		return fmt.Errorf("received packet is not a ping: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CPong)       // Predefined packet is expected to always exist.
	cpacket.(*protocol.CPacketPong).Payload = ping.Payload // And always be of the correct type.
	transmitter(cpacket)
	return nil
}

// HandleSPing handles the Ping packet.
func HandleSRequest(transmitter func(protocol.CPacket), pacFac protocol.PacketFactory, spacket protocol.SPacket) error {
	_, ok := spacket.(*protocol.SPacketRequest)
	if !ok {
		return fmt.Errorf("received packet is not a status request: %v", spacket)
	}

	cpacket, _ := pacFac.MakeCPacket(protocol.CResponse)                     // Predefined packet is expected to always exist.
	cpacket.(*protocol.CPacketResponse).Status = status.DefaultResponse(578) // And always be of the correct type.
	transmitter(cpacket)
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
