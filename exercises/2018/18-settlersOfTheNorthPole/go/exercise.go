package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 18.
type Exercise struct {
	common.BaseExercise
}

const (
	open       = '.'
	trees      = '|'
	lumberyard = '#'
)

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	grid := parse(instr)
	for range 10 {
		grid = step(grid)
	}
	return resourceValue(grid), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	const target = 1_000_000_000
	grid := parse(instr)

	// The state settles into a cycle; record the minute each state was first seen,
	// then fast-forward the remaining minutes modulo the cycle length.
	seen := map[string]int{}
	for minute := range target {
		key := string(bytesOf(grid))
		if first, ok := seen[key]; ok {
			period := minute - first
			remaining := (target - minute) % period
			for range remaining {
				grid = step(grid)
			}
			return resourceValue(grid), nil
		}
		seen[key] = minute
		grid = step(grid)
	}
	return resourceValue(grid), nil
}

// parse reads the acre map into a 2D byte grid.
func parse(instr string) [][]byte {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	grid := make([][]byte, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return grid
}

// bytesOf flattens the grid to a single byte slice for use as a map key.
func bytesOf(grid [][]byte) []byte {
	b := make([]byte, 0, len(grid)*len(grid[0]))
	for _, row := range grid {
		b = append(b, row...)
	}
	return b
}

// step advances every acre one minute by the transformation rules.
func step(grid [][]byte) [][]byte {
	h, w := len(grid), len(grid[0])
	next := make([][]byte, h)
	for y := range h {
		next[y] = make([]byte, w)
		for x := range w {
			next[y][x] = nextCell(grid, x, y)
		}
	}
	return next
}

// neighbors counts adjacent trees and lumberyards around (x, y).
func neighbors(grid [][]byte, x, y int) (int, int) {
	var treeCount, lumberCount int
	h, w := len(grid), len(grid[0])
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx < 0 || ny < 0 || nx >= w || ny >= h {
				continue
			}
			switch grid[ny][nx] {
			case trees:
				treeCount++
			case lumberyard:
				lumberCount++
			}
		}
	}
	return treeCount, lumberCount
}

func nextCell(grid [][]byte, x, y int) byte {
	t, l := neighbors(grid, x, y)
	switch grid[y][x] {
	case open:
		if t >= 3 {
			return trees
		}
	case trees:
		if l >= 3 {
			return lumberyard
		}
	case lumberyard:
		if l >= 1 && t >= 1 {
			return lumberyard
		}
		return open
	}
	return grid[y][x]
}

// resourceValue is the number of wooded acres times the number of lumberyards.
func resourceValue(grid [][]byte) int {
	var t, l int
	for _, row := range grid {
		for _, c := range row {
			switch c {
			case trees:
				t++
			case lumberyard:
				l++
			}
		}
	}
	return t * l
}
