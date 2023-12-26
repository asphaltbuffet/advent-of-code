package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 16.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	m, err := parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	m.Power(Point{0, 0}, Right)

	// m.DebugPrintContraption()
	// m.DebugPrintEnergized()

	return m.CountEnergized(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	m, err := parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	max := 0

	// test starting from top edge
	for x := 0; x < m.Width; x++ {
		c := m.Clone()

		c.Power(Point{x, 0}, Down)

		e := c.CountEnergized()

		if e > max {
			max = e
		}
	}

	// test starting from bottom edge
	for x := 0; x < m.Width; x++ {
		c := m.Clone()

		c.Power(Point{x, m.Height - 1}, Up)

		e := c.CountEnergized()

		if e > max {
			max = e
		}
	}

	// test starting from left edge
	for y := 0; y < m.Height; y++ {
		c := m.Clone()

		c.Power(Point{0, y}, Right)

		e := c.CountEnergized()

		if e > max {
			max = e
		}
	}

	// test starting from right edge
	for y := 0; y < m.Height; y++ {
		c := m.Clone()

		c.Power(Point{m.Width - 1, y}, Left)

		e := c.CountEnergized()

		if e > max {
			max = e
		}
	}

	return max, nil
}
