//go:generate stringer -type=clickMode window.go
package items

import (
	"fmt"
	"math"
	"sync"

	"go.uber.org/zap"
)

type WindowID uint8

const (
	InventoryWindow WindowID = 0
)

type clickable interface {
	GetSlot(slotID int16) Slot
	SetSlot(slotID int16, item Slot)
	GetRange(rangeType) slotRange
}

type windowMgr struct {
	WindowID  WindowID
	clickable clickable

	log        *zap.Logger
	mu         sync.Mutex
	lastAction int16
	isOpen     bool
	isUpset    bool
	// DEBT is server crashes while cursor is not empty - item on the cursor will be lost.
	//  Need to persist cursor contents and find a way to recover it after a crash as
	//  cursor contents cannot be communicated to the client.
	cursor Slot
}

type clickMode uint8

const (
	simpleClick clickMode = 0
	shftClick   clickMode = 1
	numberKey   clickMode = 2
	middleClick clickMode = 3
	drop        clickMode = 4
	dragPaint   clickMode = 5
	doubleClick clickMode = 6
)

const slotOutsideWindow = -999

const (
	leftMouseButton   = 0
	rightMouseButton  = 1
	middleMouseButton = 2

	kbdKey1 = 0
	kbdKey2 = 1
	kbdKey3 = 2
	kbdKey4 = 3
	kbdKey5 = 4
	kbdKey6 = 5
	kbdKey7 = 6
	kbdKey8 = 7
	kbdKey9 = 8

	kbdKeyQ = 0
)

func (m *windowMgr) HandleClick(actionID, slotID, mode int16, button uint8, clickedItem Slot) (*Slot, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Debug(fmt.Sprintf("button: %d; slotID: %d; clickedItem: %v", button, slotID, clickedItem), zap.Int16("mode", mode))

	m.log.Debug(fmt.Sprintf("actionID: %d, lastAction: %d", actionID, m.lastAction))

	if m.lastAction+1 != actionID {
		m.isUpset = true
		return nil, false, fmt.Errorf("action ID out of sequence, next action should be %d", m.lastAction+1)
	} else if m.isUpset {
		return nil, false, fmt.Errorf("expecting client to apologise before sending any other clicks")
	}

	// Notchian client never sends an OpenWindow packet for player inventory, we just start receiving inventory clicks.
	// It does send a CloseWindow though.
	if m.WindowID == InventoryWindow {
		m.isOpen = true
	}

	if !m.isOpen {
		return nil, false, fmt.Errorf("window ID %d is not open", m.WindowID)
	}

	var err error
	var inventoryUpdated bool
	var droppedItem Slot
	switch clickMode(mode) {
	case simpleClick:
		droppedItem, inventoryUpdated, err = m.handleMode0(slotID, button, clickedItem)
	case shftClick:
		inventoryUpdated, err = m.handleMode1(slotID, button, clickedItem)
	case numberKey:
		inventoryUpdated, err = m.handleMode2(slotID, button, clickedItem)
	case drop:
		droppedItem, inventoryUpdated, err = m.handleMode4(slotID, button, clickedItem)
	case middleClick, dragPaint, doubleClick:
		return nil, false, fmt.Errorf("mode %s not supported", clickMode(mode).String())
	default:
		return nil, false, fmt.Errorf("invalid mode %d received", mode)
	}

	if err == nil {
		m.lastAction = actionID
	}

	if droppedItem.IsPresent {
		return &droppedItem, inventoryUpdated, err
	} else {
		return nil, inventoryUpdated, err
	}
}

func (m *windowMgr) IsUpset() bool { return m.isUpset }
func (m *windowMgr) Apologise(clientActionID int16) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Debug("client apologising", zap.Bool("isUpset", m.isUpset),
		zap.Int16("lastAction", m.lastAction), zap.Any("last client action", clientActionID))

	m.isUpset = false
	if m.lastAction < clientActionID {
		m.lastAction = clientActionID
	}
}

