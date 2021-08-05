package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
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
	// DEBT this should be part of chat implementation.
	// Broadcast(message string) error
}

type server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	reg map[control.Component]control.ComponentState // component state registry

	control chan control.Command
	config  control.ServerConf

	killSignal chan os.Signal
	db         *sql.DB

	log *zap.Logger
	ps  nats.PubSub

	net *network.Network

	players *players.Roster
	world   *w.World
	sharder *w.Sharder

	// chat    // chat implementation needed
}

// NewServer wires up and provides new server instance.
func NewServer(conf control.ServerConf) (Server, error) {
	serverCtx, cancel := context.WithCancel(context.Background())
	controlChan := make(chan control.Command)

	rootLog, err := log.GetRootLogger(serverCtx, conf.LogLevels.Baseline)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate root logger: %w", err)
	}

	pubSub := nats.NewPubSub(serverCtx, log.LevelUp(log.Named(rootLog, "pubsub"), conf.LogLevels.PubSub), nats.NewNats(), controlChan)

	database, err := db.New(conf.DBURL, conf.LogLevels.DB == "DEBUG", log.LevelUp(log.Named(rootLog, "DB"), conf.LogLevels.DB))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate DB: %w", err)
	}

	roster := players.NewRoster(serverCtx,
		log.LevelUp(log.Named(rootLog, "players"), conf.LogLevels.Players),
		log.LevelUp(log.Named(rootLog, "windows"), conf.LogLevels.Players),
		pubSub, database)

	world, err := w.NewWorld(serverCtx, conf.World, log.LevelUp(log.Named(rootLog, "world"), conf.LogLevels.World), database)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate world: %w", err)
	}

	sharder := w.NewSharder(serverCtx, conf.World, log.LevelUp(log.Named(rootLog, "sharder"), conf.LogLevels.Sharder), pubSub, controlChan, world, roster)

	dispatcher := network.NewDispatcher(serverCtx,
		log.LevelUp(log.Named(rootLog, "dispatcher"), conf.LogLevels.Dispatcher),
		pubSub, auth.GetAuther(),
		roster,
		network.NewKeepAliver(controlChan, pubSub, log.LevelUp(log.Named(rootLog, "aliver"), conf.LogLevels.Dispatcher)),
		sharder,
	)

	net := network.NewNetwork(serverCtx, conf.Network, log.LevelUp(log.Named(rootLog, "network"), conf.LogLevels.Network), controlChan, pubSub, dispatcher)

	return &server{
		ctx:        serverCtx,
		cancelFunc: cancel,

		reg: make(map[control.Component]control.ComponentState),

		config: conf,

		log: rootLog,
		db:  database,
		ps:  pubSub,
		net: net,

		control:    controlChan,
		killSignal: make(chan os.Signal),

		players: roster,
		world:   world,
		sharder: sharder,
	}, nil
}

func (s *server) Start() error {
	go func() {
		if err := s.startServer(); err != nil {
			s.log.Error("failed to start the server", zap.Error(err))
			s.control <- control.Command{Signal: control.SERVER_FAIL, Message: err.Error()}
		}
		s.log.Info("server started")
	}()

	s.startControlLoop()
	return nil
}

// shutdown will unconditionally kill the server and exit the process.
func (s *server) shutdown(failed bool) {
	var code int
	if failed {
		code = 1
	}

	if err := s.Stop(); err != nil {
		s.log.Error("received error while stopping server", zap.Error(err))
		code = 1
	}

	os.Exit(code)
}

