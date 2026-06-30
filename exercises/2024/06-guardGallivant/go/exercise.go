package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 6.
type Exercise struct {
	common.BaseExercise
}

// Facing directions, ordered so that (dir+1)%4 is a right turn.
var dirs = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // up, right, down, left

// pos is a cell coordinate.
type pos struct{ r, c int }

// state is a cell plus a facing direction.
type state struct {
	r, c, dir int
}

// findStart returns the guard's starting position (the '^' cell).
func findStart(grid []string) pos {
	for r, row := range grid {
		for c := 0; c < len(row); c++ {
			if row[c] == '^' {
				return pos{r, c}
			}
		}
	}
	return pos{-1, -1}
}

// walk simulates the guard from start. extra is an optional added obstruction
// (use {-1,-1} for none). It returns the set of visited cells and whether the
// guard fell into a loop (true) versus walking off the map (false).
func walk(grid []string, start pos, extra pos) (map[pos]bool, bool) {
	rows := len(grid)
	blocked := func(r, c int) bool {
		return grid[r][c] == '#' || (r == extra.r && c == extra.c)
	}

	visited := make(map[pos]bool)
	seen := make(map[state]bool)
	r, c, dir := start.r, start.c, 0
	for {
		visited[pos{r, c}] = true
		st := state{r, c, dir}
		if seen[st] {
			return visited, true
		}
		seen[st] = true

		nr, nc := r+dirs[dir][0], c+dirs[dir][1]
		if nr < 0 || nr >= rows || nc < 0 || nc >= len(grid[nr]) {
			return visited, false
		}
		if blocked(nr, nc) {
			dir = (dir + 1) % 4
			continue
		}
		r, c = nr, nc
	}
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	grid := strings.Fields(instr)
	start := findStart(grid)
	visited, _ := walk(grid, start, pos{-1, -1})
	return len(visited), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	grid := strings.Fields(instr)
	start := findStart(grid)

	// Only cells on the original path can change the route.
	path, _ := walk(grid, start, pos{-1, -1})
	delete(path, start)

	count := 0
	for cell := range path {
		if _, looped := walk(grid, start, cell); looped {
			count++
		}
	}
	return count, nil
}
