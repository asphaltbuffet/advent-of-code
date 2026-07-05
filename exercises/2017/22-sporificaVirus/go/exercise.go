package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 22.
type Exercise struct {
	common.BaseExercise
}

type pt struct{ x, y int }

// Node states.
const (
	clean = iota
	weakened
	infected
	flagged
)

// Directions, clockwise from up: up, right, down, left.
var dirs = [4]pt{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

// parseGrid loads the starting infected nodes into a state map, centred so the
// carrier begins at the origin. It returns the map and the carrier's start.
func parseGrid(instr string) map[pt]int {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	rows := len(lines)
	cols := len(lines[0])
	grid := map[pt]int{}
	for r, line := range lines {
		for c, ch := range line {
			if ch == '#' {
				grid[pt{c - cols/2, r - rows/2}] = infected
			}
		}
	}
	return grid
}

// One runs 10000 bursts of the two-state virus and returns how many caused an
// infection.
func (e Exercise) One(instr string) (any, error) {
	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0 // up
	infections := 0

	for range 10000 {
		if grid[pos] == infected {
			dir = (dir + 1) & 3 // turn right
			delete(grid, pos)   // becomes clean
		} else {
			dir = (dir + 3) & 3 // turn left
			grid[pos] = infected
			infections++
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}
	}
	return infections, nil
}

// Two runs 10 million bursts of the four-state virus and returns how many
// bursts turned a node infected.
func (e Exercise) Two(instr string) (any, error) {
	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0
	infections := 0

	for range 10_000_000 {
		switch grid[pos] {
		case clean:
			dir = (dir + 3) & 3 // turn left
			grid[pos] = weakened
		case weakened:
			grid[pos] = infected // no turn
			infections++
		case infected:
			dir = (dir + 1) & 3 // turn right
			grid[pos] = flagged
		case flagged:
			dir = (dir + 2) & 3 // reverse
			delete(grid, pos)   // back to clean
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}
	}
	return infections, nil
}
