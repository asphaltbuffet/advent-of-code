package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 13.
type Exercise struct {
	common.BaseExercise
}

type point struct{ x, y int }

type fold struct {
	axis byte // 'x' or 'y'
	at   int
}

func parse(instr string) (map[point]bool, []fold, error) {
	dotsPart, foldsPart, ok := strings.Cut(strings.TrimRight(instr, "\n"), "\n\n")
	if !ok {
		return nil, nil, fmt.Errorf("input missing blank-line separator")
	}

	dots := map[point]bool{}
	for _, line := range strings.Split(strings.TrimSpace(dotsPart), "\n") {
		var p point
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d,%d", &p.x, &p.y); err != nil {
			return nil, nil, fmt.Errorf("parsing dot %q: %w", line, err)
		}
		dots[p] = true
	}

	var folds []fold
	for _, line := range strings.Split(strings.TrimSpace(foldsPart), "\n") {
		var f fold
		if _, err := fmt.Sscanf(strings.TrimSpace(line), "fold along %c=%d", &f.axis, &f.at); err != nil {
			return nil, nil, fmt.Errorf("parsing fold %q: %w", line, err)
		}
		folds = append(folds, f)
	}

	return dots, folds, nil
}

// apply reflects every dot across the fold line, returning the new dot set.
func apply(dots map[point]bool, f fold) map[point]bool {
	out := make(map[point]bool, len(dots))
	for p := range dots {
		switch f.axis {
		case 'x':
			if p.x > f.at {
				p.x = 2*f.at - p.x
			}
		case 'y':
			if p.y > f.at {
				p.y = 2*f.at - p.y
			}
		}
		out[p] = true
	}
	return out
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	dots, folds, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}
	if len(folds) == 0 {
		return nil, fmt.Errorf("no folds in input")
	}

	dots = apply(dots, folds[0])

	return fmt.Sprintf("%d", len(dots)), nil
}

// render draws the dot set as ASCII art (█ lit, ░ dark) with no trailing newline.
func render(dots map[point]bool) string {
	maxX, maxY := 0, 0
	for p := range dots {
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	var sb strings.Builder
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if dots[point{x, y}] {
				sb.WriteString("█")
			} else {
				sb.WriteString("░")
			}
		}
		if y != maxY {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	dots, folds, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	for _, f := range folds {
		dots = apply(dots, f)
	}

	return render(dots), nil
}