func (s *server) Stop() error {
	const serverStopTimeout = time.Second * 5
	s.log.Info("stopping server")

	s.cancelFunc()

	var stopTimeout = time.NewTimer(serverStopTimeout)
	var stopTicker = time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case <-stopTimeout.C:
			s.log.Error("reached server stop timeout", zap.Any("regState", s.reg))
			return errors.New("reached server stop timeout")
		case <-stopTicker.C:
			var stillWaiting bool // still waiting for some components to stop
			for _, state := range s.reg {
				if state != control.STOPPED && state != control.FAILED {
					stillWaiting = true
					break
				}
			}

			if !stillWaiting {
				s.log.Info("server stopped")
				return nil
			}
		}
	}

	// DEBT
	//  To implement stopping the server:
	//    - register all stoppable processes on start
	//    - signal all stoppables to stop
	//    - block in a loop until all stoppable processes have indeed stopped
	//    - have a breaker deadline timeout in case a stoppable is stuck indefinitely
	//  While stopping the server node do this:
	//    - notify players connected to this node that it is stopping and they're about to be disconnected
	//    - close all open connections
	//    - stop processing chunks
	//    - flush all chunk state to DB and unload all chunks
	//    - notify other nodes that chunks are unloaded
	//    - leave the cluster group
	//    - stop NATS node
	//    - close all goroutines
	//    - exit server

}

func (s *server) Version() string {
	return "0.0.1-SNAPSHOT"
}

func (s *server) startServer() error {
	rand.Seed(time.Now().UnixNano())

	// DEBT when starting server:
	//  - register all startable server objects
	//  - separate object starting sequence and readiness signal
	//  - retrieve readiness from every startable objects
	//  - signal overall readiness to all objects when ready
	//  - only accept network connections when everything is ready
	//  - use control loop channel to signal READY states
	//  - have a breaker timeout in case some startable object is stuck and cannot start in time

	control.RegisterCurrentConfig(s.config)

	if err := s.ps.Start(); err != nil {
		return fmt.Errorf("failed to start nats: %w", err)
	}

	if err := db.Migrate(s.db); err != nil {
		return fmt.Errorf("failed to migrate the database schema: %w", err)
	}

	if err := db.RegisterStateRecorders(s.ps, s.db); err != nil {
		return fmt.Errorf("failed to register DB state persisting handlers: %w", err)
	}

	if err := s.net.Start(context.TODO()); err != nil {
		return fmt.Errorf("failed to start network: %w", err)
	}

	if err := s.world.Load(); err != nil {
		return fmt.Errorf("failed to load world data: %w", err)
	}

	s.sharder.Start()

	if err := handlers.RegisterEventHandlersState3(log.LevelUp(log.Named(s.log, "players"), s.config.LogLevels.Players),
		s.ps, s.players, s.world); err != nil {
		return fmt.Errorf("failed to register Play state handlers: %w", err)
	}

	if err := s.players.RegisterHandlers(); err != nil {
		return fmt.Errorf("failed to register global player handlers: %w", err)
	}

	return nil
}

func (s *server) startControlLoop() {
	signal.Notify(s.killSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	var mu sync.Mutex

	for {
		select { // select over server commands channel
		case command := <-s.control:
			switch command.Signal {
			case control.COMPONENT:
				mu.Lock()
				// These transitions of the state machine are allowed:
				// <no state>       => <any state>
				// control.STARTING => control.READY, control.STOPPED, control.FAILED
				// control.READY    => control.STOPPED, control.FAILED
				// Any other state transitions are assumed to be a result race condition
				// inside the component and will be ignored.
				// On receiving a control.FAILED message - server shutdown will be initiated.
				if state, ok := s.reg[command.Component]; !ok {
					s.reg[command.Component] = command.State
				} else if state == control.STARTING {
					s.reg[command.Component] = command.State
				} else if state == control.READY && command.State != control.STARTING {
					s.reg[command.Component] = command.State
				}
				mu.Unlock()

				if command.State == control.FAILED {
					s.log.Error("component failed", zap.String("comp", string(command.Component)), zap.String("component", command.Message))
					go s.shutdown(true)
				}
			}
		case <-s.killSignal:
			s.log.Info("received interrupt signal")
			go s.shutdown(false)
		}
	}
}
