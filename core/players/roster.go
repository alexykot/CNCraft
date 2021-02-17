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
    "github.com/alexykot/cncraft/pkg/game/entities"
    "github.com/alexykot/cncraft/pkg/game/player"
)

// Roster holds the map of all players logged into this server.
type Roster struct {
    log *zap.Logger
    ps  nats.PubSub
    mu  sync.Mutex
    repo *repo

    players map[uuid.UUID]*Player // ID of the user always matches ID of the corresponding connection.
}

func NewRoster(log *zap.Logger, ps nats.PubSub, db *sql.DB) *Roster {
    return &Roster{
        log:  log,
        ps:   ps,
        repo: newRepo(db),

        players: make(map[uuid.UUID]*Player),
    }
}

func (r *Roster) AddPlayer(userID uuid.UUID, username string) *Player {
    r.mu.Lock()
    defer r.mu.Unlock()

    if existing, ok := r.players[userID]; ok {
        r.log.Error("player already exists", zap.String("id", userID.String()),
            zap.String("existing", existing.Username), zap.String("new", username))
        return nil
    }

    // TODO whole user state hardcoded until persistence is properly implemented
    p := &Player{
        ID:       userID,
        PC:       entities.NewPC(username, player.MaxHealth),
        Username: username,
        Settings: player.Settings{
            ViewDistance: 7,
            FlyingSpeed:  0.05,
            FoVModifier:  0.1,
        },
        State: player.State{
            CurrentHotbarSlot: player.Slot0,
            CurrentLocation: data.Location{
                PositionF: data.PositionF{
                    X: 15,
                    Y: 32,
                    Z: 15,
                },
            },
        },
    }

    r.players[userID] = p
    return p
}

func (r *Roster) GetPlayer(userID uuid.UUID) (*Player, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()

    p, ok := r.players[userID]
    return p, ok
}

func (r *Roster) SetPlayerPos(userID uuid.UUID, position data.PositionF) {
    p, ok := r.GetPlayer(userID)
    if !ok {
        return
    }
    p.SetPosition(position)

    r.publishPlayerPosUpdate(p)
}

// RegisterHandlers creates subscriptions to all relevant global subjects.
func (r *Roster) RegisterHandlers() error {
    if err := r.ps.Subscribe(subj.MkPlayerJoined(), r.playerJoinedHandler); err != nil {
        return fmt.Errorf("failed to subscribe for joined users: %w", err)
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

func (r *Roster) publishPlayerPosUpdate(p *Player) {
    lope := envelope.PlayerPosUpdate(&pb.PlayerPosUpdate{
        Id: p.ID.String(),
        Pos: &pb.Position{
            X: p.State.CurrentLocation.X,
            Y: p.State.CurrentLocation.Y,
            Z: p.State.CurrentLocation.Z,
        },
    })
    if err := r.ps.Publish(subj.MkPlayerPosUpdate(), lope); err != nil {
        r.log.Error("failed to publish position update", zap.Error(err))
    }
}
