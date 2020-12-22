package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
)

const CurrentProtocol = protocol.MC1_15_2

type Network struct {
	host string
	port int

	dispatcher *SPacketDispatcher

	log *zap.Logger

	ps      nats.PubSub
	control chan control.Command
}

func NewNetwork(conf control.NetworkConf, log *zap.Logger, report chan control.Command, bus nats.PubSub, disp *SPacketDispatcher) *Network {
	return &Network{
		host:       conf.Host,
		port:       conf.Port,
		dispatcher: disp,
		control:    report,
		log:        log,
		ps:         bus,
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
	if err := n.ps.Subscribe(subj.MkConnSend(conn.ID()), handler.HandlePacketSend); err != nil {
		n.log.Error("failed to subscribe to conn.send subject", zap.Error(err), zap.Any("conn", conn))
		if err = conn.Close(); err != nil {
			n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
		}
		return
	}
	if err := n.ps.Subscribe(subj.MkConnStateChange(conn.ID()), handler.HandleConnState); err != nil {
		n.log.Error("failed to subscribe to conn.state subject", zap.Error(err), zap.Any("conn", conn))
		if err = conn.Close(); err != nil {
			n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
		}
		return
	}

	// DEBT This subj us not used atm, may be still useful later, or to delete if turns out to be not useful.
	//if err := n.ps.Publish(subj.MkNewConn(), envelope.NewConn(&pb.NewConnection{Id: conn.ID().String()}, nil)); err != nil {
	//	n.log.Error("failed to publish conn.new message", zap.Error(err), zap.Any("conn", conn))
	//	if err = conn.Close(); err != nil {
	//		n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
	//	}
	//	return
	//}

	// Register packet receiving subscription
	// make sure we subscribe to receive incoming packets before pulling incoming bytes from connection
	n.dispatcher.RegisterNewConnection(conn.ID())

	var inf []byte
	for {
		inf = make([]byte, 1024)
		size, err := conn.Pull(inf)
		n.log.Sugar().Debugf("received %d bytes from connection", size)

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
	s.net.log.Debug("changed connstate", zap.String("conn", s.conn.ID().String()), zap.String("state", state.String()))
}

func (s *connHandler) HandlePacketSend(lope *envelope.E) {
	cpacket := lope.GetCpacket()
	if cpacket == nil {
		s.net.log.Error("received empty CPacket message, cannot send")
		return
	}
	bufO := buffer.NewFrom(cpacket.GetBytes())
	if bufO.Len() < 2 {
		s.net.log.Error("received CPacket with zero length buffer, cannot send")
		return
	}

	temp := buffer.New()
	temp.PushVrI(bufO.Len())

	deflated := buffer.New()
	deflated.PushUAS(s.conn.Deflate(bufO.UAS()), false)
	temp.PushUAS(deflated.UAS(), false)

	if _, err := s.conn.Push(s.conn.Encrypt(temp.UAS())); err != nil {
		s.net.log.Error("Failed to push client bound packet", zap.Error(err))
		return
	}
	s.net.log.Debug("pushed packet to client", zap.Any("packet", cpacket))
}

func (s *connHandler) HandlePacketReceive(buffBytes []byte) {
	println("received bytes:", hex.EncodeToString(buffBytes))
	lope := envelope.SPacket(&pb.SPacket{
		Bytes: buffBytes,
		State: pb.ConnState(s.conn.GetState()),
	}, nil)

	if err := s.net.ps.Publish(subj.MkConnReceive(s.conn.ID()), lope); err != nil {
		s.net.log.Error("failed to publish SPacket message", zap.Error(err))
	}

	s.net.log.Debug("received packet from client", zap.String("conn", s.conn.ID().String()))
}
