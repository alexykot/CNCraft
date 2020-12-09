package core

import (
	"net"

	"github.com/alexykot/cncraft/pkg/entities"
	"github.com/alexykot/cncraft/pkg/protocol"
)

type User struct {
	PC entities.PlayerCharacter
	Username string
}

type Connection interface {
	ID() string
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
