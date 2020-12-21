package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/core/state"
	"github.com/alexykot/cncraft/core/users"
	"github.com/alexykot/cncraft/pkg/log"
)

type Server interface {
	Start() error
	Stop() error

	Log() *zap.Logger

	// DEBT maybe move all users into substruct
	Users() []users.User
	GetUser(uuid.UUID) (users.User, bool)

	Version() string

	//	Chat() // chat implementation needed
	// TODO this should be part of chat implementation.
	//Broadcast(message string) error
}

type server struct {
	log *zap.Logger
	ps  nats.PubSub

	control chan control.Command
	signal  chan os.Signal

	dispatcher *state.SPacketDispatcher

	//chat    // chat implementation needed
	network *network.Network

	users map[uuid.UUID]users.User
}

// NewServer wires up and provides new server instance.
func NewServer(conf control.ServerConf) (Server, error) {
	logger, err := log.GetLogger(conf.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate logger: %w", err)
	}

	controlChan := make(chan control.Command)
	packetFactory := state.NewPacketFactory()
	pubSub := nats.NewPubSub(logger.Named("pubsub"), nats.NewNats(), controlChan)

	return &server{
		log:        logger.Named("core"),
		ps:         pubSub,
		control:    controlChan,
		signal:     make(chan os.Signal),
		dispatcher: state.NewDispatcher(logger.Named("dispatcher"), packetFactory, pubSub),

		network: network.NewNetwork(conf.Network, logger.Named("network"), controlChan, pubSub),

		users: make(map[uuid.UUID]users.User),
	}, nil
}

// ==== State ====
func (s *server) Start() error {
	go func() {
		if err := s.startServer(); err != nil {
			s.log.Error("failed to start the server", zap.Error(err))
			s.control <- control.Command{Signal: control.FAIL, Message: err.Error()}
		}
	}()
	s.startControlLoop()

	return nil
}

func (s *server) Stop() error {
	s.log.Info("stopping server")

	// TODO to stop server node:
	//  - notify users connected to this node that it is stopping and they're about to be disconnected
	//  - close all open connections
	//  - stop processing chunks
	//  - flush all chunk state to DB and unload all chunks
	//  - notify other nodes to reallocate unloaded chunks
	//  - leave the cluster group
	//  - stop NATS node
	//  - close all goroutines
	//  - exit server

	s.log.Info("server stopped")
	return nil
}

func (s *server) Log() *zap.Logger { return s.log }

func (s *server) Bus() nats.PubSub { return s.ps }

func (s *server) Users() []users.User {
	if len(s.users) == 0 {
		return nil
	}

	userList := make([]users.User, len(s.users), len(s.users))

	var i int
	for _, player := range s.users {
		userList[i] = player
	}

	return userList
}

func (s *server) GetUser(uuid uuid.UUID) (users.User, bool) {
	user, ok := s.users[uuid]
	return user, ok
}

func (s *server) Version() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) stopServer(after time.Duration) {
	s.log.Warn(fmt.Sprintf("stopping server in %s", after))
	if after == 0 {
		s.control <- control.Command{Signal: control.STOP}
	} else {
		// TODO schedule shutdown {after} seconds later
	}
}

func (s *server) startServer() error {
	if err := s.ps.Start(); err != nil {
		return fmt.Errorf("failed to start nats: %w", err)
	}

	if err := s.network.Start(); err != nil {
		return fmt.Errorf("failed to start network: %w", err)
	}

	if err := s.registerGlobalHandlers(); err != nil {
		return fmt.Errorf("failed register state0 handlers: %w", err)
	}

	return nil

	//s.console.Load()
	//s.command.Load()
	//s.tasking.Load()

	//state.RegisterHandlersState0(s.ps)
	//state.RegisterHandlersState1(s.ps)
	//state.RegisterHandlersState2(s.ps, join)
	//state.RegisterHandlersState3(s.ps, logger, tasking, join, quit)

	//s.ps.Subscribe(func(event apisEvent.PlayerJoinEvent) {
	//	s.log.InfoF("player %s logged in with uuid:%v", event.PlayerCharacter.Name(), event.PlayerCharacter.UUID())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has joined!", chat.Yellow, event.PlayerCharacter.Name())))
	//})

	//s.ps.Subscribe(func(event apisEvent.PlayerQuitEvent) {
	//	s.log.InfoF("%s disconnected!", event.PlayerCharacter.Name())
	//
	//	s.Broadcast(chat.Translate(fmt.Sprintf("%s%s has left!", chat.Yellow, event.PlayerCharacter.Name())))
	//})

	//s.ps.Subscribe(func(event implEvent.PlayerConnJoinEvent) {
	//	s.users.addData(event.Conn)
	//
	//	s.ps.Publish(apisEvent.PlayerJoinEvent{PlayerEvent: apisEvent.PlayerEvent{PlayerCharacter: event.Conn.PlayerCharacter}})
	//})

	//s.ps.Subscribe(func(event implEvent.PlayerConnQuitEvent) {
	//	player := s.users.playerByConn(event.Conn.Connection)
	//
	//	if player != nil {
	//		s.ps.Publish(apisEvent.PlayerQuitEvent{PlayerEvent: apisEvent.PlayerEvent{PlayerCharacter: player}})
	//	}
	//
	//	s.users.delData(event.Conn)
	//})

	//s.ps.Subscribe(func(event implEvent.PlayerPluginMessagePullEvent) {
	//	s.log.DebugF("received message on channel '%s' from player %s:%s", event.Channel, event.Conn.Name(), event.Conn.UUID())
	//
	//	switch event.Channel {
	//	case plugin.CHANNEL_BRAND:
	//		s.log.DebugF("their client's brand is '%s'", event.Message.(*plugin.Brand).Name)
	//	}
	//})
}

func (s *server) startControlLoop() {
	signal.Notify(s.signal, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// select over server commands channel
	select {
	case command := <-s.control:
		switch command.Signal {
		// stop selecting when stop is received
		case control.STOP:
			s.log.Info("received stop command")
			if err := s.Stop(); err != nil {
				s.log.Error("received error while stopping server", zap.Error(err))
			}
			return
		case control.FAIL:
			s.log.Error("internal server error", zap.String("message", command.Message))
			if err := s.Stop(); err != nil {
				s.log.Error("received error while stopping server", zap.Error(err))
			}
			return
		}
	case <-s.signal:
		s.log.Info("received interrupt signal")
		if err := s.Stop(); err != nil {
			s.log.Error("received error while stopping server", zap.Error(err))
		}
		return
	}
}

func (s *server) registerGlobalHandlers() error {
	if err := s.ps.Subscribe(subj.MkNewConn(), s.dispatcher.HandleNewConnection); err != nil {
		return fmt.Errorf("failed to subscribe to new connections: %w", err)
	}
	return nil
}
