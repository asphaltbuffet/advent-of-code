package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 15.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	return sumAllSteps(input), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	steps := strings.Split(input, ",")
	ops := make([]*Op, len(steps))

	for i, step := range steps {
		op, err := parseStep(step)
		if err != nil {
			return nil, fmt.Errorf("failed to parse step %q: %w", step, err)
		}

		ops[i] = op
	}

	return calcFocusingPower(ops), nil
}
