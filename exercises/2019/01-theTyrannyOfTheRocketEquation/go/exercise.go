package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	sum := 0

	for _, mass := range strings.Split(instr, "\n") {
		m, err := strconv.Atoi(mass)
		if err != nil {
			return nil, err
		}

		sum += calculateFuel(m, false)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	sum := 0

	for _, mass := range strings.Split(instr, "\n") {
		m, err := strconv.Atoi(mass)
		if err != nil {
			return nil, err
		}

		sum += calculateFuel(m, true)
	}

	return sum, nil
}

func calculateFuel(mass int, includeFuel bool) int {
	fuel := (mass / 3) - 2

	switch {
	case !includeFuel:
		return fuel
	case fuel <= 0:
		return 0
	default:
		return fuel + calculateFuel(fuel, includeFuel)
	}
}
