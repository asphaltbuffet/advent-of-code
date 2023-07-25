package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 1.
type Exercise struct {
	common.BaseExercise
}

const (
	north heading = iota
	east
	south
	west
)

type heading int

type point struct {
	h    heading
	x, y int
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	position := point{h: north, x: 0, y: 0}

	for _, i := range strings.Split(instr, ", ") {
		switch {
		case strings.HasPrefix(i, "L"):
			switch position.h {
			case north:
				position.h = west
			case east:
				position.h = north
			case south:
				position.h = east
			case west:
				position.h = south
			}
		case strings.HasPrefix(i, "R"):
			switch position.h {
			case north:
				position.h = east
			case east:
				position.h = south
			case south:
				position.h = west
			case west:
				position.h = north
			}
		}

		distance, err := strconv.Atoi(i[1:])
		if err != nil {
			return nil, fmt.Errorf("reading movement magnitude: %w", err)
		}

		switch position.h {
		case north:
			position.y += distance
		case east:
			position.x += distance
		case south:
			position.y -= distance
		case west:
			position.x -= distance
		}
	}

	return abs(position.x) + abs(position.y), nil
}

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	position := point{h: north, x: 0, y: 0}
	visited := make(map[string]bool)

	for _, i := range strings.Split(instr, ", ") {
		switch {
		case strings.HasPrefix(i, "L"):
			switch position.h {
			case north:
				position.h = west
			case east:
				position.h = north
			case south:
				position.h = east
			case west:
				position.h = south
			}
		case strings.HasPrefix(i, "R"):
			switch position.h {
			case north:
				position.h = east
			case east:
				position.h = south
			case south:
				position.h = west
			case west:
				position.h = north
			}
		}

		distance, err := strconv.Atoi(i[1:])
		if err != nil {
			return nil, fmt.Errorf("reading movement magnitude: %w", err)
		}

		for j := 0; j < distance; j++ {
			switch position.h {
			case north:
				position.y++
			case east:
				position.x++
			case south:
				position.y--
			case west:
				position.x--
			}

			location := fmt.Sprintf("%d,%d", position.x, position.y)
			if visited[location] {
				return abs(position.x) + abs(position.y), nil
			}

			visited[location] = true
		}
	}

	return nil, fmt.Errorf("no location visited twice")
}
