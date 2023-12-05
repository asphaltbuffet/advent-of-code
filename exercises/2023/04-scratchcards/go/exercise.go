package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 4.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// not 1007
func (e Exercise) One(instr string) (any, error) {
	var sum int

	for _, line := range strings.Split(instr, "\n") {
		card := New(line)
		sum += card.Score()
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	total := countTotalCards(strings.Split(instr, "\n"))

	return total, nil
}
