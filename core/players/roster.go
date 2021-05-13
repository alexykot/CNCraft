package players

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/game/data"
)

// Roster holds the map of all players logged into this server.
type Roster struct {
	log  *zap.Logger
	ps   nats.PubSub
	mu   sync.Mutex
	repo *repo

	players map[uuid.UUID]*Player
}

func NewRoster(log *zap.Logger, ps nats.PubSub, db *sql.DB) *Roster {
	return &Roster{
		log:  log,
		ps:   ps,
		repo: newRepo(db),

		// DEBT this is a single point of synchronisation for all currently connected players. This will break
		//  down in multi-node cluster setup, and cluster-global synchronisation will be needed instead.
		players: make(map[uuid.UUID]*Player),
	}
}

func (r *Roster) AddPlayer(connID uuid.UUID, username string) (*Player, error) {
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
	if p, isNew, err = r.repo.InitPlayer(username, connID); err != nil {
		return nil, fmt.Errorf("failed to init player: %w", err)
	}

	if isNew {
		r.log.Debug("new player joined for the first time", zap.String("name", p.Username))
		r.publishNewPlayerJoined(p)
	} else {
		r.log.Debug("rejoining player loaded", zap.String("name", p.Username))
	}

	r.players[p.ID] = p
	return p, nil
}

func (r *Roster) GetPlayerByConnID(connID uuid.UUID) (*Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.players {
		if p.ConnID == connID {
			return p, true
		}
	}

	return nil, false
}

func (r *Roster) GetPlayerIDByConnID(connID uuid.UUID) (uuid.UUID, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.players {
		if p.ConnID == connID {
			return p.ID, true
		}
	}

	return uuid.UUID{}, false
}

// SetPlayerSpatial - *bool type is used here to separate true/false valus from an absence of value to update
func (r *Roster) SetPlayerSpatial(connID uuid.UUID, position *data.PositionF, rotation *data.RotationF, onGround *bool) {
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

func (r *Roster) SetPlayerHeldItem(connID uuid.UUID, heldItem uint8) {
	p, ok := r.GetPlayerByConnID(connID)
	if !ok {
		return
	}
	p.State.Inventory.CurrentHotbarSlot = heldItem
	r.publishPlayerInventoryUpdate(p)
}

func (r *Roster) PlayerInventoryChanged(connID uuid.UUID) {
	p, ok := r.GetPlayerByConnID(connID)
	if !ok {
		return
	}
	r.publishPlayerInventoryUpdate(p)
}

// RegisterHandlers creates subscriptions to all relevant global subjects.
func (r *Roster) RegisterHandlers() error {
	if err := r.ps.Subscribe(subj.MkPlayerJoined(), r.playerJoinedHandler); err != nil {
		return fmt.Errorf("failed to subscribe for joining users: %w", err)
	}

	if err := r.ps.Subscribe(subj.MkPlayerLeft(), r.playerLeftHandler); err != nil {
		return fmt.Errorf("failed to subscribe for leaving users: %w", err)
	}

	r.log.Info("global player handlers registered")
	return nil
}

func (r *Roster) publishNewPlayerJoined(p *Player) {
	lope := envelope.NewPlayerJoined(&pb.NewPlayerJoined{
		PlayerId: p.ID.String(),
		ConnId:   p.ConnID.String(),
		Username: p.Username,
		Pos: &pb.Position{
			X: p.State.Location.X,
			Y: p.State.Location.Y,
			Z: p.State.Location.Z,
		},
	})
	if err := r.ps.Publish(subj.MkNewPlayerJoined(), lope); err != nil {
		r.log.Error("failed to publish position update", zap.Error(err))
	}
}

// TODO not sure if this is needed
func (r *Roster) playerJoinedHandler(_ *envelope.E) {}

func (r *Roster) playerLeftHandler(lope *envelope.E) {
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

func (r *Roster) publishPlayerSpatialUpdate(p *Player) {
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

func (r *Roster) publishPlayerInventoryUpdate(p *Player) {
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