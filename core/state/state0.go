package state

import (
	"github.com/alexykot/cncraft/impl/conn"
	"github.com/alexykot/cncraft/impl/protocol"
	"github.com/alexykot/cncraft/impl/protocol/server"
	"github.com/alexykot/cncraft/pkg/bus"
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

