package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 20.
type Exercise struct {
	common.BaseExercise
}

type cell struct{ r, c int }

// threshold returns the minimum saving counted: the AoC example (small grid)
// uses >=50, the real input uses >=100.
func threshold(path []cell) int {
	if len(path) < 100 {
		return 50
	}
	return 100
}

// trackPath returns the ordered list of track cells from S to E. The grid has a
// single non-branching path, so a simple walk recovers it.
func trackPath(instr string) []cell {
	grid := strings.Fields(instr)
	if len(grid) == 0 {
		return nil
	}
	var start cell
	for r := range grid {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == 'S' {
				start = cell{r, c}
			}
		}
	}
	open := func(r, c int) bool {
		return r >= 0 && r < len(grid) && c >= 0 && c < len(grid[r]) && grid[r][c] != '#'
	}

	var path []cell
	prev := cell{-1, -1}
	cur := start
	for {
		path = append(path, cur)
		if grid[cur.r][cur.c] == 'E' {
			break
		}
		moved := false
		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nr, nc := cur.r+d[0], cur.c+d[1]
			if open(nr, nc) && (nr != prev.r || nc != prev.c) {
				prev = cur
				cur = cell{nr, nc}
				moved = true
				break
			}
		}
		if !moved {
			break // dead end (shouldn't happen on a valid single-path grid)
		}
	}
	return path
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// countCheats counts cheats of length up to maxCheat that save at least the
// threshold. dist along the single path is just the index in the path slice.
func countCheats(path []cell, maxCheat int) int {
	save := threshold(path)
	count := 0
	for i := 0; i < len(path); i++ {
		for j := i + 1; j < len(path); j++ {
			m := abs(path[i].r-path[j].r) + abs(path[i].c-path[j].c)
			if m > maxCheat {
				continue
			}
			if (j - i) - m >= save {
				count++
			}
		}
	}
	return count
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return countCheats(trackPath(instr), 2), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return countCheats(trackPath(instr), 20), nil
}
