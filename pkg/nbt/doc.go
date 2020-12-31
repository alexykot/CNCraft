// Package nbt implements an encoding/decoding library for Named Binary Tag format  used by the Minecraft protocol.
// It provides API similar to "encoding/xml" package.
//
// This is a vendored in copy of the github.com/tnze/go-mc/nbt library.
//
// NBT tag type will be derived from the data type of the supplied field, following this mapping:
// struct, interface     => TagCompound
// uint8                 => TagByte
// int16, uint16         => TagShort
// int32, uint32         => TagInt
// float32               => TagFloat
// int64, uint64         => TagLong
// float64               => TagDouble
// string                => TagString
// []uint8               => TagByteArray
// []int32               => TagIntArray
// []int64               => TagLongArray
// []AnyOtherType        => TagList
// if none above matched => TagNone

package nbt
