package nats

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
