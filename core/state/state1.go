package state

import (
	"github.com/alexykot/cncraft/pkg/envelope"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/buffers"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/status"
)

// RegisterHandlersState1 registers handlers for packets transmitted/received in the Status connection state.
func RegisterHandlersState1(ps nats.PubSub, logger *zap.Logger) {
	{ // client bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.CResponse), func(envelopeIn envelope.E) {
			connID, ok := envelopeIn.GetMetaKey(nats.MetaConn)
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
			ps.Publish(network.MakeConnTopicSend(connID), nats.NewEnvelope(buf, nil))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.CPong), func(envelopeIn envelope.E) {
			connID, ok := envelopeIn.GetMetaKey(nats.MetaConn)
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
			ps.Publish(network.MakeConnTopicSend(connID), nats.NewEnvelope(buf, nil))
		})
	}

	{ // server bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.SRequest), func(envelopeIn envelope.E) {
			//packet, ok := envelopeIn.GetMessage().(protocol.SPacketRequest)
			//if !ok {
			//	// DEBT figure out logging here
			//	return
			//}

			ps.Publish(protocol.MakePacketTopic(protocol.CResponse),
				nats.NewEnvelope(protocol.CPacketResponse{Status: status.DefaultResponse(0)}, envelopeIn.GetMetaMap()))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.SPing), func(envelopeIn envelope.E) {
			packet, ok := envelopeIn.GetMessage().(protocol.SPacketPing)
			if !ok {
				// DEBT figure out logging here
				return
			}

			ps.Publish(protocol.MakePacketTopic(protocol.CPong),
				nats.NewEnvelope(protocol.CPacketPong{Ping: packet.Ping}, envelopeIn.GetMetaMap()))
		})
	}
}
