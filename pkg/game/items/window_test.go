package items

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	pItems "github.com/alexykot/cncraft/pkg/protocol/items"
)

func TestClickMgr(t *testing.T) {
	suite.Run(t, &clickMgrSuite{})
}

type clickMgrSuite struct {
	suite.Suite

	i *Inventory
}

func (s *clickMgrSuite) SetupSuite() {}

func (s *clickMgrSuite) SetupTest() {
	s.i = NewInventory(zap.NewNop())
}

const (
	rowTop1 = 9
	rowTop2 = 10
	rowTop3 = 11
	rowTop4 = 12
	rowTop5 = 13
	rowTop6 = 14
	rowTop7 = 15
	rowTop8 = 16
	rowTop9 = 17

	rowMiddle1 = 18
	rowMiddle2 = 19
	rowMiddle3 = 20
	rowMiddle4 = 21
	rowMiddle5 = 22
	rowMiddle6 = 23
	rowMiddle7 = 24
	rowMiddle8 = 25
	rowMiddle9 = 26

	rowBottom1 = 27
	rowBottom2 = 28
	rowBottom3 = 29
	rowBottom4 = 30
	rowBottom5 = 31
	rowBottom6 = 32
	rowBottom7 = 33
	rowBottom8 = 34
	rowBottom9 = 35

	hotbar1 = 36
	hotbar2 = 37
	hotbar3 = 38
	hotbar4 = 39
	hotbar5 = 40
	hotbar6 = 41
	hotbar7 = 42
	hotbar8 = 43
	hotbar9 = 44
)

type testSlot struct {
	Slot
	slotID int16
}

func tSlot(slotID int16, slot Slot) testSlot { return testSlot{slotID: slotID, Slot: slot} }

func pickaxe() Slot { return Slot{IsPresent: true, ItemID: pItems.DiamondPickaxe, ItemCount: 1} }
func showel() Slot  { return Slot{IsPresent: true, ItemID: pItems.DiamondShovel, ItemCount: 1} }

func bedrock(stackCount ...int16) Slot {
	var itemCount int16 = 64
	if len(stackCount) > 0 {
		itemCount = stackCount[0]
	}

	return Slot{IsPresent: true, ItemID: pItems.Bedrock, ItemCount: itemCount}
}
func empty() Slot { return Slot{} }

type testCase struct {
	name string

	invStart    []testSlot
	cursorStart Slot

	slotID      int16
	button      uint8
	clickedItem Slot

	invEnd       []testSlot
	cursorEnd    Slot
	dropped      Slot
	shouldChange bool

	isBeingDragged bool
}

func (s *clickMgrSuite) TestHandleClick_OkMode0() {
	testCases := []testCase{
		{
			name:         "left_click/single_item_pickup",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       nil,
			cursorStart:  empty(),
			cursorEnd:    pickaxe(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  pickaxe(),
		},
		{
			name:         "left_click/single_item_putdown_empty_slot",
			invStart:     nil,
			invEnd:       []testSlot{tSlot(hotbar2, pickaxe())},
			cursorStart:  pickaxe(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  empty(),
		},

		{
			name:         "left_click/stack_pickup",
			invStart:     []testSlot{tSlot(hotbar2, bedrock())},
			invEnd:       nil,
			cursorStart:  empty(),
			cursorEnd:    bedrock(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(),
		},
		{
			name:         "left_click/stack_putdown",
			invStart:     nil,
			invEnd:       []testSlot{tSlot(hotbar3, bedrock())},
			cursorStart:  bedrock(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar3,
			clickedItem:  empty(),
		},

		{
			name:         "left_click/halfstack_join_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(40))},
			cursorStart:  bedrock(20),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(20),
		},
		{
			name:         "left_click/overstack_join_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(40))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(64))},
			cursorStart:  bedrock(40),
			cursorEnd:    bedrock(16),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(40),
		},
		{
			name:         "left_click/replacement_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(40))},
			invEnd:       []testSlot{tSlot(hotbar2, pickaxe())},
			cursorStart:  pickaxe(),
			cursorEnd:    bedrock(40),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(40),
		},

		{
			name:         "right_click/single_item_pickup",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       nil,
			cursorStart:  empty(),
			cursorEnd:    pickaxe(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  pickaxe(),
		},
		{
			name:         "right_click/single_item_putdown_empty_slot",
			invStart:     nil,
			invEnd:       []testSlot{tSlot(hotbar2, pickaxe())},
			cursorStart:  pickaxe(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar2,
			clickedItem:  empty(),
		},

		{
			name:         "right_click/halfstack_pickup",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(40))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(20))},
			cursorStart:  empty(),
			cursorEnd:    bedrock(20),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(40),
		},
		{
			name:         "left_click/halfstack_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(20)), tSlot(hotbar3, bedrock(20))},
			cursorStart:  bedrock(20),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar3,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/stackitem_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(20)), tSlot(hotbar3, bedrock(1))},
			cursorStart:  bedrock(20),
			cursorEnd:    bedrock(19),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar3,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/stackitem_join_putdown",
			invStart:     []testSlot{tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar2, bedrock(21))},
			cursorStart:  bedrock(20),
			cursorEnd:    bedrock(19),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar2,
			clickedItem:  bedrock(20),
		},
		{
			name:         "left_click/outside_window_drop_single_item",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			cursorStart:  pickaxe(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       slotOutsideWindow,
			clickedItem:  empty(),
			dropped:      pickaxe(),
		},
		{
			name:         "right_click/outside_window_drop_single_item",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			cursorStart:  pickaxe(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       slotOutsideWindow,
			clickedItem:  empty(),
			dropped:      pickaxe(),
		},
		{
			name:         "left_click/outside_window_drop_stack_whole",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			cursorStart:  bedrock(20),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       slotOutsideWindow,
			clickedItem:  empty(),
			dropped:      bedrock(20),
		},
		{
			name:         "right_click/outside_window_drop_stack_item",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, bedrock(20))},
			cursorStart:  bedrock(20),
			cursorEnd:    bedrock(19),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       slotOutsideWindow,
			clickedItem:  empty(),
			dropped:      bedrock(1),
		},
	}
	s.runTests(simpleClick, testCases)
}

