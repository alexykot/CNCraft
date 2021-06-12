package handlers

import (
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
)

type EventHandler func(tick game.Tick, event *envelope.E) error
type TickHandler func(tick game.Tick) error

type EventProcessor interface {
	GetEventHandlers() map[pb.OneOfEvent][]EventHandler
	GetTickHandlers() []TickHandler
}

func Get() []EventProcessor {
	// TODO instantiate all processors for a shard in here.
	//  Likely pass shard and it's chunks into this thing to bind them to internal state of handlers.
	return nil
}
