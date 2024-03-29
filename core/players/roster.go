//go:generate mockgen -package mocks -source=roster.go -destination=mocks/mocks.go Roster

// Package players contains implementation for players list (Roster), the player repo that allows to load players
// from persistence and the Player details itself.
package players

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game/data"
)

// Roster handles the map of all players logged into this server.
type Roster interface {
	Start(ctx context.Context)
	AddPlayer(username string, connID, dimensionID uuid.UUID) (*Player, error)
	GetPlayerByConnID(connID uuid.UUID) (*Player, bool)
	GetPlayerIDByConnID(connID uuid.UUID) (uuid.UUID, bool)
	SetPlayerSpatial(connID uuid.UUID, position *data.PositionF, rotation *data.RotationF, onGround *bool)
	SetPlayerHeldItem(connID uuid.UUID, heldItem uint8)
	PlayerInventoryChanged(connID uuid.UUID)
}

type roster struct {
	control chan control.Command
	log     *zap.Logger
	ps      nats.PubSub
	mu      sync.Mutex
	repo    *repo

	players map[uuid.UUID]*Player
}

func NewRoster(log, windowLog *zap.Logger, ctrlChan chan control.Command, ps nats.PubSub, db *sql.DB) Roster {
	return &roster{
		control: ctrlChan,
		log:     log,
		ps:      ps,
		repo:    newRepo(windowLog, db), // DEBT this really should be instantiated separately and passed in as option.

		// DEBT this is a single point of synchronisation for all currently connected players. This will break
		//  down in multi-node cluster setup, and cluster-global synchronisation will be needed instead.
		//
		//  On GetPlayer calls rosters of different nodes in the cluster will need to synchronously request player
		//  state from the node where the given player is connected. This ideally should be done via
		//  request-response facility inside NATS, but this will need further research.
		//
		//  On player inventory updates (e.g. on PlayerInventoryChanged when player picks up something from the ground)
		//  the local roster will need to transmit async updates via player-specific channels to the remote roster
		//  that tracks that player. Remote roster will receive the update, update local memory state and transmit
		//  update for the persistent state.
		players: make(map[uuid.UUID]*Player),
	}
}

func (r *roster) AddPlayer(username string, connID, dimensionID uuid.UUID) (*Player, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// DEBT the IDs are getting recreated for every new connection, so the ID will not work for identifying
	//  the duplicates and we have to check through the whole map for usernames. Suboptimal.
	//  Afterthought: ID is actually useless and the username is in fact a globally unique ID of the player.
	var found bool
	for _, existing := range r.players {
		if existing.Username == username {
			found = true
			break
		}
	}

	if found {
		return nil, fmt.Errorf("player %s already exists", username)
	}

	var err error
	var isNew bool
	var p *Player
	if p, isNew, err = r.repo.InitPlayer(username, connID, dimensionID); err != nil {
		return nil, fmt.Errorf("failed to init player: %w", err)
	}

	if isNew {
		r.log.Debug("new player joined for the first time", zap.String("name", p.Username))
	} else {
		r.log.Debug("rejoining player loaded", zap.String("name", p.Username))
	}
	r.publishPlayerJoined(p)

	r.players[p.ID] = p
	return p, nil
}

func (r *roster) GetPlayerByConnID(connID uuid.UUID) (*Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.players {
		if p.ConnID == connID {
			return p, true
		}
	}

	return nil, false
}

func (r *roster) GetPlayerIDByConnID(connID uuid.UUID) (uuid.UUID, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.players {
		if p.ConnID == connID {
			return p.ID, true
		}
	}

	return uuid.UUID{}, false
}

// SetPlayerSpatial - pointer types are used here to separate possible default values from an absence of value to update.
func (r *roster) SetPlayerSpatial(connID uuid.UUID, position *data.PositionF, rotation *data.RotationF, onGround *bool) {
	p, ok := r.GetPlayerByConnID(connID)
	if !ok {
		return
	}

	if position != nil {
		p.SetPosition(*position)
	}

	if rotation != nil {
		p.SetRotation(*rotation)
	}

	if onGround != nil {
		p.SetOnGround(*onGround)
	}

	r.publishPlayerSpatialUpdate(p)
}

