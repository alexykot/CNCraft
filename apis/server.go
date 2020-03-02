package apis

import (
	"sync"

	apisBase "github.com/golangmc/minecraft-server/apis/base"
	"github.com/golangmc/minecraft-server/apis/cmds"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/uuid"
	implBase "github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

type Server interface {
	apisBase.State

	Logging() *logs.Logging

	Command() *cmds.CommandManager

	Tasking() *task.Tasking

	Watcher() pubsub.PubSub

	Players() []ents.Player

	ConnByUUID(uuid uuid.UUID) implBase.Connection

	PlayerByUUID(uuid uuid.UUID) ents.Player

	PlayerByConn(conn implBase.Connection) ents.Player

	ServerVersion() string

	Broadcast(message string)
}

var instance *Server
var syncOnce sync.Once

func MinecraftServer() Server {
	if instance == nil {
		panic("server is unavailable")
	}

	return *instance
}

func SetMinecraftServer(server Server) {
	syncOnce.Do(func() {
		instance = &server
	})
}
