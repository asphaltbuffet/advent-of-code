// Code generated by "stringer -type=PipeShape"; DO NOT EDIT.

package exercises

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoPipe-46]
	_ = x[VerticalPipe-124]
	_ = x[HorizontalPipe-45]
	_ = x[NECornerPipe-76]
	_ = x[NWCornerPipe-74]
	_ = x[SWCornerPipe-55]
	_ = x[SECornerPipe-70]
	_ = x[StartPipe-83]
	_ = x[InvalidPipe-88]
}

const (
	_PipeShape_name_0 = "HorizontalPipeNoPipe"
	_PipeShape_name_1 = "SWCornerPipe"
	_PipeShape_name_2 = "SECornerPipe"
	_PipeShape_name_3 = "NWCornerPipe"
	_PipeShape_name_4 = "NECornerPipe"
	_PipeShape_name_5 = "StartPipe"
	_PipeShape_name_6 = "InvalidPipe"
	_PipeShape_name_7 = "VerticalPipe"
)

var (
	_PipeShape_index_0 = [...]uint8{0, 14, 20}
)

func (i PipeShape) String() string {
	switch {
	case 45 <= i && i <= 46:
		i -= 45
		return _PipeShape_name_0[_PipeShape_index_0[i]:_PipeShape_index_0[i+1]]
	case i == 55:
		return _PipeShape_name_1
	case i == 70:
		return _PipeShape_name_2
	case i == 74:
		return _PipeShape_name_3
	case i == 76:
		return _PipeShape_name_4
	case i == 83:
		return _PipeShape_name_5
	case i == 88:
		return _PipeShape_name_6
	case i == 124:
		return _PipeShape_name_7
	default:
		return "PipeShape(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}