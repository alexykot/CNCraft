package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/buffer"
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

	// kept to reuse when we'll get to Play state
	//if err := n.ps.Publish(subj.MkNewConn(), envelope.NewConn(&pb.NewConnection{Id: conn.ID().String()}, nil)); err != nil {
	//	n.log.Error("failed to publish conn.new message", zap.Error(err), zap.Any("conn", conn))
	//	if err = conn.Close(); err != nil {
	//		n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
	//	}
	//	return
	//}

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
		packetBytes := buf.UAS()[buf.InI() : buf.InI()+packetLen]

		n.log.Sugar().Debugf("received bytes: %s", hex.EncodeToString(packetBytes))

		n.dispatcher.HandleSPacket(conn, packetBytes)
		n.log.Debug("received packet from client", zap.String("conn", conn.ID().String()))
	}
}
