package state

import (
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/protocol"
	"github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/bus"
)

func RegisterHandlersState0(ps bus.PubSub) {
	ps.Subscribe(protocol.MakePacketTopic(protocol.SHandshake), func(envelopeIn bus.Envelope) {
		connID, ok := envelopeIn.GetMeta(bus.MetaConn)
		if !ok {
			// DEBT figure out logging here
			return
		}
		packet, ok := envelopeIn.GetMessage().(server.SPacketHandshake)
		if !ok {
			// DEBT figure out logging here
			return
		}
		ps.Publish(conn.MakeConnTopicState(connID), bus.NewEnvelope(packet.NextState, nil))
	})
}

