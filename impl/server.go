package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golangmc/minecraft-server/apis"
	apisBase "github.com/golangmc/minecraft-server/apis/base"
	"github.com/golangmc/minecraft-server/apis/cmds"
	"github.com/golangmc/minecraft-server/apis/data/chat"
	"github.com/golangmc/minecraft-server/apis/ents"
	apisEvent "github.com/golangmc/minecraft-server/apis/game/event"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/apis/uuid"
	implBase "github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/conf"
	"github.com/golangmc/minecraft-server/impl/conn"
	"github.com/golangmc/minecraft-server/impl/console"
	"github.com/golangmc/minecraft-server/impl/data/plugin"
	"github.com/golangmc/minecraft-server/impl/data/system"
	"github.com/golangmc/minecraft-server/impl/data/values"
	implEvent "github.com/golangmc/minecraft-server/impl/game/event"
	"github.com/golangmc/minecraft-server/impl/game/state"
	"github.com/golangmc/minecraft-server/impl/protocol"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

type server struct {
	message chan system.Message

	console *console.Console

	logging *logs.Logging
	tasking *task.Tasking
	pubsub  pubsub.PubSub

	command *cmds.CommandManager

	network     implBase.Network
	packFactory protocol.PacketFactory

	players *playerAssociation
}

// ==== new ====
func NewServer(conf conf.ServerConfig) apis.Server {
	message := make(chan system.Message)
	tasking := task.NewTasking(values.MPT)

	join := make(chan implBase.PlayerAndConnection)
	quit := make(chan implBase.PlayerAndConnection)

	packetFactory := protocol.NewPacketFactory(tasking, join, quit)
	pubsub := pubsub.New()
	logger := logs.NewLogging("server", logs.EveryLevel...)

	// TODO not sure if it should be here or somewhere else. Probably there should be an "OnStart" or something.
	//  Maybe `Load()` is that. To review later.
	state.HandleState0(pubsub)
	state.HandleState1(pubsub)
	state.HandleState2(pubsub, join)
	state.HandleState3(pubsub, logger, tasking, join, quit)

	return &server{
		message: message,

		console: console.NewConsole(message),

		logging: logger,
		tasking: tasking,
		pubsub:  pubsub,

		command: cmds.NewCommandManager(),

		packFactory: packetFactory,
		network:     conn.NewNetwork(conf.Network.Host, conf.Network.Port, packetFactory, message, join, quit),

		players: &playerAssociation{
			uuidToData: make(map[uuid.UUID]ents.Player),
			connToUUID: make(map[implBase.Connection]uuid.UUID),
			uuidToConn: make(map[uuid.UUID]implBase.Connection),
		},
	}
}

// ==== State ====
func (s *server) Load() {
	apis.SetMinecraftServer(s)

	go s.loadServer()
	go s.readInputs()

	s.wait()
}

func (s *server) Kill() {

	s.console.Kill()
	s.command.Kill()
	s.tasking.Kill()
	s.network.Kill()

	// push the stop message to the server exit channel
	s.message <- system.Make(system.STOP, "normal stop")
	close(s.message)

	s.logging.Info(chat.DarkRed, "server stopped")
}

// ==== Server ====
func (s *server) Logging() *logs.Logging {
	return s.logging
}

func (s *server) Command() *cmds.CommandManager {
	return s.command
}

func (s *server) Tasking() *task.Tasking {
	return s.tasking
}

func (s *server) Watcher() pubsub.PubSub {
	return s.pubsub
}

func (s *server) Players() []ents.Player {
	players := make([]ents.Player, 0)

	for _, player := range s.players.uuidToData {
		players = append(players, player)
	}

	return players
}

func (s *server) ConnByUUID(uuid uuid.UUID) implBase.Connection {
	return s.players.uuidToConn[uuid]
}

func (s *server) PlayerByUUID(uuid uuid.UUID) ents.Player {
	return s.players.uuidToData[uuid]
}

func (s *server) PlayerByConn(conn implBase.Connection) ents.Player {
	uuid, con := s.players.connToUUID[conn]
	if !con {
		return nil
	}

	return s.PlayerByUUID(uuid)
}

func (s *server) ServerVersion() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) Broadcast(message string) {
	s.console.SendMessage(message)

	for _, player := range s.Players() {
		player.SendMessage(message)
	}
}

// ==== server commands ====
func (s *server) broadcastCommand(sender ents.Sender, params []string) {
	message := strings.Join(params, " ")

	for _, player := range s.Players() {
		player.SendMessage(message)
	}
}

