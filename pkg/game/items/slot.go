package items

import pItems "github.com/alexykot/cncraft/pkg/protocol/items"

type Slot struct {
	IsPresent bool
	ItemID    pItems.ItemID
	ItemCount int16
	NBT       map[string]string // DEBT not clear atm how to handle this NBT compound
}

func slotEqual(itemLeft, itemRight Slot) bool {
	return itemLeft.IsPresent == itemRight.IsPresent &&
		itemLeft.ItemID == itemRight.ItemID &&
		itemLeft.ItemCount == itemRight.ItemCount
}

type rangeType string

const top rangeType = "top"
const bottom rangeType = "bottom"
const hotbar rangeType = "hotbar"

type slotRange struct {
	start, end int16
}

func (r slotRange) InRange(slotID int16) bool {
	return slotID >= r.start && slotID <= r.end
}

func (r slotRange) GetSlots() []int16 {
	var slots []int16
	for slotID := r.start; slotID <= r.end; slotID++ {
		slots = append(slots, slotID)
	}
	return slots
}

func (r slotRange) GetEmptySlots(window clickable) []int16 {
	var emptySlots []int16
	for slotID := r.start; slotID <= r.end; slotID++ {
		if !window.GetSlot(slotID).IsPresent {
			emptySlots = append(emptySlots, slotID)
		}
	}
	return emptySlots
}

func (r slotRange) GetItemSlots(window clickable, itemID pItems.ItemID) []int16 {
	var itemSlots []int16
	for slotID := r.start; slotID <= r.end; slotID++ {
		slotItem := window.GetSlot(slotID)
		if slotItem.IsPresent && slotItem.ItemID == itemID {
			itemSlots = append(itemSlots, slotID)
		}
	}
	return itemSlots
}
