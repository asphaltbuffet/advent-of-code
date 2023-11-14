package sorting

import (
	"fmt"
	"strconv"
)

// DescribeNumber should return a string describing the number.
func DescribeNumber(f float64) string {
	return fmt.Sprintf("This is the number %0.1f", f)
}

type NumberBox interface {
	Number() int
}

// DescribeNumberBox should return a string describing the NumberBox.
func DescribeNumberBox(nb NumberBox) string {
	f := float64(nb.Number())

	return fmt.Sprintf("This is a box containing the number %0.1f", f)
}

type FancyNumber struct {
	n string
}

func (i FancyNumber) Value() string {
	return i.n
}

type FancyNumberBox interface {
	Value() string
}

// ExtractFancyNumber should return the integer value for a FancyNumber
// and 0 if any other FancyNumberBox is supplied.
func ExtractFancyNumber(fnb FancyNumberBox) int {
	switch v := fnb.(type) {
	case FancyNumber:
		if i, err := strconv.Atoi(v.Value()); err == nil {
			return i
		}

		return 0
	default:
		return 0
	}
}

// DescribeFancyNumberBox should return a string describing the FancyNumberBox.
func DescribeFancyNumberBox(fnb FancyNumberBox) string {
	return fmt.Sprintf("This is a fancy box containing the number %d.0", ExtractFancyNumber(fnb))
}

// DescribeAnything should return a string describing whatever it contains.
func DescribeAnything(i interface{}) string {
	out := "This is %sthe number %0.1f"
	var box string
	var n float64

	switch v := i.(type) {
	case int:
		box = ""
		n = float64(v)
	case float64:
		box = ""
		n = v
	case NumberBox:
		box = "a box containing "
		n = float64(v.(NumberBox).Number())
	case FancyNumberBox:
		box = "a fancy box containing "
		n = float64(ExtractFancyNumber(v.(FancyNumberBox)))
	default:
		return "Return to sender"
	}

	return fmt.Sprintf(out, box, n)
}
