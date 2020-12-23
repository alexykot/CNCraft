package network

import (
	"fmt"

	auth2 "github.com/alexykot/cncraft/core/network/auth"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// DispatcherTransmitter parses and dispatches processing for incoming server bound protocol packets.
//  Also it collects and transmits outgoing client bound packets.
type DispatcherTransmitter struct {
	log  *zap.Logger
	ps   nats.PubSub
	auth *auth2.Auther

	pacFac protocol.PacketFactory
}

func NewDispatcher(log *zap.Logger, ps nats.PubSub, pacfac protocol.PacketFactory, auth *auth2.Auther) *DispatcherTransmitter {
	return &DispatcherTransmitter{
		log:    log,
		ps:     ps,
		auth:   auth,
		pacFac: pacfac,
	}
}

//func (d *DispatcherTransmitter) Register() {
//	if err := d.ps.Publish(subj.MkNewConn(), envelope.NewConn(&pb.NewConnection{Id: conn.ID().String()}, nil)); err != nil {
//		d.log.Error("failed to publish conn.new message", zap.Error(err), zap.Any("conn", conn))
//		if err = conn.Close(); err != nil {
//			d.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
//		}
//		return
//	}
//}

func (d *DispatcherTransmitter) RegisterConnHandlers(conn Connection) error {
	transmitHandler := func(lope *envelope.E) {
		conn := conn
		log := d.log.With(zap.String("conn", conn.ID().String()))

		cpacket := lope.GetCpacket()
		if cpacket == nil {
			log.Error("failed to parse envelope - there is no CPacket inside", zap.Any("envelope", lope))
			return
		}

		if err := d.transmitBytes(conn, cpacket.GetBytes()); err != nil {
			log.Error("failed to transmit CPacket", zap.Error(err))
			return
		}
		log.Debug("transmitted CPacket", zap.String("type", fmt.Sprintf("%X", cpacket.PacketType)))
	}

	if err := d.ps.Subscribe(subj.MkConnTransmit(conn.ID()), transmitHandler); err != nil {
		return fmt.Errorf("failed to subscribe to connTransmit: %w", err)
	}

	d.log.Debug("handled new connection", zap.Any("conn", conn.ID()))
	return nil
}

func (d *DispatcherTransmitter) HandleSPacket(conn Connection, packetBytes []byte) {
	log := d.log.With(zap.String("connId", conn.ID().String()))
	sPacket, err := d.parseSPacket(conn.GetState(), packetBytes)

	if err != nil {
		log.Error("cannot handle new SPacket - could not parse bytes", zap.Error(err))
		return
	}
	if err = d.dispatchSPacket(conn, sPacket); err != nil {
		log.Error("cannot handle new packet - failed to dispatch handling", zap.Error(err))
		return
	}
	log.Debug("handled incoming packet", zap.String("type", sPacket.Type().String()), zap.Any("data", sPacket))
}

func (d *DispatcherTransmitter) parseSPacket(connState protocol.State, packetBytes []byte) (protocol.SPacket, error) {
	bufI := buffer.NewFrom(packetBytes)
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVrI())

	var pacType protocol.PacketType
	if d.checkIsStatusHandshake(connState, packetBytes) {
		// hack for Status->Login state upgrade, see checkIsStatusHandshake for details
		pacType = protocol.MakeSType(protocol.Handshake, protocolPacketID)
	} else {
		pacType = protocol.MakeSType(connState, protocolPacketID)
	}

	sPacket, err := d.pacFac.MakeSPacket(pacType)
	if err != nil {
		return nil, fmt.Errorf("failed to make SPacket: %w", err)
	}

	if err := sPacket.Pull(bufI); err != nil {
		return nil, fmt.Errorf("failed to parse buffer into SPacket `%X`: %w", int32(pacType), err)
	}
	return sPacket, nil
}

