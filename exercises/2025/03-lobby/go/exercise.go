package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 3.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var sum int
	for bank := range strings.Lines(instr) {
		sum += Largest(strings.Trim(bank, "\n"))
		// sum += LongLargest(strings.Trim(bank, "\n"), 2)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var sum int
	for bank := range strings.Lines(instr) {
		sum += LongLargest(strings.Trim(bank, "\n"), 12)
	}

	return sum, nil
}

// Largest creates the max 2-digit integer from consecutive integers in a string.
// It assumes that every character in the string is '0' -> '9'.
func Largest(b string) int {
	var l int
	var r int

	for i := 0; i < len(b)-1; i++ {
		n := int(b[i] - '0')
		if n > l {
			l = n
			r = -1
		} else if n > r {
			r = n
		}
	}

	//set r to last value if unset
	n := int(b[len(b)-1] - '0')
	if r < n {
		r = n
	}

	return l*10 + r
}

func LongLargest(b string, digits int) int {
	bank := make([]int, len(b))

	for i, r := range b {
		bank[i] = int(r - '0')
	}
	// fmt.Println("bank: ", bank)

	out := 0
	remaining := digits

	for i := 0; remaining > 0; {
		window := len(b) - i - remaining + 1
		// fmt.Printf("i=%d w=%d r= %d out=%d\n", i, window, remaining, out)

		m, idx := MaxInRange(bank[i : i+window])
		out = out*10 + m

		i += idx + 1
		remaining--
	}

	return out
}

func MaxInRange(b []int) (int, int) {
	idx := 0
	m := 0

	for i, n := range b {
		if n > m {
			m = n
			idx = i
		}
	}

	return m, idx
}
