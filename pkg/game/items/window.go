package items

import (
	"fmt"
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

	lastAction int16
	isOpen     bool
	cursor     *Slot
}

const (
	mode0 = 0
	mode1 = 1
	mode2 = 2
	mode3 = 3
	mode4 = 4
	mode5 = 5
	mode6 = 6
)

const (
	button0  = 0
	button1  = 1
	button2  = 2
	button3  = 3
	button4  = 4
	button5  = 5
	button6  = 6
	button7  = 7
	button8  = 8
	button9  = 9
	button10 = 10
)

func (m *clickMgr) HandleClick(actionID, slotID, mode int16, button uint8, clickedItem Slot) (bool, error) {
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
	switch mode {
	case mode0:
		inventoryUpdated, err = m.handleMode0(slotID, button, clickedItem)
	default:
		return false, fmt.Errorf("mode %d not supported", mode)
	}

	if err == nil {
		m.lastAction = actionID
	}

	return inventoryUpdated, err
}

func (m *clickMgr) OpenWindow()  { m.isOpen = true }
func (m *clickMgr) CloseWindow() { m.isOpen = false }

func (m *clickMgr) handleMode0(slotID int16, button uint8, clickedItem Slot) (bool, error) {
	slotItem := m.Clickable.GetSlot(slotID)

	println(fmt.Sprintf("slotID: %d", slotID))
	println(fmt.Sprintf("slotItem: %v", slotItem))
	println(fmt.Sprintf("clickedItem: %v", clickedItem))

	switch button {
	case button0:
		if slotItem.IsPresent {
			if !clickedItem.IsPresent {
				return false, fmt.Errorf("non-empty slot was clicked, therefore clicked item must be present")
			}

			if slotItem.ItemID != clickedItem.ItemID {
				return false, fmt.Errorf("clicked item is different from the item in the slot")
			}

			if slotItem.ItemCount != clickedItem.ItemCount {
				return false, fmt.Errorf("clicked item count is different from the item count in the slot")
			}

			println(fmt.Sprintf("picked from slotID: %d, item: %v", slotID, m.cursor))
			m.cursor = &slotItem
			m.Clickable.SetSlot(slotID, Slot{})
			return true, nil
		} else if m.cursor != nil {
			if clickedItem.IsPresent {
				return false, fmt.Errorf("empty slot was clicked with a cursor item, clickedItem item must be not present")
			}

			println(fmt.Sprintf("put to slotID: %d, item: %v", slotID, m.cursor))
			m.Clickable.SetSlot(slotID, *m.cursor)
			m.cursor = nil
			return true, nil
		}
		return false, nil
	case button1:
		if slotItem.IsPresent {
			if !clickedItem.IsPresent {
				return false, fmt.Errorf("non-empty slot was clicked, therefore clicked item must be present")
			}

			if slotItem.ItemID != clickedItem.ItemID {
				return false, fmt.Errorf("clicked item is different from the item in the slot")
			}

			if slotItem.ItemCount != clickedItem.ItemCount {
				return false, fmt.Errorf("clicked item count is different from the item count in the slot")
			}

			m.cursor = &slotItem
			m.Clickable.SetSlot(slotID, Slot{})
			return true, nil
		} else if m.cursor != nil {
			if clickedItem.IsPresent {
				return false, fmt.Errorf("empty slot was clicked with a cursor item, clickedItem item must be not present")
			}
			m.Clickable.SetSlot(slotID, *m.cursor)
			println("item put down")

			m.cursor = nil
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("button %d not supported for mode 0", button)

	}
}
