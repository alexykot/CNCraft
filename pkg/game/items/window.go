//go:generate stringer -type=clickMode window.go
package items

import (
	"fmt"
	"math"
	"sync"

	"go.uber.org/zap"
)

type WindowID int8

const (
	CursorWindow    WindowID = -1 // Special value only used in CSetSlot to set cursor contents
	InventoryWindow WindowID = 0
)

const CursorSlot = -1

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
	// DEBT if server crashes while cursor is not empty - item on the cursor will be lost.
	//  Need to persist cursor contents and find a way to recover it after a crash.
	//  It is possible to set cursor contents using SetSlot packet.
	cursor Slot

	dragged    Slot // dragged.IsPresent=true indicates that dragging is in progress
	dragSlots  []int16
	dragPlaced map[int16]int16
}

type clickMode uint8
type button uint8

const (
	simpleClick clickMode = 0
	shftClick   clickMode = 1
	numberKey   clickMode = 2
	middleClick clickMode = 3
	drop        clickMode = 4
	drag        clickMode = 5
	doubleClick clickMode = 6
)

const slotOutsideWindow = -999

const (
	leftMouseButton   button = 0
	rightMouseButton  button = 1
	middleMouseButton button = 2

	kbdKey1 button = 0
	kbdKey2 button = 1
	kbdKey3 button = 2
	kbdKey4 button = 3
	kbdKey5 button = 4
	kbdKey6 button = 5
	kbdKey7 button = 6
	kbdKey8 button = 7
	kbdKey9 button = 8

	kbdKeyQ button = 0

	startLeftMouseDrag   button = 0
	startRightMouseDrag  button = 4
	startMiddleMouseDrag button = 8
	addLeftDragSlot      button = 1
	addRightDragSlot     button = 5
	addMiddleDragSlot    button = 9
	endLeftMouseDrag     button = 2
	endRightMouseDrag    button = 6
	endMiddleMouseDrag   button = 10
)

