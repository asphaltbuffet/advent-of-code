package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 3.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) []string {
	return strings.Split(strings.TrimSpace(instr), "\n")
}

// countTrees walks the given slope from the top-left, wrapping horizontally
// (the map repeats to the right), and returns how many trees ('#') are hit.
func countTrees(grid []string, right, down int) int {
	trees, x := 0, 0
	w := len(grid[0])
	for y := 0; y < len(grid); y += down {
		if grid[y][x%w] == '#' {
			trees++
		}
		x += right
	}
	return trees
}

// One counts trees hit descending the right-3, down-1 slope.
func (e Exercise) One(instr string) (any, error) {
	grid := parse(instr)
	return fmt.Sprintf("%d", countTrees(grid, 3, 1)), nil
}

// slopes are the five trajectories checked in part two.
var slopes = [][2]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}

// Two multiplies the tree counts across all five slopes.
func (e Exercise) Two(instr string) (any, error) {
	grid := parse(instr)
	product := 1
	for _, s := range slopes {
		product *= countTrees(grid, s[0], s[1])
	}
	return fmt.Sprintf("%d", product), nil
}
