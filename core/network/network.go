package network

import (
	"fmt"
	"net"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
)

const CurrentProtocol = protocol.MC1_15_2

type Network struct {
	host string
	port int

	log     *zap.Logger
	packFac PacketFactory

	ps      nats.PubSub
	control chan control.Command
}

func NewNetwork(conf control.NetworkConf, packFac PacketFactory, log *zap.Logger, report chan control.Command, bus nats.PubSub) *Network {
	return &Network{
		host:    conf.Host,
		port:    conf.Port,
		control: report,
		log:     log,
		ps:      bus,
		packFac: packFac,
	}
}

func (n *Network) Start() error {
	if err := n.startListening(); err != nil {
		return fmt.Errorf("failed to start listening on %s:%d: %w", n.host, n.port, err)
	}
	n.log.Info("started listening", zap.String("host", n.host), zap.Int("port", n.port))

	return nil
}

func (n *Network) Stop() {
	// TODO gracefully close all connections here
}

func (n *Network) startListening() error {
	ser, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("address resolution failed: %w", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("failed to bind TCP: %w", err)
	}

	// TODO add supervisor and ensure server gracefully stops if the connections goroutine fails.
	go func() {
		for {
			conn, err := tcp.AcceptTCP()

			if err != nil {
				n.log.Error("failed to accept a TCP connection",
					zap.String("host", n.host), zap.Int("port", n.port), zap.Error(err))
				n.control <- control.Command{Signal: control.FAIL, Message: err.Error()}
				break
			}

			n.log.Debug("received TCP connection",
				zap.String("from", conn.RemoteAddr().String()), zap.Any("conn", conn))

			_ = conn.SetNoDelay(true)
			_ = conn.SetKeepAlive(true)

			go n.handleNewConnection(NewConnection(conn))
		}
	}()

	return nil
}

func (n *Network) handleNewConnection(conn Connection) {
	n.log.Debug("new connection", zap.Any("address", conn.Address().String()))

	handler := &connHandler{conn: conn, net: n}

	if err := n.ps.Subscribe(MkConnSubjSend(conn.ID()), handler.HandlePacketSend); err != nil {
		n.log.Error("failed to subscribe to conn.send subject", zap.Error(err), zap.Any("conn", conn))
		if err = conn.Close(); err != nil {
			n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
		}
	}
	if err := n.ps.Subscribe(MkConnSubjStateChange(conn.ID()), handler.HandleConnState); err != nil {
		n.log.Error("failed to subscribe to conn.state subject", zap.Error(err), zap.Any("conn", conn))
		if err = conn.Close(); err != nil {
			n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
		}
	}

	var inf []byte
	for {
		inf = make([]byte, 1024)
		size, err := conn.Pull(inf)

		if err != nil && err.Error() == "EOF" {
			// TODO broadcast player leaving

			break
		} else if err != nil || size == 0 {
			_ = conn.Close()

			// TODO broadcast player leaving
			break
		}

		buf := buffer.NewFrom(conn.Decrypt(inf[:size]))

		// decompression
		// decryption

		if buf.UAS()[0] == 0xFE { // LEGACY PING
			continue
		}

		packetLen := buf.PullVrI()

		handler.HandlePacketReceive(buf.UAS()[buf.InI() : buf.InI()+packetLen])
	}
}

type connHandler struct {
	conn Connection
	net  *Network
}

func (s *connHandler) HandleConnState(envelope *envelope.E) {
	state, err := protocol.IntToState(int(envelope.GetConnState().State))
	if err != nil {
		s.net.log.Error("failed to parse state", zap.Error(err))
		return
	}
	s.conn.SetState(state)
}

func (s *connHandler) HandlePacketSend(lope *envelope.E) {
	cpacket := lope.GetCpacket()
	if cpacket == nil {
		s.net.log.Error("received empty CPacket message, cannot send")
		return
	}
	bufO := buffer.NewFrom(cpacket.GetBytes())

	if bufO.Len() > 1 {
		temp := buffer.New()
		temp.PushVrI(bufO.Len())

		comp := buffer.New()
		comp.PushUAS(s.conn.Deflate(bufO.UAS()), false)
		temp.PushUAS(comp.UAS(), false)

		if _, err := s.conn.Push(s.conn.Encrypt(temp.UAS())); err != nil {
			s.net.log.Error("Failed to push client bound packet", zap.Error(err))
		} else {
			s.net.log.Debug("pushed packet to client")
		}
	}
}

func (s *connHandler) HandlePacketReceive(buffBytes []byte) {
	// TODO relocate this to subscribing side
	//bufI := buffer.NewFrom(buffBytes)
	//protocolPacketID := protocol.ProtocolPacketID(bufI.PullVrI())
	//id := protocol.MakeID(protocol.ServerBound, s.conn.GetState(), protocolPacketID)
	//bufI.Len()
	//incomingPacket, err := s.net.packFac.MakeSPacket(id)
	//if err != nil {
	//	s.net.log.Warn("unable to decode packet", zap.Int("packet_id", int(id)), zap.Error(err))
	//	return
	//}
	//
	//// populate incoming packet
	//if err := incomingPacket.Pull(bufI); err != nil {
	//	s.net.log.Warn("malformed packet", zap.Int("packet_id", int(id)), zap.Error(err))
	//	return
	//}

	lope := envelope.NewWithSPacket(&pb.SPacket{Bytes: buffBytes}, nil)
	if err := s.net.ps.Publish(MkConnSubjReceive(s.conn.ID()), lope); err != nil {
		s.net.log.Error("failed to publish SPacket message", zap.Error(err))
	}
}

// MkConnSubjSend creates a subject name string for given connection ID for receiving server bound packets
func MkConnSubjReceive(connID uuid.UUID) string {
	return "conn." + connID.String() + ".receive"
}

// MkConnSubjSend creates a subject name string for given connection ID for sending client bound packets
func MkConnSubjSend(connID uuid.UUID) string {
	return "conn." + connID.String() + ".send"
}

// MkConnSubjSend creates a subject name string for given connection ID for handling connection state changes
func MkConnSubjStateChange(connID uuid.UUID) string {
	return "conn." + connID.String() + ".state"
}
