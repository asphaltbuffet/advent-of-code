package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 11.
type Exercise struct {
	common.BaseExercise
}

const (
	floor    = '.'
	empty    = 'L'
	occupied = '#'
)

// grid is the seat layout with a row stride.
type grid struct {
	cells []byte
	w, h  int
}

func parseGrid(instr string) grid {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	g := grid{h: len(lines), w: len(lines[0])}
	g.cells = make([]byte, g.w*g.h)
	for y, line := range lines {
		copy(g.cells[y*g.w:], line)
	}
	return g
}

func (g grid) at(x, y int) byte { return g.cells[y*g.w+x] }

var dirs = [8][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}

// adjacentOccupied counts occupied seats in the 8 cells around (x, y).
func adjacentOccupied(g grid, x, y int) int {
	n := 0
	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && nx < g.w && ny >= 0 && ny < g.h && g.at(nx, ny) == occupied {
			n++
		}
	}
	return n
}

// visibleOccupied counts the first visible seat in each of the 8 directions that
// is occupied, seeing over floor.
func visibleOccupied(g grid, x, y int) int {
	n := 0
	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]
		for nx >= 0 && nx < g.w && ny >= 0 && ny < g.h {
			if c := g.at(nx, ny); c != floor {
				if c == occupied {
					n++
				}
				break
			}
			nx += d[0]
			ny += d[1]
		}
	}
	return n
}

// stabilize runs the seating rules to a fixed point and returns the number of
// occupied seats. neighbors counts crowding; a seat empties when its neighbor
// count reaches threshold.
func stabilize(g grid, neighbors func(grid, int, int) int, threshold int) int {
	cur := make([]byte, len(g.cells))
	copy(cur, g.cells)
	next := make([]byte, len(g.cells))
	work := grid{cur, g.w, g.h}

	for {
		changed := false
		for y := 0; y < g.h; y++ {
			for x := 0; x < g.w; x++ {
				c := work.at(x, y)
				nc := c
				switch c {
				case empty:
					if neighbors(work, x, y) == 0 {
						nc = occupied
					}
				case occupied:
					if neighbors(work, x, y) >= threshold {
						nc = empty
					}
				}
				next[y*g.w+x] = nc
				if nc != c {
					changed = true
				}
			}
		}
		copy(cur, next)
		if !changed {
			break
		}
	}

	count := 0
	for _, c := range cur {
		if c == occupied {
			count++
		}
	}
	return count
}

// One stabilizes with adjacent-neighbor rules (threshold 4).
func (e Exercise) One(instr string) (any, error) {
	g := parseGrid(instr)
	return fmt.Sprintf("%d", stabilize(g, adjacentOccupied, 4)), nil
}

// Two stabilizes with line-of-sight rules (threshold 5).
func (e Exercise) Two(instr string) (any, error) {
	g := parseGrid(instr)
	return fmt.Sprintf("%d", stabilize(g, visibleOccupied, 5)), nil
}