func (s *server) stopServerCommand(sender ents.Sender, params []string) {
	if _, ok := sender.(*console.Console); !ok {
		s.logging.ErrorF("non console sender %s tried to stop the server", sender.Name())
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
		s.logging.Warn(chat.Gold, "stopping server in ", chat.Green, util.FormatTime(after))

		// schedule shutdown {after} seconds later
		s.tasking.AfterTime(after, time.Second, func(task *task.Task) {
			s.Kill()
		})

	}
}

func (s *server) versionCommand(sender ents.Sender, params []string) {
	sender.SendMessage(s.ServerVersion())
}

// ==== internal ====
func (s *server) loadServer() {
	s.console.Load()
	s.command.Load()
	s.tasking.Load()
	s.network.Load()

	s.command.Register("vers", s.versionCommand)
	s.command.Register("send", s.broadcastCommand)
	s.command.Register("stop", s.stopServerCommand)

	s.pubsub.Subscribe(func(event apisEvent.PlayerJoinEvent) {
		s.logging.InfoF("player %s logged in with uuid:%v", event.Player.Name(), event.Player.UUID())

		s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has joined!", chat.Yellow, event.Player.Name())))
	})
	s.pubsub.Subscribe(func(event apisEvent.PlayerQuitEvent) {
		s.logging.InfoF("%s disconnected!", event.Player.Name())

		s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has left!", chat.Yellow, event.Player.Name())))
	})

	s.pubsub.Subscribe(func(event implEvent.PlayerConnJoinEvent) {
		s.players.addData(event.Conn)

		s.pubsub.Publish(apisEvent.PlayerJoinEvent{PlayerEvent: apisEvent.PlayerEvent{Player: event.Conn.Player}})
	})
	s.pubsub.Subscribe(func(event implEvent.PlayerConnQuitEvent) {
		player := s.players.playerByConn(event.Conn.Connection)

		if player != nil {
			s.pubsub.Publish(apisEvent.PlayerQuitEvent{PlayerEvent: apisEvent.PlayerEvent{Player: player}})
		}

		s.players.delData(event.Conn)
	})

	s.pubsub.Subscribe(func(event implEvent.PlayerPluginMessagePullEvent) {
		s.logging.DebugF("received message on channel '%s' from player %s:%s", event.Channel, event.Conn.Name(), event.Conn.UUID())

		switch event.Channel {
		case plugin.CHANNEL_BRAND:
			s.logging.DebugF("their client's brand is '%s'", event.Message.(*plugin.Brand).Name)
		}
	})
}

func (s *server) readInputs() {
	for {
		// read input from console
		text := strings.Trim(<-s.console.IChannel, " ")
		if len(text) == 0 {
			continue
		}

		args := strings.Split(text, " ")
		if len(args) == 0 {
			continue
		}

		if command := s.command.Search(args[0]); command != nil {

			err := apisBase.Attempt(func() {
				(*command).Evaluate(s.console, args[1:])
			})

			if err != nil {
				s.logging.Error(
					chat.Red, "failed to evaluate ",
					chat.DarkGray, "`",
					chat.White, (*command).Name(),
					chat.DarkGray, "`",
					chat.Red, ": ", err.Error()[8:])
			}

			continue
		}

		s.console.SendMessage(text)
	}
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
			s.logging.Error("internal server error: ", command.Message)
			s.logging.Error("stopping server")
			return
		}
	}

	s.wait()
}

// ==== players ====
type playerAssociation struct {
	uuidToData map[uuid.UUID]ents.Player

	connToUUID map[implBase.Connection]uuid.UUID
	uuidToConn map[uuid.UUID]implBase.Connection
}

func (p *playerAssociation) addData(data implBase.PlayerAndConnection) {
	p.uuidToData[data.Player.UUID()] = data.Player

	p.connToUUID[data.Connection] = data.Player.UUID()
	p.uuidToConn[data.Player.UUID()] = data.Connection
}

func (p *playerAssociation) delData(data implBase.PlayerAndConnection) {
	player := p.playerByConn(data.Connection)

	uuid := p.connToUUID[data.Connection]

	delete(p.connToUUID, data.Connection)
	delete(p.uuidToConn, uuid)

	if player != nil {
		delete(p.uuidToData, player.UUID())
	}
}

func (p *playerAssociation) playerByUUID(uuid uuid.UUID) ents.Player {
	return p.uuidToData[uuid]
}

func (p *playerAssociation) playerByConn(conn implBase.Connection) ents.Player {
	uuid, con := p.connToUUID[conn]

	if !con {
		return nil
	}

	data, con := p.uuidToData[uuid]

	if !con {
		return nil
	}

	return data
}
