package conn

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"

	"github.com/golangmc/minecraft-server/apis/rand"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conn/crypto"
	"github.com/golangmc/minecraft-server/impl/protocol"
)

type connection struct {
	new bool
	tcp *net.TCPConn
	id  string

	state protocol.State

	certify Certify
	compact Compact
}

func NewConnection(conn *net.TCPConn) base.Connection {
	return &connection{
		new: true,
		tcp: conn,
		id:  uuid.New().String(),

		certify: Certify{},
		compact: Compact{},
	}
}

func (c *connection) Address() net.Addr {
	return c.tcp.RemoteAddr()
}

func (c *connection) ID() string {
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
	c.certify.data = rand.RandomByteArray(4)
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

func (c *connection) Stop() (err error) {
	err = c.tcp.Close()
	return
}

func (c *connection) SendPacket(packet protocol.CPacket) {
	bufO := NewBuffer()
	temp := NewBuffer()

	// write buffer
	bufO.PushVrI(int32(packet.ID()))
	packet.Push(bufO, c)

	temp.PushVrI(bufO.Len())
	temp.PushUAS(bufO.UAS(), false)

	_, _ = c.tcp.Write(c.Encrypt(temp.UAS()))
}
