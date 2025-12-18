package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 8.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	jj, err := NewJunctions(instr)
	if err != nil {
		return nil, err
	}

	var wires int
	if len(*jj) < 100 {
		wires = 10
	} else {
		wires = 1000
	}

	circuits := jj.CreateCircuits(wires)
	// multiply size of 3 largest circuits
	var prod int
	prod = circuits[0] * circuits[1] * circuits[2]

	return prod, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	jj, err := NewJunctions(instr)
	if err != nil {
		return nil, err
	}

	return jj.EndCircuits(), nil
}
