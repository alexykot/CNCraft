package items

type Slot struct {
	IsPresent bool
	ItemID    int32
	ItemCount int8
	NBT       map[string]string // DEBT not clear atm how to handle this NBT compound
}
