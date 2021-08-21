//go:generate mockgen -package mocks -source=pubsub.go -destination=mocks/mocks.go PubSub

// Package nats bootstraps internal NATS server and provides a pubsub interface for handling async message delivery.
// Expected types of subscriptions:
// per entity:
//  - connection receiving server bound packets, subj per connection
//  - connection sending client bound packets, subj per connection
//  - player events, subj per player
//  - chunkster events, subj per chunkster
//
// global:
//  - announcing joining users, global subj for whole cluster
//  - announcing leaving users, global subj for whole cluster
//  - chat broadcast messages, global subj for whole cluster
package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	natsd "github.com/nats-io/nats-server/v2/server"
	natsc "github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
)

// DEBT probably should have and propagate a global bootstrap timeout across whole server
const natsStartTimeout = 500 * time.Millisecond

type PubSub interface {
	Start(ctx context.Context)
	Publish(subj subj.Subj, messages ...*envelope.E) error
	Subscribe(subj subj.Subj, handleFunc func(message *envelope.E)) error
	Unsubscribe(subj subj.Subj)
}

type natsServer interface {
	ConfigureLogger()
	Start()
	Shutdown()
	ReadyForConnections(time.Duration) bool
}

type pubsub struct {
	subs    map[subj.Subj][]*natsc.Subscription
	natsd   natsServer
	client  *natsc.Conn
	log     *zap.Logger
	control chan control.Command
}

func NewNATSServer() (natsServer, error) {
	opts := &natsd.Options{
		Port:   4222, // NATS default port.
		NoSigs: true,
	}

	server, err := natsd.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate NATS server: %w", err)
	}

	return server, nil
}

func NewPubSub(log *zap.Logger, control chan control.Command, nats natsServer) PubSub {
	return &pubsub{
		natsd:   nats,
		control: control,
		log:     log,
		subs:    map[subj.Subj][]*natsc.Subscription{},
	}
}

func (ps *pubsub) Start(ctx context.Context) {
	ps.signal(control.STARTING, nil)

	if err := ps.startServer(ctx); err != nil {
		ps.signal(control.FAILED, fmt.Errorf("failed to start NATS server: %w", err)) // signal NATS has failed
		return
	}

	if err := ps.startClient(ctx); err != nil {
		ps.signal(control.FAILED, fmt.Errorf("failed to start NATS client: %w", err)) // signal NATS has failed
		return
	}

	ps.signal(control.READY, nil)
	ps.log.Info("pubsub started")
}

func (ps *pubsub) startServer(ctx context.Context) error {
	go ps.handleContextCancel(ctx) // make sure NATS shutdown is called whenever context is cancelled

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ps.signal(control.FAILED, fmt.Errorf("nats panicked: %v", r)) // signal NATS has failed
			}
		}()
		// This needs to be called manually for some reason. Without it NATS will run completely silently.
		ps.natsd.ConfigureLogger()

		ps.natsd.Start()
	}()

	// This will block until NATS server is ready for client connections,
	// or until provided timeout runs out, whichever comes earlier.
	if ok := ps.natsd.ReadyForConnections(natsStartTimeout); !ok {
		ps.natsd.Shutdown()
		return errors.New("failed to start NATS server within the timeout")
	}

	ps.log.Debug("started NATS server")
	return nil
}

func (ps *pubsub) startClient(_ context.Context) error {
	var err error
	opts := []natsc.Option{
		natsc.DontRandomize(),
	}

	if ps.client, err = natsc.Connect(natsc.DefaultURL, opts...); err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	ps.log.Debug("connected NATS client")
	return nil
}

func (ps *pubsub) Publish(subject subj.Subj, lopes ...*envelope.E) error {
	for _, lope := range lopes {
		bytes, err := lope.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if err = ps.client.Publish(string(subject), bytes); err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}
	}
	ps.log.Debug("published into subj", zap.String("subj", string(subject)))
	return nil
}

func (ps *pubsub) Subscribe(subject subj.Subj, handleFunc func(*envelope.E)) error {
	sub, err := ps.client.Subscribe(string(subject), ps.makeHandler(handleFunc))
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}
	ps.subs[subject] = append(ps.subs[subject], sub)

	ps.log.Debug("subscribed to subj", zap.String("subj", string(subject)))
	return nil
}

func (ps *pubsub) Unsubscribe(subject subj.Subj) {
	subList, ok := ps.subs[subject]
	if !ok {
		return
	}
	for _, sub := range subList {
		if err := sub.Drain(); err != nil {
			ps.log.Warn("failed to drain subscription", zap.Error(err), zap.String("subj", string(subject)))
		}
	}
}

func (ps *pubsub) makeHandler(handleFunc func(*envelope.E)) natsc.MsgHandler {
	return func(msg *natsc.Msg) {
		lope := envelope.Empty()
		if err := lope.Unmarshal(msg.Data); err != nil {
			ps.log.Error("failed to unmarshal incoming message data", zap.Error(err))
		}

		handleFunc(lope)
	}
}

func (ps *pubsub) handleContextCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			defer func() {
				if r := recover(); r != nil {
					ps.signal(control.FAILED, fmt.Errorf("NATS server panicked while shutting down: %v", r))
				}
			}()
			ps.natsd.Shutdown()
			ps.signal(control.STOPPED, nil)
			return
		}
	}
}

func (ps *pubsub) signal(state control.ComponentState, err error) {
	ps.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.PUBSUB,
		State:     state,
		Err:       err,
	}
}
