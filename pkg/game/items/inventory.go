package items

import (
	pItems "github.com/alexykot/cncraft/pkg/protocol/items"
)

type Inventory struct {
	CurrentHotbarSlot HotBarSlot

	RowTop    [9]Slot
	RowMiddle [9]Slot
	RowBottom [9]Slot
	RowHotbar [9]Slot

	Armor   [4]Slot
	Offhand Slot

	Craft  [4]Slot
	Result Slot
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

// ToArray converts Inventory into correctly numbered array of slots for marshalling into a packet.
// TODO Possibly later move this into centralised Window implementation for all types of windows.
func (i Inventory) ToArray() []Slot {
	result := make([]Slot, 46, 46)
	result[0] = i.Result
	result[1] = i.Craft[0]
	result[2] = i.Craft[1]
	result[3] = i.Craft[2]
	result[4] = i.Craft[3]

	result[5] = i.Armor[0]
	result[6] = i.Armor[1]
	result[7] = i.Armor[2]
	result[8] = i.Armor[3]

	ii := 9
	for _, item := range i.RowTop {
		result[ii] = item
		ii++
	}
	for _, item := range i.RowMiddle {
		result[ii] = item
		ii++
	}
	for _, item := range i.RowBottom {
		result[ii] = item
		ii++
	}
	for _, item := range i.RowHotbar {
		result[ii] = item
		ii++
	}
	result[45] = i.Offhand

	return result
}

func (i *Inventory) AssignSlot(slotNumber, itemID, itemCount int) {
	var item Slot
	if itemID != int(pItems.Air) {
		item.IsPresent = true
		item.ItemID = itemID
		item.ItemCount = itemCount
	}

	if slotNumber == 0 {
		i.Result = item
	} else if slotNumber > 0 && slotNumber < 5 {
		i.Craft[slotNumber-1] = item
	} else if slotNumber > 4 && slotNumber < 9 {
		i.Armor[slotNumber-5] = item
	} else if slotNumber > 8 && slotNumber < 18 {
		i.RowTop[slotNumber-9] = item
	} else if slotNumber > 17 && slotNumber < 27 {
		i.RowMiddle[slotNumber-18] = item
	} else if slotNumber > 26 && slotNumber < 36 {
		i.RowBottom[slotNumber-27] = item
	} else if slotNumber > 35 && slotNumber < 45 {
		i.RowHotbar[slotNumber-36] = item
	} else if slotNumber == 45 {
		i.Offhand = item
	}
}
