package plugin

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/game/data"
)

type Message interface {
	Chan() Channel

	buffer.BufferPush
	buffer.BufferPull
}

func GetMessageForChannel(channel Channel) Message {
	switch channel {
	case ChannelBrand:
		return &Brand{}
	case ChannelDebugPaths:
		return &DebugPaths{}
	case ChannelDebugNeighbors:
		return &DebugNeighbors{}
	}
	return nil
}

type Channel string

const (
	ChannelBrand          Channel = "minecraft:brand"
	ChannelDebugPaths     Channel = "minecraft:debug/paths"
	ChannelDebugNeighbors Channel = "minecraft:debug/neighbors_update"
)

// look, they're like cute little packets :D

type Brand struct {
	Name string
}

func (b *Brand) Chan() Channel {
	return ChannelBrand
}

func (b *Brand) Push(writer buffer.B) {
	writer.PushTxt(b.Name)
}

func (b *Brand) Pull(reader buffer.B) {
	b.Name = reader.PullTxt()
}

type DebugPaths struct { // unused? honestly why did I do this
	UnknownValue1 int32
	UnknownValue2 float32
	PathEntity    PathEntity
}

type PathEntity struct {
	Index   int
	Target  PathPoint
	PSetLen int
	PSet    []PathPoint
	OSetLen int
	OSet    []PathPoint
	CSetLen int
	CSet    []PathPoint
}

func (p *PathEntity) Push(writer buffer.B) {
	writer.PushI32(int32(p.Index))

	p.Target.Push(writer)

	writer.PushI32(int32(p.PSetLen))
	for _, point := range p.PSet {
		point.Push(writer)
	}

	writer.PushI32(int32(p.OSetLen))
	for _, point := range p.OSet {
		point.Push(writer)
	}

	writer.PushI32(int32(p.CSetLen))
	for _, point := range p.CSet {
		point.Push(writer)
	}
}

func (p *PathEntity) Pull(reader buffer.B) {
	p.Index = int(reader.PullI32())

	target := PathPoint{}
	target.Pull(reader)

	p.Target = target

	p.PSet = make([]PathPoint, 0)
	p.PSetLen = int(reader.PullI32())

	for i := 0; i < p.PSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.PSet = append(p.PSet, point)
	}

	p.OSet = make([]PathPoint, 0)
	p.OSetLen = int(reader.PullI32())

	for i := 0; i < p.OSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.OSet = append(p.OSet, point)
	}

	p.CSet = make([]PathPoint, 0)
	p.CSetLen = int(reader.PullI32())

	for i := 0; i < p.CSetLen; i++ {
		point := PathPoint{}
		point.Pull(reader)

		p.CSet = append(p.CSet, point)
	}
}

type PathPoint struct {
	X int32
	Y int32
	Z int32

	DistanceOrigin float32
	Cost           float32
	CostMalus      float32
	Visited        bool
	NodeType       NodeType
	DistanceTarget float32
}

func (p *PathPoint) Push(writer buffer.B) {
	writer.PushI32(p.X)
	writer.PushI32(p.Y)
	writer.PushI32(p.Z)
	writer.PushF32(p.DistanceOrigin)
	writer.PushF32(p.Cost)
	writer.PushF32(p.CostMalus)
	writer.PushBit(p.Visited)
	writer.PushI32(int32(p.NodeType))
	writer.PushF32(p.DistanceTarget)
}

func (p *PathPoint) Pull(reader buffer.B) {
	p.X = reader.PullI32()
	p.Y = reader.PullI32()
	p.Z = reader.PullI32()
	p.DistanceOrigin = reader.PullF32()
	p.Cost = reader.PullF32()
	p.CostMalus = reader.PullF32()
	p.Visited = reader.PullBit()
	p.NodeType = NodeType(reader.PullI32())
	p.DistanceTarget = reader.PullF32()
}

type NodeType int

const (
	BLOCKED NodeType = iota
	OPEN
	WALKABLE
	TRAPDOOR
	FENCE
	LAVA
	WATER
	RAIL
	DANGER_FIRE
	DAMAGE_FIRE
	DANGER_CACTUS
	DAMAGE_CACTUS
	DANGER_OTHER
	DAMAGE_OTHER
	DOOR_OPEN
	DOOR_WOOD_CLOSED
	DOOR_IRON_CLOSED
)

func (d *DebugPaths) Chan() Channel {
	return ChannelDebugPaths
}

func (d *DebugPaths) Push(writer buffer.B) {
	writer.PushI32(d.UnknownValue1)
	writer.PushF32(d.UnknownValue2)
	d.PathEntity.Push(writer)
}

func (d *DebugPaths) Pull(reader buffer.B) {
	d.UnknownValue1 = reader.PullI32()
	d.UnknownValue2 = reader.PullF32()

	entity := PathEntity{}
	entity.Pull(reader)

	d.PathEntity = entity
}

type DebugNeighbors struct {
	Time     int64
	Location data.PositionI
}

func (d *DebugNeighbors) Chan() Channel {
	return ChannelDebugNeighbors
}

func (d *DebugNeighbors) Push(writer buffer.B) {
	writer.PushVrL(d.Time)
	d.Location.Push(writer)
}

func (d *DebugNeighbors) Pull(reader buffer.B) {
	d.Time = reader.PullVrL()
	d.Location.Pull(reader)
}