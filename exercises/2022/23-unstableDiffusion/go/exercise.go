package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 23.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 4034
func (e Exercise) One(instr string) (any, error) {
	elfCoords, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return diffuse(elfCoords, 1)
}

// Two returns the answer to the second part of the exercise.
// answer: 960
func (e Exercise) Two(instr string) (any, error) {
	elfCoords, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return diffuse(elfCoords, 2)
}

func parse(instr string) (map[point]string, error) {
	parsed := map[point]string{}

	for row, line := range strings.Split(instr, "\n") {
		for col, c := range line {
			switch c {
			case '#':
				parsed[point{col, row}] = "#"
			case '.':
			// do nothing
			default:
				return nil, fmt.Errorf("invalid character: %q", c)
			}
		}
	}

	return parsed, nil
}
