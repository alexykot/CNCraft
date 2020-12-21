package buffer

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/data/tags"
)

/*
 Language used:
	- Len = Length
	- Arr = Array

	- Bit = Boolean
	- Byt = Byte
	- Int = int
	- VrI = VarInt
	- Srt = Short
	- Txt = String
*/

// TODO
//  - break this interface into parts
//  - make sure possible reading errors are checked and returned
//  - rename methods for clarity and get rid of the special dictionary
type B interface {
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
	Push(writer B)
}

type BufferPull interface {
	Pull(reader B)
}

type BufferCodec interface {
	BufferPush
	BufferPull
}

type buffer struct {
	iIndex int32
	oIndex int32

	bArray []byte
}

func (b *buffer) String() string {
	return fmt.Sprintf("Buffer[%d](i: %d, o: %d)%v", b.Len(), b.iIndex, b.oIndex, b.bArray)
}

// new
func New() B {
	return NewFrom(make([]byte, 0))
}

func NewFrom(bArray []byte) B {
	return &buffer{bArray: bArray}
}

// server_data
func (b *buffer) Len() int32 {
	return int32(len(b.bArray))
}

func (b *buffer) SAS() []int8 {
	return asSArray(b.bArray)
}

func (b *buffer) UAS() []byte {
	return b.bArray
}

func (b *buffer) InI() int32 {
	return b.iIndex
}

func (b *buffer) InO() int32 {
	return b.oIndex
}

func (b *buffer) SkpAll() {
	b.SkpLen(b.Len() - 1)
}

func (b *buffer) SkpLen(delta int32) {
	b.iIndex += delta
}

// pull
func (b *buffer) PullBit() bool {
	return b.pullNext() != 0
}

func (b *buffer) PullByt() byte {
	return b.pullNext()
}

func (b *buffer) PullI16() int16 {
	return int16(binary.BigEndian.Uint16(b.pullSize(4)))
}

func (b *buffer) PullU16() uint16 {
	return uint16(b.pullNext())<<8 | uint16(b.pullNext())
}

func (b *buffer) PullI32() int32 {
	return int32(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullI64() int64 {
	return int64(b.PullU64())
}

func (b *buffer) PullU64() uint64 {
	return binary.BigEndian.Uint64(b.pullSize(8))
}

func (b *buffer) PullF32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b.pullSize(4)))
}

func (b *buffer) PullF64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b.pullSize(8)))
}

func (b *buffer) PullVrI() int32 {
	return int32(b.pullVariable(5))
}

func (b *buffer) PullVrL() int64 {
	return b.pullVariable(10)
}

func (b *buffer) PullTxt() string {
	return string(b.PullUAS())
}

func (b *buffer) PullUAS() []byte {
	size := b.PullVrI()
	array := b.bArray[b.iIndex : b.iIndex+size]

	b.iIndex += size

	return array
}

func (b *buffer) PullSAS() []int8 {
	return asSArray(b.PullUAS())
}

func (b *buffer) PullUID() uuid.UUID {
	id, _ := bitsToUUID(b.PullI64(), b.PullI64())
	return id
}

func (b *buffer) PullPos() data.PositionI {
	val := b.PullU64()

	x := int64(val) >> 38
	y := int64(val) & 0xFFF
	z := int64(val) << 26 >> 38

	return data.PositionI{
		X: x,
		Y: y,
		Z: z,
	}
}

func (b *buffer) PullNbt() *tags.NbtCompound {
	typ := tags.Typ(b.PullByt())

	fmt.Println("==type")
	fmt.Println(typ)

	if typ != tags.TAG_Compound {
		panic("root tag must be compound") // probably shouldn't panic?
	}

	name := b.PullTxt()
	if len(name) != 0 {
		panic("root compound should have an empty name")
	}

	fmt.Println("==name")
	fmt.Println(name)

	tag := &tags.NbtCompound{}
	b.pullNbt(tag)

	return tag
}

// push
func (b *buffer) PushBit(data bool) {
	if data {
		b.pushNext(byte(0x01))
	} else {
		b.pushNext(byte(0x00))
	}
}

func (b *buffer) PushByt(data byte) {
	b.pushNext(data)
}

func (b *buffer) PushI16(data int16) {
	b.pushNext(
		byte(data)>>8,
		byte(data))
}

func (b *buffer) PushI32(data int32) {
	b.pushNext(
		byte(data>>24),
		byte(data>>16),
		byte(data>>8),
		byte(data))
}

