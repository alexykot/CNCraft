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

	log *zap.Logger
	ctx context.Context

	dispatcher *DispatcherTransmitter

	ps      nats.PubSub
	control chan control.Command
}

func NewNetwork(ctx context.Context, conf control.NetworkConf, log *zap.Logger, ctrlChan chan control.Command, bus nats.PubSub, disp *DispatcherTransmitter) *Network {
	return &Network{
		ctx:        ctx,
		host:       conf.Host,
		port:       conf.Port,
		dispatcher: disp,
		control:    ctrlChan,
		log:        log,
		ps:         bus,
	}
}

func (n *Network) Start() {
	n.signal(control.STARTING, nil)

	if err := n.startListening(n.ctx); err != nil {
		n.signal(control.FAILED, fmt.Errorf("failed to start network: failed to start listening on %s:%d: %w", n.host, n.port, err))
		return
	}

	if err := n.dispatcher.Start(n.ctx); err != nil {
		n.signal(control.FAILED, fmt.Errorf("failed to start network: failed to start dispatcher: %w", err))
		return
	}

	n.signal(control.READY, nil)
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

	// Context cancellation will signal server stopping and needs to be handled correctly.
	// Context cancellation cannot be handled in the infinite for{} loop of the tcpListener because
	// tcpListener.AcceptTCP() will block indefinitely until a new connection will appear, so will not allow
	// to handle context.Done() signal timely.
	go func() {
		defer func() {
			if r := recover(); r != nil { // stop the server if the TCP listener goroutine dies
				n.signal(control.FAILED, fmt.Errorf("TCP listener panicked: %v", r))
			}
		}()

		select {
		case <-ctx.Done():
			if err := tcpListener.Close(); err != nil {
				n.signal(control.FAILED, err)
				return
			}
			n.log.Info("TCP listener closed")
			n.signal(control.STOPPED, nil)
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil { // stop the server if the TCP listener goroutine dies
				n.signal(control.FAILED, fmt.Errorf("TCP listener panicked: %v", r))
			}
		}()

		for {
			conn, err := tcpListener.AcceptTCP() // this blocking call will wait until some connections will appear on the wire
			if err != nil {
				n.signal(control.FAILED, fmt.Errorf("failed to accept a TCP connection on %s:%d: %w", n.host, n.port, err))
				return
			}

			n.log.Debug("received TCP connection",
				zap.String("from", conn.RemoteAddr().String()), zap.Any("conn", conn))

			_ = conn.SetNoDelay(true)
			_ = conn.SetKeepAlive(true)

			go n.handleNewConnection(ctx, NewConnection(conn))
		}
	}()

	n.log.Info("started TCP listener", zap.String("host", n.host), zap.Int("port", n.port))

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

	go func() {
		defer func() {
			if r := recover(); r != nil { // stop the server if the TCP listener goroutine dies
				n.signal(control.FAILED, fmt.Errorf("TCP listener panicked: %v", r))
			}
		}()

		select {
		case <-ctx.Done():
			if err := conn.Close(); err != nil {
				n.log.Error("failed while closing connection", zap.String("conn", conn.ID().String()), zap.Error(err))
				return
			}
			n.log.Info("connection closed", zap.String("conn", conn.ID().String()))
		}
	}()

	for {
		bufIn := buffer.New()
		size, err := conn.Receive(bufIn) // this blocking call will wait until some bytes will appear on the wire
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

func (n *Network) signal(state control.ComponentState, err error) {
	n.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.NETWORK,
		State:     state,
		Err:       err,
	}
}
