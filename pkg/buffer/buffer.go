package buffer

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/google/uuid"
)

type BPush interface {
	Push(writer B)
}

type BPull interface {
	Pull(reader B)
}

// TODO
//  - break this interface into parts
//  - make sure possible reading errors are checked and returned
type B interface {
	io.ReadWriter

	Len() int

	Bytes() []byte

	IndexI() int32

	IndexO() int32

	SkipAll()

	SkipLen(delta int)

	// pull
	PullBool() bool

	PullByte() byte

	PullInt16() int16

	PullUint16() uint16

	PullInt32() int32

	PullInt64() int64

	PullUint64() uint64

	PullFloat32() float32

	PullFloat64() float64

	PullVarInt() int32

	PullVarLong() int64

	PullString() string

	PullBytes() []byte

	PullUUID() uuid.UUID

	//PullNbt() *nbt.NbtCompound  // DEBT this is not going to work like that. Figure out idiomatic Go interface here.

	// push
	PushBool(data bool)

	PushByte(data byte)

	PushInt16(data int16)

	PushInt32(data int32)

	PushInt64(data int64)

	PushUint64(data uint64)

	PushFloat32(data float32)

	PushFloat64(data float64)

	PushVarInt(data int32)

	PushVarLong(data int64)

	PushString(data string)

	PushBytes(data []byte, prefixWithLen bool)

	PushUUID(data uuid.UUID)

	//PushNbt(data *tags.NbtCompound)  // DEBT this is not going to work like that. Figure out idiomatic Go interface here.
}

type buffer struct {
	iIndex int32 // TODO figure out what this means
	oIndex int32 // TODO figure out what this means

	bArray []byte
}

func (b *buffer) String() string {
	return fmt.Sprintf("Buffer[%d](i: %d, o: %d)%v", b.Len(), b.iIndex, b.oIndex, b.bArray)
}

// new
func New() B {
	return NewFrom(make([]byte, 0, 1024))
}

func NewFrom(bArray []byte) B {
	return &buffer{bArray: bArray}
}

// stdlib ReadWriter interface
func (b *buffer) Read(target []byte) (n int, err error) {
	if b.Len() <= len(target) {
		copy(target, b.bArray)
		return b.Len(), nil
	} else {
		copy(target, b.bArray[:len(target)])
		return len(target), nil
	}
}

func (b *buffer) Write(data []byte) (n int, err error) {
	b.pushNext(data...)
	return len(data), nil
}

// server_data
func (b *buffer) Len() int {
	return len(b.bArray)
}

func (b *buffer) Bytes() []byte {
	return b.bArray
}

func (b *buffer) IndexI() int32 {
	return b.iIndex
}

func (b *buffer) IndexO() int32 {
	return b.oIndex
}

func (b *buffer) SkipAll() {
	b.SkipLen(b.Len() - 1)
}

func (b *buffer) SkipLen(delta int) {
	b.iIndex = b.iIndex + int32(delta)
}

// pull
func (b *buffer) PullBool() bool {
	return b.pullNext() != 0
}

func (b *buffer) PullByte() byte {
	return b.pullNext()
}

func (b *buffer) PullInt16() int16 {
	return int16(binary.BigEndian.Uint16(b.pullSize(4)))
}

func (b *buffer) PullUint16() uint16 {
	return uint16(b.pullNext())<<8 | uint16(b.pullNext())
}

func (b *buffer) PullInt32() int32 {
	return int32(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullInt64() int64 {
	return int64(b.PullUint64())
}

func (b *buffer) PullUint64() uint64 {
	return binary.BigEndian.Uint64(b.pullSize(8))
}

func (b *buffer) PullFloat32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullFloat64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b.pullSize(8)))
}

func (b *buffer) PullVarInt() int32 {
	return int32(b.pullVariable(5))
}

func (b *buffer) PullVarLong() int64 {
	return b.pullVariable(10)
}

