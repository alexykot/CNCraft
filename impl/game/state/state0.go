package state

import (
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/protocol"
	"github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

func HandleState0(ps pubsub.PubSub) {
	ps.Subscribe(protocol.MakePacketTopic(protocol.SHandshake), func(envelopeIn pubsub.Envelope) {
		connID, ok := envelopeIn.GetMeta(pubsub.MetaConn)
		if !ok {
			// DEBT figure out logging here
			return
		}
		packet, ok := envelopeIn.GetMessage().(server.SPacketHandshake)
		if !ok {
			// DEBT figure out logging here
			return
		}
		ps.Publish(conn.MakeConnTopicState(connID), pubsub.NewEnvelope(packet.NextState, nil))
	})
}