func (m *windowMgr) HandleClick(actionID, slotID, mode int16, keyPress uint8, clickedItem Slot) (*Slot, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Debug(fmt.Sprintf("button: %d; slotID: %d; clickedItem: %v", keyPress, slotID, clickedItem), zap.Int16("mode", mode))

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

	if m.dragged.IsPresent && clickMode(mode) != drag {
		m.dragged = Slot{}
		m.dragSlots = nil
		m.dragPlaced = nil
		return nil, false, fmt.Errorf("received non-dragging mode %s while drag is active", clickMode(mode).String())
	}

	var err error
	var inventoryUpdated bool
	var droppedItem Slot
	switch clickMode(mode) {
	case simpleClick:
		droppedItem, inventoryUpdated, err = m.handleMode0(slotID, button(keyPress), clickedItem)
	case shftClick:
		inventoryUpdated, err = m.handleMode1(slotID, button(keyPress), clickedItem)
	case numberKey:
		inventoryUpdated, err = m.handleMode2(slotID, button(keyPress), clickedItem)
	case drop:
		droppedItem, inventoryUpdated, err = m.handleMode4(slotID, button(keyPress), clickedItem)
	case drag:
		inventoryUpdated, err = m.handleMode5(slotID, button(keyPress), clickedItem)
	case middleClick, doubleClick:
		return nil, false, fmt.Errorf("mode %s not supported", clickMode(mode).String())
	default:
		return nil, false, fmt.Errorf("invalid mode %d received", mode)
	}

	if err == nil {
		m.lastAction = actionID
	} else {
		m.log.Debug("click cannot be handled", zap.Error(err))
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

func (m *windowMgr) GetCursor() Slot   { return m.cursor }
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

func (m *windowMgr) handleMode0(slotID int16, button button, clickedItem Slot) (Slot, bool, error) {
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

func (m *windowMgr) handleMode1(slotID int16, button button, clickedItem Slot) (bool, error) {
	slotItem := m.clickable.GetSlot(slotID)

	if clickedItem.IsPresent {
		return false, fmt.Errorf("clickedItem should not be present in mode 1")
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

func (m *windowMgr) handleMode2(slotID int16, button button, clickedItem Slot) (bool, error) {
	switch button { // Pick up that can.
	case kbdKey1, kbdKey2, kbdKey3, kbdKey4, kbdKey5, kbdKey6, kbdKey7, kbdKey8, kbdKey9:
		// Okay, you can go.
	default:
		return false, fmt.Errorf("button %d not supported for mode 2", button)
	}

	if clickedItem.IsPresent {
		return false, fmt.Errorf("clickedItem should not be present in mode 2")
	}

	hotbarRange := m.clickable.GetRange(hotbar)
	hotbarSlots := hotbarRange.GetSlots()
	if len(hotbarSlots) != int(kbdKey9)+1 {
		return false, fmt.Errorf("unexpected number of hotbar slots")
	}

	if slotID == hotbarSlots[button] {
		return false, nil // swapping a hotbar slot with itself
	}

	slotItem := m.clickable.GetSlot(slotID)
	hotbarItem := m.clickable.GetSlot(hotbarSlots[button])
	m.clickable.SetSlot(hotbarSlots[button], slotItem)
	m.clickable.SetSlot(slotID, hotbarItem)
	return true, nil
}

func (m *windowMgr) handleMode4(slotID int16, button button, _ Slot) (Slot, bool, error) {
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

func (m *windowMgr) handleMode5(slotID int16, button button, _ Slot) (bool, error) {
	switch button {
	case startLeftMouseDrag, startRightMouseDrag:
		if slotID != slotOutsideWindow {
			m.dragSlots = nil
			m.dragged = Slot{}
			return false, fmt.Errorf("slotID must be -999 to start the drag")
		}

		if !m.cursor.IsPresent {
			m.dragSlots = nil
			m.dragged = Slot{}
			return false, fmt.Errorf("nothing to drag - cursor is empty")
		}

		if m.dragged.IsPresent {
			m.dragSlots = nil
			m.dragged = Slot{}
			return false, fmt.Errorf("was already dragging")
		}

		m.dragged = m.cursor
		m.dragPlaced = make(map[int16]int16)
		return false, nil

	case addLeftDragSlot, addRightDragSlot:
		if !m.dragged.IsPresent || m.dragPlaced == nil {
			m.dragSlots = nil
			m.dragged = Slot{}
			return false, fmt.Errorf("is not dragging")
		}

		slotItem := m.clickable.GetSlot(slotID)
		if slotItem.IsPresent && slotItem.ItemID != m.dragged.ItemID {
			m.dragSlots = nil
			m.dragged = Slot{}
			return false, fmt.Errorf("dragging over non-matching item")
		}

		m.dragSlots = append(m.dragSlots, slotID)
		var perSlot int16
		switch button {
		case addLeftDragSlot:
			perSlot = int16(math.Floor(float64(m.dragged.ItemCount) / float64(len(m.dragSlots))))
		case addRightDragSlot:
			perSlot = 1
		}

		dragged := m.dragged
		for _, dragSlotID := range m.dragSlots {
			if dragged.ItemCount < 1 {
				continue // nothing to place anymore, all distributed
			}

			dragSlotItem := m.clickable.GetSlot(dragSlotID)
			if !dragSlotItem.IsPresent { // empty slot is now slot of this type, with zero count
				dragSlotItem = m.dragged
				dragSlotItem.ItemCount = 0
			}

			// if we have already painted - need to take away what was just added to not double-paint into same slot
			dragSlotItem.ItemCount = dragSlotItem.ItemCount - m.dragPlaced[dragSlotID]

			// update itemcount with as much as can be distributed into this slot, but no more than stackable
			oldCount := dragSlotItem.ItemCount
			dragSlotItem.ItemCount = int16(math.Min(float64(dragSlotItem.ItemCount+perSlot), float64(dragSlotItem.ItemID.MaxStack())))

			// account for how much was painted into this slot
			m.dragPlaced[dragSlotID] = dragSlotItem.ItemCount - oldCount

			// reduce available to distribute in this painting round
			dragged.ItemCount = dragged.ItemCount - (dragSlotItem.ItemCount - oldCount)
			m.clickable.SetSlot(dragSlotID, dragSlotItem) // save paint-over slot
		}

		m.cursor.ItemCount = dragged.ItemCount
		return true, nil
	case endLeftMouseDrag, endRightMouseDrag:
		var err error
		if slotID != slotOutsideWindow {
			err = fmt.Errorf("slotID must be -999 to end the drag")
		}

		if !m.dragged.IsPresent {
			err = fmt.Errorf("was not dragging")
		}

		m.dragSlots = nil
		m.dragPlaced = nil
		m.dragged = Slot{}
		return false, err
	default:
		return false, fmt.Errorf("button %d not supported for dragging", button)
	}
}