func (s *clickMgrSuite) TestHandleClick_OkMode1() {
	testCases := []testCase{
		{
			name:         "left_click/single_item_moveup",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/single_item_movedown",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/item_moveup_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(rowTop2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/item_movedown_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name: "left_click/item_moveup_row_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, pickaxe()),
				tSlot(rowTop2, pickaxe()),
				tSlot(rowTop3, pickaxe()),
				tSlot(rowTop4, pickaxe()),
				tSlot(rowTop5, pickaxe()),
				tSlot(rowTop6, pickaxe()),
				tSlot(rowTop7, pickaxe()),
				tSlot(rowTop8, pickaxe()),
				tSlot(rowTop9, pickaxe()),
				tSlot(hotbar1, pickaxe()),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, pickaxe()),
				tSlot(rowTop2, pickaxe()),
				tSlot(rowTop3, pickaxe()),
				tSlot(rowTop4, pickaxe()),
				tSlot(rowTop5, pickaxe()),
				tSlot(rowTop6, pickaxe()),
				tSlot(rowTop7, pickaxe()),
				tSlot(rowTop8, pickaxe()),
				tSlot(rowTop9, pickaxe()),
				tSlot(rowMiddle1, pickaxe()),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name: "left_click/stack_moveup_multistack_row_mostly_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, bedrock(63)),
				tSlot(rowTop2, bedrock(63)),
				tSlot(rowTop3, bedrock(63)),
				tSlot(rowTop4, bedrock(63)),
				tSlot(rowTop5, bedrock(63)),
				tSlot(rowTop6, bedrock(63)),
				tSlot(rowTop7, bedrock(63)),
				tSlot(rowTop8, bedrock(63)),
				tSlot(rowTop9, bedrock(63)),
				tSlot(hotbar1, bedrock(64)),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, bedrock(64)),
				tSlot(rowTop2, bedrock(64)),
				tSlot(rowTop3, bedrock(64)),
				tSlot(rowTop4, bedrock(64)),
				tSlot(rowTop5, bedrock(64)),
				tSlot(rowTop6, bedrock(64)),
				tSlot(rowTop7, bedrock(64)),
				tSlot(rowTop8, bedrock(64)),
				tSlot(rowTop9, bedrock(64)),
				tSlot(rowMiddle1, bedrock(55)),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name: "left_click/stack_moveup_multistack_range_mostly_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, bedrock(63)),
				tSlot(rowTop2, bedrock(63)),
				tSlot(rowTop3, bedrock(63)),
				tSlot(rowTop4, bedrock(63)),
				tSlot(rowTop5, bedrock(63)),
				tSlot(rowTop6, bedrock(63)),
				tSlot(rowTop7, bedrock(63)),
				tSlot(rowTop8, bedrock(63)),
				tSlot(rowTop9, bedrock(63)),
				tSlot(rowMiddle1, bedrock(63)),
				tSlot(rowMiddle2, bedrock(63)),
				tSlot(rowMiddle3, bedrock(63)),
				tSlot(rowMiddle4, bedrock(63)),
				tSlot(rowMiddle5, bedrock(63)),
				tSlot(rowMiddle6, bedrock(63)),
				tSlot(rowMiddle7, bedrock(63)),
				tSlot(rowMiddle8, bedrock(63)),
				tSlot(rowMiddle9, bedrock(63)),
				tSlot(rowBottom1, bedrock(63)),
				tSlot(rowBottom2, bedrock(63)),
				tSlot(rowBottom3, bedrock(63)),
				tSlot(rowBottom4, bedrock(63)),
				tSlot(rowBottom5, bedrock(63)),
				tSlot(rowBottom6, bedrock(63)),
				tSlot(rowBottom7, bedrock(63)),
				tSlot(rowBottom8, bedrock(63)),
				tSlot(rowBottom9, bedrock(63)),
				tSlot(hotbar1, bedrock(64)),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, bedrock(64)),
				tSlot(rowTop2, bedrock(64)),
				tSlot(rowTop3, bedrock(64)),
				tSlot(rowTop4, bedrock(64)),
				tSlot(rowTop5, bedrock(64)),
				tSlot(rowTop6, bedrock(64)),
				tSlot(rowTop7, bedrock(64)),
				tSlot(rowTop8, bedrock(64)),
				tSlot(rowTop9, bedrock(64)),
				tSlot(rowMiddle1, bedrock(64)),
				tSlot(rowMiddle2, bedrock(64)),
				tSlot(rowMiddle3, bedrock(64)),
				tSlot(rowMiddle4, bedrock(64)),
				tSlot(rowMiddle5, bedrock(64)),
				tSlot(rowMiddle6, bedrock(64)),
				tSlot(rowMiddle7, bedrock(64)),
				tSlot(rowMiddle8, bedrock(64)),
				tSlot(rowMiddle9, bedrock(64)),
				tSlot(rowBottom1, bedrock(64)),
				tSlot(rowBottom2, bedrock(64)),
				tSlot(rowBottom3, bedrock(64)),
				tSlot(rowBottom4, bedrock(64)),
				tSlot(rowBottom5, bedrock(64)),
				tSlot(rowBottom6, bedrock(64)),
				tSlot(rowBottom7, bedrock(64)),
				tSlot(rowBottom8, bedrock(64)),
				tSlot(rowBottom9, bedrock(64)),
				tSlot(hotbar1, bedrock(37)),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/item_movedown_row_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       leftMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/empty_slot_up",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: false,
			button:       leftMouseButton,
			slotID:       hotbar2,
			clickedItem:  empty(),
		},
		{
			name:         "left_click/empty_slot_down",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: false,
			button:       leftMouseButton,
			slotID:       rowTop2,
			clickedItem:  empty(),
		},

		// right click, same behaviour
		{
			name:         "right_click/single_item_moveup",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/single_item_movedown",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/item_moveup_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(rowTop2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/item_movedown_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name: "right_click/item_moveup_row_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, pickaxe()),
				tSlot(rowTop2, pickaxe()),
				tSlot(rowTop3, pickaxe()),
				tSlot(rowTop4, pickaxe()),
				tSlot(rowTop5, pickaxe()),
				tSlot(rowTop6, pickaxe()),
				tSlot(rowTop7, pickaxe()),
				tSlot(rowTop8, pickaxe()),
				tSlot(rowTop9, pickaxe()),
				tSlot(hotbar1, pickaxe()),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, pickaxe()),
				tSlot(rowTop2, pickaxe()),
				tSlot(rowTop3, pickaxe()),
				tSlot(rowTop4, pickaxe()),
				tSlot(rowTop5, pickaxe()),
				tSlot(rowTop6, pickaxe()),
				tSlot(rowTop7, pickaxe()),
				tSlot(rowTop8, pickaxe()),
				tSlot(rowTop9, pickaxe()),
				tSlot(rowMiddle1, pickaxe()),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name: "right_click/stack_moveup_multistack_row_mostly_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, bedrock(63)),
				tSlot(rowTop2, bedrock(63)),
				tSlot(rowTop3, bedrock(63)),
				tSlot(rowTop4, bedrock(63)),
				tSlot(rowTop5, bedrock(63)),
				tSlot(rowTop6, bedrock(63)),
				tSlot(rowTop7, bedrock(63)),
				tSlot(rowTop8, bedrock(63)),
				tSlot(rowTop9, bedrock(63)),
				tSlot(hotbar1, bedrock(64)),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, bedrock(64)),
				tSlot(rowTop2, bedrock(64)),
				tSlot(rowTop3, bedrock(64)),
				tSlot(rowTop4, bedrock(64)),
				tSlot(rowTop5, bedrock(64)),
				tSlot(rowTop6, bedrock(64)),
				tSlot(rowTop7, bedrock(64)),
				tSlot(rowTop8, bedrock(64)),
				tSlot(rowTop9, bedrock(64)),
				tSlot(rowMiddle1, bedrock(55)),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name: "right_click/stack_moveup_multistack_range_mostly_occupied",
			invStart: []testSlot{
				tSlot(rowTop1, bedrock(63)),
				tSlot(rowTop2, bedrock(63)),
				tSlot(rowTop3, bedrock(63)),
				tSlot(rowTop4, bedrock(63)),
				tSlot(rowTop5, bedrock(63)),
				tSlot(rowTop6, bedrock(63)),
				tSlot(rowTop7, bedrock(63)),
				tSlot(rowTop8, bedrock(63)),
				tSlot(rowTop9, bedrock(63)),
				tSlot(rowMiddle1, bedrock(63)),
				tSlot(rowMiddle2, bedrock(63)),
				tSlot(rowMiddle3, bedrock(63)),
				tSlot(rowMiddle4, bedrock(63)),
				tSlot(rowMiddle5, bedrock(63)),
				tSlot(rowMiddle6, bedrock(63)),
				tSlot(rowMiddle7, bedrock(63)),
				tSlot(rowMiddle8, bedrock(63)),
				tSlot(rowMiddle9, bedrock(63)),
				tSlot(rowBottom1, bedrock(63)),
				tSlot(rowBottom2, bedrock(63)),
				tSlot(rowBottom3, bedrock(63)),
				tSlot(rowBottom4, bedrock(63)),
				tSlot(rowBottom5, bedrock(63)),
				tSlot(rowBottom6, bedrock(63)),
				tSlot(rowBottom7, bedrock(63)),
				tSlot(rowBottom8, bedrock(63)),
				tSlot(rowBottom9, bedrock(63)),
				tSlot(hotbar1, bedrock(64)),
			},
			invEnd: []testSlot{
				tSlot(rowTop1, bedrock(64)),
				tSlot(rowTop2, bedrock(64)),
				tSlot(rowTop3, bedrock(64)),
				tSlot(rowTop4, bedrock(64)),
				tSlot(rowTop5, bedrock(64)),
				tSlot(rowTop6, bedrock(64)),
				tSlot(rowTop7, bedrock(64)),
				tSlot(rowTop8, bedrock(64)),
				tSlot(rowTop9, bedrock(64)),
				tSlot(rowMiddle1, bedrock(64)),
				tSlot(rowMiddle2, bedrock(64)),
				tSlot(rowMiddle3, bedrock(64)),
				tSlot(rowMiddle4, bedrock(64)),
				tSlot(rowMiddle5, bedrock(64)),
				tSlot(rowMiddle6, bedrock(64)),
				tSlot(rowMiddle7, bedrock(64)),
				tSlot(rowMiddle8, bedrock(64)),
				tSlot(rowMiddle9, bedrock(64)),
				tSlot(rowBottom1, bedrock(64)),
				tSlot(rowBottom2, bedrock(64)),
				tSlot(rowBottom3, bedrock(64)),
				tSlot(rowBottom4, bedrock(64)),
				tSlot(rowBottom5, bedrock(64)),
				tSlot(rowBottom6, bedrock(64)),
				tSlot(rowBottom7, bedrock(64)),
				tSlot(rowBottom8, bedrock(64)),
				tSlot(rowBottom9, bedrock(64)),
				tSlot(hotbar1, bedrock(37)),
			},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/item_movedown_row_occupied",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar2, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       rightMouseButton,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/empty_slot_up",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: false,
			button:       rightMouseButton,
			slotID:       hotbar2,
			clickedItem:  empty(),
		},
		{
			name:         "right_click/empty_slot_down",
			invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: false,
			button:       rightMouseButton,
			slotID:       rowTop2,
			clickedItem:  empty(),
		},
	}
	s.runTests(shftClick, testCases)
}

