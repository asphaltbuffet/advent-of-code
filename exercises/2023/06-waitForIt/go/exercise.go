package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 6.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	races := parseRaces(instr)
	total := 1

	for _, r := range races {
		n := r.CountFasterTimes()
		total *= n
	}

	return total, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	r := parseBigRace(instr)

	n := r.CountFasterTimes()

	return n, nil
}
