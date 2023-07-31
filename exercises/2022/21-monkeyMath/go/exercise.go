package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 21.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	raw := parse(instr)

	result, err := calc("root", raw, make(map[string]int))
	if err != nil {
		return nil, fmt.Errorf("calculating part 1 answer: %w", err)
	}

	return result, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	return nil, nil
}
