package player

import (
	"github.com/alexykot/cncraft/pkg/buffers"
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
	buffers.BufferPush
}

type PlayerInfoAddPlayer struct {
	Player entities.PlayerCharacter
}

func (p *PlayerInfoAddPlayer) Push(writer buffers.Buffer) {
	profile := p.Player.GetProfile()
	writer.PushUID(profile.UUID)
	writer.PushTxt(profile.Name)

	writer.PushVrI(int32(len(profile.Properties)))

	for _, prop := range profile.Properties {
		writer.PushTxt(prop.Name)
		writer.PushTxt(prop.Value)

		if prop.Signature == nil {
			writer.PushBit(false)
		} else {
			writer.PushBit(true)
			writer.PushTxt(*prop.Signature)
		}
	}

	writer.PushVrI(int32(p.Player.GetGameMode()))

	writer.PushVrI(0) // update this to the player's actual ping

	writer.PushBit(false) // update this to be whether the player has a custom display name or not, write that name as json if they do
}

type PlayerInfoUpdateLatency struct {}
