package objects

import "time"

type BlockID uint32

func (b BlockID) String() string { return blockNamesMap[b] }
func (b BlockID) ID() uint32     { return uint32(b) }

// TODO Need to setup automated code generation from Notchian data export and provide here detailed data about every block.
//  This will need to consider the tool used for digging.
func (b BlockID) IsDiggable(tool ItemID) bool {
	switch tool {
	case ItemWoodenPickaxe, ItemStonePickaxe, ItemIronPickaxe, ItemDiamondPickaxe, ItemGoldenPickaxe:
		fallthrough
	case ItemWoodenShovel, ItemStoneShovel, ItemIronShovel, ItemDiamondShovel, ItemGoldenShovel:
		fallthrough
	case ItemWoodenAxe, ItemStoneAxe, ItemIronAxe, ItemDiamondAxe, ItemGoldenAxe:
		fallthrough
	case ItemWoodenHoe, ItemStoneHoe, ItemIronHoe, ItemDiamondHoe, ItemGoldenHoe:
		return true
	}

	return false
}

// TODO This will need to be exported as well, and will need to account for the actual tool used.
func (b BlockID) DigTime(tool ItemID) time.Duration {
	switch tool {
	case ItemWoodenPickaxe, ItemStonePickaxe, ItemIronPickaxe, ItemDiamondPickaxe, ItemGoldenPickaxe:
		fallthrough
	case ItemWoodenShovel, ItemStoneShovel, ItemIronShovel, ItemDiamondShovel, ItemGoldenShovel:
		fallthrough
	case ItemWoodenAxe, ItemStoneAxe, ItemIronAxe, ItemDiamondAxe, ItemGoldenAxe:
		fallthrough
	case ItemWoodenHoe, ItemStoneHoe, ItemIronHoe, ItemDiamondHoe, ItemGoldenHoe:
		return time.Second * 1
	}

	return time.Second * 86400
}
