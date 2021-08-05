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

	natsd "github.com/nats-io/nats-server/server"
	natsc "github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/pkg/envelope"
)

// DEBT probably should have and propagate a global bootstrap timeout across whole server
const natsStartTimeout = 500 * time.Millisecond

type PubSub interface {
	Start()
	Publish(subj string, messages ...*envelope.E) error
	Subscribe(subj string, handleFunc func(message *envelope.E)) error
	Unsubscribe(subj string)
}

type pubsub struct {
	ctx     context.Context
	subs    map[string][]*natsc.Subscription
	natsd   *natsd.Server
	client  *natsc.Conn
	log     *zap.Logger
	control chan control.Command
}

func NewNats() *natsd.Server {
	opts := &natsd.Options{
		Port: 4222, // NATS default port.
	}

	return natsd.New(opts)
}

func NewPubSub(ctx context.Context, log *zap.Logger, nats *natsd.Server, control chan control.Command) PubSub {
	return &pubsub{
		ctx:     ctx,
		natsd:   nats,
		control: control,
		log:     log,
		subs:    map[string][]*natsc.Subscription{},
	}
}

func (ps *pubsub) Start() {
	ps.signal(control.STARTING, nil)

	if err := ps.startServer(); err != nil {
		ps.signal(control.FAILED, fmt.Errorf("failed to start NATS server: %w", err)) // signal NATS has failed
		return
	}

	if err := ps.startClient(); err != nil {
		ps.signal(control.FAILED, fmt.Errorf("failed to start NATS client: %w", err)) // signal NATS has failed
		return
	}

	ps.signal(control.READY, nil)
	ps.log.Info("pubsub started")
}

func (ps *pubsub) startServer() error {
	go ps.handleContextCancel() // make sure NATS shutdown is called whenever context is cancelled

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

func (ps *pubsub) startClient() error {
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

func (ps *pubsub) Publish(subject string, lopes ...*envelope.E) error {
	for _, lope := range lopes {
		bytes, err := lope.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if err = ps.client.Publish(subject, bytes); err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}
	}
	ps.log.Debug("published into subj", zap.String("subj", subject))
	return nil
}

func (ps *pubsub) Subscribe(subject string, handleFunc func(*envelope.E)) error {
	sub, err := ps.client.Subscribe(subject, ps.makeHandler(handleFunc))
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}
	ps.subs[subject] = append(ps.subs[subject], sub)

	ps.log.Debug("subscribed to subj", zap.String("subj", subject))
	return nil
}

func (ps *pubsub) Unsubscribe(subject string) {
	subList, ok := ps.subs[subject]
	if !ok {
		return
	}
	for _, sub := range subList {
		if err := sub.Drain(); err != nil {
			ps.log.Warn("failed to drain subscription", zap.Error(err), zap.String("subject", subject))
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

func (ps *pubsub) handleContextCancel() {
	for {
		select {
		case <-ps.ctx.Done():
			ps.signal(control.STOPPED, nil)
			ps.natsd.Shutdown()
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
