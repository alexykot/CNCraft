//go:generate stringer -type=DiggingAction status.go

package player

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
)

type ClientStatusAction int

const (
	Respawn ClientStatusAction = iota
	Request
)

type DiggingAction int32

const (
	StartedDigging DiggingAction = iota
	CancelledDigging
	FinishedDigging
	DropItemStack
	DropItem
	ShootArrowFinishEating
	SwapItemInHand
)

func DiggingActionFromPb(pbAction pb.PlayerDigging_Action) DiggingAction {
	return DiggingAction(pbAction)
}

func (d *DiggingAction) Pull(reader *buffer.Buffer) error {
	val := DiggingAction(reader.PullVarInt())
	switch val {
	case StartedDigging, CancelledDigging, FinishedDigging, DropItemStack, DropItem, ShootArrowFinishEating, SwapItemInHand:
		*d = val
		return nil
	}

	return fmt.Errorf("digging action index %d not allowed", int(val))
}

func (d *DiggingAction) Push(writer *buffer.Buffer) {
	writer.PushVarInt(int32(*d))
}
