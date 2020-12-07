package state

import (
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/data/status"
	"github.com/golangmc/minecraft-server/impl/protocol"
	"github.com/golangmc/minecraft-server/impl/protocol/client"
	"github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/bus"
)

/**
 * status
 */

func RegisterHandlersState1(ps bus.PubSub) {
	{ // client bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.CResponse), func(envelopeIn bus.Envelope) {
			connID, ok := envelopeIn.GetMeta(bus.MetaConn)
			if !ok {
				// DEBT figure out logging here
				return
			}
			packet, ok := envelopeIn.GetMessage().(client.CPacketResponse)
			if !ok {
				// DEBT figure out logging here
				return
			}
			buf := conn.NewBuffer()
			packet.Push(buf, nil)
			ps.Publish(conn.MakeConnTopicSend(connID), bus.NewEnvelope(buf, nil))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.CPong), func(envelopeIn bus.Envelope) {
			connID, ok := envelopeIn.GetMeta(bus.MetaConn)
			if !ok {
				// DEBT figure out logging here
				return
			}
			packet, ok := envelopeIn.GetMessage().(client.CPacketPong)
			if !ok {
				// DEBT figure out logging here
				return
			}
			buf := conn.NewBuffer()
			packet.Push(buf, nil)
			ps.Publish(conn.MakeConnTopicSend(connID), bus.NewEnvelope(buf, nil))
		})
	}

	{ // server bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.SRequest), func(envelopeIn bus.Envelope) {
			//packet, ok := envelopeIn.GetMessage().(server.SPacketRequest)
			//if !ok {
			//	// DEBT figure out logging here
			//	return
			//}

			ps.Publish(protocol.MakePacketTopic(protocol.CResponse),
				bus.NewEnvelope(client.CPacketResponse{Status: status.DefaultResponse()}, envelopeIn.GetAllMeta()))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.SPing), func(envelopeIn bus.Envelope) {
			packet, ok := envelopeIn.GetMessage().(server.SPacketPing)
			if !ok {
				// DEBT figure out logging here
				return
			}

			ps.Publish(protocol.MakePacketTopic(protocol.CPong),
				bus.NewEnvelope(client.CPacketPong{Ping: packet.Ping}, envelopeIn.GetAllMeta()))
		})
	}
}
