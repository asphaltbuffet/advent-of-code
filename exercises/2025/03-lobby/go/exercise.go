package exercises

import (
	"fmt"
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
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
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
