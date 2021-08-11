package objects

// TODO Need to setup automated code generation from Notichan data export and provide here detailed data about every block.
var stackable16 = map[ItemID]struct{}{
	ItemEgg:      {},
	ItemSnowball: {},
}

var stackable64 = map[ItemID]struct{}{
	ItemDirt:    {},
	ItemBedrock: {},
}

func (b ItemID) MaxStack() int16 {
	if _, ok := stackable64[b]; ok {
		return 64
	}

	if _, ok := stackable16[b]; ok {
		return 16
	}

	return 1
}
