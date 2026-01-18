package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 12.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
		}
	}()

	f := Parse(instr)
	count := 0
	for _, r := range f.regions {
		// simple check if area of all presents is <= total area
		pArea := r.area
		for i, p := range r.presents {
			pArea -= f.presents[i] * p
		}

		if pArea >= 0 {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
		}
	}()

	// there is no part 2 for this day
	return 0, nil
}
