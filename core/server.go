package core

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/db"
	"github.com/alexykot/cncraft/core/handlers"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/network"
	"github.com/alexykot/cncraft/core/players"
	w "github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/log"
	"github.com/alexykot/cncraft/pkg/protocol/auth"
)

type Server interface {
	Start() error
	Stop() error

	Version() string

	//	Chat() // chat implementation needed
	// TODO this should be part of chat implementation.
	// Broadcast(message string) error
}

type server struct {
	control chan control.Command
	config  control.ServerConf

	signal chan os.Signal
	db     *sql.DB

	log *zap.Logger
	ps  nats.PubSub

	network *network.Network

	players *players.Roster
	world   *w.World

	// chat    // chat implementation needed
}

// NewServer wires up and provides new server instance.
func NewServer(conf control.ServerConf) (Server, error) {
	rootLog, err := log.GetRootLogger(conf.LogLevels.Baseline)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate root logger: %w", err)
	}

	controlChan := make(chan control.Command)

	pubSub := nats.NewPubSub(log.LevelUp(log.Named(rootLog, "pubsub"), conf.LogLevels.PubSub), nats.NewNats(), controlChan)

	database, err := db.New(conf.DBURL, conf.LogLevels.DB == "DEBUG", log.LevelUp(log.Named(rootLog, "DB"), conf.LogLevels.DB))
	if err != nil {
		return nil, fmt.Errorf("could not instantiate DB: %w", err)
	}

	roster := players.NewRoster(
		log.LevelUp(log.Named(rootLog, "players"), conf.LogLevels.Players),
		log.LevelUp(log.Named(rootLog, "windows"), conf.LogLevels.Players),
		pubSub, database)

	dispatcher := network.NewDispatcher(log.LevelUp(log.Named(rootLog, "dispatcher"), conf.LogLevels.Dispatcher),
		pubSub, auth.GetAuther(), roster,
		network.NewKeepAliver(controlChan, pubSub, log.LevelUp(log.Named(rootLog, "aliver"), conf.LogLevels.Dispatcher)))

	net := network.NewNetwork(conf.Network, log.LevelUp(log.Named(rootLog, "network"), conf.LogLevels.Network),
		controlChan, pubSub, dispatcher)

	world, err := w.NewWorld(conf.WorldID, log.LevelUp(log.Named(rootLog, "world"), conf.LogLevels.World), database)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate world: %w", err)
	}

	return &server{
		config:  conf,
		log:     rootLog,
		ps:      pubSub,
		control: controlChan,
		signal:  make(chan os.Signal),
		network: net,
		players: roster,
		world:   world,
		db:      database,
	}, nil
}

func (s *server) Start() error {
	go func() {
		if err := s.startServer(); err != nil {
			s.log.Error("failed to start the server", zap.Error(err))
			s.control <- control.Command{Signal: control.FAIL, Message: err.Error()}
		}
		s.log.Info("server started")
	}()

	s.startControlLoop()
	return nil
}

func (s *server) Stop() error {
	s.log.Info("stopping server")

	// TODO to stop server node:
	//  - notify players connected to this node that it is stopping and they're about to be disconnected
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

// TODO use context when starting all subsystems and use cancellation to signal graceful shutdown.
func (s *server) startServer() error {
	rand.Seed(time.Now().UnixNano())

	control.RegisterCurrentConfig(s.config)

	if err := s.ps.Start(); err != nil {
		return fmt.Errorf("failed to start nats: %w", err)
	}

	// TODO make network not accept connections until bootstrap is fully done.
	//  Probably use control loop channel to signal READY state.
	if err := s.network.Start(context.TODO()); err != nil {
		return fmt.Errorf("failed to start network: %w", err)
	}

	if err := s.players.RegisterHandlers(); err != nil {
		return fmt.Errorf("failed to register global player handlers: %w", err)
	}

	if err := s.world.Load(); err != nil {
		return fmt.Errorf("failed to load world data: %w", err)
	}

	if err := handlers.RegisterEventHandlersState3(log.LevelUp(log.Named(s.log, "players"), s.config.LogLevels.Players),
		s.ps, s.players, s.world); err != nil {
		return fmt.Errorf("failed to register Play state handlers: %w", err)
	}

	if err := db.Migrate(s.db); err != nil {
		return fmt.Errorf("failed to migrate the database schema: %w", err)
	}

	if err := db.RegisterDBRecorders(s.ps, s.db); err != nil {
		return fmt.Errorf("failed to register DB state persisting handlers: %w", err)
	}

	return nil
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
