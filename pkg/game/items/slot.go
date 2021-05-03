package items

type Slot struct {
	IsPresent bool
	ItemID    int
	ItemCount int
	NBT       map[string]string // DEBT not clear atm how to handle this NBT compound
}