func (m *windowMgr) LastAction() int16 { return m.lastAction }
func (m *windowMgr) OpenWindow() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isOpen = true
}

func (m *windowMgr) CloseWindow() Slot {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isOpen = false
	droppedItem := m.cursor
	m.cursor = Slot{}
	return droppedItem
}

func (m *windowMgr) handleMode0(slotID int16, button uint8, clickedItem Slot) (Slot, bool, error) {
	slotItem := m.clickable.GetSlot(slotID)
	if !slotEqual(slotItem, clickedItem) {
		return Slot{}, false, fmt.Errorf("slot contents not equal to clickedItem supplied")
	}

	switch button {
	case leftMouseButton:
		m.log.Debug("left click", zap.Int("mode", 0))
		if m.cursor.IsPresent { // put down something
			m.log.Debug("put down", zap.Int("mode", 0))

			if slotItem.IsPresent {
				m.log.Debug("slot not empty", zap.Int("mode", 0))
				if m.cursor.ItemID == slotItem.ItemID { // same things in cursor and slot - join
					m.log.Debug("join", zap.Int("mode", 0))
					slotNewCount := m.cursor.ItemCount + slotItem.ItemCount
					cursorNewCount := slotNewCount - slotItem.ItemID.MaxStack()
					if cursorNewCount > 0 { // total doesn't fit in one stack, something left on cursor
						slotItem.ItemCount = slotItem.ItemID.MaxStack()
						m.cursor.ItemCount = cursorNewCount
					} else { // total fits in one stack, cursor empty
						slotItem.ItemCount = slotNewCount
						m.cursor = Slot{}
					}
					m.clickable.SetSlot(slotID, slotItem)
				} else { // different things in cursor and slot - swap
					m.log.Debug("swap", zap.Int("mode", 0))
					m.clickable.SetSlot(slotID, m.cursor)
					m.cursor = slotItem
				}
			} else {
				if slotID == slotOutsideWindow { // click outside window, drop all from cursor
					m.log.Debug("dropping all from cursor", zap.Int("mode", 0))
					dropped := m.cursor
					m.cursor = Slot{}
					return dropped, true, nil
				} else { // slot empty, put all down
					m.log.Debug("slot empty", zap.Int("mode", 0))
					m.clickable.SetSlot(slotID, m.cursor)
					m.cursor = Slot{}
				}
			}
			return Slot{}, true, nil
		} else { // pick up something
			m.log.Debug("pick up", zap.Int("mode", 0))
			if slotItem.IsPresent {
				m.log.Debug("picked item", zap.Int("mode", 0))
				m.cursor = slotItem
				m.clickable.SetSlot(slotID, Slot{})
				return Slot{}, true, nil
			} else {
				m.log.Debug("nothing to pick up", zap.Int("mode", 0))
				return Slot{}, false, nil
			}
		}
	case rightMouseButton:
		m.log.Debug("right click", zap.Int("mode", 0))
		if m.cursor.IsPresent { // put down something
			m.log.Debug("put down", zap.Int("mode", 0))

			if slotItem.IsPresent {
				m.log.Debug("slot not empty", zap.Int("mode", 0))
				if m.cursor.ItemID == slotItem.ItemID { // same things in cursor and slot - join one item from cursor stack
					m.log.Debug("join one item", zap.Int("mode", 0))
					slotNewCount := int16(math.Min(float64(slotItem.ItemCount+1), float64(slotItem.ItemID.MaxStack())))
					moved := slotNewCount - slotItem.ItemCount
					cursorNewCount := m.cursor.ItemCount - moved
					slotItem.ItemCount = slotNewCount
					m.cursor.ItemCount = cursorNewCount
					if m.cursor.ItemCount < 1 { // nothing left on the cursor
						m.cursor = Slot{}
					}
					m.clickable.SetSlot(slotID, slotItem)
				} else { // different things in cursor and slot
					if slotItem.ItemCount == 1 { // only one item on cursor - swap
						m.log.Debug("swap", zap.Int("mode", 0))
						m.clickable.SetSlot(slotID, m.cursor)
						m.cursor = slotItem
					} else { // more than one item on cursor - can't do anything
						m.log.Debug("can't swap, ignoring", zap.Int("mode", 0))
						return Slot{}, false, nil
					}
				}
			} else {
				if slotID == slotOutsideWindow { // click outside window, drop one item from cursor
					// TODO this should return dropped items
					m.log.Debug("dropping one item from cursor", zap.Int("mode", 0))
					dropped := m.cursor
					dropped.ItemCount = 1
					m.cursor.ItemCount = m.cursor.ItemCount - dropped.ItemCount
					if m.cursor.ItemCount < 1 {
						m.cursor = Slot{}
					}
					return dropped, true, nil
				} else { // slot empty, put one item down
					m.log.Debug("slot empty - put one down", zap.Int("mode", 0))
					var newSlotItem Slot
					newSlotItem.IsPresent = true
					newSlotItem.ItemID = m.cursor.ItemID
					newSlotItem.ItemCount = 1
					m.cursor.ItemCount -= 1
					if m.cursor.ItemCount == 0 { // nothing left on the cursor
						m.cursor = Slot{}
					}
					m.clickable.SetSlot(slotID, newSlotItem)
				}
			}
			return Slot{}, true, nil
		} else { // pick up something
			m.log.Debug("pick up", zap.Int("mode", 0))
			if slotItem.IsPresent {
				m.log.Debug("picked half-stack", zap.Int("mode", 0))

				var pickupItem Slot
				pickupItem.IsPresent = true
				pickupItem.ItemID = slotItem.ItemID
				pickupItem.ItemCount = int16(math.Ceil(float64(slotItem.ItemCount) / 2))

				slotItem.ItemCount = slotItem.ItemCount - pickupItem.ItemCount

				// there was only one item in the stack, or unstackable item
				if slotItem.ItemCount == 0 {
					slotItem.IsPresent = false
					slotItem.ItemID = 0
				}
				m.clickable.SetSlot(slotID, slotItem)
				m.cursor = pickupItem
				return Slot{}, true, nil
			} else {
				m.log.Debug("nothing to pick up", zap.Int("mode", 0))
				return Slot{}, false, nil
			}
		}
	default:
		return Slot{}, false, fmt.Errorf("button %d is invalid for mode 0", button)
	}
}

