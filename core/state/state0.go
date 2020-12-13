package state

import (
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// RegisterHandlersState1 registers handlers for packets transmitted/received in the Handshake connection state.
func RegisterHandlersState0(ps bus.PubSub, logger *zap.Logger) {
	ps.Subscribe(protocol.MakePacketTopic(protocol.SHandshake), func(envelopeIn bus.Envelope) {
		connID, ok := envelopeIn.GetMeta(bus.MetaConn)
		if !ok {
			// DEBT figure out logging here
			return
		}
		packet, ok := envelopeIn.GetMessage().(protocol.SPacketHandshake)
		if !ok {
			// DEBT figure out logging here
			return
		}
		ps.Publish(network.MakeConnTopicState(connID), bus.NewEnvelope(packet.NextState, nil))
	})
}

