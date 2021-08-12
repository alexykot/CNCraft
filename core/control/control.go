package control

type signal int

const (
	COMPONENT signal = iota
)

type Component string

const (
	PUBSUB  Component = "pubsub"
	NETWORK Component = "network"
	// Not making DISPATCHER a separate component for now, as it does not have any own async event loops
	// and is part of the network bootstrap process, so covered by the NETWORK component.
	// DISPATCHER Component = "dispatcher"
	KEEPALIVER Component = "keepaliver"
	WORLD      Component = "world"
	EVENTS     Component = "events"
	SHARDER    Component = "sharder"
	ROSTER     Component = "roster"
	DB         Component = "db"
)

type ComponentState int

const (
	STARTING ComponentState = iota
	READY
	STOPPED
	FAILED
)

type Command struct {
	Signal    signal
	Component Component
	State     ComponentState
	Err       error
}
