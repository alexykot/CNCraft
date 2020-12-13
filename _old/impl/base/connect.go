package base

import (
	"net"

	"github.com/golangmc/minecraft-server/impl/protocol"
)

type Connection interface {
	ID() string
	Address() net.Addr

	GetState() protocol.State
	SetState(state protocol.State)

	Encrypt(data []byte) (output []byte)
	Decrypt(data []byte) (output []byte)

	CertifyName() string
	CertifyData() []byte
	CertifyValues(name string)
	CertifyUpdate(secret []byte)

	Deflate(data []byte) (output []byte)
	Inflate(data []byte) (output []byte)

	Pull(data []byte) (len int, err error)
	Push(data []byte) (len int, err error)

	Stop() (err error)

	SendPacket(packet protocol.CPacket)
}
