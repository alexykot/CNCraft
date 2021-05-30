package player

import (
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/game/entities"
)

type PlayerInfoAction int32

const (
	AddPlayer PlayerInfoAction = iota
	UpdateGameMode
	UpdateLatency
	UpdateDisplayName
	RemovePlayer
)

type PlayerInfo interface {
	buffer.BPush
}

type PlayerInfoAddPlayer struct {
	Player entities.PlayerCharacter
}

func (p *PlayerInfoAddPlayer) Push(writer *buffer.Buffer) {
	//profile := p.Player.GetProfile()
	//writer.PushUUID(profile.UUID)
	//writer.PushString(profile.Name)
	//
	//writer.PushVarInt(int32(len(profile.Properties)))
	//
	//for _, prop := range profile.Properties {
	//	writer.PushString(prop.Name)
	//	writer.PushString(prop.Value)
	//
	//	if prop.Signature == nil {
	//		writer.PushBool(false)
	//	} else {
	//		writer.PushBool(true)
	//		writer.PushString(*prop.Signature)
	//	}
	//}
	//
	//writer.PushVarInt(int32(p.Player.GetGameMode()))
	//
	//writer.PushVarInt(0) // update this to the player's actual ping
	//
	//writer.PushBool(false) // update this to be whether the player has a custom display name or not, write that name as json if they do
}

type PlayerInfoUpdateLatency struct{}
