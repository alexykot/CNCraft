package core

import (
	"context"
	"database/sql"
	"errors"
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
	serverCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	ctrlChan := make(chan control.Command)

	rootLog, err := log.GetRootLogger(serverCtx, conf.LogLevels.Baseline)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate root logger: %w", err)
	}

	pubSub := nats.NewPubSub(serverCtx, log.LevelUp(log.Named(rootLog, "pubsub"), conf.LogLevels.PubSub), nats.NewNats(), ctrlChan)

	database, err := db.New(conf.DBURL, conf.LogLevels.DB == "DEBUG", log.LevelUp(log.Named(rootLog, "DB"), conf.LogLevels.DB))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate DB: %w", err)
	}

	roster := players.NewRoster(serverCtx, ctrlChan,
		log.LevelUp(log.Named(rootLog, "players"), conf.LogLevels.Players),
		log.LevelUp(log.Named(rootLog, "windows"), conf.LogLevels.Players),
		pubSub, database)

	world, err := w.NewWorld(serverCtx, conf.World, log.LevelUp(log.Named(rootLog, "world"), conf.LogLevels.World), database)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not instantiate world: %w", err)
	}

	sharder := w.NewSharder(serverCtx, conf.World, log.LevelUp(log.Named(rootLog, "sharder"), conf.LogLevels.Sharder), pubSub, ctrlChan, world, roster)

	dispatcher := network.NewDispatcher(serverCtx,
		log.LevelUp(log.Named(rootLog, "dispatcher"), conf.LogLevels.Dispatcher),
		pubSub, auth.GetAuther(),
		roster,
		network.NewKeepAliver(ctrlChan, pubSub, log.LevelUp(log.Named(rootLog, "aliver"), conf.LogLevels.Dispatcher)),
		sharder,
	)

	net := network.NewNetwork(serverCtx, conf.Network, log.LevelUp(log.Named(rootLog, "network"), conf.LogLevels.Network), ctrlChan, pubSub, dispatcher)

	return &server{
		ctx:        serverCtx,
		cancelFunc: cancel,

		reg: make(map[control.Component]control.ComponentState),

		config: conf,

		log: rootLog,
		db:  database,
		ps:  pubSub,
		net: net,

		control:    ctrlChan,
		killSignal: make(chan os.Signal),

		players: roster,
		world:   world,
		sharder: sharder,
	}, nil
}

func (s *server) Start() error {
	go s.startServer()

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

	for {
		select {
		case <-stopTimeout.C:
			s.log.Error("reached server stop timeout", zap.Any("regState", s.reg))
			return errors.New("reached server stop timeout")
		case command := <-s.control:
			if command.Signal == control.COMPONENT && (command.State == control.STOPPED || command.State == control.FAILED) {
				s.log.Info("component stopped", zap.Any("comp", command.Component))
				s.reg[command.Component] = command.State
			}

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

func (s *server) startServer() {
	rand.Seed(time.Now().UnixNano())

	// DEBT when starting server:
	//  - only accept network connections when everything is ready
	//  - have a breaker timeout in case some startable object is stuck and cannot start in time

	control.RegisterCurrentConfig(s.config)

	s.ps.Start()

	db.Init(s.control, s.ps, s.db)

	s.net.Start()

	s.world.Load(s.control)

	s.sharder.Start()

	handlers.RegisterEventHandlersState3(s.control, log.LevelUp(log.Named(s.log, "players"), s.config.LogLevels.Players),
		s.ps, s.players, s.world)

	s.players.RegisterHandlers()

	s.log.Info("server started")
}

func (s *server) startControlLoop() {
	var serverFailed bool

controlLoop:
	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("server context cancelled, shutdown initiated")
			break controlLoop

		case command := <-s.control:
			switch command.Signal {
			case control.COMPONENT:
				// These transitions of the state machine are allowed:
				// <no state>       => <any state>
				// control.STARTING => control.READY, control.STOPPED, control.FAILED
				// control.READY    => control.STOPPED, control.FAILED
				// Any other state transitions are assumed to be a result race condition
				// inside the component and will be ignored.
				// On receiving a control.FAILED message - server shutdown will be initiated.
				//
				// This map does not need to be thread safe as it's only ever accessed here and inside the shutdown
				// sequence, and never - in parallel.
				if state, ok := s.reg[command.Component]; !ok {
					s.reg[command.Component] = command.State
				} else if state == control.STARTING {
					s.reg[command.Component] = command.State
				} else if state == control.READY && command.State != control.STARTING {
					s.reg[command.Component] = command.State
				}

				if command.State == control.FAILED {
					s.log.Error("component failed", zap.String("comp", string(command.Component)), zap.Error(command.Err))
					serverFailed = true
					break controlLoop
				}
			}
		}
	}
	s.shutdown(serverFailed)
}
