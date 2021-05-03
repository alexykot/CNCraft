// This is a test client intended for protocol reverse engineering only. Not a production client, not even a stub of it.
package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/protocol"
)

const serverHost = "127.0.0.1"
const serverPort = 25565
const username = "Kolsar"

var conn network.Connection
var connState protocol.State

func main() {
	err := startNetwork()
	if err != nil {
		panic(err)
	}

	go handleReceive()

	// Send SHandshake
	spacket, _ := protocol.GetPacketFactory().MakeSPacket(protocol.SHandshake)
	handshake, _ := spacket.(*protocol.SPacketHandshake)
	handshake.Version = protocol.Version
	handshake.Host = serverHost
	handshake.Port = serverPort
	handshake.NextState = protocol.Login

	bufHandshake := buffer.New()
	handshake.Push(bufHandshake)
	sendBuffer(bufHandshake)
	println("SHandshake sent")
	println(fmt.Sprintf("SHandshake bytes: %X", bufHandshake.Bytes()))

	// Send SLoginStart
	spacket, _ = protocol.GetPacketFactory().MakeSPacket(protocol.SLoginStart)
	loginStart, _ := spacket.(*protocol.SPacketLoginStart)
	loginStart.Username = username

	connState = protocol.Login
	bufLoginStart := buffer.New()
	loginStart.Push(bufLoginStart)
	sendBuffer(bufLoginStart)
	println("SLoginStart sent")
	println(fmt.Sprintf("SLoginStart bytes: %X", bufLoginStart.Bytes()))

	wait()
}

func wait() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		break
	}
}

func startNetwork() error {
	addr, err := net.ResolveTCPAddr("tcp", serverHost+":"+strconv.Itoa(serverPort))
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial TCP: %w", err)
	}
	conn = network.NewConnection(tcpConn)
	connState = protocol.Handshake

	return nil
}

func sendBuffer(bufOut buffer.B) {
	if _, err := conn.Transmit(bufOut); err != nil {
		panic(fmt.Errorf("failed to transmit buffer: %w", err))
	}
}

func handleReceive() {
	for {
		bufIn := buffer.New()
		size, err := conn.Receive(bufIn)
		if err != nil && err.Error() == "EOF" {
			break
		} else if err != nil || size == 0 {
			_ = conn.Close()
			break
		}

		packetLen := bufIn.PullVarInt()
		// println(fmt.Sprintf("received size: %d, declared packet length: %d", size, packetLen))
		packetBytes := bufIn.Bytes()[bufIn.IndexI() : bufIn.IndexI()+packetLen]

		handleCPacket(packetBytes)
	}
}

func handleCPacket(packetBytes []byte) {
	bufI := buffer.NewFrom(packetBytes)
	protocolPacketID := protocol.ProtocolPacketID(bufI.PullVarInt())
	if protocolPacketID == 0x19 { // server sends CDisconnectPlay in Handshake state
		connState = protocol.Play
	}

	pacType := protocol.MakeCType(connState, protocolPacketID)
	cPacket, err := protocol.GetPacketFactory().MakeCPacket(pacType)
	println(fmt.Sprintf("received packet; type %d/%X, %s; size %d", connState, protocolPacketID, pacType.String(), len(packetBytes)))
	if err != nil {
		println(fmt.Sprintf("failed to make CPacket: %v", err))
		println()
		return
	}
	println()

	switch pacType {
	case protocol.CDisconnectLogin:
		disconnect := cPacket.(*protocol.CPacketDisconnectLogin)
		disconnect.Pull(bufI)
		println(fmt.Sprintf("received disconnect, reason: %s", disconnect.Reason))

	case protocol.CDisconnectPlay:
		disconnect := cPacket.(*protocol.CPacketDisconnectPlay)
		disconnect.Pull(bufI)
		println(fmt.Sprintf("received disconnect, reason: %s", disconnect.Reason))

	case protocol.CLoginSuccess:
		connState = protocol.Play

	case protocol.CWindowConfirmation:
		winConfirm := cPacket.(*protocol.CPacketWindowConfirmation)
		winConfirm.Pull(bufI)
		println(fmt.Sprintf("Window Confirmation: %v", winConfirm))

	case protocol.CChunkData:
		spacket, _ := protocol.GetPacketFactory().MakeSPacket(protocol.SHandshake)
		handshake, _ := spacket.(*protocol.SPacketHandshake)
		handshake.Version = protocol.Version
		handshake.Host = serverHost
		handshake.Port = serverPort
		handshake.NextState = protocol.Login

		bufHandshake := buffer.New()
		handshake.Push(bufHandshake)
		sendBuffer(bufHandshake)
		println("SHandshake sent")
		println(fmt.Sprintf("SHandshake bytes: %X", bufHandshake.Bytes()))
	}
}

func prettyPrintBytes(bytes []byte) {
	for i, byte := range bytes {
		if i > 0 && i%64 == 0 {
			println()
		}
		print(fmt.Sprintf("%02X", byte))
	}
	println()

	// var i int
	// for i = 64; i < len(bytes); i = i + 64 {
	// 	println(fmt.Sprintf("%X", bytes[i-64:i]))
	// }
	//
	// if i-64 < len(bytes) {
	// 	println(fmt.Sprintf("%X", bytes[i-64:len(bytes)-1]))
	// }
}
