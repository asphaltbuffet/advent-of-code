package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 10.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	m, start, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	n, err := findPathLength(m, start)
	if err != nil {
		return nil, err
	}

	return n / 2, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}
