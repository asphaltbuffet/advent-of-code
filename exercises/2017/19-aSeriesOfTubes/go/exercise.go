package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 19.
type Exercise struct {
	common.BaseExercise
}

type cell struct {
	r, c     int
	isLetter bool
	ch       byte
}

// walk follows the routing diagram from the top-row entry and returns the
// letters collected in order, the total number of steps taken, and the sequence
// of cells visited.
func walk(instr string) (string, int, []cell) {
	// Preserve leading spaces — they position the path. Only trim the trailing
	// newline so the final row isn't empty.
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")

	at := func(r, c int) byte {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[r]) {
			return ' '
		}
		return grid[r][c]
	}

	turn := func(r, c, dr, _ int) (int, int) {
		if dr != 0 { // moving vertically -> turn horizontally
			if at(r, c-1) != ' ' {
				return 0, -1
			}
			return 0, 1
		}
		// moving horizontally -> turn vertically
		if at(r-1, c) != ' ' {
			return -1, 0
		}
		return 1, 0
	}

	// Entry: the only '|' in the top row.
	r, c := 0, strings.IndexByte(grid[0], '|')
	dr, dc := 1, 0 // moving down

	var letters strings.Builder
	var path []cell
	steps := 0
	for at(r, c) != ' ' {
		ch := at(r, c)
		isLetter := ch >= 'A' && ch <= 'Z'
		if isLetter {
			letters.WriteByte(ch)
		}
		path = append(path, cell{r, c, isLetter, ch})
		if ch == '+' {
			dr, dc = turn(r, c, dr, dc)
		}
		r, c = r+dr, c+dc
		steps++
	}
	return letters.String(), steps, path
}

// One returns the letters collected along the path.
func (e Exercise) One(instr string) (any, error) {
	letters, _, _ := walk(instr)
	return letters, nil
}

// Two returns the total number of steps taken along the path.
func (e Exercise) Two(instr string) (any, error) {
	_, steps, _ := walk(instr)
	return steps, nil
}
