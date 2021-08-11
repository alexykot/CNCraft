package events

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

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
	sync.RWMutex

	chunks     []level.Chunk
	activeDigs map[data.PositionI]activeDig // block positions and active dig details
	roster     *players.Roster
}

type activeDig struct {
	startTime   game.Tick // tick time when digging started
	diggerCount int       // number of players simultaneously digging the block
}

func newDigger(chunks []level.Chunk, roster *players.Roster) Handler {
	return &digger{
		chunks: chunks,
		roster: roster,
	}
}

func (d *digger) Name() string { return "digger" }

func (d *digger) GetTickHandler() TickHandler {
	return nil // TODO this needs to handle ongoing active digs
}

func (d *digger) GetEventHandlers() map[pb.OneOfEvent]EventHandler {
	return map[pb.OneOfEvent]EventHandler{
		pb.Event_PlayerDigging: d.handlePlayerDiggingEvent,
	}
}

func (d *digger) handlePlayerDiggingEvent(tick game.Tick, event *envelope.E) (map[uuid.UUID][]*envelope.E, error) {
	shardEvent := event.GetShardEvent()
	if shardEvent == nil {
		return nil, errors.New("provided event is not a shardEvent")
	}

	playerDigging := shardEvent.GetPlayerDigging()
	if playerDigging == nil {
		return nil, errors.New("provided event is not a playerDigging event")
	}

	playerID, err := uuid.FromBytes([]byte(playerDigging.PlayerId))
	if err != nil {
		return nil, fmt.Errorf("PlayerId invalid: %w", err)
	}

	digAction := player.DiggingActionFromPb(playerDigging.Action)

	res, err := d.handleDig(tick, digAction, playerID, data.PositionFFromPb(playerDigging.Pos))
	if err != nil {
		return nil, fmt.Errorf("failed to handle player digging: %w", err)
	}
	return res, nil
}

func (d *digger) handleDig(tick game.Tick, digAction player.DiggingAction, playerID uuid.UUID, blockPosF data.PositionF) (map[uuid.UUID][]*envelope.E, error) {
	switch digAction {
	case player.StartedDigging:
		return d.handleStartedDigging(tick, playerID, blockPosF)
	case player.FinishedDigging:
		return d.handleFinishedDigging(tick, playerID, blockPosF)
	case player.CancelledDigging:
		return d.handleCancelledDigging(tick, playerID, blockPosF)
	default:
		return nil, fmt.Errorf("unsupported digging action %s", digAction.String())
	}
}

// handleStartedDigging handles digging starts sent by clients. It accounts for multiple clients potentially trying
// to dig the same block at the same time.
func (d *digger) handleStartedDigging(tick game.Tick, playerID uuid.UUID, blockPosF data.PositionF) (map[uuid.UUID][]*envelope.E, error) {
	blockPosI := blockPosF.ToInt()
	block, err := d.getBlockAtCoords(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("failed to find block at coords %s: %w", blockPosI.String(), err)
	}

	_, isLegal, err := d.digIsLegal(playerID, block, blockPosF)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if dig is legal for player %s, coords %s: %w", playerID, blockPosF.String(), err)
	} else if !isLegal {
		return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(false, blockPosI, block, player.StartedDigging)}}, nil
	}

	d.Lock()
	dig, ok := d.activeDigs[blockPosI]
	if !ok || dig.diggerCount == 0 { // new digging effort starting
		dig.startTime = tick
		dig.diggerCount = 1
	} else { // player joining ongoing digging effort
		dig.diggerCount++
	}
	d.activeDigs[blockPosI] = dig
	d.Unlock()

	return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(true, blockPosI, block, player.StartedDigging)}}, nil
}

