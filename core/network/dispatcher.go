package network

import (
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/buffer"
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

func (d *SPacketDispatcher) HandleSPacket(conn Connection, packetBytes []byte) {
	log := d.log.With(zap.String("connId", conn.ID().String()))
	pacType, sPacket, err := d.parseSPacket(conn.GetState(), packetBytes)

	if err != nil {
		log.Error("cannot handle new packet - could not parse spacket", zap.Error(err))
		return
	}
	if err = d.receiveSPacket(conn, pacType, sPacket); err != nil {
		log.Error("cannot handle new packet - failed to dispatch handling", zap.Error(err))
		return
	}
	log.Debug("handled incoming packet", zap.String("type", pacType.String()), zap.Any("data", sPacket))
}

func (d *SPacketDispatcher) parseSPacket(connState protocol.State, packetBytes []byte) (protocol.PacketType, protocol.SPacket, error) {
	println("spacket bytes:", hex.EncodeToString(packetBytes))

	bufI := buffer.NewFrom(packetBytes)

	id := bufI.PullVrI()
	println(fmt.Sprintf("spacket id: %X", id))
	protocolPacketID := protocol.ProtocolPacketID(id)

	pacType := protocol.MakeSType(connState, protocolPacketID)
	println(fmt.Sprintf("spacket type: %X", pacType.Value()))

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
func (d *SPacketDispatcher) receiveSPacket(conn Connection, pacType protocol.PacketType, spacket protocol.SPacket) error {
	transmitter := func(cpacket protocol.CPacket) {
		conn := conn
		d.transmitCPacket(conn, cpacket)
	}

	switch pacType {
	case protocol.SHandshake:
		stateSetter := func(state protocol.State) {
			conn := conn
			conn.SetState(state)
			d.log.Debug("changed connstate", zap.String("conn", conn.ID().String()), zap.String("state", state.String()))
		}

		if err := handlers.HandleSHandshake(stateSetter, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	case protocol.SPing:
		if err := handlers.HandleSPing(transmitter, d.pacFac, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	case protocol.SRequest:
		if err := handlers.HandleSRequest(transmitter, d.pacFac, spacket); err != nil {
			return fmt.Errorf("failed to handle handshake packet: %w", err)
		}
	default:
		return fmt.Errorf("unhandled packet type: %d", pacType)
	}

	return nil
}

func (d *SPacketDispatcher) transmitCPacket(conn Connection, cpacket protocol.CPacket) {
	println(fmt.Sprintf("cpacket type: %X", cpacket.Type().Value()))

	bufO := buffer.New()
	bufO.PushVrI(int32(cpacket.ProtocolID()))
	cpacket.Push(bufO)
	if bufO.Len() < 2 {
		d.log.Error("received CPacket with zero length buffer, cannot send")
	}

	println("cpacket bytes:", hex.EncodeToString(bufO.UAS()))

	temp := buffer.New()
	temp.PushVrI(bufO.Len())

	deflated := buffer.New()
	deflated.PushUAS(conn.Deflate(bufO.UAS()), false)
	temp.PushUAS(deflated.UAS(), false)

	if _, err := conn.Push(conn.Encrypt(temp.UAS())); err != nil {
		d.log.Error("Failed to push client bound packet", zap.Error(err))
		return
	}
	d.log.Debug("pushed packet to conn", zap.String("conn", conn.ID().String()),
		zap.String("type", cpacket.Type().String()))
}
