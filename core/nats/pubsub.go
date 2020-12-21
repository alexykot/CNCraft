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
	"errors"
	"fmt"
	"time"

	natsd "github.com/nats-io/nats-server/server"
	natsc "github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/pkg/envelope"
)

// DEBT centralising all meta names in one place may not be a great idea
const (
	MetaConnID = "conn_id"
)

const natsStartTimeout = 500 * time.Millisecond

type PubSub interface {
	Start() error
	Publish(subj string, messages ...*envelope.E) error
	Subscribe(subj string, handleFunc func(message *envelope.E)) error
	Unsubscribe(subj string)
}

type pubsub struct {
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

func NewPubSub(log *zap.Logger, nats *natsd.Server, control chan control.Command) PubSub {
	return &pubsub{
		natsd:   nats,
		control: control,
		log:     log,
		subs:    map[string][]*natsc.Subscription{},
	}
}

func (ps *pubsub) Start() error {
	if err := ps.startServer(); err != nil {
		return fmt.Errorf("failed to start NATS server: %w", err)
	}

	if err := ps.startClient(); err != nil {
		return fmt.Errorf("failed to start NATS client: %w", err)
	}

	return nil
}

func (ps *pubsub) startServer() error {
	go func() {
		defer func() {
			message := "nats stopped unexpectedly"
			if r := recover(); r != nil {
				message = fmt.Sprintf("nats panicked: %v", r)
			}
			// stop the server if NATS exits for any reason
			ps.control <- control.Command{Signal: control.FAIL, Message: message}
		}()
		// This needs to be called manually for some reason. Without it NATS will run completely silently.
		ps.natsd.ConfigureLogger()
		ps.natsd.Start()
	}()

	// This will block until NATS server is ready for client connections,
	// or until provided timeout runs out, whichever comes earlier.
	if ok := ps.natsd.ReadyForConnections(natsStartTimeout); !ok {
		return errors.New("failed to start NATS server within the timeout")
	}

	ps.log.Info("started NATS server")
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

	ps.log.Info("connected NATS client")
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
	return nil
}

func (ps *pubsub) Subscribe(subject string, handleFunc func(*envelope.E)) error {
	sub, err := ps.client.Subscribe(subject, ps.makeHandler(handleFunc))
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}
	ps.subs[subject] = append(ps.subs[subject], sub)

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
		lope := envelope.NewEmpty()
		if err := lope.Unmarshal(msg.Data); err != nil {
			ps.log.Error("failed to unmarshal incoming message data", zap.Error(err))
		}

		handleFunc(lope)
	}
}
