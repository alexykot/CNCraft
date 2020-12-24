package network

import (
	"fmt"
	"net"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/protocol"
)

type Connection interface {
	ID() uuid.UUID
	Address() net.Addr

	GetState() protocol.State
	SetState(protocol.State)

	EnableEncryption(secret []byte) error
	EnableCompression()

	Receive(bufIn buffer.B) (len int, err error)
	Transmit(bufOut buffer.B) (len int, err error)

	Close() error
}

type connection struct {
	tcp *net.TCPConn
	id  uuid.UUID

	state protocol.State

	aes crypter
	zip compressor
}

func NewConnection(conn *net.TCPConn) Connection {
	return &connection{
		tcp: conn,
		id:  uuid.New(),

		aes: crypter{},
		zip: compressor{},
	}
}

func (c *connection) Address() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *connection) ID() uuid.UUID {
	return c.id
}

func (c *connection) GetState() protocol.State {
	return c.state
}

func (c *connection) SetState(state protocol.State) {
	c.state = state
}

// Close closes underlying TCP connection.
func (c *connection) Close() (err error) {
	return c.tcp.Close()
}

func (c *connection) EnableEncryption(secret []byte) error {
	if err := c.aes.Enable(secret); err != nil {
		return fmt.Errorf("failed to enable AES encryption: %w", err)
	}
	return nil
}
func (c *connection) EnableCompression() {
	c.zip.Enable()
}

func (c *connection) Receive(bufIn buffer.B) (len int, err error) {
	data := make([]byte, 1024)

	readLen, err := c.pull(data)
	if err != nil {
		return 0, err
	}
	if c.aes.enabled {
		data = c.aes.Decrypt(data[:readLen])
	}

	if c.zip.enabled {
		bufIn.PushUAS(c.zip.Inflate(data), false)
	} else {
		bufIn.PushUAS(data, false)
	}
	return readLen, nil
}

func (c *connection) Transmit(bufOut buffer.B) (len int, err error) {
	temp := buffer.New()
	temp.PushVrI(bufOut.Len())

	if c.zip.enabled {
		deflated := buffer.New()
		deflated.PushUAS(c.zip.Deflate(bufOut.UAS()), false)
	} else {
		temp.PushUAS(bufOut.UAS(), false)
	}

	if c.aes.enabled {
		return c.push(c.aes.Encrypt(temp.UAS()))
	} else {
		return c.push(temp.UAS())
	}
}

func (c *connection) pull(data []byte) (int, error) {
	readLen, err := c.tcp.Read(data)
	if err != nil {
		return 0, fmt.Errorf("failed to read from TCP: %w", err)
	}
	return readLen, nil
}

func (c *connection) push(data []byte) (int, error) {
	wroteLen, err := c.tcp.Write(data)
	if err != nil {
		return 0, fmt.Errorf("failed to write to TCP: %w", err)
	}
	return wroteLen, nil
}
