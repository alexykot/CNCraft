package items

import (
    "testing"

    "github.com/stretchr/testify/suite"

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
    s.i = NewInventory()
}

const (
    rowTop1 = 9
    rowTop2 = 10
    rowtop3 = 11
    rowtop4 = 12
    rowtop5 = 13
    rowtop6 = 14
    rowtop7 = 15
    rowtop8 = 16
    rowtop9 = 17

    rowMiddle1 = 9
    rowMiddle2 = 10
    rowMiddle3 = 11
    rowMiddle4 = 12
    rowMiddle5 = 13
    rowMiddle6 = 14
    rowMiddle7 = 15
    rowMiddle8 = 16
    rowMiddle9 = 17
/**/
    rowBottom1 = 9
    rowBottom2 = 10
    rowBottom3 = 11
    rowBottom4 = 12
    rowBottom5 = 13
    rowBottom6 = 14
    rowBottom7 = 15
    rowBottom8 = 16
    rowBottom9 = 17

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
func pickaxe() Slot {
    return Slot{IsPresent: true, ItemID: int16(pItems.DiamondPickaxe), ItemCount: 1}
}
func bedrock(stackCount ...int16) Slot {
    var itemCount int16 = 64
    if len(stackCount) > 0 {
        itemCount = stackCount[0]
    }

    return Slot{IsPresent: true, ItemID: int16(pItems.Bedrock), ItemCount: itemCount}
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
    shouldChange bool
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
            clickedItem:  pickaxe(),
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
            clickedItem:  pickaxe(),
        },
        {
            name:         "left_click/item_moveup_occipied",
            invStart:     []testSlot{tSlot(rowTop1, pickaxe()), tSlot(hotbar1, pickaxe())},
            invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(rowTop2, pickaxe())},
            cursorStart:  empty(),
            cursorEnd:    empty(),
            shouldChange: true,
            button:       leftMouseButton,
            slotID:       hotbar1,
            clickedItem:  pickaxe(),
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
            clickedItem:  pickaxe(),
        },
        {
            name:         "left_click/item_moveup_row_occipied",
            invStart:     []testSlot{
                tSlot(rowTop1, pickaxe()),
                tSlot(rowTop2, pickaxe()),
                tSlot(rowtop3, pickaxe()),
                tSlot(rowtop4, pickaxe()),
                tSlot(rowtop5, pickaxe()),
                tSlot(rowtop6, pickaxe()),
                tSlot(rowtop7, pickaxe()),
                tSlot(rowtop8, pickaxe()),
                tSlot(rowtop9, pickaxe()),
                tSlot(hotbar1, pickaxe()),
            },
            invEnd:       []testSlot{tSlot(rowTop1, pickaxe()), tSlot(rowTop2, pickaxe())},
            cursorStart:  empty(),
            cursorEnd:    empty(),
            shouldChange: true,
            button:       leftMouseButton,
            slotID:       hotbar1,
            clickedItem:  pickaxe(),
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
            clickedItem:  pickaxe(),
        },
    }
    s.runTests(shftClick, testCases)
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

            hasChanged, err := s.i.HandleClick(actionID, test.slotID, int16(mode), test.button, test.clickedItem)
            s.Require().NoError(err)
            s.Equal(test.shouldChange, hasChanged, "inventory has not changed")

            invCompare(test.invEnd, s.i.ToArray(), s.Require().Equal)
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
