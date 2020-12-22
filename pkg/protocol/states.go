package protocol

import (
	"fmt"
)

// State is one of four states of the client-server connection in Minecraft protocol
type State int

const (
	Handshake State = 0
	Status    State = 1
	Login     State = 2
	Play      State = 3
)

func IntToState(s int) (State, error) {
	state := State(s)
	var err error
	if state != Handshake && state != Status && state != Login && state != Play {
		err = fmt.Errorf("no state defined for %d", s)
	}
	return state, err
}

func (state State) String() string {
	switch state {
	case Handshake:
		return "Handshake"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Play:
		return "Play"
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}

func (state State) Next() State {
	switch state {
	case Handshake:
		return Status
	case Status:
		return Login
	case Login:
		return Play
	case Play:
		return Handshake
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}
