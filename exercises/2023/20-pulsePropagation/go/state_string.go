// Code generated by "stringer -type=State"; DO NOT EDIT.

package exercises

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Off-0]
	_ = x[On-1]
}

const _State_name = "OffOn"

var _State_index = [...]uint8{0, 3, 5}

func (i State) String() string {
	if i < 0 || i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}