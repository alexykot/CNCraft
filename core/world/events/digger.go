package events

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/players"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/level"
	"github.com/alexykot/cncraft/pkg/game/player"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// maxDigDistance - very crude way to determine if digging is within legal distance.
// more details at https://wiki.vg/index.php?title=Protocol&oldid=16676#Player_Digging
const maxDigDistance = 7.5

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
		return nil, errors.New("provided event is not a shardEvent")
	}

	playerDigging := shardEvent.GetPlayerDigging()
	if playerDigging == nil {
		return nil, errors.New("provided event is not a playerDigging event")
	}

	blockPosF := data.PositionF{
		X: playerDigging.Pos.X,
		Y: playerDigging.Pos.Y,
		Z: playerDigging.Pos.Z,
	}
	blockPosI := blockPosF.ToInt()

	playerID, err := uuid.FromBytes([]byte(playerDigging.PlayerId))
	if err != nil {
		return nil, fmt.Errorf("PlayerId invalid: %w", err)
	}

	digAction := player.DiggingActionFromPb(playerDigging.Action)

	block, err := d.getBlockAtCoords(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("failed to find block at coords %s: %w", blockPosI.String(), err)
	}

	isLegal, err := d.digIsLegal(playerID, block, blockPosF)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if dig is legal for player %s, coords %s: %w", playerID, blockPosF.String(), err)
	} else if !isLegal {
		return map[uuid.UUID][]*envelope.E{
			playerID: {d.response(false, blockPosI, block, digAction)},
		}, nil
	}

	return nil, nil
}

func (d *digger) digIsLegal(playerID uuid.UUID, block level.Block, digCoordF data.PositionF) (bool, error) {
	pl, ok := d.roster.GetPlayerByConnID(playerID)
	if !ok {
		return false, fmt.Errorf("player %s not found", playerID.String())
	}

	playerLoc := pl.GetLocation()
	xDistance := math.Abs(playerLoc.PositionF.X - digCoordF.X)
	yDistance := math.Abs(playerLoc.PositionF.Y - digCoordF.Y)
	zDistance := math.Abs(playerLoc.PositionF.Z - digCoordF.Z)
	// DEBT likely not the correct algorithm according to Notchian server, but good enough for now.
	if xDistance > maxDigDistance || yDistance > maxDigDistance || zDistance > maxDigDistance {
		return false, nil
	}

	tool := pl.GetState().Inventory.GetCurrentTool()
	if !block.ID().IsDiggable(tool.ItemID) {
		return false, nil
	}

	return true, nil
}

func (d *digger) getBlockAtCoords(digCoordI data.PositionI) (level.Block, error) {
	chunk, err := d.getChunkAtCoords(digCoordI)
	if err != nil {
		return nil, fmt.Errorf("no chunk available for given coords, x:y:z %s", digCoordI.String())
	}

	block, err := chunk.GetGlobalBlock(digCoordI)
	if err != nil {
		return nil, fmt.Errorf("block not found in the chunk, x:y:z %s", digCoordI.String())
	}

	return block, nil
}

func (d *digger) getChunkAtCoords(coord data.PositionI) (level.Chunk, error) {
	chunkID := level.FindChunkID(coord)
	for _, chunk := range d.chunks {
		if chunk.ID() == chunkID {
			return chunk, nil
		}
	}

	// DEBT ideally add shardID to this error
	return nil, fmt.Errorf("chunk %s not found in shard", chunkID.String())
}

func (d *digger) response(ack bool, blockPos data.PositionI, block level.Block, action player.DiggingAction) *envelope.E {
	cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CAcknowledgePlayerDigging)
	nack := cpacket.(*protocol.CPacketAcknowledgePlayerDigging)

	nack.Location = blockPos
	nack.Block = block.ID()
	nack.Status = action
	nack.Successful = ack

	return envelope.MkCpacketEnvelope(nack)
}
