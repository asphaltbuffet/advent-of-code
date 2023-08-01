package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 22.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	mm, path := parse(instr)

	b := newBoard(mm)

	// b.debugPrint()

	b.move(path)

	row := b.location.position.y + 1
	col := b.location.position.x + 1

	return 1000*row + 4*col + b.facing, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	return nil, nil
}
