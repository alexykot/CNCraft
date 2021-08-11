package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
)

const keepAliveInterval = time.Second * 5
const keepAliveMaxLag = time.Second * 30

// KeepAliver maintains the list of connections to keep alive. Also it controls the keepalive responses from clients
//  and triggers disconnection for clients that failed to respond.
type KeepAliver struct {
	sync.RWMutex

	log     *zap.Logger
	control chan control.Command
	ps      nats.PubSub

	theyLive map[uuid.UUID]int64 // latest timestamp of the keepalive response received from the given connection
}

func NewKeepAliver(control chan control.Command, ps nats.PubSub, log *zap.Logger) *KeepAliver {
	return &KeepAliver{
		control:  control,
		ps:       ps,
		log:      log,
		theyLive: make(map[uuid.UUID]int64),
	}
}

func (k *KeepAliver) Start(ctx context.Context) error {
	if err := k.ps.Subscribe(subj.MkConnClosed(), k.connClosedHandler); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subj.MkConnClosed().String(), err)
	}
	k.log.Debug("registered conn closing handler")

	go k.tick(ctx)
	k.log.Debug("started keepalive ticker")
	k.signal(control.READY, nil)

	return nil
}

func (k *KeepAliver) AddAliveConn(connID uuid.UUID) {
	k.Lock()
	defer k.Unlock()

	if _, ok := k.theyLive[connID]; ok {
		return
	}

	k.theyLive[connID] = time.Now().Unix() // assume initially client was last seen right now
}

func (k *KeepAliver) connClosedHandler(lope *envelope.E) {
	closeConn := lope.GetCloseConn()
	if closeConn == nil {
		k.log.Error("failed to parse envelope: there is no closeConn inside", zap.Any("envelope", lope))
		return
	}

	connID, err := uuid.Parse(closeConn.ConnId)
	if err != nil {
		k.log.Error("failed to parse conn ID as UUID", zap.String("id", closeConn.ConnId))
		return
	}
	k.log.Debug("connection closed", zap.String("conn", closeConn.ConnId))

	k.pronounceDead(connID, false)
}

func (k *KeepAliver) tick(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			k.signal(control.FAILED, fmt.Errorf("keepaliver panicked: %v", r))
		}
	}()

	aliveTicker := time.NewTicker(keepAliveInterval)
	for {
		select {
		case <-ctx.Done():
			k.signal(control.STOPPED, nil)
			return
		case nowTime := <-aliveTicker.C:
			nowTimeNix := nowTime.Unix()
			k.RLock()
			for connID, lastReceived := range k.theyLive {
				if nowTimeNix-lastReceived <= int64(keepAliveMaxLag.Seconds()) {
					k.transmitKeepAlive(connID, nowTime)
				} else {
					k.pronounceDead(connID, true)
					delete(k.theyLive, connID)
				}
			}
			k.RUnlock()
		}
	}
}

// pronounceDead removes given connection from the list of connections that need to be kept alive, and optionally
// broadcasts this event.
//
// DEBT this does not actually close the connection itself since it doesn't have access to the Connection object.
//  This may cause resource leaks.
func (k *KeepAliver) pronounceDead(connID uuid.UUID, publish bool) {
	k.Lock()
	if _, ok := k.theyLive[connID]; ok {
		delete(k.theyLive, connID)
	}
	k.Unlock()

	// Publish will be true only if the keepaliver itself recognised the failure and needs to broadcast it.
	// If some other part has recognised it first (TCP connection closed, or clean
	// client disconnect packet sent) - no need to publish as this was already received via ConnClosed channel.
	if publish {
		lope := envelope.CloseConn(&pb.CloseConn{
			ConnId: connID.String(),
			State:  pb.ConnState_PLAY, // keepAliver is only active in the Play state
		})

		if err := k.ps.Publish(subj.MkConnClosed(), lope); err != nil {
			k.log.Error("failed to publish CloseConn", zap.Error(err), zap.String("conn", connID.String()))
			return
		}
		k.log.Debug("connection pronounced dead", zap.String("conn", connID.String()))
	}
}

func (k *KeepAliver) receiveKeepAlive(connID uuid.UUID, liveID int64) {
	k.Lock()
	if _, ok := k.theyLive[connID]; !ok {
		return
	}

	k.theyLive[connID] = liveID
	k.Unlock()
}

func (k *KeepAliver) transmitKeepAlive(connID uuid.UUID, timeNow time.Time) {
	cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CKeepAlive)
	keepAlive := cpacket.(*protocol.CPacketKeepAlive)
	keepAlive.KeepAliveID = timeNow.Unix()

	bufOut := buffer.New()
	keepAlive.Push(bufOut)

	lope := envelope.CPacket(&pb.CPacket{
		Bytes:      bufOut.Bytes(),
		PacketType: keepAlive.Type().Value(),
	})

	if err := k.ps.Publish(subj.MkConnTransmit(connID), lope); err != nil {
		k.log.Error("failed to publish conn keepalive CPacket", zap.Error(err), zap.String("conn", connID.String()))
	}
}

func (k *KeepAliver) signal(state control.ComponentState, err error) {
	k.control <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.KEEPALIVER,
		State:     state,
		Err:       err,
	}
}
