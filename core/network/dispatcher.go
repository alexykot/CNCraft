package network

import (
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// SPacketDispatcher parses and dispatches processing for incoming server bound protocol packets.
type SPacketDispatcher struct {
	log    *zap.Logger
	pacFac protocol.PacketFactory
	ps     nats.PubSub
}

func NewDispatcher(log *zap.Logger, pacfac protocol.PacketFactory, ps nats.PubSub) *SPacketDispatcher {
	return &SPacketDispatcher{
		log:    log,
		pacFac: pacfac,
		ps:     ps,
	}
}

func (d *SPacketDispatcher) RegisterNewConnection(connId uuid.UUID) {
	handler := func(lope *envelope.E) {
		connId := connId
		log := d.log.With(zap.String("connId", connId.String()))
		pacType, sPacket, err := d.parseSPacket(lope)

		if err != nil {
			log.Error("cannot handle new packet - could not parse spacket", zap.Error(err))
			return
		}
		if err = d.dispatchPacketHandling(connId, pacType, sPacket); err != nil {
			log.Error("cannot handle new packet - failed to dispatch handling", zap.Error(err))
			return
		}
		log.Debug("handled incoming packet", zap.String("type", pacType.String()), zap.Any("data", sPacket))
	}

	if err := d.ps.Subscribe(subj.MkConnReceive(connId), handler); err != nil {
		d.log.Error("cannot handle new connection: cannot subscribe to connection receive subject", zap.Error(err))
		return
	}
	d.log.Debug("handled new connection", zap.Any("ID", connId))
}

func (d *SPacketDispatcher) parseSPacket(lope *envelope.E) (protocol.PacketType, protocol.SPacket, error) {

	println("parsing packet")

	spacketPb := lope.GetSpacket()
	if spacketPb == nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot parse SPacket: there is no SPacket in the envelope %v", lope)
	}
	println("packet bytes:", hex.EncodeToString(spacketPb.Bytes))

	bufI := buffer.NewFrom(spacketPb.Bytes)

	id := bufI.PullVrI()
	println(fmt.Sprintf("packet id: %X", id))
	protocolPacketID := protocol.ProtocolPacketID(id)
	state, err := protocol.IntToState(int(spacketPb.State))
	println(fmt.Sprintf("conn state: %d", state))
	if err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot parse SPacket connection state: %w", err)
	}

	pacType := protocol.MakeSType(state, protocolPacketID)
	println(fmt.Sprintf("packet type: %v", pacType))

	sPacket, err := d.pacFac.MakeSPacket(pacType)
	if err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot make SPacket: %w", err)
	}

	if err := sPacket.Pull(bufI); err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot pasrse buffer into SPacket `%d`: %w", pacType, err)
	}
	return pacType, sPacket, nil
}

// DispatchStatePacketHandling parses incoming server bound packet envelopes and dispatches packet handlers.
func (d *SPacketDispatcher) dispatchPacketHandling(connID uuid.UUID, pactype protocol.PacketType, spacket protocol.SPacket) error {
	switch pactype {
	case protocol.SHandshake:
		if err := handlers.HandleSHandshake(d.ps, connID, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	case protocol.SPing:
		if err := handlers.HandleSPing(d.ps, d.pacFac, connID, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	case protocol.SRequest:
		if err := handlers.HandleSHandshake(d.ps, connID, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	default:
		return fmt.Errorf("unhandled packet type: %d", pactype)
	}

	return nil
}
