package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 11.
type Exercise struct {
	common.BaseExercise
}

// hexDeltas maps each hex direction to its cube-coordinate step. Cube
// coordinates keep x+y+z == 0, so distance is a single formula.
var hexDeltas = map[string][3]int{
	"n":  {0, 1, -1},
	"s":  {0, -1, 1},
	"ne": {1, 0, -1},
	"sw": {-1, 0, 1},
	"nw": {-1, 1, 0},
	"se": {1, -1, 0},
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// hexDist returns the number of steps from the origin to a cube coordinate.
func hexDist(x, y, z int) int {
	return (abs(x) + abs(y) + abs(z)) / 2
}

// walk follows the path and returns the final distance from origin and the
// greatest distance reached at any point along the way.
func walk(instr string) (int, int) {
	var furthest int
	x, y, z := 0, 0, 0
	for step := range strings.SplitSeq(strings.TrimSpace(instr), ",") {
		d := hexDeltas[strings.TrimSpace(step)]
		x, y, z = x+d[0], y+d[1], z+d[2]
		if dist := hexDist(x, y, z); dist > furthest {
			furthest = dist
		}
	}
	return hexDist(x, y, z), furthest
}

// One returns the fewest steps back to the origin from the path's end.
func (e Exercise) One(instr string) (any, error) {
	final, _ := walk(instr)
	return final, nil
}

// Two returns the furthest distance from the origin reached during the walk.
func (e Exercise) Two(instr string) (any, error) {
	_, furthest := walk(instr)
	return furthest, nil
}

// hexToPixel converts cube coordinates to pointy-top hex pixel coordinates.
