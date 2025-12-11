package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 5.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	inputs := strings.Split(instr, "\n\n")

	inv, err := LoadInventory(inputs[0])
	if err != nil {
		return nil, err
	}

	ing, err := LoadIngredients(inputs[1])
	if err != nil {
		return nil, err
	}

	var count int
	for _, i := range ing {
		if inv.IsFresh(i) {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}
