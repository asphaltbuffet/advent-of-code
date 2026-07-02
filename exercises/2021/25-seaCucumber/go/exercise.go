package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 25.
type Exercise struct {
	common.BaseExercise
}

// grid is a toroidal seafloor: rows of '.', '>' (east herd) and 'v' (south herd).
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

// step advances one move: the east herd moves, then the south herd moves against
// the already-updated grid. It reports whether any cucumber moved.
func step(g *grid) bool {
	moved := false

	// East herd.
	next := make([]byte, len(g.cells))
	copy(next, g.cells)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.at(x, y) != '>' {
				continue
			}
			nx := (x + 1) % g.w
			if g.at(nx, y) == '.' {
				next[y*g.w+x] = '.'
				next[y*g.w+nx] = '>'
				moved = true
			}
		}
	}
	g.cells = next

	// South herd (against the east-updated grid).
	next = make([]byte, len(g.cells))
	copy(next, g.cells)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.at(x, y) != 'v' {
				continue
			}
			ny := (y + 1) % g.h
			if g.at(x, ny) == '.' {
				next[y*g.w+x] = '.'
				next[ny*g.w+x] = 'v'
				moved = true
			}
		}
	}
	g.cells = next

	return moved
}

// One counts the steps until no sea cucumber can move.
func (e Exercise) One(instr string) (any, error) {
	g := parseGrid(instr)
	count := 1
	for step(&g) {
		count++
	}
	return fmt.Sprintf("%d", count), nil
}

// Two is the free star for completing all other days; there is no puzzle here.
func (e Exercise) Two(_ string) (any, error) {
	return "Merry Christmas!", nil
}
