package network

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/world"

	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/auth"
)

// Dispatcher is a dispatcher-transmitter interface. Implementation of this interface is expected to handle the incoming
// packets from the network for all given connections, and transmit the outgoing packets to the given connections.
type Dispatcher interface {
	Init(ctx context.Context) error
	RegisterNewConn(conn Connection) error
	HandleSPacket(conn Connection, packetBytes []byte)
}

// dispatcherTransmitter parses and dispatches processing for incoming server bound protocol packets.
//  Also it collects and transmits outgoing client bound packets and handles disconnections.
type dispatcherTransmitter struct {
	log     *zap.Logger
	ps      nats.PubSub
	auth    auth.A
	roster  players.Roster
	aliver  *KeepAliver
	sharder *world.Sharder

	// map of mutexes intended to control access to individual connections. Each connection needs to be thread-safe,
	// but unrelated connections may be processed in parallel.
	connMu map[uuid.UUID]*sync.Mutex

	// mutex protecting editing the connMu map itself.
	connMapMu sync.Mutex
}

func NewDispatcher(log *zap.Logger, ps nats.PubSub, auth auth.A, roster players.Roster, aliver *KeepAliver, sharder *world.Sharder) Dispatcher {
	return &dispatcherTransmitter{
		log:     log,
		ps:      ps,
		auth:    auth,
		roster:  roster,
		aliver:  aliver,
		sharder: sharder,

		connMu: make(map[uuid.UUID]*sync.Mutex),
	}
}

func (d *dispatcherTransmitter) Init(ctx context.Context) error {
	if err := d.aliver.Start(ctx); err != nil {
		return fmt.Errorf("failed to start the keepaliver: %w", err)
	}

	if err := d.ps.Subscribe(subj.MkConnClosed(), d.connClosedHandler); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subj.MkConnClosed().String(), err)
	}
	d.log.Debug("registered conn closing handler")

	if err := d.ps.Subscribe(subj.MkConnBroadcast(), d.packetBroadcastHandler); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subj.MkConnBroadcast().String(), err)
	}
	d.log.Debug("registered conn broadcast handler")

	d.log.Info("dispatcher started")
	return nil
}

func (d *dispatcherTransmitter) RegisterNewConn(conn Connection) error {
	d.connMu[conn.ID()] = &sync.Mutex{}

	if err := d.ps.Subscribe(subj.MkConnTransmit(conn.ID()), d.getTransmitHandler(conn)); err != nil {
		return fmt.Errorf("failed to subscribe to connTransmit: %w", err)
	}

	d.log.Info("new connection opened", zap.Any("conn", conn.ID()))
	return nil
}

func (d *dispatcherTransmitter) HandleSPacket(conn Connection, packetBytes []byte) {
	log := d.log.With(zap.String("conn", conn.ID().String()))
	d.connMu[conn.ID()].Lock()
	defer d.connMu[conn.ID()].Unlock()

	sPacket, err := d.parseSPacket(conn.GetState(), packetBytes)
	if err != nil {
		log.Error("cannot handle new SPacket: could not parse bytes", zap.Error(err))
		return
	}
	log.Debug("handling SPacket", zap.String("type", sPacket.Type().String()))

	if err = d.dispatchSPacket(conn, sPacket); err != nil {
		if errors.Is(err, handlers.InvalidLoginErr) {
			log.Info("invalid login attempt, evicting user", zap.Error(err))
			d.auth.LoginFailure(conn.ID())
			if err := d.forceDisconnect(conn.GetState(), conn.ID()); err != nil {
				log.Error("failed to trigger disconnect", zap.Error(err))
			}
		} else {
			log.Error("cannot handle new packet: failed to dispatch handling", zap.Error(err))
		}
		return
	}
}