func (d *digger) handleFinishedDigging(tick game.Tick, playerID uuid.UUID, blockPosF data.PositionF) (map[uuid.UUID][]*envelope.E, error) {
	blockPosI := blockPosF.ToInt()
	block, err := d.getBlockAtCoords(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("failed to find block at coords %s: %w", blockPosI.String(), err)
	}

	digDuration, isLegal, err := d.digIsLegal(playerID, block, blockPosF)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if dig is legal for player %s, coords %s: %w", playerID, blockPosF.String(), err)
	} else if !isLegal {
		return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(false, blockPosI, block, player.FinishedDigging)}}, nil
	}

	d.Lock()
	dig, ok := d.activeDigs[blockPosI]
	if !ok || dig.diggerCount == 0 { // no digging was actually happening on this block, NAck.
		d.Unlock()
		return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(false, blockPosI, block, player.FinishedDigging)}}, nil
	}

	// enough time has passed to dig out the given block with given tool, Ack and broadcast world state update
	if tick.AsTime().Sub(dig.startTime.AsTime()) >= digDuration {
		delete(d.activeDigs, blockPosI) // block digged successfully, all digging now stops
		d.Unlock()

		return map[uuid.UUID][]*envelope.E{playerID: {
			d.ackResponse(true, blockPosI, block, player.FinishedDigging),
			// TODO broadcast placed block disappeared
			// TODO broadcast block entity spawned
		}}, nil
	} else { // enough time has passed to dig out the given block with given tool, NAck
		d.Unlock()
		return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(false, blockPosI, block, player.FinishedDigging)}}, nil
	}
}

// handleCancelledDigging handles digging cancellations sent by clients. It does not consider/handle abandoned digs,
// e.g. where client connection is lost. Those are handled during regular ticks where overtime digs are quietly removed.
func (d *digger) handleCancelledDigging(_ game.Tick, playerID uuid.UUID, blockPosF data.PositionF) (map[uuid.UUID][]*envelope.E, error) {
	blockPosI := blockPosF.ToInt()

	d.Lock()
	dig, ok := d.activeDigs[blockPosI]
	if ok && dig.diggerCount > 0 { // there is actual digging ongoing
		dig.diggerCount--
		if dig.diggerCount == 0 { // nobody is digging this block anymore
			delete(d.activeDigs, blockPosI)
		} else {
			d.activeDigs[blockPosI] = dig // somebody is still digging it it seems
		}
	}
	d.Unlock()

	block, err := d.getBlockAtCoords(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("failed to find block at coords %s: %w", blockPosI.String(), err)
	}

	return map[uuid.UUID][]*envelope.E{playerID: {d.ackResponse(true, blockPosI, block, player.CancelledDigging)}}, nil
}

func (d *digger) digIsLegal(playerID uuid.UUID, block level.Block, blockPosF data.PositionF) (digDuration time.Duration, isLegal bool, err error) {
	pl, ok := d.roster.GetPlayerByConnID(playerID)
	if !ok {
		return 0, false, fmt.Errorf("player %s not found", playerID.String())
	}

	playerLoc := pl.GetLocation()
	xDistance := math.Abs(playerLoc.PositionF.X - blockPosF.X)
	yDistance := math.Abs(playerLoc.PositionF.Y - blockPosF.Y)
	zDistance := math.Abs(playerLoc.PositionF.Z - blockPosF.Z)
	// DEBT likely not the correct algorithm according to Notchian server, but good enough for now.
	if xDistance > maxDigDistance || yDistance > maxDigDistance || zDistance > maxDigDistance {
		return 0, false, nil
	}

	tool := pl.GetState().Inventory.GetCurrentTool()
	if !block.ID().IsDiggable(tool.ItemID) {
		return 0, false, nil
	}

	return block.ID().DigTime(tool.ItemID), true, nil
}

func (d *digger) getBlockAtCoords(blockPosI data.PositionI) (level.Block, error) {
	chunk, err := d.getChunkAtCoords(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("no chunk available for given coords, x:y:z %s", blockPosI.String())
	}

	block, err := chunk.GetGlobalBlock(blockPosI)
	if err != nil {
		return nil, fmt.Errorf("block not found in the chunk, x:y:z %s", blockPosI.String())
	}

	return block, nil
}

func (d *digger) getChunkAtCoords(blockPosI data.PositionI) (level.Chunk, error) {
	chunkID := level.FindChunkID(blockPosI)
	for _, chunk := range d.chunks {
		if chunk.ID() == chunkID {
			return chunk, nil
		}
	}

	// DEBT ideally add shardID to this error
	return nil, fmt.Errorf("chunk %s not found in shard", chunkID.String())
}

func (d *digger) ackResponse(ack bool, blockPosI data.PositionI, block level.Block, action player.DiggingAction) *envelope.E {
	cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CAcknowledgePlayerDigging)
	nack := cpacket.(*protocol.CPacketAcknowledgePlayerDigging)

	nack.Location = blockPosI
	nack.Block = block.ID()
	nack.Status = action
	nack.Successful = ack

	return envelope.MkCpacketEnvelope(nack)
}
