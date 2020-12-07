package bus

import (
	"sync"
)

// DEBT centralising all meta names in one place may not be a great idea
const (
	MetaConn = "conn"
)

type Envelope interface {
	GetMessage() interface{}
	GetMeta(string) (string, bool)
	GetAllMeta() map[string]string
}

type envelope struct {
	message interface{}
	meta map[string]string
}

func NewEnvelope(message interface{}, meta map[string]string) Envelope {
	return envelope{message: message, meta: meta}
}

func (e envelope) GetMessage() interface{} {return e.message}
func (e envelope) GetAllMeta() map[string]string {return e.meta}
func (e envelope) GetMeta(key string) (string, bool) {
	val, ok := e.meta[key]
	return val, ok
}

type PubSub interface {
	Publish(topic string, message ...Envelope)
	Subscribe(topic string, handleFunc func(message Envelope)) Handler
}

type Handler interface {
	UnSub()
}


func New() PubSub {
	return &pubsub{
		locker: sync.Mutex{},
		topics: make(map[string][]*handler),
	}
}

type pubsub struct {
	locker sync.Mutex
	topics map[string][]*handler
}

type handler struct {
	topic string
	watch *pubsub // DEBT not sure why this is needed, need to look at UnSub() closely.

	function func(message Envelope)
}

func (w *pubsub) publish(topic string, envelope Envelope) {
	handlers, _ := w.topics[topic]
	for _, handler := range handlers {
		handler.function(envelope)
	}
}

func (w *pubsub) Publish(topic string, messages ...Envelope) {
	for _, value := range messages {
		w.publish(topic, value)
	}
}

func (w *pubsub) Subscribe(topic string, handleFunc func(Envelope)) Handler {
	handler := &handler{
		topic:    topic,
		watch:    w,
		function: handleFunc,
	}

	w.locker.Lock()
	w.topics[topic] = append(w.topics[topic], handler)
	w.locker.Unlock()

	return handler
}

func (h *handler) UnSub() {
	handlers := h.watch.topics[h.topic]
	if handlers == nil {
		return
	}

	for i, elem := range handlers {
		if elem == h {
			h.watch.topics[h.topic] = append(h.watch.topics[h.topic][:i], h.watch.topics[h.topic][i+1:]...)
		}
	}
}