func (d *dispatcherTransmitter) parseSPacket(connState protocol.State, packetBytes []byte) (protocol.SPacket, error) {
	bufI := buffer.NewFrom(packetBytes)
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVarInt())

	var pacType protocol.PacketType
	if d.checkIsStatusHandshake(connState, packetBytes) {
		// hack for Status->Login state upgrade, see checkIsStatusHandshake for details
		pacType = protocol.MakeSType(protocol.Handshake, protocolPacketID)
	} else {
		pacType = protocol.MakeSType(connState, protocolPacketID)
	}

	sPacket, err := protocol.GetPacketFactory().MakeSPacket(pacType)
	if err != nil {
		return nil, fmt.Errorf("failed to make SPacket from pacType %d/%X, %s: %w", connState, protocolPacketID, pacType.String(), err)
	}

	if err := sPacket.Pull(bufI); err != nil {
		return nil, fmt.Errorf("failed to parse buffer into SPacket, pacType %d/%X, %s: %w", connState, protocolPacketID, pacType.String(), err)
	}
	return sPacket, nil
}

// dispatchSPacket dispatches handling for the provided packet according to it's type.
func (d *dispatcherTransmitter) dispatchSPacket(conn Connection, sPacket protocol.SPacket) error {
	var err error
	var cPackets []protocol.CPacket

	debugStateSetter := func(state protocol.State) { // only needed to add the debug log line
		conn := conn
		conn.SetState(state)
		d.log.Debug("changed connState", zap.String("conn", conn.ID().String()), zap.String("state", state.String()))
	}

	switch sPacket.Type() {
	case protocol.SHandshake:
		err = handlers.HandleSHandshake(debugStateSetter, sPacket)
	case protocol.SRequest:
		if cPackets, err = handlers.HandleSRequest(sPacket); err != nil {
			return fmt.Errorf("failed to handle SRequest packet: %w", err)
		}
	case protocol.SPing:
		cPackets, err = handlers.HandleSPing(sPacket)
	case protocol.SLoginStart:
		cPackets, err = handlers.HandleSLoginStart(d.auth, d.ps, debugStateSetter, d.aliver.AddAliveConn, conn.ID(), sPacket)
	case protocol.SEncryptionResponse:
		cPackets, err = handlers.HandleSEncryptionResponse(
			d.auth, d.ps, debugStateSetter, conn.EnableEncryption, conn.EnableCompression, d.aliver.AddAliveConn, conn.ID(), sPacket)
	case protocol.SPluginMessage:
		player, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}
		err = handlers.HandleSPluginMessage(d.log, player, sPacket)
	case protocol.SClientSettings:
		player, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}
		err = handlers.HandleSClientSettings(player, sPacket)
	case protocol.SKeepAlive:
		err = handlers.HandleSKeepAlive(d.aliver.receiveKeepAlive, conn.ID(), sPacket)
	case protocol.SPlayerPosition, protocol.SPlayerMovement:
		err = handlers.HandleSPlayerSpatial(d.roster.SetPlayerSpatial, conn.ID(), sPacket)
	case protocol.SEntityAction:
		err = handlers.HandleSEntityAction(sPacket)
	case protocol.SAnimation:
		err = handlers.HandleSAnimation(sPacket)
	case protocol.SHeldItemChange:
		err = handlers.HandleSHeldItemChange(d.roster.SetPlayerHeldItem, conn.ID(), sPacket)
	case protocol.SClickWindow:
		thisPlayer, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}
		var inventoryUpdated bool
		inventoryUpdated, cPackets, err = handlers.HandleSClickWindow(thisPlayer.State.Inventory, d.log, sPacket)
		if inventoryUpdated {
			d.roster.PlayerInventoryChanged(conn.ID())
		}
	case protocol.SPlayerDigging:
		thisPlayer, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}

		err = handlers.HandleSPlayerDigging(d.ps, d.sharder, thisPlayer, sPacket)
	case protocol.SCloseWindow:
		thisPlayer, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}
		err = handlers.HandleSCloseWindow(thisPlayer, sPacket)
	case protocol.SWindowConfirmation:
		thisPlayer, ok := d.roster.GetPlayerByConnID(conn.ID())
		if !ok {
			err = fmt.Errorf("player %s not found ", conn.ID())
			break
		}
		err = handlers.HandleSWindowConfirmation(thisPlayer.State.Inventory, sPacket)
	default:
		return nil
		// DEBT turn this error back on once all expected packets are handled
		// return fmt.Errorf("unhandled packet type: %X", int32(pacType))
	}

	if err != nil {
		return fmt.Errorf("failed to handle %s packet: %w", sPacket.Type().String(), err)
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

func (d *dispatcherTransmitter) connClosedHandler(lope *envelope.E) {
	closeConn := lope.GetCloseConn()
	if closeConn == nil {
		d.log.Error("failed to parse envelope: there is no closeConn inside", zap.Any("envelope", lope))
		return
	}

	connID, err := uuid.Parse(closeConn.ConnId)
	if err != nil {
		d.log.Error("failed to parse conn ID as UUID", zap.String("id", closeConn.ConnId))
		return
	}
	d.log.Debug("connection closed", zap.String("conn", closeConn.ConnId))

	if playerID, ok := d.roster.GetPlayerIDByConnID(connID); ok {
		playerLeftLope := envelope.PlayerLeft(&pb.PlayerLeft{PlayerId: playerID.String()})
		if err := d.ps.Publish(subj.MkPlayerLeft(), playerLeftLope); err != nil {
			d.log.Error("failed to publish player left message", zap.Error(err))
		}
	}

	d.connMapMu.Lock()
	delete(d.connMu, connID)
	d.connMapMu.Unlock()
}

// DEBT There is a need for global transmissions, to broadcast world state updates, chat messages etc.
//  Broadcasting is done via dedicated channel, which is handled by the broadcast handler. This implementation
//  of the handler is very dumb and will not work at any meaningful scale. This will need to become a separate
//  Broadcaster component that will select a subset of relevant players that actually need to receive the message
//  (e.g. only those subscribed to relevant chat channels, or close enough to the updated chunk).
//  Broadcaster of every cluster node will unicast messages to every relevant player connected to that node.
func (d *dispatcherTransmitter) packetBroadcastHandler(lope *envelope.E) {
	if cPacked := lope.GetCpacket(); cPacked == nil {
		d.log.Error("failed to parse broadcast envelope: there is no cPacked inside", zap.Any("envelope", lope))
		return
	}

	d.connMapMu.Lock()
	for connID := range d.connMu {
		if err := d.ps.Publish(subj.MkConnTransmit(connID), lope); err != nil {
			d.log.Error("failed to publish player left message", zap.Error(err))
		}
	}
	d.connMapMu.Unlock()
}

func (d *dispatcherTransmitter) getTransmitHandler(conn Connection) func(lope *envelope.E) {
	return func(lope *envelope.E) {
		conn := conn
		log := d.log.With(zap.String("conn", conn.ID().String()))

		d.connMu[conn.ID()].Lock()
		defer d.connMu[conn.ID()].Unlock()

		pbCPacket := lope.GetCpacket()
		if pbCPacket == nil {
			log.Error("failed to parse envelope: there is no CPacket inside", zap.Any("envelope", lope))
			return
		}
		pacType := protocol.PacketType(pbCPacket.PacketType)
		log.Debug("transmitting CPacket", zap.String("type", pacType.String()))

		bufOut := buffer.New()
		bufOut.PushVarInt(int32(pacType.ProtocolID()))

		packetBytes := bytes.Join([][]byte{
			bufOut.Bytes(),
			pbCPacket.GetBytes(),
		}, nil)

		if err := d.transmitBytes(conn, packetBytes); err != nil {
			if errors.Is(err, ErrTCPWriteFail) { // we've noticed a dead connection before KeepAliver triggered
				log.Info("closing failed connection", zap.Error(err)) // assuming connection dead and client gone
				_ = conn.Close()                                      // errors here don't really matter, 'cus connection is already dead anyway

				lope := envelope.CloseConn(&pb.CloseConn{
					ConnId: conn.ID().String(),
					State:  pb.ConnState(conn.GetState()),
				})

				if err := d.ps.Publish(subj.MkConnClosed(), lope); err != nil {
					d.log.Error("failed to publish CloseConn", zap.Error(err), zap.String("conn", conn.ID().String()))
				}
			} else {
				log.Error("failed to transmit bytes", zap.Error(err))
			}
			return
		}

		// Cleanly disconnecting the player.
		if pacType == protocol.CDisconnectLogin || pacType == protocol.CDisconnectPlay {
			_ = conn.Close() // errors here don't really matter, 'cus connection is already dead anyway

			lope := envelope.CloseConn(&pb.CloseConn{
				ConnId: conn.ID().String(),
				State:  pb.ConnState(conn.GetState()),
			})

			if err := d.ps.Publish(subj.MkConnClosed(), lope); err != nil {
				d.log.Error("failed to publish CloseConn", zap.Error(err), zap.String("conn", conn.ID().String()))
			}
			return
		}
	}
}

func (d *dispatcherTransmitter) transmitCPacket(conn Connection, cpacket protocol.CPacket) error {
	bufOut := buffer.New()
	bufOut.PushVarInt(int32(cpacket.ProtocolID()))
	cpacket.Push(bufOut)

	d.log.Debug("transmitting packet", zap.String("conn", conn.ID().String()),
		zap.String("type", cpacket.Type().String()))

	if err := d.transmitBuffer(conn, bufOut); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	return nil
}

func (d *dispatcherTransmitter) transmitBytes(conn Connection, packetBytes []byte) error {
	d.log.Debug("transmitting bytes", zap.String("conn", conn.ID().String()),
		zap.String("bytes", hex.EncodeToString(packetBytes)))

	if err := d.transmitBuffer(conn, buffer.NewFrom(packetBytes)); err != nil {
		return fmt.Errorf("failed to transmit buffer: %w", err)
	}

	return nil
}

func (d *dispatcherTransmitter) transmitBuffer(conn Connection, bufOut *buffer.Buffer) error {
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

// TODO add chat message here to tell user why they were disconnected
// DEBT looks like this is not actually dropping the TCP connection, need to add that as well.
func (d *dispatcherTransmitter) forceDisconnect(connState protocol.State, connID uuid.UUID) error {
	d.log.Info("evicting player", zap.String("conn", connID.String()))

	bufOut := buffer.New()
	var pacType protocol.PacketType

	switch connState {
	case protocol.Login:
		cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CDisconnectLogin)
		disconnect := cpacket.(*protocol.CPacketDisconnectLogin)
		disconnect.Push(bufOut)
		pacType = disconnect.Type()
	case protocol.Play:
		cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CDisconnectPlay)
		disconnect := cpacket.(*protocol.CPacketDisconnectPlay)
		disconnect.Push(bufOut)
		pacType = disconnect.Type()
	default:
		return fmt.Errorf("cannot trigger disconnect on conn %s for state %d", connID.String(), connState)
	}

	packetLope := envelope.CPacket(&pb.CPacket{
		Bytes:      bufOut.Bytes(),
		PacketType: pacType.Value(),
	})

	if err := d.ps.Publish(subj.MkConnTransmit(connID), packetLope); err != nil {
		return fmt.Errorf("failed to publish conn disconnect CPacket: %w", err)
	}

	if playerID, ok := d.roster.GetPlayerIDByConnID(connID); ok {
		playerLeftLope := envelope.PlayerLeft(&pb.PlayerLeft{PlayerId: playerID.String()})
		if err := d.ps.Publish(subj.MkPlayerLeft(), playerLeftLope); err != nil {
			return fmt.Errorf("failed to publish player left message: %w", err)
		}
	}

	return nil
}

// checkIsStatusHandshake checks if the packet looks like a Handshake packet. This is needed because in Status
// connection state there is no way in the protocol to correctly signal upgrade to Login state, so the Notchian client
// sends a SHandshake packet, which belongs to Handshake state and it's packetID collides with the SRequest packet
// from the Status state. So we have to hack around this by checking the packet size, and if the connState is Status,
// packetID is 0x00 and size is bigger than 1 byte - assume it is an SHandshake, not SRequest as it normally would be.
func (d *dispatcherTransmitter) checkIsStatusHandshake(connState protocol.State, packetBytes []byte) bool {
	if connState != protocol.Status { // if the connState is not Status - this hack does not apply.
		return false
	} else if len(packetBytes) < 6 { // 6 bytes is absolute minimum possible length of the SHandshake packet.
		return false
	}

	protocolID := buffer.NewFrom(packetBytes).PullVarInt()
	if protocolID != 0x00 { // SHandshake packet protocol ID is 0x00, same as the ID of SRequest packet.
		return false
	} else if len(packetBytes) == 1 { // SRequest packet has no fields so it's length is always 1 byte.
		return false
	}

	return true
}
