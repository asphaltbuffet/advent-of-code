package exercises

import (
	"fmt"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 6.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var sum int

	nums, ops, w := LoadHomework(instr)

	for i := range w {
		// get numbers
		vals := []int{}
		for j := i; j < len(nums); j += w {
			n, err := strconv.Atoi(nums[j])
			if err != nil {
				return nil, fmt.Errorf("P%d: %w", i, err)
			}

			vals = append(vals, n)
		}

		sum += op[ops[i]](vals...)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var sum int

	pp, err := RTLParse(instr)
	if err != nil {
		return nil, err
	}

	for _, p := range pp {
		sum += op[p.Operator](p.Numbers...)
	}

	return sum, nil
}
