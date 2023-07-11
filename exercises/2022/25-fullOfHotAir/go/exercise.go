package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 25.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 2==221=-002=0-02-000
func (c Exercise) One(instr string) (any, error) {
	sum := 0

	for _, line := range strings.Split(instr, "\n") {
		sum += Decode(line)
	}

	return Encode(sum), nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	return nil, nil
}
