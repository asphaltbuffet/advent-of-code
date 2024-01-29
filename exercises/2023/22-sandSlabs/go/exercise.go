package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 22.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	bricks := parseInput(instr)

	bricksBelow, bricksAbove := getBrickOrder(bricks)

	canDisintegrate := getNonSupportingBricks(bricks, bricksBelow, bricksAbove)

	return len(canDisintegrate), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	bricks := parseInput(instr)

	bricksBelow, bricksAbove := getBrickOrder(bricks)

	canDisintegrate := getNonSupportingBricks(bricks, bricksBelow, bricksAbove)

	total := countDisintegratable(bricks, bricksBelow, bricksAbove, canDisintegrate)

	return total, nil
}
