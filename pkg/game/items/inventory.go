package items

import (
	"go.uber.org/zap"

	pItems "github.com/alexykot/cncraft/pkg/protocol/items"
)

type Inventory struct {
	windowMgr

	CurrentHotbarSlot uint8

	RowTop    [9]Slot
	RowMiddle [9]Slot
	RowBottom [9]Slot
	RowHotbar [9]Slot

	Armor   [4]Slot
	Offhand Slot

	Craft  [4]Slot
	Result Slot
}

func NewInventory(windowLog *zap.Logger) *Inventory {
	inv := &Inventory{
		windowMgr: windowMgr{
			WindowID: InventoryWindow,
			log:      windowLog,
		},
	}
	inv.clickable = inv // TODO don't like this Ouroboros wiring
	return inv
}

// ToArray converts Inventory into correctly numbered array of slots for marshalling into a packet.
func (i *Inventory) ToArray() []Slot {
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

func (i *Inventory) GetSlot(slotID int16) Slot {
	if slotID < 0 {
		return Slot{}
	}

	items := i.ToArray()
	if slotID > int16(len(items)-1) {
		return Slot{}
	}

	return items[slotID]
}

func (i *Inventory) SetSlot(slotID int16, item Slot) {
	item.IsPresent = item.ItemID != pItems.Air

	if slotID == 0 {
		i.Result = item
	} else if slotID > 0 && slotID < 5 {
		i.Craft[slotID-1] = item
	} else if slotID > 4 && slotID < 9 {
		i.Armor[slotID-5] = item
	} else if slotID > 8 && slotID < 18 {
		i.RowTop[slotID-9] = item
	} else if slotID > 17 && slotID < 27 {
		i.RowMiddle[slotID-18] = item
	} else if slotID > 26 && slotID < 36 {
		i.RowBottom[slotID-27] = item
	} else if slotID > 35 && slotID < 45 {
		i.RowHotbar[slotID-36] = item
	} else if slotID == 45 {
		i.Offhand = item
	}
}

func (i *Inventory) reset() {
	wipe := make([]Slot, 46, 46)
	for ii := range wipe {
		i.SetSlot(int16(ii), Slot{})
	}
}

func (i *Inventory) GetRange(rangeType rangeType) slotRange {
	var slots slotRange
	switch rangeType {
	case top:
		slots = slotRange{9, 35}
	case bottom:
		slots = slotRange{36, 44}
	case hotbar:
		slots = slotRange{36, 44}
	}
	return slots
}
