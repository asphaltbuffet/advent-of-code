package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 24.
type Exercise struct {
	common.BaseExercise
}

// hex is an axial hex coordinate.
type hex struct{ q, r int }

// hexDirs maps each direction token to its axial step.
var hexDirs = map[string]hex{
	"e":  {1, 0},
	"w":  {-1, 0},
	"ne": {1, -1},
	"sw": {-1, 1},
	"nw": {0, -1},
	"se": {0, 1},
}

// walk follows a direction string and returns the destination tile.
func walk(line string) hex {
	var pos hex
	for i := 0; i < len(line); {
		var dir string
		switch line[i] {
		case 'e', 'w':
			dir = line[i : i+1]
			i++
		default: // n/s prefix, two-char direction
			dir = line[i : i+2]
			i += 2
		}
		d := hexDirs[dir]
		pos.q += d.q
		pos.r += d.r
	}
	return pos
}

// blackTiles returns the set of black tiles after applying all flip paths.
func blackTiles(instr string) map[hex]bool {
	black := map[hex]bool{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		t := walk(strings.TrimSpace(line))
		if black[t] {
			delete(black, t)
		} else {
			black[t] = true
		}
	}
	return black
}

// One counts black tiles after the initial flips.
func (e Exercise) One(instr string) (any, error) {
	return fmt.Sprintf("%d", len(blackTiles(instr))), nil
}

// step runs one day of the hex cellular automaton: a black tile stays black only
// with 1 or 2 black neighbors, and a white tile turns black with exactly 2.
func step(black map[hex]bool) map[hex]bool {
	counts := map[hex]int{}
	for t := range black {
		for _, d := range hexDirs {
			counts[hex{t.q + d.q, t.r + d.r}]++
		}
	}

	next := map[hex]bool{}
	for t, n := range counts {
		if black[t] {
			if n == 1 || n == 2 {
				next[t] = true
			}
		} else if n == 2 {
			next[t] = true
		}
	}
	return next
}

// Two runs 100 days of the flipping rules and returns the black-tile count.
func (e Exercise) Two(instr string) (any, error) {
	black := blackTiles(instr)
	for day := 0; day < 100; day++ {
		black = step(black)
	}
	return fmt.Sprintf("%d", len(black)), nil
}
