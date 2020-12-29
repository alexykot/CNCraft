package network

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/users"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/auth"
)

// DispatcherTransmitter parses and dispatches processing for incoming server bound protocol packets.
//  Also it collects and transmits outgoing client bound packets.
type DispatcherTransmitter struct {
	log   *zap.Logger
	ps    nats.PubSub
	auth  auth.A
	tally *users.Roster

	connMu map[uuid.UUID]*sync.Mutex
}

func NewDispatcher(log *zap.Logger, ps nats.PubSub, auth auth.A, tally *users.Roster) *DispatcherTransmitter {
	return &DispatcherTransmitter{
		log:   log,
		ps:    ps,
		auth:  auth,
		tally: tally,

		connMu: make(map[uuid.UUID]*sync.Mutex),
	}
}

//func (d *DispatcherTransmitter) Register() {
//	if err := d.ps.Publish(subj.MkNewConn(), envelope.NewConn(&pb.NewConnection{Id: conn.ID().String()})); err != nil {
//		d.log.Error("failed to publish conn.new message", zap.Error(err), zap.Any("conn", conn))
//		if err = conn.Close(); err != nil {
//			d.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
//		}
//		return
//	}
//}

func (d *DispatcherTransmitter) RegisterNewConn(conn Connection) error {
	d.connMu[conn.ID()] = &sync.Mutex{}

	transmitHandler := func(lope *envelope.E) {
		conn := conn
		log := d.log.With(zap.String("conn", conn.ID().String()))

		d.connMu[conn.ID()].Lock()
		defer d.connMu[conn.ID()].Unlock()

		pbCPacket := lope.GetCpacket()
		if pbCPacket == nil {
			log.Error("failed to parse envelope - there is no CPacket inside", zap.Any("envelope", lope))
			return
		}
		pacType := protocol.PacketType(pbCPacket.PacketType)
		log.Debug("transmitting CPacket", zap.String("type", pacType.String()))

		bufOut := buffer.New()
		bufOut.PushVrI(int32(pacType.ProtocolID()))

		packetBytes := bytes.Join([][]byte{
			bufOut.UAS(),
			pbCPacket.GetBytes(),
		}, nil)

		if err := d.transmitBytes(conn, packetBytes); err != nil {
			log.Error("failed to transmit CPacket", zap.Error(err))
			return
		}
	}

	if err := d.ps.Subscribe(subj.MkConnTransmit(conn.ID()), transmitHandler); err != nil {
		return fmt.Errorf("failed to subscribe to connTransmit: %w", err)
	}

	d.log.Debug("handled new connection", zap.Any("conn", conn.ID()))
	return nil
}

func (d *DispatcherTransmitter) HandleSPacket(conn Connection, packetBytes []byte) {
	log := d.log.With(zap.String("connId", conn.ID().String()))
	d.connMu[conn.ID()].Lock()
	defer d.connMu[conn.ID()].Unlock()

	sPacket, err := d.parseSPacket(conn.GetState(), packetBytes)
	if err != nil {
		log.Error("cannot handle new SPacket - could not parse bytes", zap.Error(err))
		return
	}
	log.Debug("handling SPacket", zap.String("type", sPacket.Type().String()), zap.Any("data", sPacket))

	if err = d.dispatchSPacket(conn, sPacket); err != nil {
		if errors.Is(err, handlers.InvalidLoginErr) {
			// Send CDisconnect packet here.
			log.Info("invalid login attempt, evicting user", zap.Error(err))
			d.auth.LoginFailure(conn.ID())
			if err := conn.Close(); err != nil {
				log.Error("error while closing connection", zap.Error(err))
			}
		} else {
			log.Error("cannot handle new packet - failed to dispatch handling", zap.Error(err))
		}
		return
	}
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

	sPacket, err := protocol.GetPacketFactory().MakeSPacket(pacType)
	if err != nil {
		return nil, fmt.Errorf("failed to make SPacket: %w", err)
	}

	if err := sPacket.Pull(bufI); err != nil {
		return nil, fmt.Errorf("failed to parse buffer into SPacket `%X`: %w", int32(pacType), err)
	}
	return sPacket, nil
}