func (b *buffer) PushI64(data int64) {
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

func (b *buffer) PushF32(data float32) {
	b.PushI32(int32(math.Float32bits(data)))
}

func (b *buffer) PushF64(data float64) {
	b.PushI64(int64(math.Float64bits(data)))
}

func (b *buffer) PushVrI(data int32) {
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

func (b *buffer) PushVrL(data int64) {
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

func (b *buffer) PushTxt(data string) {
	b.PushUAS([]byte(data), true)
}

func (b *buffer) PushUAS(data []byte, prefixWithLen bool) {
	if prefixWithLen {
		b.PushVrI(int32(len(data)))
	}

	b.pushNext(data...)
}

func (b *buffer) PushSAS(data []int8, prefixWithLen bool) {
	b.PushUAS(asUArray(data), prefixWithLen)
}

func (b *buffer) PushUID(data uuid.UUID) {
	msb, lsb := bitsFromUUID(data)

	b.PushI64(msb)
	b.PushI64(lsb)
}

func (b *buffer) PushPos(data data.PositionI) {
	b.PushI64(((data.X & 0x3FFFFFF) << 38) | ((data.Z & 0x3FFFFFF) << 12) | (data.Y & 0xFFF))
}

func (b *buffer) PushNbt(data *tags.NbtCompound) {
	if data == nil {
		b.PushByt(0)
	} else {
		b.PushByt(byte(data.Type()))

		b.pushNext(0, 0)

		b.pushNbt(data)
	}
}

func (b *buffer) pullNext() byte {

	if b.iIndex >= b.Len() {
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

func asSArray(bytes []byte) []int8 {
	array := make([]int8, 0)

	for _, b := range bytes {
		array = append(array, int8(b))
	}

	return array
}

func asUArray(bytes []int8) []byte {
	array := make([]byte, 0)

	for _, b := range bytes {
		array = append(array, byte(b))
	}

	return array
}

var typeToInst = map[tags.Typ]func() tags.Nbt{
	tags.TAG_End: func() tags.Nbt {
		return &tags.NbtEnd{}
	},
	tags.TAG_Byte: func() tags.Nbt {
		return &tags.NbtByt{}
	},
	tags.TAG_Short: func() tags.Nbt {
		return &tags.NbtI16{}
	},
	tags.TAG_Int: func() tags.Nbt {
		return &tags.NbtI32{}
	},
	tags.TAG_Long: func() tags.Nbt {
		return &tags.NbtI64{}
	},
	tags.TAG_Float: func() tags.Nbt {
		return &tags.NbtF32{}
	},
	tags.TAG_Double: func() tags.Nbt {
		return &tags.NbtF64{}
	},
	tags.TAG_Byte_Array: func() tags.Nbt {
		return &tags.NbtArrByt{}
	},
	tags.TAG_String: func() tags.Nbt {
		return &tags.NbtTxt{}
	},
	tags.TAG_List: func() tags.Nbt {
		return &tags.NbtArrAny{}
	},
	tags.TAG_Compound: func() tags.Nbt {
		return &tags.NbtCompound{}
	},
	tags.TAG_Int_Array: func() tags.Nbt {
		return &tags.NbtArrI32{}
	},
	tags.TAG_Long_Array: func() tags.Nbt {
		return &tags.NbtArrI64{}
	},
}

func (b *buffer) pullNbt(data tags.Nbt) {
	switch data.Type() {
	case tags.TAG_End:
		// nothing
		break
	case tags.TAG_Byte:
		data.(*tags.NbtByt).Value = int8(b.PullByt())
		break
	case tags.TAG_Short:
		panic("unimplemented")
		break
	case tags.TAG_Int:
		panic("unimplemented")
		break
	case tags.TAG_Long:
		panic("unimplemented")
		break
	case tags.TAG_Float:
		panic("unimplemented")
		break
	case tags.TAG_Double:
		panic("unimplemented")
		break
	case tags.TAG_Byte_Array:
		panic("unimplemented")
		break
	case tags.TAG_String:
		panic("unimplemented")
		break
	case tags.TAG_List:
		panic("unimplemented")
		break
	case tags.TAG_Compound:
		value := make(map[string]tags.Nbt)

		fmt.Println("reading compound")

		for {
			typ := tags.Typ(b.PullByt())
			if typ == tags.TAG_End {
				fmt.Println("encountered break")
				break
			}

			fmt.Println("===type")
			fmt.Println(typ)

			name := b.pullNbtTxt()

			fmt.Println("===name")
			fmt.Println(name)

			inst := typeToInst[typ]()
			b.pullNbt(inst)

			value[name] = inst
		}

		data.(*tags.NbtCompound).Value = value
		break
	case tags.TAG_Int_Array:
		panic("unimplemented")
		break
	case tags.TAG_Long_Array:
		value := make([]int64, b.PullI32())

		for i := 0; i < len(value); i++ {
			value[i] = b.PullI64()
		}

		data.(*tags.NbtArrI64).Value = value
		break
	}
}

func (b *buffer) pushNbt(data tags.Nbt) {
	switch data.Type() {
	case tags.TAG_End:
		// nothing
		break
	case tags.TAG_Byte:
		b.PushByt(byte(data.(*tags.NbtByt).Value))
		break
	case tags.TAG_Short:
		panic("unimplemented")
		break
	case tags.TAG_Int:
		panic("unimplemented")
		break
	case tags.TAG_Long:
		panic("unimplemented")
		break
	case tags.TAG_Float:
		panic("unimplemented")
		break
	case tags.TAG_Double:
		panic("unimplemented")
		break
	case tags.TAG_Byte_Array:
		panic("unimplemented")
		break
	case tags.TAG_String:
		panic("unimplemented")
		break
	case tags.TAG_List:
		panic("unimplemented")
		break
	case tags.TAG_Compound:
		for name, tag := range data.(*tags.NbtCompound).Value {
			b.PushByt(byte(tag.Type()))

			if tag.Type() == tags.TAG_End {
				continue
			}

			b.pushNbtTxt(name)
			b.pushNbt(tag)
		}

		b.PushByt(0)
		break
	case tags.TAG_Int_Array:
		panic("unimplemented")
		break
	case tags.TAG_Long_Array:
		value := data.(*tags.NbtArrI64).Value

		b.PushI32(int32(len(value)))

		for _, value := range value {
			b.PushI64(value)
		}

		break
	}
}

func (b *buffer) pullNbtTxt() string {
	size := b.PullI16()
	data := b.pullSize(int(size))

	return string(data)
}

func (b *buffer) pushNbtTxt(data string) {
	b.PushI16(int16(len(data)))
	b.PushUAS([]byte(data), false)
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
