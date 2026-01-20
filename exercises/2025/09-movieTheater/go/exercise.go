package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 9.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	f, _ := NewFloor(instr)

	return f.BiggestRect(), nil
}

// Two returns the answer to the second part of the exercise.
// LOW: 123274025
// HIGH: 4653414735
// HIGH: 4455008748
func (e Exercise) Two(instr string) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
		}
	}()

	floor, _ := NewFloor(instr)

	bfr := floor.BoundedRect()

	return bfr, nil
}
