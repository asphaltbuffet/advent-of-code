package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 24.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	hailStones, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	hsCount := len(hailStones)
	intersectCount := 0

	for i := 0; i < hsCount-1; i++ {
		for j := i + 1; j < hsCount; j++ {
			if hailStones[i].Intersects(hailStones[j]) {
				intersectCount++
			}
		}
	}

	return intersectCount, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	hailStones, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	// potential velocities for each dimension.
	tmpVel := make(map[Dimension][]int)

	// Iterate through all pairs of hailstone positions to update potential velocities.
	for i := 0; i < len(hailStones)-1; i++ {
		for j := i + 1; j < len(hailStones); j++ {
			for _, d := range []Dimension{X, Y, Z} {
				updatePotential(d, hailStones[i], hailStones[j], tmpVel)
			}
		}
	}

	// if there is exactly one potential velocity in each dimension, we can solve for position.
	if len(tmpVel[X]) != 1 || len(tmpVel[Y]) != 1 || len(tmpVel[Z]) != 1 {
		return nil, fmt.Errorf("no solution found")
	}

	// final position based on the calculated velocity and the first two hailstone positions.
	rockVel := Vector3D{float64(tmpVel[X][0]), float64(tmpVel[Y][0]), float64(tmpVel[Z][0])}
	pos, err := getPosition(rockVel, hailStones[0], hailStones[1])

	return int(pos.X + pos.Y + pos.Z), err
}
