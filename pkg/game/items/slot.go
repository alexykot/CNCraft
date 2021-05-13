package items

type Slot struct {
	IsPresent bool
	ItemID    int16
	ItemCount int16
	NBT       map[string]string // DEBT not clear atm how to handle this NBT compound
}
