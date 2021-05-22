package items

var stackable16 = map[ItemID]struct{}{
	Egg:      {},
	Snowball: {},
}

var stackable64 = map[ItemID]struct{}{
	Dirt:    {},
	Bedrock: {},
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
