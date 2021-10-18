// This is a test client intended for protocol reverse engineering only. Not a production client, not even a stub of it.
package main

import (
    "context"
    "fmt"
    "net"

    "github.com/alexykot/cncraft/core/control"
    "github.com/alexykot/cncraft/core/nats/subj"
    "github.com/alexykot/cncraft/core/network"
    "github.com/alexykot/cncraft/pkg/envelope"
)

// nubSub - noop pubsub implementation.
type nubSub struct {}

func (ps *nubSub) Start(_ context.Context)                                  {}
func (ps *nubSub) Publish(_ subj.Subj, _ ...*envelope.E) error              { return nil }
func (ps *nubSub) Subscribe(_ subj.Subj, _ func(message *envelope.E)) error { return nil }
func (ps *nubSub) Unsubscribe(_ subj.Subj)                                  {}

func newDispatcher(ctrl chan control.Command, host string, port int) (network.Dispatcher, error) {
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        return nil, fmt.Errorf("failed to resolve address: %w", err)
    }

    return &dispatcher{serverAddr: addr, ctrl: ctrl}, nil
}

// dispatcher - simplest dispatcher
type dispatcher struct {
    ctrl chan control.Command
    serverAddr *net.TCPAddr
    clientConn network.Connection
    serverConn network.Connection
}

func (d *dispatcher) Init(_ context.Context) error {
    return nil
}

func (d *dispatcher) RegisterNewConn(conn network.Connection) error {
    d.clientConn = conn

    tcpConn, err := net.DialTCP("tcp", nil, d.serverAddr)/**/
    if err != nil {
        d.signal(control.FAILED, err)

        return fmt.Errorf("failed to dial TCP: %w", err)
    }
    d.serverConn = network.NewConnection(tcpConn)

    return nil
}

func (d *dispatcher) HandleSPacket(conn network.Connection, packetBytes []byte) {

}

func (d *dispatcher) signal(state control.ComponentState, err error) {
    d.ctrl <- control.Command{
        Signal:    control.COMPONENT,
        Component: control.NETWORK,
        State:     state,
        Err:       err,
    }
}
