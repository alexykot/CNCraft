package events

import (
	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/level"
)

type digger struct {
	chunks []level.Chunk
	roster *players.Roster
}

func newDigger(chunks []level.Chunk, roster *players.Roster) Handler {
	return &digger{
		chunks: chunks,
		roster: roster,
	}
}

func (d *digger) Name() string { return "digger" }

func (d *digger) GetTickHandler() TickHandler { return nil }

func (d *digger) GetEventHandlers() map[pb.OneOfEvent]EventHandler {
	return map[pb.OneOfEvent]EventHandler{
		pb.Event_PlayerDigging: d.handlePlayerDiggingEvent,
	}
}

func (d *digger) handlePlayerDiggingEvent(_ game.Tick, event *envelope.E) (map[uuid.UUID][]*envelope.E, error) {
	shardEvent := event.GetShardEvent()
	if shardEvent == nil {
		return nil, nil // silently ignore irrelevant events
	}

	playerDigging := shardEvent.GetPlayerDigging()
	if playerDigging == nil {
		return nil, nil // silently ignore irrelevant events
	}

	return nil, nil
}

func (d *digger) digIsLegal(playerID uuid.UUID, digCoord data.PositionF) (bool, error) {
	return true, nil
}
