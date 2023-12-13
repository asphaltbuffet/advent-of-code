package exercises

import (
	"slices"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 9.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var sum int

	for _, line := range strings.Split(instr, "\n") {
		history := lineToIntSlice(line)
		sum += calculateReductions(history)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var sum int

	for _, line := range strings.Split(instr, "\n") {
		history := lineToIntSlice(line)
		slices.Reverse(history)
		sum += calculateReductions(history)
	}

	return sum, nil
}
