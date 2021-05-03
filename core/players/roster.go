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

	players map[uuid.UUID]*Player // ID of the user always matches ID of the corresponding connection.
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

func (r *Roster) AddPlayer(userID uuid.UUID, username string) (*Player, error) {
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
	if r.players[userID], isNew, err = r.repo.InitPlayer(userID, username); err != nil {
		return nil, fmt.Errorf("failed to init player: %w", err)
	}

	if isNew {
		r.log.Debug("new player joined for the first time", zap.Any("player", r.players[userID]))
		r.publishNewPlayerJoined(r.players[userID])
	} else {
		r.log.Debug("rejoining player loaded", zap.Any("player", r.players[userID]))
	}

	return r.players[userID], nil
}

func (r *Roster) GetPlayer(userID uuid.UUID) (*Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.players[userID]
	return p, ok
}

func (r *Roster) SetPlayerSpatial(userID uuid.UUID, position *data.PositionF, rotation *data.RotationF, onGround *bool) {
	p, ok := r.GetPlayer(userID)
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

func (r *Roster) playerJoinedHandler(lope *envelope.E) {
	// joined := lope.GetPlayerJoined()
	// if joined == nil {
	//	r.log.Error("failed to parse envelope - no JoinedPlayer inside", zap.Any("envelope", lope))
	//	return
	// }
	//
	// userId, err := uuid.Parse(joined.Id)
	// if err != nil {
	//	r.log.Error("failed to parse user ID as UUID", zap.Any("id", joined.Id))
	//	return
	// }

	// TODO implement this
}

func (r *Roster) playerLeftHandler(lope *envelope.E) {
	left := lope.GetPlayerLeft()
	if left == nil {
		r.log.Error("failed to parse envelope - no PlayerLeft inside", zap.Any("envelope", lope))
		return
	}

	userID, err := uuid.Parse(left.Id)
	if err != nil {
		r.log.Error("failed to parse user ID as UUID", zap.Any("id", left.Id))
	}

	r.log.Debug("player leaving", zap.Any("player", r.players[userID]))
	delete(r.players, userID)
}

func (r *Roster) publishPlayerSpatialUpdate(p *Player) {
	lope := envelope.PlayerSpatialUpdate(&pb.PlayerSpatialUpdate{
		Id: p.ID.String(),
		Pos: &pb.Position{
			X: p.State.CurrentLocation.X,
			Y: p.State.CurrentLocation.Y,
			Z: p.State.CurrentLocation.Z,
		},
		Rot: &pb.Rotation{
			Yaw:   p.State.CurrentLocation.Yaw,
			Pitch: p.State.CurrentLocation.Pitch,
		},
		OnGround: p.State.CurrentLocation.OnGround,
	})
	if err := r.ps.Publish(subj.MkPlayerSpatialUpdate(), lope); err != nil {
		r.log.Error("failed to publish position update", zap.Error(err))
	}
}

func (r *Roster) publishNewPlayerJoined(p *Player) {
	lope := envelope.NewPlayerJoined(&pb.NewPlayerJoined{
		Id:       p.ID.String(),
		Username: p.Username,
		Pos: &pb.Position{
			X: p.State.CurrentLocation.X,
			Y: p.State.CurrentLocation.Y,
			Z: p.State.CurrentLocation.Z,
		},
	})
	if err := r.ps.Publish(subj.MkNewPlayerJoined(), lope); err != nil {
		r.log.Error("failed to publish position update", zap.Error(err))
	}
}
