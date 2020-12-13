package ents

import (
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/game"

	"github.com/golangmc/minecraft-server/impl/protocol/client"

	apisBase "github.com/golangmc/minecraft-server/apis/base"
	implBase "github.com/golangmc/minecraft-server/impl/base"
)

type player struct {
	entityLiving

	prof *game.Profile

	online bool

	conn implBase.Connection

	mode game.GameMode
}

func NewPlayer(prof *game.Profile, conn implBase.Connection) ents.Player {
	player := &player{
		prof:         prof,
		entityLiving: newEntityLiving(),
	}

	player.SetName(prof.Name)
	player.SetUUID(prof.UUID)

	player.SetConn(conn)

	return player
}

func (p *player) SendMessage(message ...interface{}) {
	packet := client.CPacketChatMessage{
		Message:         *msgs.New(apisBase.ConvertToString(message...)),
		MessagePosition: msgs.NormalChat,
	}

	p.conn.SendPacket(&packet)
}

func (p *player) GetGameMode() game.GameMode {
	return p.mode
}

func (p *player) SetGameMode(mode game.GameMode) {
	p.mode = mode
}

func (p *player) GetIsOnline() bool {
	return p.online
}

func (p *player) SetIsOnline(state bool) {
	p.online = state
}

func (p *player) GetProfile() *game.Profile {
	return p.prof
}

func (p *player) SetConn(conn implBase.Connection) {
	p.conn = conn
}
