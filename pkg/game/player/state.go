package player

import "github.com/alexykot/cncraft/pkg/game/data"

type State struct {
	CurrentHotbarSlot HotBarSlot
	CurrentLocation   data.Location
}

type HotBarSlot byte

const (
	Slot0 HotBarSlot = iota
	Slot1
	Slot2
	Slot3
	Slot4
	Slot5
	Slot6
	Slot7
	Slot8
)
