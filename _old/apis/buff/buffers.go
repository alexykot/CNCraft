package buff

import (
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/data/tags"
	"github.com/golangmc/minecraft-server/apis/uuid"
)

// DEBT break this interface into parts, make sure possible reading errors are checked and returned
type Buffer interface {
	Len() int32

	SAS() []int8

	UAS() []byte

	InI() int32

	InO() int32

	SkpAll()

	SkpLen(delta int32)

	// pull
	PullBit() bool

	PullByt() byte

	PullI16() int16

	PullU16() uint16

	PullI32() int32

	PullI64() int64

	PullU64() uint64

	PullF32() float32

	PullF64() float64

	PullVrI() int32

	PullVrL() int64

	PullTxt() string

	PullUAS() []byte

	PullSAS() []int8

	PullUID() uuid.UUID

	PullPos() data.PositionI

	PullNbt() *tags.NbtCompound

	// push
	PushBit(data bool)

	PushByt(data byte)

	PushI16(data int16)

	PushI32(data int32)

	PushI64(data int64)

	PushF32(data float32)

	PushF64(data float64)

	PushVrI(data int32)

	PushVrL(data int64)

	PushTxt(data string)

	PushUAS(data []byte, prefixWithLen bool)

	PushSAS(data []int8, prefixWithLen bool)

	PushUID(data uuid.UUID)

	PushPos(data data.PositionI)

	PushNbt(data *tags.NbtCompound)
}

type BufferPush interface {
	Push(writer Buffer)
}

type BufferPull interface {
	Pull(reader Buffer)
}

type BufferCodec interface {
	BufferPush
	BufferPull
}
