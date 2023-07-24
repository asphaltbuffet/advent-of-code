package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

type cube struct {
	x, y, z int
}

// Exercise for Advent of Code 2022 day 18.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 4456
func (c Exercise) One(instr string) (any, error) {
	cubes, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	cubeMap := make(map[cube]bool)
	for _, c := range cubes {
		cubeMap[c] = true
	}

	totalExposure := len(cubes) * 6

	adjacent := []cube{
		{0, 0, -1},
		{0, 0, 1},
		{0, -1, 0},
		{0, 1, 0},
		{-1, 0, 0},
		{1, 0, 0},
	}

	// check all cubes for adjacent cubes
	for _, c := range cubes {
		for _, a := range adjacent {
			x := c.x + a.x
			y := c.y + a.y
			z := c.z + a.z

			if cubeMap[cube{x, y, z}] {
				totalExposure--
			}
		}
	}

	return totalExposure, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	return nil, nil
}

func parse(instr string) ([]cube, error) {
	var cubes []cube

	for _, line := range strings.Split(instr, "\n") {
		c := cube{}

		_, err := fmt.Sscanf(line, "%d,%d,%d", &c.x, &c.y, &c.z)
		if err != nil {
			return nil, fmt.Errorf("parsing cube from %q: %w", line, err)
		}

		cubes = append(cubes, c)
	}

	return cubes, nil
}
