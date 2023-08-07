package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 24.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 264
func (e Exercise) One(instr string) (any, error) {
	basin, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	return startToEnd(basin, 0)
}

// Two returns the answer to the second part of the exercise.
// answer: 789
func (e Exercise) Two(instr string) (any, error) {
	basin, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	first, err := startToEnd(basin, 0)
	if err != nil {
		return nil, fmt.Errorf("calculating first trip: %w", err)
	}

	second, err := endToStart(basin, first)
	if err != nil {
		return nil, fmt.Errorf("calculating second trip: %w", err)
	}

	third, err := startToEnd(basin, second)
	if err != nil {
		return nil, fmt.Errorf("calculating third trip: %w", err)
	}

	return third, nil
}

func startToEnd(b *basin, elapsed int) (int, error) {
	return calcPath(b.winds, b.start, b.end, b.totalRows, b.totalCols, elapsed)
}

func endToStart(b *basin, elapsed int) (int, error) {
	return calcPath(b.winds, b.end, b.start, b.totalRows, b.totalCols, elapsed)
}
