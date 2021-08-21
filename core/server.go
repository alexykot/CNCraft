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
	"github.com/alexykot/cncraft/core/world"
	"github.com/alexykot/cncraft/pkg/log"
	"github.com/alexykot/cncraft/pkg/protocol/auth"
)

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	reg map[control.Component]control.ComponentState // component state registry

	control chan control.Command
	config  control.ServerConf

	db *sql.DB

	log *zap.Logger
	ps  nats.PubSub

	net *network.Network

	roster  players.Roster
	world   *world.World
	sharder *world.Sharder
}

// NewServer wires up and provides new server instance.
func NewServer(config control.ServerConf) (*Server, error) {
	var err error
	srv := &Server{
		reg:     make(map[control.Component]control.ComponentState),
		control: make(chan control.Command),
		config:  config,
	}

	srv.ctx, srv.cancelFunc = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	if srv.log, err = log.GetRoot(srv.config.Log.Baseline); err != nil {
		srv.cancelFunc()
		return nil, fmt.Errorf("could not instantiate root logger: %w", err)
	}

	natsd, err := nats.NewNATSServer()
	if err != nil {
		srv.cancelFunc()
		return nil, fmt.Errorf("could not instantiate NATS server: %w", err)
	}
	srv.ps = nats.NewPubSub(log.NamedLevelUp(srv.log, "pubsub", srv.config.Log.PubSub), srv.control, natsd)

	if srv.db, err = db.New(log.NamedLevelUp(srv.log, "DB", srv.config.Log.DB),
		srv.config.DBURL, srv.config.Log.DB == "DEBUG"); err != nil {
		srv.cancelFunc()
		return nil, fmt.Errorf("could not instantiate DB: %w", err)
	}

	srv.roster = players.NewRoster(
		log.NamedLevelUp(srv.log, "players", srv.config.Log.Players),
		log.NamedLevelUp(srv.log, "windows", srv.config.Log.Players),
		srv.control, srv.ps, srv.db)

	if srv.world, err = world.NewWorld(log.NamedLevelUp(srv.log, "world", srv.config.Log.World),
		srv.config.World, srv.db); err != nil {
		srv.cancelFunc()
		return nil, fmt.Errorf("could not instantiate world: %w", err)
	}

	srv.sharder = world.NewSharder(log.NamedLevelUp(srv.log, "sharder", srv.config.Log.Sharder), srv.control, srv.config.World, srv.ps, srv.world, srv.roster)

	dispatcher := network.NewDispatcher(
		log.NamedLevelUp(srv.log, "dispatcher", srv.config.Log.Dispatcher),
		srv.ps, auth.GetAuther(),
		srv.roster,
		network.NewKeepAliver(log.NamedLevelUp(srv.log, "aliver", srv.config.Log.Dispatcher), srv.control, srv.ps),
		srv.sharder,
	)
	srv.net = network.NewNetwork(log.NamedLevelUp(srv.log, "network", srv.config.Log.Network), srv.control, srv.config.Net, srv.ps, dispatcher)

	return srv, nil
}

func (s *Server) Start() error {
	go s.startServer()

	s.startControlLoop()
	return nil
}

// shutdown will unconditionally kill the server and exit the process.
func (s *Server) shutdown(failed bool) {
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

func (s *Server) Stop() error {
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

func (s *Server) Version() string {
	return "0.0.1-SNAPSHOT"
}

func (s *Server) startServer() {
	rand.Seed(time.Now().UnixNano())

	// DEBT when starting server:
	//  - only accept network connections when everything is ready
	//  - have a breaker timeout in case some startable object is stuck and cannot start in time

	control.RegisterCurrentConfig(s.config)

	s.ps.Start(s.ctx)

	db.Init(s.ctx, s.control, s.ps, s.db)

	s.net.Start(s.ctx)

	s.world.Load(s.ctx, s.control)

	s.sharder.Start(s.ctx)

	handlers.RegisterEventHandlersState3(log.NamedLevelUp(s.log, "players", s.config.Log.Players),
		s.control, s.ps, s.roster, s.world)

	s.roster.Start(s.ctx)

	s.log.Info("server started")
}

func (s *Server) startControlLoop() {
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
				// Any other state transitions are assumed to be a result of possible race conditions
				// inside the component and will be ignored.
				//
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
