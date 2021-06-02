package control

type Signal int

const (
	SERVER_STOP Signal = iota
	SERVER_FAIL
	SHARD_FAIL
	SHARD_STOP
)

type Command struct {
	Signal  Signal
	Message string
}
