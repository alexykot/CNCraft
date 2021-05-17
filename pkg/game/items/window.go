//go:generate stringer -type=clickMode window.go
package items

import (
	"fmt"
	"math"
	"sync"
)

type windowID uint8

const (
	InventoryWindow windowID = 0
)

type clickable interface {
	GetSlot(slotID int16) Slot
	SetSlot(slotID int16, item Slot)
}

type clickMgr struct {
	WindowID  windowID
	Clickable clickable

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

func (m *clickMgr) HandleClick(actionID, slotID, mode int16, button uint8, clickedItem Slot) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// println()
	// println(fmt.Sprintf("actionID: %d, lastAction: %d", actionID, m.lastAction))

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
	default:
		return false, fmt.Errorf("mode %s not supported", clickMode(mode).String())
	}

	if err == nil {
		m.lastAction = actionID
	}

	return inventoryUpdated, err
}

func (m *clickMgr) OpenWindow()  { m.isOpen = true }
func (m *clickMgr) CloseWindow() { m.isOpen = false }

// DEBT this all does not yet consider item stackability
func (m *clickMgr) handleMode0(slotID int16, button uint8, clickedItem Slot) (bool, error) {
	slotItem := m.Clickable.GetSlot(slotID)

	if !slotEqual(slotItem, clickedItem) {
		return false, fmt.Errorf("slot contents not equal to clickedItem supplied")
	}

	println(fmt.Sprintf("button: %d; slotID: %d; slotItem: %v; clickedItem: %v", button, slotID, slotItem, clickedItem))

	switch button {
	case leftMouseButton:
		println("left click")
		if m.cursor.IsPresent { // put down something
			println("put down")

			if slotItem.IsPresent {
				println("slot not empty")
				if m.cursor.ItemID == slotItem.ItemID { // same things in cursor and slot - join
					println("join")
					slotNewCount := m.cursor.ItemCount + slotItem.ItemCount
					cursorNewCount := slotNewCount - getMaxStack(slotItem.ItemID)
					if cursorNewCount > 0 { // total doesn't fit in one stack, something left on cursor
						slotItem.ItemCount = getMaxStack(slotItem.ItemID)
						m.cursor.ItemCount = cursorNewCount
					} else { // total fits in one stack, cursor empty
						slotItem.ItemCount = slotNewCount
						m.cursor = Slot{}
					}
					m.Clickable.SetSlot(slotID, slotItem)
				} else { // different things in cursor and slot - swap
					println("swap")
					m.Clickable.SetSlot(slotID, m.cursor)
					m.cursor = slotItem
				}
			} else {
				println("slot empty")
				m.Clickable.SetSlot(slotID, m.cursor)
				m.cursor = Slot{}
			}
			return true, nil
		} else { // pick up something
			println("pick up")
			if slotItem.IsPresent {
				println("picked item")
				m.cursor = slotItem // TODO the contents of the cursor should be persisted as well
				m.Clickable.SetSlot(slotID, Slot{})
				return true, nil
			} else {
				println("nothing to pick up")
				return false, nil
			}
		}
	case rightMouseButton:
		println("right click")
		if m.cursor.IsPresent { // put down something
			println("put down")

			if slotItem.IsPresent {
				println("slot not empty")
				if m.cursor.ItemID == slotItem.ItemID { // same things in cursor and slot - join one item from cursor stack
					println("join one item")
					slotNewCount := int16(math.Min(float64(slotItem.ItemCount + 1), float64(getMaxStack(slotItem.ItemID))))
					moved := slotNewCount - slotItem.ItemCount
					cursorNewCount := m.cursor.ItemCount - moved
					slotItem.ItemCount = slotNewCount
					m.cursor.ItemCount = cursorNewCount
					if m.cursor.ItemCount == 0 { // nothing left on the cursor
						m.cursor = Slot{}
					}
					m.Clickable.SetSlot(slotID, slotItem)
				} else { // different things in cursor and slot
					if slotItem.ItemCount == 1 { // only one item on cursor - swap
						println("swap")
						m.Clickable.SetSlot(slotID, m.cursor)
						m.cursor = slotItem
					} else { // more than one item on cursor - can't do anything
						println("can't swap, ignoring")
						return false, nil
					}
				}
			} else { // slot empty, put one item down
				println("slot empty - put one down")
				var newSlotItem Slot
				newSlotItem.IsPresent = true
				newSlotItem.ItemID = m.cursor.ItemID
				newSlotItem.ItemCount = 1
				m.cursor.ItemCount -= 1
				if m.cursor.ItemCount == 0 { // nothing left on the cursor
					m.cursor = Slot{}
				}
				m.Clickable.SetSlot(slotID, newSlotItem)
			}
			return true, nil
		} else { // pick up something
			println("pick up")
			if slotItem.IsPresent {
				println("picked half-stack")

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
				m.Clickable.SetSlot(slotID, slotItem)
				m.cursor = pickupItem // TODO the contents of the cursor should be persisted as well
				return true, nil
			} else {
				println("nothing to pick up")
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("button %d not supported for mode 0", button)

	}
}

