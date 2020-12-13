package state

import (
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/buffers"
	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/status"
)

// RegisterHandlersState1 registers handlers for packets transmitted/received in the Status connection state.
func RegisterHandlersState1(ps bus.PubSub, logger *zap.Logger) {
	{ // client bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.CResponse), func(envelopeIn bus.Envelope) {
			connID, ok := envelopeIn.GetMeta(bus.MetaConn)
			if !ok {
				// DEBT figure out logging here
				return
			}
			packet, ok := envelopeIn.GetMessage().(protocol.CPacketResponse)
			if !ok {
				// DEBT figure out logging here
				return
			}
			buf := buffers.NewBuffer()
			packet.Push(buf)
			ps.Publish(network.MakeConnTopicSend(connID), bus.NewEnvelope(buf, nil))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.CPong), func(envelopeIn bus.Envelope) {
			connID, ok := envelopeIn.GetMeta(bus.MetaConn)
			if !ok {
				// DEBT figure out logging here
				return
			}
			packet, ok := envelopeIn.GetMessage().(protocol.CPacketPong)
			if !ok {
				// DEBT figure out logging here
				return
			}
			buf := buffers.NewBuffer()
			packet.Push(buf)
			ps.Publish(network.MakeConnTopicSend(connID), bus.NewEnvelope(buf, nil))
		})
	}

	{ // server bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.SRequest), func(envelopeIn bus.Envelope) {
			//packet, ok := envelopeIn.GetMessage().(protocol.SPacketRequest)
			//if !ok {
			//	// DEBT figure out logging here
			//	return
			//}

			ps.Publish(protocol.MakePacketTopic(protocol.CResponse),
				bus.NewEnvelope(protocol.CPacketResponse{Status: status.DefaultResponse()}, envelopeIn.GetAllMeta()))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.SPing), func(envelopeIn bus.Envelope) {
			packet, ok := envelopeIn.GetMessage().(protocol.SPacketPing)
			if !ok {
				// DEBT figure out logging here
				return
			}

			ps.Publish(protocol.MakePacketTopic(protocol.CPong),
				bus.NewEnvelope(protocol.CPacketPong{Ping: packet.Ping}, envelopeIn.GetAllMeta()))
		})
	}
}