// dispatchSPacket dispatches handling for the provided packet according to it's type.
func (d *DispatcherTransmitter) dispatchSPacket(conn Connection, sPacket protocol.SPacket) error {
	var err error
	var cPackets []protocol.CPacket

	debugStateSetter := func(state protocol.State) { // only needed to add the debug log line
		conn := conn
		conn.SetState(state)
		d.log.Debug("changed connState", zap.String("conn", conn.ID().String()), zap.String("state", state.String()))
	}

	pacType := sPacket.Type()
	switch pacType {
	case protocol.SHandshake:
		err = handlers.HandleSHandshake(debugStateSetter, sPacket)
	case protocol.SRequest:
		if cPackets, err = handlers.HandleSRequest(sPacket); err != nil {
			return fmt.Errorf("failed to handle SRequest packet: %w", err)
		}
	case protocol.SPing:
		cPackets, err = handlers.HandleSPing(sPacket)
	case protocol.SLoginStart:
		cPackets, err = handlers.HandleSLoginStart(d.auth, d.ps, debugStateSetter, conn.ID(), sPacket)
	case protocol.SEncryptionResponse:
		cPackets, err = handlers.HandleSEncryptionResponse(
			d.auth, d.ps, debugStateSetter, conn.EnableEncryption, conn.EnableCompression, conn.ID(), sPacket)
	case protocol.SPluginMessage:
		cPackets, err = handlers.HandleSPluginMessage(d.log, d.tally, conn.ID(), sPacket)
	case protocol.SClientSettings:
		cPackets, err = handlers.HandleSClientSettings(d.tally, conn.ID(), sPacket)
	default:
		return fmt.Errorf("unhandled packet type: %X", int32(pacType))
	}

	if err != nil {
		return fmt.Errorf("failed to handle %s packet: %w", pacType.String(), err)
	}
	if cPackets != nil {
		for _, cPacket := range cPackets {
			if err := d.transmitCPacket(conn, cPacket); err != nil {
				return fmt.Errorf("failed to transmit %s packet: %w", cPacket.Type().String(), err)
			}
		}
	}

	return nil
}

func (d *DispatcherTransmitter) transmitBuffer(conn Connection, bufOut buffer.B) error {
	if bufOut.Len() < 2 {
		return fmt.Errorf("buffer data is too short")
	}

	count, err := conn.Transmit(bufOut)
	if err != nil {
		return fmt.Errorf("failed to push client bound data: %w", err)
	}

	d.log.Debug("transmitted bytes", zap.String("conn", conn.ID().String()), zap.Int("count", count))
	return nil
}

func (d *DispatcherTransmitter) transmitBytes(conn Connection, packetBytes []byte) error {
	d.log.Debug("pushing bytes to conn", zap.String("conn", conn.ID().String()),
		zap.String("bytes", hex.EncodeToString(packetBytes)))

	if err := d.transmitBuffer(conn, buffer.NewFrom(packetBytes)); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	return nil
}

func (d *DispatcherTransmitter) transmitCPacket(conn Connection, cpacket protocol.CPacket) error {
	bufOut := buffer.New()
	bufOut.PushVrI(int32(cpacket.ProtocolID()))
	cpacket.Push(bufOut)

	d.log.Debug("pushing packet to conn", zap.String("conn", conn.ID().String()),
		zap.String("type", cpacket.Type().String()), zap.String("bytes", hex.EncodeToString(bufOut.UAS())))

	if err := d.transmitBuffer(conn, bufOut); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	return nil
}

// checkIsStatusHandshake checks if the packet looks like a Handshake packet. This is needed because in Status
// connection state there is no way in the protocol to correctly signal upgrade to Login state, so the Notchian client
// sends a SHandshake packet, which belongs to Handshake state and it's packetID collides with the SRequest packet
// from the Status state. So we have to hack around this by checking the packet size, and if the connState is Status,
// packetID is 0x00  and size is bigger than 1 byte - assume it is an SHandshake, not SRequest as it normally would be.
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
