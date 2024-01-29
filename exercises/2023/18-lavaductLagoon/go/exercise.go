package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 18.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	steps, _, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	boundary, err := steps.GetBoundaryPoints()
	if err != nil {
		return nil, err
	}

	return shoelace(boundary), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	_, steps, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	boundary, err := steps.GetBoundaryPoints()
	if err != nil {
		return nil, err
	}

	return shoelace(boundary), nil
}
