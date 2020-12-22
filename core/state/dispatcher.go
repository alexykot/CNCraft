package state

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/state/handlers"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/protocol"
)

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

func (d *SPacketDispatcher) Register() error {
	if err := d.ps.Subscribe(subj.MkNewConn(), d.HandleNewConnection); err != nil {
		return fmt.Errorf("failed to subscribe to new connections: %w", err)
	}
	d.log.Info("dispatcher handlers registered")
	return nil
}

func (d *SPacketDispatcher) parseSPacket(lope *envelope.E) (protocol.PacketType, protocol.SPacket, error) {
	spacketPb := lope.GetSpacket()
	if spacketPb == nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot parse SPacket: there is no SPacket in the envelope %v", lope)
	}

	bufI := buffer.NewFrom(spacketPb.Bytes)
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVrI())
	state, err := protocol.IntToState(int(spacketPb.State))
	if err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot parse SPacket connection state: %w", err)
	}

	pacType := protocol.MakeSType(state, protocolPacketID)
	bufI.Len()

	sPacket, err := d.pacFac.MakeSPacket(pacType)
	if err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot make SPacket: %w", err)
	}

	if err := sPacket.Pull(bufI); err != nil {
		return protocol.Unspecified, nil, fmt.Errorf("cannot pasrse buffer into SPacket `%d`: %w", pacType, err)
	}
	return pacType, sPacket, nil
}

func (d *SPacketDispatcher) HandleNewConnection(lope *envelope.E) {
	newConn := lope.GetNewConn()
	if newConn == nil {
		d.log.Error("cannot handle new connection - there is no NewConn in the envelope", zap.Any("envelope", lope))
		return
	}
	connId, err := uuid.Parse(newConn.Id)
	if err != nil {
		d.log.Error("cannot handle new connection - cannot parse connection ID", zap.Error(err))
		return
	}

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

	if err = d.ps.Subscribe(subj.MkConnReceive(connId), handler); err != nil {
		d.log.Error("cannot handle new connection: cannot subscribe to connection receive subject", zap.Error(err))
		return
	}
	d.log.Debug("handled new connection", zap.Any("ID", newConn.Id))
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
