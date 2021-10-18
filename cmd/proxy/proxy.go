// This is a test client intended for protocol reverse engineering only. Not a production client, not even a stub of it.
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "go.uber.org/zap"

    "github.com/alexykot/cncraft/core/control"
    "github.com/alexykot/cncraft/core/network"
    "github.com/alexykot/cncraft/pkg/buffer"
    "github.com/alexykot/cncraft/pkg/log"
    "github.com/alexykot/cncraft/pkg/protocol"
)

// Address of the real vanilla Minecraft server
const serverHost = "127.0.0.1"
const serverPort = 25565

// Address where the proxy itself will start
const proxyHost = "127.0.0.1"
const proxyPort = 25566

var serverConn network.Connection
var connState protocol.State

const logLevel = "DEBUG"

var l *zap.Logger

func main() {
    var err error

    globalCtx := context.Background()

    if l, err = log.GetRoot(logLevel); err != nil {
        panic(err)
    }

    ctrl := make(chan control.Command)
    disp, err := newDispatcher(ctrl, serverHost, serverPort)
    if err != nil {
        panic(err)
    }

    if err = startProxyServer(globalCtx, ctrl, disp, proxyHost, proxyPort); err != nil {
        panic(err)
    }

    wait(globalCtx, ctrl)
}

func wait(ctx context.Context, ctrl chan control.Command) {
    ctx, cancelFunc := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
    defer cancelFunc()

    for {
        select {
        case <-ctx.Done():
            return
        case cmd := <-ctrl:
            if cmd.State == control.FAILED {
                l.Error("component failed", zap.Error(cmd.Err), zap.String("comp", string(cmd.Component)))
                return
            } else {
                l.Info("signal from component", zap.Int("state", int(cmd.State)), zap.String("comp", string(cmd.Component)))
            }
        }
    }
}

func startProxyServer(ctx context.Context, ctrl chan control.Command, disp network.Dispatcher, proxyHost string, proxyPort int) error {
    n := network.NewNetwork(l, ctrl, control.NetworkConf{
        Host: proxyHost,
        Port: proxyPort,
    }, &nubSub{}, disp)

    n.Start(ctx)

    return nil
}

func sendBuffer(bufOut *buffer.Buffer) {
    if _, err := serverConn.Transmit(bufOut); err != nil {
        panic(fmt.Errorf("failed to transmit buffer: %w", err))
    }
}

func handleReceiveClientBound() {
    for {
        bufIn := buffer.New()
        size, err := serverConn.Receive(bufIn)
        if err != nil && err.Error() == "EOF" {
            break
        } else if err != nil || size == 0 {
            _ = serverConn.Close()
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