func (b *buffer) PullString() string {
	return string(b.PullBytes())
}

func (b *buffer) PullBytes() []byte {
	size := b.PullVarInt()
	array := b.bArray[b.iIndex : b.iIndex+size]

	b.iIndex += size

	return array
}

func (b *buffer) PullUUID() uuid.UUID {
	id, _ := bitsToUUID(b.PullInt64(), b.PullInt64())
	return id
}

// push
func (b *buffer) PushBool(data bool) {
	if data {
		b.pushNext(byte(0x01))
	} else {
		b.pushNext(byte(0x00))
	}
}

func (b *buffer) PushByte(data byte) {
	b.pushNext(data)
}

func (b *buffer) PushInt16(data int16) {
	b.pushNext(
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushInt32(data int32) {
	b.pushNext(
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushInt64(data int64) {
	b.pushNext(
		byte(data>>56),
		byte(data>>48),
		byte(data>>40),
		byte(data>>32),
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushUint64(data uint64) {
	b.pushNext(
		byte(data>>56),
		byte(data>>48),
		byte(data>>40),
		byte(data>>32),
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushFloat32(data float32) {
	b.PushInt32(int32(math.Float32bits(data)))
}

func (b *buffer) PushFloat64(data float64) {
	b.PushInt64(int64(math.Float64bits(data)))
}

func (b *buffer) PushVarInt(data int32) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *buffer) PushVarLong(data int64) {
	for {
		temp := data & 0x7F
		data >>= 7

		if data != 0 {
			temp |= 0x80
		}

		b.pushNext(byte(temp))

		if data == 0 {
			break
		}
	}
}

func (b *buffer) PushString(data string) {
	b.PushBytes([]byte(data), true)
}

func (b *buffer) PushBytes(data []byte, prefixWithLen bool) {
	if prefixWithLen {
		b.PushVarInt(int32(len(data)))
	}

	b.pushNext(data...)
}

func (b *buffer) PushUUID(data uuid.UUID) {
	msb, lsb := bitsFromUUID(data)

	b.PushInt64(msb)
	b.PushInt64(lsb)
}

func (b *buffer) pullNext() byte {
	if b.iIndex >= int32(b.Len()) {
		return 0
		// panic("reached end of buffer")
	}

	next := b.bArray[b.iIndex]
	b.iIndex++

	if b.oIndex > 0 {
		b.oIndex--
	}

	return next
}

func (b *buffer) pullSize(next int) []byte {
	bytes := make([]byte, next)

	for i := 0; i < next; i++ {
		bytes[i] = b.pullNext()
	}

	return bytes
}

func (b *buffer) pushNext(bArray ...byte) {
	b.oIndex += int32(len(bArray))
	b.bArray = append(b.bArray, bArray...)
}

func (b *buffer) pullVariable(max int) int64 {
	var num int
	var res int64

	for {
		tmp := int64(b.pullNext())
		res |= (tmp & 0x7F) << uint(num*7)

		if num++; num > max {
			panic("VarInt > " + strconv.Itoa(max))
		}

		if tmp&0x80 != 0x80 {
			break
		}
	}

	return res
}

func bitsToUUID(msb, lsb int64) (data uuid.UUID, err error) {
	mBytes := make([]byte, 8)
	lBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(mBytes, uint64(msb))
	binary.BigEndian.PutUint64(lBytes, uint64(lsb))

	return uuid.FromBytes(append(mBytes, lBytes...))
}

func bitsFromUUID(uuid uuid.UUID) (msb, lsb int64) {
	bytes, _ := uuid.MarshalBinary()

	msb = 0
	lsb = 0

	for i := 0; i < 8; i++ {
		msb = (msb << 0x08) | int64(bytes[i]&0xFF)
	}

	for i := 8; i < 16; i++ {
		lsb = (lsb << 0x08) | int64(bytes[i]&0xFF)
	}

	return
}
