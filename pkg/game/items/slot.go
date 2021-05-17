package items

type Slot struct {
	IsPresent bool
	ItemID    int16 // TODO pItems.ItemID ?
	ItemCount int16
	NBT       map[string]string // DEBT not clear atm how to handle this NBT compound
}

func slotEqual(itemLeft, itemRight Slot) bool {
	return itemLeft.IsPresent == itemRight.IsPresent &&
		itemLeft.ItemID == itemRight.ItemID &&
		itemLeft.ItemCount == itemRight.ItemCount
}

func getMaxStack(itemID int16) int16 {
	return 64 // DEBT replace this with per-item stackability count
}
