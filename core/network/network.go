package network

import (
	"fmt"
	"net"
	"strconv"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/pkg/buffers"
	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/protocol"
)

type Network struct {
	host string
	port int

	log     *zap.Logger
	packFac PacketFactory

	bus     bus.PubSub
	control chan control.Command
}

func NewNetwork(host string, port int, packFac PacketFactory, log *zap.Logger, report chan control.Command, bus bus.PubSub) *Network {
	return &Network{
		host:    host,
		port:    port,
		control: report,
		log:     log,
		bus:     bus,
		packFac: packFac,
	}
}

func (n *Network) Load() {
	if err := n.startListening(); err != nil {
		n.log.Error("failed to start listening",
			zap.String("host", n.host), zap.Int("port", n.port), zap.Error(err))
		n.control <- control.Command{Signal: control.FAIL}
		return
	}

	n.log.Info("started listening", zap.String("host", n.host), zap.Int("port", n.port))
}

func (n *Network) Kill() {
	// TODO gracefully close all connections here
}

func (n *Network) startListening() error {
	ser, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("address resolution failed [%v]", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("failed to bind [%v]", err)
	}

	// TODO add supervisor and ensure server gracefully stops if the connections goroutine fails.
	go func() {
		for {
			conn, err := tcp.AcceptTCP()

			if err != nil {
				n.log.Error("failed to accept a TCP connection",
					zap.String("host", n.host), zap.Int("port", n.port), zap.Error(err))
				n.control <- control.Command{Signal: control.FAIL}
				break
			}

			_ = conn.SetNoDelay(true)
			_ = conn.SetKeepAlive(true)

			go handleConnect(n, NewConnection(conn))
		}
	}()

	return nil
}

func handleConnect(n *Network, conn Connection) {
	n.log.Debug("new connection", zap.String("address", conn.Address().String()))

	var inf []byte

	subscribeConn(n, conn)

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

		buf := buffers.NewBufferWith(conn.Decrypt(inf[:size]))

		// decompression
		// decryption

		if buf.UAS()[0] == 0xFE { // LEGACY PING
			continue
		}

		packetLen := buf.PullVrI()
		bufI := buffers.NewBufferWith(buf.UAS()[buf.InI() : buf.InI()+packetLen])
		handleReceive(n, conn, bufI)
	}
}

type connHandler struct {
	conn Connection
	net  *Network
}

func (s *connHandler) HandlePacketSend(envelope bus.Envelope) {
	bufO, ok := envelope.GetMessage().(buffers.Buffer)
	if !ok {
		s.net.log.Error("failed to cast message to buffer")
		return
	}

	if bufO.Len() > 1 {
		temp := buffers.NewBuffer()
		temp.PushVrI(bufO.Len())

		comp := buffers.NewBuffer()
		comp.PushUAS(s.conn.Deflate(bufO.UAS()), false)
		temp.PushUAS(comp.UAS(), false)

		if _, err := s.conn.Push(s.conn.Encrypt(temp.UAS())); err != nil {
			s.net.log.Error("Failed to push client bound packet", zap.Error(err))
		} else {
			s.net.log.Debug("pushed packet to client")
		}
	}
}

func (s *connHandler) HandleConnState(envelope bus.Envelope) {
	state, ok := envelope.GetMessage().(protocol.State)
	if !ok {
		s.net.log.Error("Failed to cast message to protocol.state")
		return
	}
	s.conn.SetState(state)
}

func subscribeConn(net *Network, conn Connection) {
	s := &connHandler{
		conn: conn,
		net:  net,
	}

	net.bus.Subscribe(MakeConnTopicSend(conn.ID()), s.HandlePacketSend)
	net.bus.Subscribe(MakeConnTopicState(conn.ID()), s.HandleConnState)
}

func handleReceive(net *Network, conn Connection, bufI buffers.Buffer) {
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVrI())

	id := protocol.MakeID(protocol.ServerBound, conn.GetState(), protocolPacketID)

	incomingPacket, err := net.packFac.MakeSPacket(id)
	if err != nil {
		net.log.Warn("unable to decode packet", zap.Int("packet_id", int(id)), zap.Error(err))
		return
	}

	// populate incoming packet
	if err := incomingPacket.Pull(bufI); err != nil {
		net.log.Warn("malformed packet", zap.Int("packet_id", int(id)), zap.Error(err))
		return
	}

	net.bus.Publish(protocol.MakePacketTopic(incomingPacket.ID()),
		bus.NewEnvelope(incomingPacket, map[string]string{bus.MetaConn: conn.ID()}))

	// TODO this double publishing is weird and some subscribers actually expect messages both at once.
	//  Those subscribers need to be refactored, really they are trying to handle incoming packet and immediately
	//  send back the response over the provided connection. Instead it should publish response in a topic and
	//  there should be separate connection subscribers listening for client bound packets.
	//  This will need to also understand what packet goes to what client. Maybe topic per packet per client?
	//
	// net.bus.Publish(incomingPacket, conn)
}

func MakeConnTopicSend(connID string) string {
	return "conn." + connID + ".send"
}

func MakeConnTopicState(connID string) string {
	return "conn." + connID + ".state"
}
