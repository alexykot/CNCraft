//go:generate stringer -type=clickMode window.go
package items

import (
	"fmt"
	"math"
	"sync"

	"go.uber.org/zap"
)

type windowID uint8

const (
	InventoryWindow windowID = 0
)

type clickable interface {
	GetSlot(slotID int16) Slot
	SetSlot(slotID int16, item Slot)
	GetRange(whatRange) slotRange
}

type windowMgr struct {
	WindowID  windowID
	clickable clickable

	log        *zap.Logger
	mu         sync.Mutex
	lastAction int16
	isOpen     bool
	cursor     Slot
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

const (
	leftMouseButton   = 0
	rightMouseButton  = 1
	middleMouseButton = 2
)

func (m *windowMgr) HandleClick(actionID, slotID, mode int16, button uint8, clickedItem Slot) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Debug(fmt.Sprintf("actionID: %d, lastAction: %d", actionID, m.lastAction))

	if m.lastAction+1 != actionID {
		return false, fmt.Errorf("action ID out of sequence, next action should be %d", m.lastAction+1)
	}

	// Notchian client never sends an OpenWindow packet for player inventory, we just start receiving inventory clicks.
	// It does send a CloseWindow though.
	if m.WindowID == InventoryWindow {
		m.isOpen = true
	}

	if !m.isOpen {
		return false, fmt.Errorf("window ID %d is not open", m.WindowID)
	}

	var err error
	var inventoryUpdated bool
	switch clickMode(mode) {
	case simpleClick:
		inventoryUpdated, err = m.handleMode0(slotID, button, clickedItem)
	case shftClick:
		inventoryUpdated, err = m.handleMode1(slotID, button, clickedItem)
	case numberKey, middleClick, drop, dragPaint, doubleClick:
		return false, fmt.Errorf("mode %s not supported", clickMode(mode).String())
	default:
		return false, fmt.Errorf("invalid mode %d received", mode)
	}

	if err == nil {
		m.lastAction = actionID
	}

	return inventoryUpdated, err
}

func (m *windowMgr) OpenWindow()  { m.isOpen = true }
func (m *windowMgr) CloseWindow() { m.isOpen = false } // TODO do something with cursor contents here

func (m *windowMgr) handleMode0(slotID int16, button uint8, clickedItem Slot) (bool, error) {
	slotItem := m.clickable.GetSlot(slotID)

	if !slotEqual(slotItem, clickedItem) {
		return false, fmt.Errorf("slot contents not equal to clickedItem supplied")
	}

	m.log.Debug(fmt.Sprintf("button: %d; slotID: %d; slotItem: %v; clickedItem: %v", button, slotID, slotItem, clickedItem),
		zap.Int("mode", 0))

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
				m.log.Debug("slot empty", zap.Int("mode", 0))
				m.clickable.SetSlot(slotID, m.cursor)
				m.cursor = Slot{}
			}
			return true, nil
		} else { // pick up something
			m.log.Debug("pick up", zap.Int("mode", 0))
			if slotItem.IsPresent {
				m.log.Debug("picked item", zap.Int("mode", 0))
				m.cursor = slotItem // TODO the contents of the cursor should be persisted as well
				m.clickable.SetSlot(slotID, Slot{})
				return true, nil
			} else {
				m.log.Debug("nothing to pick up", zap.Int("mode", 0))
				return false, nil
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
					if m.cursor.ItemCount == 0 { // nothing left on the cursor
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
						return false, nil
					}
				}
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
			return true, nil
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
				m.cursor = pickupItem // TODO the contents of the cursor should be persisted as well
				return true, nil
			} else {
				m.log.Debug("nothing to pick up", zap.Int("mode", 0))
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("button %d is invalid for mode 0", button)
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
