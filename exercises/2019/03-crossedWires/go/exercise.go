package exercises

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 3.
type Exercise struct {
	common.BaseExercise
}

type point struct{ x, y int }

// traceWire returns a map of each point visited to the step count of the first visit.
func traceWire(path string) map[point]int {
	visited := make(map[point]int)
	cur := point{0, 0}
	steps := 0

	for seg := range strings.SplitSeq(strings.TrimSpace(path), ",") {
		if len(seg) == 0 {
			continue
		}

		dir := seg[0]
		dist, err := strconv.Atoi(seg[1:])
		if err != nil {
			continue
		}

		var dx, dy int
		switch dir {
		case 'R':
			dx = 1
		case 'L':
			dx = -1
		case 'U':
			dy = 1
		case 'D':
			dy = -1
		}

		for range dist {
			cur.x += dx
			cur.y += dy
			steps++
			if _, seen := visited[cur]; !seen {
				visited[cur] = steps
			}
		}
	}

	return visited
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("expected 2 wire paths, got %d", len(lines))
	}

	wire1 := traceWire(lines[0])
	wire2 := traceWire(lines[1])

	minDist := math.MaxInt64
	for p := range wire1 {
		if _, ok := wire2[p]; ok {
			if d := abs(p.x) + abs(p.y); d < minDist {
				minDist = d
			}
		}
	}

	if minDist == math.MaxInt64 {
		return nil, errors.New("no intersections found")
	}

	return minDist, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("expected 2 wire paths, got %d", len(lines))
	}

	wire1 := traceWire(lines[0])
	wire2 := traceWire(lines[1])

	minSteps := math.MaxInt64
	for p, s1 := range wire1 {
		if s2, ok := wire2[p]; ok {
			if total := s1 + s2; total < minSteps {
				minSteps = total
			}
		}
	}

	if minSteps == math.MaxInt64 {
		return nil, errors.New("no intersections found")
	}

	return minSteps, nil
}
