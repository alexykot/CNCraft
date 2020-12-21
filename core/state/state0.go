package state

import (
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
)

// RegisterHandlersState1 registers handlers for packets transmitted/received in the Handshake connection state.
func RegisterHandlersState0(ps nats.PubSub, logger *zap.Logger) {
	//err := ps.Subscribe(protocol.MakePacketTopic(protocol.SHandshake), func(envelopeIn *envelope.E) {
	//	connID, ok := envelopeIn.GetMetaKey(nats.MetaConnID)
	//	if !ok {
	//		// DEBT figure out logging here
	//		return
	//	}
	//	packet, ok := envelopeIn.GetMessage().(protocol.SPacketHandshake)
	//	if !ok {
	//		// DEBT figure out logging here
	//		return
	//	}
	//	ps.Publish(network.MkConnSubjStateChange(connID), nats.NewEnvelope(packet.NextState, nil))
	//})
}