func (m *windowMgr) handleMode1(slotID int16, button uint8, clickedItem Slot) (bool, error) {
	slotItem := m.clickable.GetSlot(slotID)

	if !slotEqual(slotItem, clickedItem) {
		return false, fmt.Errorf("slot contents not equal to clickedItem supplied")
	}

	m.log.Debug(fmt.Sprintf("button: %d; slotID: %d; slotItem: %v; clickedItem: %v", button, slotID, slotItem, clickedItem),
		zap.Int("mode", 1))

	switch button {
	case leftMouseButton, rightMouseButton:
		if !slotItem.IsPresent {
			return false, nil
		}

		topRange := m.clickable.GetRange(top)
		bottomRange := m.clickable.GetRange(bottom)

		var targetRange slotRange
		if topRange.InRange(slotID) {
			targetRange = bottomRange // target range is opposite of where the click happened
		} else if bottomRange.InRange(slotID) {
			targetRange = topRange // target range is opposite of where the click happened
		} else {
			// TODO handle shift+click when slotID is not in one of the standard ranges
			return false, fmt.Errorf("slotID out of range")
		}

		sameItemSlots := targetRange.GetItemSlots(m.clickable, slotItem.ItemID)
		var hasChanged bool
		if len(sameItemSlots) > 0 {
			for _, sameItemSlotID := range sameItemSlots {
				sameItem := m.clickable.GetSlot(sameItemSlotID)
				newItemCount := int16(math.Min(float64(sameItem.ItemCount+slotItem.ItemCount), float64(sameItem.ItemID.MaxStack())))

				hasChanged = sameItem.ItemCount != newItemCount
				slotItem.ItemCount = slotItem.ItemCount - (newItemCount - sameItem.ItemCount)
				sameItem.ItemCount = newItemCount

				if slotItem.ItemCount <= 0 {
					slotItem = Slot{}
				}
				m.clickable.SetSlot(sameItemSlotID, sameItem)
				m.clickable.SetSlot(slotID, slotItem)
				if !slotItem.IsPresent {
					return true, nil // everything distributed, nothing else to do
				}
			}
		}

		// either no same item slots were found, or there was not enough space in those stacks to distribute everything (or item is unstackable)
		emptySlots := targetRange.GetEmptySlots(m.clickable)
		if len(emptySlots) > 0 {
			println(fmt.Sprintf("empty slots available, put remaining %d items into the first empty slot %d", slotItem.ItemCount, emptySlots[0]))
			m.clickable.SetSlot(emptySlots[0], slotItem) // place remainder of the item stack into first empty slot
			m.clickable.SetSlot(slotID, Slot{})
			return true, nil
		}

		// we get here if there was not enough spaces in the sameItemSlots to distribute everything,
		// and no empty slots are available
		println("not enough space to distribute everything, something left behind")
		return hasChanged, nil
	default:
		return false, fmt.Errorf("button %d not supported for mode 0", button)
	}
}

