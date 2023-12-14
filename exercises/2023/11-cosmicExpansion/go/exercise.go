package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 11.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	ex := expandImage(strings.Split(instr, "\n"))

	sum := sumDistances(ex, 1)

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	ex := expandImage(strings.Split(instr, "\n"))

	sum := sumDistances(ex, 999999)

	return sum, nil
}
