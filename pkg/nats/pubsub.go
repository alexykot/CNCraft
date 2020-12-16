package nats

import (
	"sync"

	natsd "github.com/nats-io/nats-server/server"
)

// DEBT centralising all meta names in one place may not be a great idea
const (
	MetaConn = "conn"
)

type PubSub interface {
	Start() error
	Publish(topic string, message ...Envelope)
	Subscribe(topic string, handleFunc func(message Envelope))
}

type pubsub struct {
	locker sync.Mutex
	topics []string
	nats *natsd.Server
}

func NewNats() *natsd.Server {
	opts := &natsd.Options{
		Port: 4222,
	}

	return natsd.New(opts)
}

func New(nats *natsd.Server) PubSub {
	return &pubsub{
		locker: sync.Mutex{},
		topics: []string{},
		nats: nats,
	}
}

func (w *pubsub) Start() error {
	// DEBT supervise this goroutine and control NATS health. Stop server on NATS crash.
	go func() {
		w.nats.Start()
	}()

	return nil
}

func (w *pubsub) publish(topic string, envelope Envelope) {

}

func (w *pubsub) Publish(topic string, messages ...Envelope) {
	for _, value := range messages {
		w.publish(topic, value)
	}
}

func (w *pubsub) Subscribe(topic string, handleFunc func(Envelope)) {
}
