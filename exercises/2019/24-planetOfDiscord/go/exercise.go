package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 24.
type Exercise struct {
	common.BaseExercise
}

// parseGrid reads a 5x5 bug grid and returns a uint32 bitmask.
// Bit i (row*5+col) is set if that tile is a bug (#).
func parseGrid(instr string) uint32 {
	var state uint32
	i := 0
	for _, ch := range strings.TrimSpace(instr) {
		switch ch {
		case '#':
			state |= 1 << i
			i++
		case '.':
			i++
		}
		// skip newlines and other chars
	}
	return state
}

// stepGrid advances the cellular automaton one minute.
//
//nolint:gocognit // neighbor-counting requires many boundary checks
func stepGrid(state uint32) uint32 {
	var next uint32
	for row := range 5 {
		for col := range 5 {
			i := row*5 + col
			neighbors := 0
			if row > 0 && (state>>(i-5))&1 == 1 {
				neighbors++
			}
			if row < 4 && (state>>(i+5))&1 == 1 {
				neighbors++
			}
			if col > 0 && (state>>(i-1))&1 == 1 {
				neighbors++
			}
			if col < 4 && (state>>(i+1))&1 == 1 {
				neighbors++
			}
			bug := (state>>i)&1 == 1
			if bug && neighbors == 1 {
				next |= 1 << i
			} else if !bug && (neighbors == 1 || neighbors == 2) {
				next |= 1 << i
			}
		}
	}
	return next
}

// One returns the answer to the first part of the exercise.
// Simulate until a layout repeats; return its biodiversity rating.
func (e Exercise) One(instr string) (any, error) {
	state := parseGrid(instr)
	seen := map[uint32]bool{}
	for {
		if seen[state] {
			return int(state), nil
		}
		seen[state] = true
		state = stepGrid(state)
	}
}

// countNeighborsRecursive counts bugs adjacent to tile (row,col) at the given level.
// levels is the current map of level->bitmask.
//
//nolint:gocognit // recursive Plutonian grid requires many per-edge cases
func countNeighborsRecursive(levels map[int]uint32, level, row, col int) int {
	count := 0

	type dir struct{ dr, dc int }
	dirs := []dir{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, d := range dirs {
		nr, nc := row+d.dr, col+d.dc

		switch {
		case nr == 2 && nc == 2:
			// stepping into center -> recurse into child level
			child := levels[level+1]
			switch {
			case d.dr == 1: // came from above center, child row 0
				for c := range 5 {
					if (child>>(0*5+c))&1 == 1 {
						count++
					}
				}
			case d.dr == -1: // came from below center, child row 4
				for c := range 5 {
					if (child>>(4*5+c))&1 == 1 {
						count++
					}
				}
			case d.dc == 1: // came from left of center, child col 0
				for r := range 5 {
					if (child>>(r*5+0))&1 == 1 {
						count++
					}
				}
			case d.dc == -1: // came from right of center, child col 4
				for r := range 5 {
					if (child>>(r*5+4))&1 == 1 {
						count++
					}
				}
			}
		case nr < 0:
			// top border -> parent tile (1,2) = index 7
			parent := levels[level-1]
			if (parent>>7)&1 == 1 {
				count++
			}
		case nr > 4:
			// bottom border -> parent tile (3,2) = index 17
			parent := levels[level-1]
			if (parent>>17)&1 == 1 {
				count++
			}
		case nc < 0:
			// left border -> parent tile (2,1) = index 11
			parent := levels[level-1]
			if (parent>>11)&1 == 1 {
				count++
			}
		case nc > 4:
			// right border -> parent tile (2,3) = index 13
			parent := levels[level-1]
			if (parent>>13)&1 == 1 {
				count++
			}
		default:
			// normal same-level neighbor
			same := levels[level]
			if (same>>(nr*5+nc))&1 == 1 {
				count++
			}
		}
	}
	return count
}

// stepRecursive advances the recursive multi-level grid one step.
//
//nolint:gocognit // level-range expansion and per-cell rules require multiple conditions
func stepRecursive(levels map[int]uint32) map[int]uint32 {
	// determine range of active levels
	minLevel, maxLevel := 0, 0
	first := true
	for lv := range levels {
		if levels[lv] == 0 {
			continue
		}
		if first {
			minLevel, maxLevel = lv, lv
			first = false
		} else {
			if lv < minLevel {
				minLevel = lv
			}
			if lv > maxLevel {
				maxLevel = lv
			}
		}
	}
	// expand range by one to allow growth
	minLevel--
	maxLevel++

	next := map[int]uint32{}
	for lv := minLevel; lv <= maxLevel; lv++ {
		var mask uint32
		for row := range 5 {
			for col := range 5 {
				if row == 2 && col == 2 {
					continue // center is always empty
				}
				i := row*5 + col
				n := countNeighborsRecursive(levels, lv, row, col)
				bug := (levels[lv]>>i)&1 == 1
				if bug && n == 1 {
					mask |= 1 << i
				} else if !bug && (n == 1 || n == 2) {
					mask |= 1 << i
				}
			}
		}
		if mask != 0 {
			next[lv] = mask
		}
	}
	return next
}

// Two returns the answer to the second part of the exercise.
// Simulate the recursive Plutonian grid for 200 minutes (10 for example input containing '?').
func (e Exercise) Two(instr string) (any, error) {
	steps := 200
	if strings.Contains(instr, "?") {
		steps = 10
	}

	// parse initial grid at level 0; treat '?' as '.'
	cleaned := strings.ReplaceAll(instr, "?", ".")
	initial := parseGrid(cleaned)
	initial &^= 1 << 12 // ensure center is clear

	levels := map[int]uint32{0: initial}

	for range steps {
		levels = stepRecursive(levels)
	}

	total := 0
	for _, mask := range levels {
		for mask != 0 {
			total += int(mask & 1)
			mask >>= 1
		}
	}
	return total, nil
}
