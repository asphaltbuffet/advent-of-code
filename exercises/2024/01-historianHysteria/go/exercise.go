package exercises

import (
	"fmt"
	"slices"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	a, b, err := parse(instr)
	if err != nil {
		return nil, err
	}

	if len(a) != len(b) {
		return nil, fmt.Errorf("mismatched lengths: %d and %d", len(a), len(b))
	}

	slices.Sort(a)
	slices.Sort(b)

	sum := 0

	for i := 0; i < len(a); i++ {
		sum += abs(a[i] - b[i])
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	a, b, err := parse(instr)
	if err != nil {
		return nil, err
	}

	if len(a) != len(b) {
		return nil, fmt.Errorf("mismatched lengths: %d and %d", len(a), len(b))
	}

	counts := make(map[int]int)

	for _, l := range b {
		counts[l]++
	}

	// calc similarity score
	sim := 0
	for _, loc := range a {
		sim += loc * counts[loc]
	}

	return sim, nil
}
