package core

type Signal int

const (
	// stops the server entirely
	STOP Signal = iota
	FAIL
)

type Command struct {
	Signal Signal
	Message string
}