// dispatchSPacket dispatches handling for the provided packet according to it's type.
func (d *DispatcherTransmitter) dispatchSPacket(conn Connection, spacket protocol.SPacket) error {
	var err error
	var cpacket protocol.CPacket

	pacType := spacket.Type()
	switch pacType {
	case protocol.SHandshake:
		stateSetter := func(state protocol.State) {
			conn := conn
			conn.SetState(state)
			d.log.Debug("changed connState", zap.String("conn", conn.ID().String()), zap.String("state", state.String()))
		}

		err = handlers.HandleSHandshake(stateSetter, spacket)
	case protocol.SRequest:
		if cpacket, err = handlers.HandleSRequest(d.pacFac, spacket); err != nil {
			return fmt.Errorf("failed to handle SRequest packet: %w", err)
		}
	case protocol.SPing:
		cpacket, err = handlers.HandleSPing(d.pacFac, spacket)
	case protocol.SLoginStart:
		cpacket, err = handlers.HandleSLoginStart(d.auth.AddNewbie, d.pacFac, conn.ID(), spacket)
	default:
		return fmt.Errorf("unhandled packet type: %X", int32(pacType))
	}

	if err != nil {
		return fmt.Errorf("failed to handle %s packet: %w", pacType.String(), err)
	}
	if cpacket != nil {
		if err := d.transmitCPacket(conn, cpacket); err != nil {
			return fmt.Errorf("failed to transmit %s packet: %w", cpacket.Type().String(), err)
		}
	}

	return nil
}

func (d *DispatcherTransmitter) transmitBuffer(conn Connection, bufOut buffer.B) error {
	if bufOut.Len() < 2 {
		return fmt.Errorf("buffer data is too short")
	}

	temp := buffer.New()
	temp.PushVrI(bufOut.Len())

	deflated := buffer.New()
	deflated.PushUAS(conn.Deflate(bufOut.UAS()), false)
	temp.PushUAS(deflated.UAS(), false)

	if _, err := conn.Push(conn.Encrypt(temp.UAS())); err != nil {
		return fmt.Errorf("failed to push client bound data: %w", err)
	}
	return nil
}

func (d *DispatcherTransmitter) transmitBytes(conn Connection, packetBytes []byte) error {
	if err := d.transmitBuffer(conn, buffer.NewFrom(packetBytes)); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	d.log.Debug("pushed bytes to conn", zap.String("conn", conn.ID().String()), zap.Int("len", len(packetBytes)))
	return nil
}

func (d *DispatcherTransmitter) transmitCPacket(conn Connection, cpacket protocol.CPacket) error {
	bufOut := buffer.New()
	bufOut.PushVrI(int32(cpacket.ProtocolID()))
	cpacket.Push(bufOut)

	if err := d.transmitBuffer(conn, bufOut); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	d.log.Debug("pushed packet to conn", zap.String("conn", conn.ID().String()),
		zap.String("type", cpacket.Type().String()))
	return nil
}

// checkIsStatusHandshake checks if the packet looks like a Handshake packet. This is needed because in Status
// connection mode there is no way in the protocol to correctly signal upgrade to login mode, so the Notchian client
// sends a SHandshake packet, which belongs to Handshake state and it's packetID collides with the SRequest packet.
// So we have to hack around this by checking the packet size, if the connState is Handshake and packetID is 0x00.
func (d *DispatcherTransmitter) checkIsStatusHandshake(connState protocol.State, packetBytes []byte) bool {
	if connState != protocol.Status { // if the connState is not Status - this hack does not apply.
		return false
	} else if len(packetBytes) < 6 { // 6 bytes is absolute minimum possible length of the SHandshake packet.
		return false
	}

	protocolID := buffer.NewFrom(packetBytes).PullVrI()
	if protocolID != 0x00 { // SHandshake packet protocol ID is 0x00, same as the ID of SRequest packet.
		return false
	} else if len(packetBytes) == 1 { // SRequest packet has no fields so it's length is always 1 byte.
		return false
	}

	return true
}
