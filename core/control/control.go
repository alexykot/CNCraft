package control

type Signal int

const (
	STOP Signal = iota
	FAIL
)

type Command struct {
	Signal  Signal
	Message string
}
