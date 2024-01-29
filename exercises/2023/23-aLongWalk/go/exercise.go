package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 23.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	tm := parseInput(input)

	visited := map[Point]int{tm.Start: 0}
	tm.walk(tm.Start, Point{0, 1}, visited)

	return visited[tm.End], nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	tm := parseInput(input)

	tm.walk(tm.Start, Point{0, 1}, map[Point]int{tm.Start: 0})
	tm.getJunctions()

	paths := tm.getPaths()
	visited := make([]bool, len(tm.Junctions))
	visited[paths[tm.Start][0].index] = true

	return tm.getLongestPath(paths, tm.Start, tm.End, 0, visited), nil
}
