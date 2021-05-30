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

	Receive(bufIn *buffer.Buffer) (len int, err error)
	Transmit(bufOut *buffer.Buffer) (len int, err error)

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

func (c *connection) Receive(bufIn *buffer.Buffer) (len int, err error) {
	data := make([]byte, 2097151) // max possible packet size according to https://wiki.vg/Protocol#Packet_format

	readLen, err := c.pull(data)
	if err != nil {
		return 0, err
	}
	if c.aes.enabled {
		data = c.aes.Decrypt(data[:readLen])
	}

	if c.zip.enabled {
		bufIn.PushBytes(c.zip.Inflate(data), false)
	} else {
		bufIn.PushBytes(data, false)
	}
	return readLen, nil
}

func (c *connection) Transmit(bufOut *buffer.Buffer) (len int, err error) {
	temp := buffer.New()
	temp.PushVarInt(getPacketLength(bufOut))
	temp.PushBytes(bufOut.Bytes(), false)
	return c.push(temp.Bytes())

	// TODO need to make it work without encryption/compression first
	// if c.zip.enabled {
	// 	deflated := buffer.New()
	// 	deflated.PushBytes(c.zip.Deflate(bufOut.Bytes()), false)
	// } else {
	// 	temp.PushBytes(bufOut.Bytes(), false)
	// }
	//
	// if c.aes.enabled {
	// 	return c.push(c.aes.Encrypt(temp.Bytes()))
	// } else {
	// 	return c.push(temp.Bytes())
	// }
}

func getPacketLength(bufOut *buffer.Buffer) int32 {
	length := bufOut.Len()
	// packetType := bufOut.PullVarInt()
	// if packetType == 0x20 {
	// 	length = length + 1
	// 	println(fmt.Sprintf("increased length to %d", length))
	// }

	return int32(length)
}

func (c *connection) pull(data []byte) (int, error) {
	readLen, err := c.tcp.Read(data)
	if err != nil {
		return 0, newNetworkError(ErrTCPReadFail, err)
	}
	return readLen, nil
}

func (c *connection) push(data []byte) (int, error) {
	wroteLen, err := c.tcp.Write(data)
	if err != nil {
		return 0, newNetworkError(ErrTCPWriteFail, err)
	}
	return wroteLen, nil
}
