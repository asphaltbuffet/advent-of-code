package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 10.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	mm := ParseMachines(instr)
	// fmt.Println(mm)

	sum := 0
	for _, m := range mm {
		sum += m.GetButtonPresses()
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}
