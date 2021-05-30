package buffer

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
)

type BPush interface {
	Push(writer *Buffer)
}

type BPull interface {
	Pull(reader *Buffer)
}

type Buffer struct {
	iIndex int32
	oIndex int32

	bArray []byte
}

func (b *Buffer) String() string {
	return fmt.Sprintf("Buffer[%d](i: %d, o: %d)%v", b.Len(), b.iIndex, b.oIndex, b.bArray)
}

func New() *Buffer {
	return NewFrom(make([]byte, 0, 1024*1024))
}

func NewFrom(bArray []byte) *Buffer {
	return &Buffer{bArray: bArray}
}

// stdlib ReadWriter interface
func (b *Buffer) Read(target []byte) (n int, err error) {
	if b.Len() <= len(target) {
		copy(target, b.bArray)
		return b.Len(), nil
	} else {
		copy(target, b.bArray[:len(target)])
		return len(target), nil
	}
}

func (b *Buffer) Write(data []byte) (n int, err error) {
	b.pushNext(data...)
	return len(data), nil
}

// server_data
func (b *Buffer) Len() int {
	return len(b.bArray)
}

func (b *Buffer) Bytes() []byte {
	return b.bArray
}

func (b *Buffer) IndexI() int32 {
	return b.iIndex
}

func (b *Buffer) IndexO() int32 {
	return b.oIndex
}

func (b *Buffer) SkipAll() {
	b.SkipLen(int32(b.Len() - 1))
}

func (b *Buffer) SkipLen(delta int32) {
	b.iIndex = b.iIndex + delta
}

// pull
func (b *Buffer) PullBool() bool {
	return b.pullNext() != 0
}

func (b *Buffer) PullByte() byte {
	return b.pullNext()
}

func (b *Buffer) PullInt16() int16 {
	return int16(binary.BigEndian.Uint16(b.pullSize(2)))
}

func (b *Buffer) PullUint16() uint16 {
	return uint16(b.pullNext())<<8 | uint16(b.pullNext())
}

func (b *Buffer) PullInt32() int32 {
	return int32(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *Buffer) PullInt64() int64 {
	return int64(b.PullUint64())
}

func (b *Buffer) PullUint64() uint64 {
	return binary.BigEndian.Uint64(b.pullSize(8))
}

func (b *Buffer) PullFloat32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *Buffer) PullFloat64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b.pullSize(8)))
}

func (b *Buffer) PullVarInt() int32 {
	return int32(b.pullVariable(5))
}

func (b *Buffer) PullVarLong() int64 {
	return b.pullVariable(10)
}

func (b *Buffer) PullString() string {
	return string(b.PullBytes())
}

func (b *Buffer) PullBytes() []byte {
	size := b.PullVarInt()
	array := b.bArray[b.iIndex : b.iIndex+size]

	b.iIndex += size

	return array
}

func (b *Buffer) PullUUID() uuid.UUID {
	id, _ := bitsToUUID(b.PullInt64(), b.PullInt64())
	return id
}

// push
func (b *Buffer) PushBool(data bool) {
	if data {
		b.pushNext(byte(0x01))
	} else {
		b.pushNext(byte(0x00))
	}
}

func (b *Buffer) PushByte(data byte) {
	b.pushNext(data)
}

func (b *Buffer) PushInt16(data int16) {
	b.pushNext(
		byte(data>>8),
		byte(data))
}

func (b *Buffer) PushInt32(data int32) {
	b.pushNext(
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *Buffer) PushInt64(data int64) {
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

func (b *Buffer) PushUint16(data uint16) {
	b.pushNext(
		byte(data>>8),
		byte(data))
}

func (b *Buffer) PushUint64(data uint64) {
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

func (b *Buffer) PushFloat32(data float32) {
	b.PushInt32(int32(math.Float32bits(data)))
}

func (b *Buffer) PushFloat64(data float64) {
	b.PushInt64(int64(math.Float64bits(data)))
}

func (b *Buffer) PushVarInt(data int32) {
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

func (b *Buffer) PushVarLong(data int64) {
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

func (b *Buffer) PushString(data string) {
	b.PushBytes([]byte(data), true)
}

func (b *Buffer) PushBytes(data []byte, prefixWithLen bool) {
	if prefixWithLen {
		b.PushVarInt(int32(len(data)))
	}

	b.pushNext(data...)
}

func (b *Buffer) PushUUID(data uuid.UUID) {
	msb, lsb := bitsFromUUID(data)

	b.PushInt64(msb)
	b.PushInt64(lsb)
}

func (b *Buffer) pullNext() byte {
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

func (b *Buffer) pullSize(next int) []byte {
	bytes := make([]byte, next)

	for i := 0; i < next; i++ {
		bytes[i] = b.pullNext()
	}

	return bytes
}

func (b *Buffer) pushNext(bArray ...byte) {
	b.oIndex += int32(len(bArray))
	b.bArray = append(b.bArray, bArray...)
}

func (b *Buffer) pullVariable(max int) int64 {
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
