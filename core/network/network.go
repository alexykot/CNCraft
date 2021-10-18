package network

import (
	"context"
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
)

type Network struct {
	host string
	port int

	log *zap.Logger

	dispatcher Dispatcher

	ps      nats.PubSub
	control chan control.Command
}

func NewNetwork(log *zap.Logger, ctrlChan chan control.Command, conf control.NetworkConf, bus nats.PubSub, disp Dispatcher) *Network {
	return &Network{
		host:       conf.Host,
		port:       conf.Port,
		dispatcher: disp,
		control:    ctrlChan,
		log:        log,
		ps:         bus,
	}
}

func (n *Network) Start(ctx context.Context) {
	n.signal(control.STARTING, nil)

	if err := n.startListening(ctx); err != nil {
		n.signal(control.FAILED, fmt.Errorf("failed to start network: failed to start listening on %s:%d: %w", n.host, n.port, err))
		return
	}

	if err := n.dispatcher.Init(ctx); err != nil {
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
		_ = conn.Close() // errors here don't really matter, 'cus connection is already dead anyway
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil { // stop the server if the TCP listener goroutine dies
				n.signal(control.FAILED, fmt.Errorf("TCP listener panicked: %v", r))
			}
		}()

		select {
		case <-ctx.Done(): // close connection when context cancelled, i.e. when server is shutting down
			_ = conn.Close() // errors here don't really matter
			n.log.Info("connection closed", zap.String("conn", conn.ID().String()))
		}
	}()

	for {
		bufIn := buffer.New()
		size, err := conn.Receive(bufIn) // this blocking call will wait until some bytes will appear on the wire

		if err != nil && (err.Error() == "EOF" || size == 0) {
			n.log.Debug("connection lost", zap.String("conn", conn.ID().String()))
			_ = conn.Close() // errors here don't really matter, 'cus connection is already dead anyway

			lope := envelope.CloseConn(&pb.CloseConn{
				ConnId: conn.ID().String(),
				State:  pb.ConnState(conn.GetState()),
			})

			if err := n.ps.Publish(subj.MkConnClosed(), lope); err != nil {
				n.log.Error("failed to publish CloseConn", zap.Error(err), zap.String("conn", conn.ID().String()))
				return
			}
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
