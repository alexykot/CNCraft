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
	log     *zap.Logger
	control chan control.Command
	ps      nats.PubSub

	mu       sync.Mutex
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

func (k *KeepAliver) Start(ctx context.Context) {
	go k.tick(ctx)
}

func (k *KeepAliver) AddAliveConn(connID uuid.UUID) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, ok := k.theyLive[connID]; ok {
		return
	}

	k.theyLive[connID] = time.Now().Unix() // assume initially client was last seen right now
}

func (k *KeepAliver) DropDeadConn(connID uuid.UUID) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.pronounceDead(connID)

	if _, ok := k.theyLive[connID]; !ok {
		return
	}

	delete(k.theyLive, connID)
}

func (k *KeepAliver) tick(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil { // stop the server if the keepAliver goroutine dies
			k.control <- control.Command{Signal: control.FAIL, Message: fmt.Sprintf("keepaliver panicked: %v", r)}
		}
	}()

	aliveTicker := time.NewTicker(keepAliveInterval)
	for {
		select {
		case <-ctx.Done():
			k.log.Info("keepAliver stopped")
			return
		case nowTime := <-aliveTicker.C:
			nowTimeNix := nowTime.Unix()
			for connID, lastReceived := range k.theyLive {
				if nowTimeNix-lastReceived <= int64(keepAliveMaxLag.Seconds()) {
					k.transmitKeepAlive(connID, nowTime)
				} else {
					k.pronounceDead(connID)
				}
			}
		}
	}

	k.control <- control.Command{Signal: control.FAIL, Message: fmt.Sprintf("keepAliver stopped unexpectedly")}
}

func (k *KeepAliver) pronounceDead(connID uuid.UUID) {
	lope := envelope.CloseConn(&pb.CloseConn{
		Id:    connID.String(),
		State: pb.ConnState_PLAY, // keepAliver is only active in the Play state
	})

	if err := k.ps.Publish(subj.MkConnClose(), lope); err != nil {
		k.log.Error("failed to publish CloseConn", zap.Error(err), zap.String("conn", connID.String()))
	}
	k.log.Debug("connection pronounced dead", zap.String("conn", connID.String()))
}

func (k *KeepAliver) receiveKeepAlive(connID uuid.UUID, liveID int64) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, ok := k.theyLive[connID]; !ok {
		return
	}

	k.theyLive[connID] = liveID
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
