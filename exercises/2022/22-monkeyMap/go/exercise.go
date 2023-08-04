package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 22.
type Exercise struct {
	common.BaseExercise
}

// this is the dimensions of a block, hacking this in for now
var blockSize = 4

// One returns the answer to the first part of the exercise.
// answer: 103224
func (e Exercise) One(instr string) (any, error) {
	if strings.Count(instr, "\n") > 50 {
		blockSize = 50
	}

	board := parse(instr, blockSize)

	state := followPath(board, false)

	return 1000*(state.row+1) + 4*(state.col+1) + int(state.face), nil
}

// Two returns the answer to the second part of the exercise.
// answer: 189097
func (e Exercise) Two(instr string) (any, error) {
	if strings.Count(instr, "\n") > 50 {
		blockSize = 50
	}

	board := parse(instr, blockSize)

	state := followPath(board, true)

	return 1000*(state.row+1) + 4*(state.col+1) + int(state.face), nil
}
