package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 5.
type Exercise struct {
	common.BaseExercise
}

type point struct {
	x, y int
}

type segment struct {
	a, b point
}

func parse(instr string) ([]segment, error) {
	var segs []segment

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var s segment
		if _, err := fmt.Sscanf(line, "%d,%d -> %d,%d", &s.a.x, &s.a.y, &s.b.x, &s.b.y); err != nil {
			return nil, fmt.Errorf("parsing segment %q: %w", line, err)
		}

		segs = append(segs, s)
	}

	return segs, nil
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// countOverlaps rasterizes every segment into a grid and returns how many cells
// are covered by at least two lines. When diagonals is false, 45° segments are
// skipped (part one). The uniform stepping — one unit of sign(dx)/sign(dy) per
// step for max(|dx|,|dy|) steps — handles horizontal, vertical, and diagonal
// lines the same way.
func countOverlaps(segs []segment, diagonals bool) int {
	grid := map[point]int{}

	for _, s := range segs {
		dx := s.b.x - s.a.x
		dy := s.b.y - s.a.y

		if !diagonals && dx != 0 && dy != 0 {
			continue // skip diagonals for part one
		}

		steps := abs(dx)
		if abs(dy) > steps {
			steps = abs(dy)
		}

		sx, sy := sign(dx), sign(dy)
		p := s.a
		for i := 0; i <= steps; i++ {
			grid[p]++
			p.x += sx
			p.y += sy
		}
	}

	count := 0
	for _, n := range grid {
		if n >= 2 {
			count++
		}
	}

	return count
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	segs, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", countOverlaps(segs, false)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	segs, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", countOverlaps(segs, true)), nil
}
