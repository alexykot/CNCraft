// Code generated by "stringer -type=clickMode window.go"; DO NOT EDIT.

package items

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[simpleClick-0]
	_ = x[shftClick-1]
	_ = x[numberKey-2]
	_ = x[middleClick-3]
	_ = x[drop-4]
	_ = x[drag-5]
	_ = x[doubleClick-6]
}

const _clickMode_name = "simpleClickshftClicknumberKeymiddleClickdropdragdoubleClick"

var _clickMode_index = [...]uint8{0, 11, 20, 29, 40, 44, 48, 59}

func (i clickMode) String() string {
	if i >= clickMode(len(_clickMode_index)-1) {
		return "clickMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _clickMode_name[_clickMode_index[i]:_clickMode_index[i+1]]
}