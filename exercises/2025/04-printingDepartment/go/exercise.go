package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 4.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	f, err := NewFloor(instr)
	if err != nil {
		return nil, fmt.Errorf("create floor: %w", err)
	}

	return f.CountRolls(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	f, err := NewFloor(instr)
	if err != nil {
		return nil, fmt.Errorf("create floor: %w", err)
	}

	return f.RemoveRolls(), nil
}