func (s *clickMgrSuite) TestHandleClick_OkMode2() {
	testCases := []testCase{
		{
			name:         "number_press/single_item_swap",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(rowTop1, showel())},
			invEnd:       []testSlot{tSlot(hotbar1, showel()), tSlot(rowTop1, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey1,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/empty_slot_swap_to_hotbar",
			invStart:     []testSlot{tSlot(rowTop1, showel())},
			invEnd:       []testSlot{tSlot(hotbar9, showel())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey9,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/empty_slot_swap_from_hotbar",
			invStart:     []testSlot{tSlot(hotbar3, showel())},
			invEnd:       []testSlot{tSlot(rowTop1, showel())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey3,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/stack_swap",
			invStart:     []testSlot{tSlot(hotbar1, bedrock()), tSlot(rowTop1, showel())},
			invEnd:       []testSlot{tSlot(hotbar1, showel()), tSlot(rowTop1, bedrock())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey1,
			slotID:       rowTop1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/single_item_swap_in_hotbar",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe()), tSlot(hotbar5, showel())},
			invEnd:       []testSlot{tSlot(hotbar1, showel()), tSlot(hotbar5, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey5,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/swap_to_empty_in_hotbar",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar5, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey5,
			slotID:       hotbar1,
			clickedItem:  empty(),
		},
		{
			name:         "number_press/swap_from_empty_in_hotbar",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{tSlot(hotbar5, pickaxe())},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKey1,
			slotID:       hotbar5,
			clickedItem:  empty(),
		},
	}
	s.runTests(numberKey, testCases)
}

func (s *clickMgrSuite) TestHandleClick_OkMode4() {
	testCases := []testCase{
		{
			name:         "drop/Q_press",
			invStart:     []testSlot{tSlot(hotbar1, pickaxe())},
			invEnd:       []testSlot{},
			cursorStart:  empty(),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       kbdKeyQ,
			slotID:       hotbar1,
			clickedItem:  empty(),
			dropped:      pickaxe(),
		},
	}

	s.runTests(drop, testCases)
}

func (s *clickMgrSuite) TestHandleClick_OkMode5() {
	testCases := []testCase{
		{
			name:         "drag/start_emptyslots",
			invStart:     []testSlot{},
			invEnd:       []testSlot{},
			cursorStart:  bedrock(20),
			cursorEnd:    empty(),
			shouldChange: true,
			button:       startLeftMouseDrag,
			slotID:       slotOutsideWindow,
			clickedItem:  empty(),
			dropped:      empty(),
		},
	}

	s.runTests(drag, testCases)
}

func (s *clickMgrSuite) runTests(mode clickMode, testCases []testCase) {
	var actionID int16
	for _, test := range testCases {
		s.Run(test.name, func() {
			actionID++

			s.i.cursor = test.cursorStart
			s.i.reset()
			for _, item := range test.invStart {
				s.i.SetSlot(item.slotID, item.Slot)
			}

			dropped, hasChanged, err := s.i.HandleClick(actionID, test.slotID, int16(mode), test.button, test.clickedItem)
			s.Require().NoError(err)
			s.Equal(test.shouldChange, hasChanged, "inventory has not changed")

			invCompare(test.invEnd, s.i.ToArray(), s.Require().Equal)

			if test.dropped.IsPresent {
				s.Require().NotNil(dropped)
				s.Equal(test.dropped.IsPresent, dropped.IsPresent)
				s.Equal(test.dropped.ItemID, dropped.ItemID)
				s.Equal(test.dropped.ItemCount, dropped.ItemCount)
			} else {
				s.Nil(dropped)
			}
		})
	}
}

func invCompare(expect []testSlot, actual []Slot, equaliser func(interface{}, interface{}, ...interface{})) {
	expectFull := make([]Slot, 46, 46)
	for _, item := range expect {
		expectFull[item.slotID] = item.Slot
	}

	equaliser(len(expectFull), len(actual))

	for i, expectItem := range expectFull {
		// println(fmt.Sprintf("%d. expected: %v; actual: %v", i, expectItem, actual[i]))
		equaliser(expectItem.IsPresent, actual[i].IsPresent)
		equaliser(expectItem.ItemID, actual[i].ItemID)
		equaliser(expectItem.ItemCount, actual[i].ItemCount)
	}
}
