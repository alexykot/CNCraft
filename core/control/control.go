package control

type signal int

const (
	COMPONENT signal = iota

	SERVER_FAIL
)

type Component string

const (
	PUBSUB  Component = "pubsub"
	NETWORK Component = "network"
	WORLD   Component = "world"
	SHARDER Component = "sharder"
	ROSTER  Component = "roster"
	DB      Component = "db"
)

type ComponentState int

const (
	STARTING ComponentState = iota
	READY
	FAILED
	STOPPED
)

type Command struct {
	Signal    signal
	Component Component
	State     ComponentState
	Message   string
}
