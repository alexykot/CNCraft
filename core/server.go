package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/log"
	"github.com/alexykot/cncraft/pkg/protocol"
)

type Server interface {
	Load()
	Kill()

	Log() *zap.Logger

	Users() []User

	GetUser(uuid.UUID) (User, bool)
	GetConnection(uuid.UUID) (network.Connection, bool)

	Version() string

	//	Chat() // chat implementation needed

	// TODO this should be part of chat implementation.
	//Broadcast(message string) error
}

type server struct {
	log *zap.Logger
	bus bus.PubSub

	control chan control.Command

	//chat    // chat implementation needed
	network network.Network

	packFactory protocol.PacketFactory
	users       map[uuid.UUID]User
	connections map[uuid.UUID]network.Connection
}

// NewServer wires up and provides new server instance.
func NewServer(conf ServerConfig) (Server, error) {
	logger, err := log.GetLogger(conf.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate logger: %w", err)
	}

	packetFactory := protocol.NewPacketFactory(logger)
	ps := bus.New()

	return &server{
		log:     logger,
		bus:     ps,
		control: make(chan control.Command),

		packFactory: packetFactory,
		//network:     conn.NewNetwork(conf.Network.Host, conf.Network.Port, packetFactory, message, join, quit),

		users: make(map[uuid.UUID]User),
	}, nil
}

// ==== State ====
func (s *server) Load() {
	go s.loadServer()

	s.wait()
}

func (s *server) Kill() {
	// TODO to stop server node:
	//  - close all open connections
	//  - stop processing chunks
	//  - flush all chunk state to DB and unload all chunks
	//  - notify other nodes to reallocate unloaded chunks
	//  - leave the cluster group
	//  - stop NATS node
	//  - close all goroutines
	//  - exit server

	s.log.Info("server stopped")
}

func (s *server) Log() *zap.Logger { return s.log }

func (s *server) Bus() bus.PubSub { return s.bus }

func (s *server) Users() []User {
	if len(s.users) == 0 {
		return nil
	}

	userList := make([]User, len(s.users), len(s.users))

	var i int
	for _, player := range s.users {
		userList[i] = player
	}

	return userList
}

func (s *server) GetConnection(uuid uuid.UUID) (network.Connection, bool) {
	conn, ok := s.connections[uuid]
	return conn, ok
}

func (s *server) GetUser(uuid uuid.UUID) (User, bool) {
	user, ok := s.users[uuid]
	return user, ok
}

func (s *server) Version() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) stopServer(after time.Duration) {
	s.log.Warn(fmt.Sprintf("stopping server in %s", after))
	if after == 0 {
		s.Kill()
	} else {
		// TODO schedule shutdown {after} seconds later
	}
}

func (s *server) loadServer() {
	//state.RegisterHandlersState0(s.bus)
	//state.RegisterHandlersState1(s.bus)
	//state.RegisterHandlersState2(s.bus, join)
	//state.RegisterHandlersState3(s.bus, logger, tasking, join, quit)

	//s.console.Load()
	//s.command.Load()
	//s.tasking.Load()
	//s.network.Load()
	//
	//s.command.Register("vers", s.versionCommand)
	//s.command.Register("send", s.broadcastCommand)
	//s.command.Register("stop", s.stopServerCommand)
	//
	//s.bus.Subscribe(func(event apisEvent.PlayerJoinEvent) {
	//	s.log.InfoF("player %s logged in with uuid:%v", event.PlayerCharacter.Name(), event.PlayerCharacter.UUID())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has joined!", chat.Yellow, event.PlayerCharacter.Name())))
	//})
	//s.bus.Subscribe(func(event apisEvent.PlayerQuitEvent) {
	//	s.log.InfoF("%s disconnected!", event.PlayerCharacter.Name())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has left!", chat.Yellow, event.PlayerCharacter.Name())))
	//})
	//
	//s.bus.Subscribe(func(event implEvent.PlayerConnJoinEvent) {
	//	s.users.addData(event.Conn)
	//
	//	s.bus.Publish(apisEvent.PlayerJoinEvent{PlayerEvent: apisEvent.PlayerEvent{PlayerCharacter: event.Conn.PlayerCharacter}})
	//})
	//s.bus.Subscribe(func(event implEvent.PlayerConnQuitEvent) {
	//	player := s.users.playerByConn(event.Conn.Connection)
	//
	//	if player != nil {
	//		s.bus.Publish(apisEvent.PlayerQuitEvent{PlayerEvent: apisEvent.PlayerEvent{PlayerCharacter: player}})
	//	}
	//
	//	s.users.delData(event.Conn)
	//})
	//
	//s.bus.Subscribe(func(event implEvent.PlayerPluginMessagePullEvent) {
	//	s.log.DebugF("received message on channel '%s' from player %s:%s", event.Channel, event.Conn.Name(), event.Conn.UUID())
	//
	//	switch event.Channel {
	//	case plugin.CHANNEL_BRAND:
	//		s.log.DebugF("their client's brand is '%s'", event.Message.(*plugin.Brand).Name)
	//	}
	//})
}

func (s *server) wait() {
	// TODO somehow signal stopping of the whole cluster or detect that the server is the last one.
	//  If is the last one - notify in chat.

	// select over server commands channel
	select {
	case command := <-s.control:
		switch command.Signal {
		// stop selecting when stop is received
		case control.STOP:
			return
		case control.FAIL:
			s.log.Error("internal server error: ", zap.Any("message", command.Message))
			s.log.Error("stopping server")
			return
		}
	}
}