func (m *windowMgr) handleMode2(slotID int16, button uint8, clickedItem Slot) (bool, error) {
	switch button { // Pick up that can.
	case kbdKey1, kbdKey2, kbdKey3, kbdKey4, kbdKey5, kbdKey6, kbdKey7, kbdKey8, kbdKey9:
		// Okay, you can go.
	default:
		return false, fmt.Errorf("button %d not supported for mode 2", button)
	}

	hotbarRange := m.clickable.GetRange(hotbar)
	hotbarSlots := hotbarRange.GetSlots()
	if len(hotbarSlots) != kbdKey9+1 {
		return false, fmt.Errorf("unexpected number of hotbar slots")
	}

	if slotID == hotbarSlots[button] {
		return false, nil // swapping a hotbar slot with itself
	}

	slotItem := m.clickable.GetSlot(slotID)
	if !slotEqual(slotItem, clickedItem) {
		return false, fmt.Errorf("slot contents not equal to clickedItem supplied")
	}

	hotbarItem := m.clickable.GetSlot(hotbarSlots[button])
	m.clickable.SetSlot(hotbarSlots[button], slotItem)
	m.clickable.SetSlot(slotID, hotbarItem)
	return true, nil
}

func (m *windowMgr) handleMode4(slotID int16, button uint8, _ Slot) (Slot, bool, error) {
	if button != kbdKeyQ {
		return Slot{}, false, fmt.Errorf("button %d not supported for mode 4", button)
	}

	slotItem := m.clickable.GetSlot(slotID)
	if !slotItem.IsPresent {
		m.log.Debug(fmt.Sprintf("slot %d, item not present, nothing to drop", slotID))
		return Slot{}, false, nil // nothing to drop
	}

	droppedItem := slotItem
	droppedItem.ItemCount = 1
	slotItem.ItemCount = slotItem.ItemCount - droppedItem.ItemCount
	m.log.Debug(fmt.Sprintf("slot %d, %d item dropped, %d items left", slotID, droppedItem.ItemCount, slotItem.ItemCount))

	if slotItem.ItemCount == 0 {
		m.log.Debug(fmt.Sprintf("slot %d, no items left", slotID))
		slotItem = Slot{}
	}

	m.clickable.SetSlot(slotID, slotItem)

	return droppedItem, true, nil
}
