package conn

import (
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/system"
	"github.com/golangmc/minecraft-server/impl/protocol"
)

type network struct {
	host string
	port int

	logger  *logs.Logging
	packFac protocol.PacketFactory

	pubsub bus.PubSub

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection

	report chan system.Message
}

func NewNetwork(host string, port int, packFac protocol.PacketFactory, report chan system.Message,
	join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Network {
	return &network{
		host: host,
		port: port,

		join: join,
		quit: quit,

		report: report,

		logger:  logs.NewLogging("network", logs.EveryLevel...),
		packFac: packFac,
	}
}

func (n *network) Load() {
	if err := n.startListening(); err != nil {
		n.report <- system.Make(system.FAIL, err)
		return
	}
}

func (n *network) Kill() {

}

func (n *network) startListening() error {
	ser, err := net.ResolveTCPAddr("tcp", n.host+":"+strconv.Itoa(n.port))
	if err != nil {
		return fmt.Errorf("address resolution failed [%v]", err)
	}

	tcp, err := net.ListenTCP("tcp", ser)
	if err != nil {
		return fmt.Errorf("failed to bind [%v]", err)
	}

	n.logger.InfoF("listening on %s:%d", n.host, n.port)

	go func() {
		for {
			con, err := tcp.AcceptTCP()

			if err != nil {
				n.report <- system.Make(system.FAIL, err)
				break
			}

			_ = con.SetNoDelay(true)
			_ = con.SetKeepAlive(true)

			go handleConnect(n, NewConnection(con))
		}
	}()

	return nil
}

func handleConnect(net *network, conn base.Connection) {
	net.logger.DebugF("New Connection from &6%v", conn.Address())

	var inf []byte

	subscribeConn(net, conn)

	for {
		inf = make([]byte, 1024)
		size, err := conn.Pull(inf)

		if err != nil && err.Error() == "EOF" {
			net.quit <- base.PlayerAndConnection{
				Player:     nil,
				Connection: conn,
			}

			break
		} else if err != nil || size == 0 {
			_ = conn.Stop()

			net.quit <- base.PlayerAndConnection{
				Player:     nil,
				Connection: conn,
			}
			break
		}

		buf := NewBufferWith(conn.Decrypt(inf[:size]))

		// decompression
		// decryption

		if buf.UAS()[0] == 0xFE { // LEGACY PING
			continue
		}

		packetLen := buf.PullVrI()
		bufI := NewBufferWith(buf.UAS()[buf.InI() : buf.InI()+packetLen])
		handleReceive(net, conn, bufI)
	}
}

type connHandler struct {
	conn base.Connection
	net  *network
}

func (s *connHandler) HandlePacketSend(envelope bus.Envelope) {
	bufO, ok := envelope.GetMessage().(buff.Buffer)
	if !ok {
		s.net.logger.Error("Failed to cast message to buffer")
		return
	}

	if bufO.Len() > 1 {
		temp := NewBuffer()
		temp.PushVrI(bufO.Len())

		comp := NewBuffer()
		comp.PushUAS(s.conn.Deflate(bufO.UAS()), false)
		temp.PushUAS(comp.UAS(), false)

		if _, err := s.conn.Push(s.conn.Encrypt(temp.UAS())); err != nil {
			s.net.logger.Error("Failed to push client bound packet: %v", err)
		} else {
			s.net.logger.DebugF("Pushed packet to client")
		}
	}
}

func (s *connHandler) HandleConnState(envelope bus.Envelope) {
	state, ok := envelope.GetMessage().(protocol.State)
	if !ok {
		s.net.logger.Error("Failed to cast message to protocol.state")
		return
	}
	s.conn.SetState(state)
}

func subscribeConn(net *network, conn base.Connection) {
	s := &connHandler{
		conn: conn,
		net:  net,
	}

	net.pubsub.Subscribe(MakeConnTopicSend(conn.ID()), s.HandlePacketSend)
	net.pubsub.Subscribe(MakeConnTopicState(conn.ID()), s.HandleConnState)
}

func handleReceive(net *network, conn base.Connection, bufI buff.Buffer) {
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVrI())

	id := protocol.MakeID(protocol.ServerBound, conn.GetState(), protocolPacketID)

	incomingPacket, err := net.packFac.GetSPacket(id)
	if err != nil {
		net.logger.WarnF("unable to decode %v packet with ID: %d", conn.GetState(), protocolPacketID)
		return
	}

	net.logger.DebugF("GET packet: %d | %v | %v", incomingPacket.ID(), reflect.TypeOf(incomingPacket), conn.GetState())

	// populate incoming packet
	if err := incomingPacket.Pull(bufI, conn); err != nil {
		net.logger.WarnF("malformed packet: %d", incomingPacket.ID())
		return
	}

	net.pubsub.Publish(protocol.MakePacketTopic(incomingPacket.ID()),
		bus.NewEnvelope(incomingPacket, map[string]string{bus.MetaConn: conn.ID()}))

	// TODO this double publishing is weird and some subscribers actually expect messages both at once.
	//  Those subscribers need to be refactored, really they are trying to handle incoming packet and immediately
	//  send back the response over the provided connection. Instead it should publish response in a topic and
	//  there should be separate connection subscribers listening for client bound packets.
	//  This will need to also understand what packet goes to what client. Maybe topic per packet per client?
	// net.pubsub.Publish(incomingPacket, conn)
}

func MakeConnTopicSend(connID string) string {
	return "conn." + connID + ".send"
}

func MakeConnTopicState(connID string) string {
	return "conn." + connID + ".state"
}
