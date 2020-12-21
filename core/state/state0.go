package state

import (
	"github.com/alexykot/cncraft/pkg/envelope"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// RegisterHandlersState1 registers handlers for packets transmitted/received in the Handshake connection state.
func RegisterHandlersState0(ps nats.PubSub, logger *zap.Logger) {
	ps.Subscribe(protocol.MakePacketTopic(protocol.SHandshake), func(envelopeIn envelope.E) {
		connID, ok := envelopeIn.GetMetaKey(nats.MetaConn)
		if !ok {
			// DEBT figure out logging here
			return
		}
		packet, ok := envelopeIn.GetMessage().(protocol.SPacketHandshake)
		if !ok {
			// DEBT figure out logging here
			return
		}
		ps.Publish(network.MakeConnTopicState(connID), nats.NewEnvelope(packet.NextState, nil))
	})
}

