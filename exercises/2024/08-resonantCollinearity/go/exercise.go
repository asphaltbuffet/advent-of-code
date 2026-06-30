package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 8.
type Exercise struct {
	common.BaseExercise
}

type pt struct{ r, c int }

// parse returns the grid dimensions and antennas grouped by frequency.
func parse(instr string) (rows, cols int, antennas map[byte][]pt) {
	grid := strings.Fields(instr)
	rows = len(grid)
	if rows > 0 {
		cols = len(grid[0])
	}
	antennas = make(map[byte][]pt)
	for r, row := range grid {
		for c := 0; c < len(row); c++ {
			if ch := row[c]; ch != '.' {
				antennas[ch] = append(antennas[ch], pt{r, c})
			}
		}
	}
	return rows, cols, antennas
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	rows, cols, antennas := parse(instr)
	inBounds := func(p pt) bool { return p.r >= 0 && p.r < rows && p.c >= 0 && p.c < cols }

	nodes := make(map[pt]bool)
	for _, pts := range antennas {
		for i := 0; i < len(pts); i++ {
			for j := i + 1; j < len(pts); j++ {
				a, b := pts[i], pts[j]
				dr, dc := b.r-a.r, b.c-a.c
				for _, n := range []pt{{a.r - dr, a.c - dc}, {b.r + dr, b.c + dc}} {
					if inBounds(n) {
						nodes[n] = true
					}
				}
			}
		}
	}
	return len(nodes), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	rows, cols, antennas := parse(instr)
	inBounds := func(p pt) bool { return p.r >= 0 && p.r < rows && p.c >= 0 && p.c < cols }

	nodes := make(map[pt]bool)
	for _, pts := range antennas {
		for i := 0; i < len(pts); i++ {
			for j := i + 1; j < len(pts); j++ {
				a, b := pts[i], pts[j]
				dr, dc := b.r-a.r, b.c-a.c
				// Walk both directions from a at every integer multiple.
				for n := a; inBounds(n); n = (pt{n.r - dr, n.c - dc}) {
					nodes[n] = true
				}
				for n := a; inBounds(n); n = (pt{n.r + dr, n.c + dc}) {
					nodes[n] = true
				}
			}
		}
	}
	return len(nodes), nil
}
