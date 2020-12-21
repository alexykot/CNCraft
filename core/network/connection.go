package network

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/network/crypto"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/protocol"
)

type Connection interface {
	ID() uuid.UUID
	Address() net.Addr

	GetState() protocol.State
	SetState(protocol.State)

	Encrypt([]byte) []byte
	Decrypt([]byte) []byte

	CertifyName() string
	CertifyData() []byte
	CertifyValues(name string)
	CertifyUpdate(secret []byte)

	Deflate([]byte) []byte
	Inflate([]byte) []byte

	Pull(data []byte) (len int, err error)
	Push(data []byte) (len int, err error)

	Close() error

	SendPacket(protocol.CPacket)
}

type connection struct {
	tcp *net.TCPConn
	id  uuid.UUID

	state protocol.State

	certify Certify
	compact Compact
}

func NewConnection(conn *net.TCPConn) Connection {
	return &connection{
		tcp: conn,
		id:  uuid.New(),

		certify: Certify{},
		compact: Compact{},
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

type Certify struct {
	name string

	used bool
	data []byte

	encrypt cipher.Stream
	decrypt cipher.Stream
}

func (c *connection) Encrypt(data []byte) (output []byte) {
	if !c.certify.used {
		return data
	}

	output = make([]byte, len(data))
	c.certify.encrypt.XORKeyStream(output, data)

	return
}

func (c *connection) Decrypt(data []byte) (output []byte) {
	if !c.certify.used {
		return data
	}

	output = make([]byte, len(data))
	c.certify.decrypt.XORKeyStream(output, data)

	return
}

func (c *connection) CertifyName() string {
	return c.certify.name
}

func (c *connection) CertifyData() []byte {
	return c.certify.data
}

func (c *connection) CertifyUpdate(secret []byte) {
	encrypt, decrypt, err := crypto.NewEncryptAndDecrypt(secret)

	c.certify.encrypt = encrypt
	c.certify.decrypt = decrypt

	if err != nil {
		panic(fmt.Errorf("failed to enable encryption for user: %s\n%v", c.CertifyName(), err))
	}

	c.certify.used = true
	c.certify.data = secret
}

func (c *connection) CertifyValues(name string) {
	c.certify.name = name
	c.certify.data = randomByteArray(4)
}

type Compact struct {
	used bool
	size int32
}

func (c *connection) Deflate(data []byte) (output []byte) {
	if !c.compact.used {
		return data
	}

	var out bytes.Buffer

	writer, _ := zlib.NewWriterLevel(&out, zlib.BestCompression)
	_, _ = writer.Write(data)
	_ = writer.Close()

	output = out.Bytes()

	return
}

func (c *connection) Inflate(data []byte) (output []byte) {
	if !c.compact.used {
		return data
	}

	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	_, _ = io.Copy(&out, reader)

	output = out.Bytes()

	return
}

func (c *connection) Pull(data []byte) (len int, err error) {
	len, err = c.tcp.Read(data)
	return
}

func (c *connection) Push(data []byte) (len int, err error) {
	len, err = c.tcp.Write(data)
	return
}

func (c *connection) Close() (err error) {
	err = c.tcp.Close()
	return
}

func (c *connection) SendPacket(packet protocol.CPacket) {
	bufO := buffer.New()
	temp := buffer.New()

	// write buffer
	bufO.PushVrI(int32(packet.ID()))
	packet.Push(bufO)

	temp.PushVrI(bufO.Len())
	temp.PushUAS(bufO.UAS(), false)

	_, _ = c.tcp.Write(c.Encrypt(temp.UAS()))
}

func randomByteArray(len int) []byte {
	array := make([]byte, len)
	_, _ = rand.Read(array)

	return array
}