func (r *roster) SetPlayerHeldItem(connID uuid.UUID, heldItem uint8) {
	p, ok := r.GetPlayerByConnID(connID)
	if !ok {
		return
	}
	p.State.Inventory.CurrentHotbarSlot = heldItem
	r.publishPlayerInventoryUpdate(p)
}

func (r *roster) PlayerInventoryChanged(connID uuid.UUID) {
	p, ok := r.GetPlayerByConnID(connID)
	if !ok {
		return
	}
	r.publishPlayerInventoryUpdate(p)
}

func (r *roster) Start(_ context.Context) {
	// DEBT When cluster mode will be developed - this will also need to start a context watching goroutine
	//  and unsubscribe from the player channels.
	//
	// For now Roster does not signal readiness as it is ready as soon as it's started, and has nothing to stop.
	if err := r.registerHandlers(); err != nil {
		r.signal(control.FAILED, err)
	}

	r.log.Info("roster started")
}

// registerHandlers creates subscriptions to all relevant global subjects.
func (r *roster) registerHandlers() error {
	if err := r.ps.Subscribe(subj.MkPlayerLeft(), r.playerLeftHandler); err != nil {
		return fmt.Errorf("failed to start roster: failed to subscribe for leaving users: %w", err)
	}

	r.log.Debug("global player handlers registered")
	return nil
}

func (r *roster) publishPlayerJoined(p *Player) {
	lope := envelope.PlayerJoined(&pb.PlayerJoined{
		PlayerId:    p.ID.String(),
		ConnId:      p.ConnID.String(),
		DimensionId: p.State.Dimension.String(),
		Username:    p.Username,
		Pos: &pb.Position{
			X: p.State.Location.X,
			Y: p.State.Location.Y,
			Z: p.State.Location.Z,
		},
	})
	if err := r.ps.Publish(subj.MkPlayerJoined(), lope); err != nil {
		r.log.Error("failed to publish position update", zap.Error(err))
	}
}

func (r *roster) playerLeftHandler(lope *envelope.E) {
	left := lope.GetPlayerLeft()
	if left == nil {
		r.log.Error("failed to parse envelope - no PlayerLeft inside", zap.Any("envelope", lope))
		return
	}

	playerID, err := uuid.Parse(left.PlayerId)
	if err != nil {
		r.log.Error("failed to parse conn ID as UUID", zap.Any("id", left.PlayerId))
	}

	r.log.Debug("player leaving", zap.String("name", r.players[playerID].Username))
	delete(r.players, playerID)
}

func (r *roster) publishPlayerSpatialUpdate(p *Player) {
	lope := envelope.PlayerSpatialUpdate(&pb.PlayerSpatialUpdate{
		PlayerId: p.ID.String(),
		Pos: &pb.Position{
			X: p.State.Location.X,
			Y: p.State.Location.Y,
			Z: p.State.Location.Z,
		},
		Rot: &pb.Rotation{
			Yaw:   p.State.Location.Yaw,
			Pitch: p.State.Location.Pitch,
		},
		OnGround: p.State.Location.OnGround,
	})
	if err := r.ps.Publish(subj.MkPlayerSpatialUpdate(), lope); err != nil {
		r.log.Error("failed to publish position update", zap.Error(err))
	}
}

func (r *roster) publishPlayerInventoryUpdate(p *Player) {
	update := &pb.PlayerInventoryUpdate{PlayerId: p.ID.String(), CurrentHotbar: int32(p.State.Inventory.CurrentHotbarSlot)}
	for i, item := range p.State.Inventory.ToArray() {
		if item.IsPresent {
			update.Inventory = append(update.Inventory, &pb.InventoryItem{
				SlotId:    int32(i),
				ItemId:    int32(item.ItemID),
				ItemCount: int32(item.ItemCount),
			})
		}
	}

	if err := r.ps.Publish(subj.MkPlayerInventoryUpdate(), envelope.PlayerInventoryUpdate(update)); err != nil {
		r.log.Error("failed to publish inventory update", zap.Error(err))
	}
}

func (r *roster) signal(state control.ComponentState, err error) {
	r.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.ROSTER,
		State:     state,
		Err:       err,
	}
}
