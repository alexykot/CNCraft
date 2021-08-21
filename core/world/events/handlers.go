package events

import (
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type EventHandler func(tick game.Tick, event *envelope.E) (map[subj.Subj][]*envelope.E, error)
type TickHandler func(tick game.Tick) (map[subj.Subj][]*envelope.E, error)

type Handler interface {
	Name() string
	GetTickHandler() TickHandler
	GetEventHandlers() map[pb.OneOfEvent]EventHandler
}

func NewHandlers(chunks []level.Chunk, roster players.Roster) []Handler {
	return []Handler{
		newDigger(chunks, roster),
	}
}
