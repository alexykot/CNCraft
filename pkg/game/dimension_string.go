// Code generated by "stringer -type=Dimension settings.go"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nether - -1]
	_ = x[Overworld-0]
	_ = x[TheEnd-1]
}

const _Dimension_name = "NetherOverworldTheEnd"

var _Dimension_index = [...]uint8{0, 6, 15, 21}

func (i Dimension) String() string {
	i -= -1
	if i < 0 || i >= Dimension(len(_Dimension_index)-1) {
		return "Dimension(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _Dimension_name[_Dimension_index[i]:_Dimension_index[i+1]]
}
