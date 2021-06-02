package network

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/pkg/buffer"
)

type Network struct {
	host string
	port int

	dispatcher *DispatcherTransmitter

	log *zap.Logger

	ps      nats.PubSub
	control chan control.Command
}

func NewNetwork(conf control.NetworkConf, log *zap.Logger, report chan control.Command, bus nats.PubSub, disp *DispatcherTransmitter) *Network {
	return &Network{
		host:       conf.Host,
		port:       conf.Port,
		dispatcher: disp,
		control:    report,
		log:        log,
		ps:         bus,
	}
}

func (n *Network) Start(ctx context.Context) error {
	if err := n.startListening(ctx); err != nil {
		return fmt.Errorf("failed to start listening on %s:%d: %w", n.host, n.port, err)
	}
	n.log.Info("started TCP listener", zap.String("host", n.host), zap.Int("port", n.port))

	if err := n.dispatcher.Start(ctx); err != nil {
		return fmt.Errorf("failed to start dispatcher: %w", err)
	}

	return nil
}

func (n *Network) Stop() {
	// TODO gracefully close all connections here
}

func (n *Network) startListening(ctx context.Context) error {
	tcpAddress, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("address resolution failed: %w", err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		return fmt.Errorf("failed to bind TCP: %w", err)
	}

	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil { // stop the server if the TCP listener goroutine dies
				n.control <- control.Command{Signal: control.SERVER_FAIL, Message: fmt.Sprintf("TCP listener panicked: %v", r)}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				n.log.Info("TCP listener closing")
				return
			default:
				conn, err := tcpListener.AcceptTCP()
				if err != nil {
					n.log.Error("failed to accept a TCP connection",
						zap.String("host", n.host), zap.Int("port", n.port), zap.Error(err))
					n.control <- control.Command{Signal: control.SERVER_FAIL, Message: err.Error()}
					return
				}

				n.log.Debug("received TCP connection",
					zap.String("from", conn.RemoteAddr().String()), zap.Any("conn", conn))

				_ = conn.SetNoDelay(true)
				_ = conn.SetKeepAlive(true)

				go n.handleNewConnection(ctx, NewConnection(conn))
			}
		}

		n.control <- control.Command{Signal: control.SERVER_FAIL, Message: fmt.Sprintf("TCP listener stopped unexpectedly")}
	}(ctx)

	return nil
}

func (n *Network) handleNewConnection(ctx context.Context, conn Connection) {
	n.log.Debug("new connection", zap.Any("address", conn.Address().String()))

	if err := n.dispatcher.RegisterNewConn(conn); err != nil {
		n.log.Error("failed to register conn subscriptions", zap.Error(err), zap.Any("conn", conn))
		if err = conn.Close(); err != nil {
			n.log.Error("error while closing failed connection", zap.Error(err), zap.Any("conn", conn))
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			n.log.Info("connection closing")
			return
		default:
			bufIn := buffer.New()
			size, err := conn.Receive(bufIn)
			if err != nil && err.Error() == "EOF" {
				// TODO broadcast player disconnect if conn.GetState() == Play.

				break
			} else if err != nil || size == 0 {
				_ = conn.Close()
				// TODO broadcast player disconnect if conn.GetState() == Play.
				break
			}

			// decompression
			// decryption

			if bufIn.Bytes()[0] == 0xFE { // LEGACY PING
				continue
			}

			n.log.Debug("reading blob of bytes", zap.Int("len", bufIn.Len()), zap.String("conn", conn.ID().String()))
			var packetCount int
			for {
				packetLen := bufIn.PullVarInt()
				if packetLen == 0 {
					break // no more packets in this blob
				}

				packetBytes := bufIn.Bytes()[bufIn.IndexI() : bufIn.IndexI()+packetLen]
				bufIn.SkipLen(packetLen)
				packetCount++

				n.log.Debug(fmt.Sprintf("read a packet %d in blob", packetCount),
					zap.Int("packetLen", int(packetLen)), zap.String("bytes", fmt.Sprintf("%X", packetBytes)))
				n.dispatcher.HandleSPacket(conn, packetBytes)
			}
		}

	}
}
