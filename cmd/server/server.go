package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/golangmc/minecraft-server/apis"
	apisBase "github.com/golangmc/minecraft-server/apis/base"
	"github.com/golangmc/minecraft-server/apis/data/chat"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	implBase "github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/console"
	"github.com/golangmc/minecraft-server/impl/data/system"
	"github.com/golangmc/minecraft-server/impl/data/values"
	"github.com/golangmc/minecraft-server/impl/protocol"

	"github.com/golangmc/minecraft-server/pkg/bus"
	"github.com/golangmc/minecraft-server/pkg/entities"
	"github.com/golangmc/minecraft-server/server/state"
)

type Server interface {
	apisBase.State

	Logging() *zap.Logger

	PlayerList() []entities.Player

	GetPlayer(id uuid.UUID) (entities.Player, bool)

	Version() string

	Broadcast(message string) error
}

type server struct {
	message chan system.Message

	log    *zap.Logger
	pubsub bus.PubSub

	//network     implBase.Network
	packFactory protocol.PacketFactory
	players     map[uuid.UUID]entities.Player
}

// NewServer wires up and provides new server instance.
func NewServer(conf conf.ServerConfig) Server {
	message := make(chan bus.Envelope)

	join := make(chan implBase.PlayerAndConnection)
	quit := make(chan implBase.PlayerAndConnection)

	packetFactory := protocol.NewPacketFactory()
	ps := bus.New()

	// TODO not sure if it should be here or somewhere else. Probably there should be an "OnStart" or something.
	//  Maybe `Load()` is that. To review later.
	//state.RegisterHandlersState0(ps)
	//state.RegisterHandlersState1(ps)
	//state.RegisterHandlersState2(ps, join)
	//state.RegisterHandlersState3(ps, logger, tasking, join, quit)

	return &server{
		message: message,

		log:    logger,
		pubsub: ps,

		packFactory: packetFactory,
		//network:     conn.NewNetwork(conf.Network.Host, conf.Network.Port, packetFactory, message, join, quit),

		players: make(map[uuid.UUID]entities.Player),
	}
}

// ==== State ====
func (s *server) Load() {
	apis.SetMinecraftServer(s)

	go s.loadServer()

	s.wait()
}

func (s *server) Kill() {

	// push the stop message to the server exit channel
	s.message <- system.Make(system.STOP, "normal stop")
	close(s.message)

	s.log.Info(chat.DarkRed, "server stopped")
}

func (s *server) Log() *zap.Logger {
	return s.log
}

func (s *server) Watcher() bus.PubSub {
	return s.pubsub
}

func (s *server) Players() []entities.Player {
	if len(s.players) == 0 {
		return nil
	}

	playerList := make([]entities.Player, len(s.players), len(s.players))

	var i int
	for _, player := range s.players {
		playerList[i] = player
	}

	return playerList
}

func (s *server) GetConn(uuid uuid.UUID) implBase.Connection {
	return s.players.uuidToConn[uuid]
}

func (s *server) GetPlayer(uuid uuid.UUID) entities.Player {
	return s.players[uuid]
}

func (s *server) ServerVersion() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) Broadcast(message string) {
	for _, player := range s.Players() {
		player.SendMessage(message)
	}
}

// ==== server commands ====
func (s *server) broadcastCommand(sender entities.Sender, params []string) {
	message := strings.Join(params, " ")

	for _, player := range s.Players() {
		player.SendMessage(message)
	}
}

func (s *server) stopServerCommand(sender entities.Sender, params []string) {
	if _, ok := sender.(*console.Console); !ok {
		s.log.Error("non console sender %s tried to stop the server", zap.String("sender", sender.Name()))
		return
	}

	var after int64 = 0

	if len(params) > 0 {
		param, err := strconv.Atoi(params[0])

		if err != nil {
			panic(err)
		}

		if param <= 0 {
			panic(fmt.Errorf("value must be a positive whole number. [1..]"))
		}

		after = int64(param)
	}

	if after == 0 {

		s.Kill()

	} else {

		// inform future shutdown
		s.log.Warn(chat.Gold, "stopping server in ", chat.Green, util.FormatTime(after))

		// schedule shutdown {after} seconds later
		s.tasking.AfterTime(after, time.Second, func(task *task.Task) {
			s.Kill()
		})

	}
}

func (s *server) versionCommand(sender entities.Sender, params []string) {
	sender.SendMessage(s.ServerVersion())
}

// ==== internal ====
func (s *server) loadServer() {
	//s.console.Load()
	//s.command.Load()
	//s.tasking.Load()
	//s.network.Load()
	//
	//s.command.Register("vers", s.versionCommand)
	//s.command.Register("send", s.broadcastCommand)
	//s.command.Register("stop", s.stopServerCommand)
	//
	//s.pubsub.Subscribe(func(event apisEvent.PlayerJoinEvent) {
	//	s.log.InfoF("player %s logged in with uuid:%v", event.Player.Name(), event.Player.UUID())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has joined!", chat.Yellow, event.Player.Name())))
	//})
	//s.pubsub.Subscribe(func(event apisEvent.PlayerQuitEvent) {
	//	s.log.InfoF("%s disconnected!", event.Player.Name())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has left!", chat.Yellow, event.Player.Name())))
	//})
	//
	//s.pubsub.Subscribe(func(event implEvent.PlayerConnJoinEvent) {
	//	s.players.addData(event.Conn)
	//
	//	s.pubsub.Publish(apisEvent.PlayerJoinEvent{PlayerEvent: apisEvent.PlayerEvent{Player: event.Conn.Player}})
	//})
	//s.pubsub.Subscribe(func(event implEvent.PlayerConnQuitEvent) {
	//	player := s.players.playerByConn(event.Conn.Connection)
	//
	//	if player != nil {
	//		s.pubsub.Publish(apisEvent.PlayerQuitEvent{PlayerEvent: apisEvent.PlayerEvent{Player: player}})
	//	}
	//
	//	s.players.delData(event.Conn)
	//})
	//
	//s.pubsub.Subscribe(func(event implEvent.PlayerPluginMessagePullEvent) {
	//	s.log.DebugF("received message on channel '%s' from player %s:%s", event.Channel, event.Conn.Name(), event.Conn.UUID())
	//
	//	switch event.Channel {
	//	case plugin.CHANNEL_BRAND:
	//		s.log.DebugF("their client's brand is '%s'", event.Message.(*plugin.Brand).Name)
	//	}
	//})
}

func (s *server) wait() {
	// select over server commands channel
	select {
	case command := <-s.message:
		switch command.Command {
		// stop selecting when stop is received
		case system.STOP:
			return
		case system.FAIL:
			s.log.Error("internal server error: ", command.Message)
			s.log.Error("stopping server")
			return
		}
	}

	s.wait()
}
