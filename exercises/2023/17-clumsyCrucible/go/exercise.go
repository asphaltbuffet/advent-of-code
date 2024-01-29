package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 17.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	defer common.Close()

	const minMove = 1
	const maxMove = 3

	city, start, stop := parseCity(input)

	return countHeatLoss(city, start, stop, minMove, maxMove), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	defer common.Close()

	const minMove = 4
	const maxMove = 10

	city, start, stop := parseCity(input)

	return countHeatLoss(city, start, stop, minMove, maxMove), nil
}

func parseCity(input string) (map[Point]int, Point, Point) {
	lines := strings.Split(input, "\n")
	city := make(map[Point]int, len(lines)*len(lines[0]))

	width := len(lines[0])
	height := len(lines)

	start, end := Point{0, 0}, Point{width - 1, height - 1}

	for y, line := range lines {
		for x, char := range line {
			city[Point{x, y}] = int(char - '0')
		}
	}

	return city, start, end
}
