package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 17.
type Exercise struct {
	common.BaseExercise
}

type target struct {
	xMin, xMax, yMin, yMax int
}

func parse(instr string) (target, error) {
	var t target
	_, err := fmt.Sscanf(strings.TrimSpace(instr), "target area: x=%d..%d, y=%d..%d",
		&t.xMin, &t.xMax, &t.yMin, &t.yMax)
	if err != nil {
		return t, fmt.Errorf("parsing target %q: %w", instr, err)
	}
	return t, nil
}

// launch simulates one shot and returns the peak height and whether it hit the
// target. Drag pulls vx toward zero each step; gravity subtracts one from vy.
func launch(vx, vy int, t target) (peak int, hit bool) {
	x, y := 0, 0
	peak = 0
	for {
		x += vx
		y += vy
		if vx > 0 {
			vx--
		} else if vx < 0 {
			vx++
		}
		vy--

		if y > peak {
			peak = y
		}

		if x >= t.xMin && x <= t.xMax && y >= t.yMin && y <= t.yMax {
			return peak, true
		}

		// Past the target with no way back: overshot x, or falling below it.
		if x > t.xMax || y < t.yMin {
			return peak, false
		}
	}
}

// search sweeps the bounded velocity space and returns the best peak height and
// the number of velocities that hit. vx ranges 1..xMax (positive targets here);
// vy ranges yMin..-yMin (above -yMin overshoots below on the way down).
func search(t target) (bestPeak, hits int) {
	for vx := 1; vx <= t.xMax; vx++ {
		for vy := t.yMin; vy <= -t.yMin; vy++ {
			peak, hit := launch(vx, vy, t)
			if hit {
				hits++
				if peak > bestPeak {
					bestPeak = peak
				}
			}
		}
	}
	return bestPeak, hits
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	t, err := parse(instr)
	if err != nil {
		return nil, err
	}

	best, _ := search(t)

	return fmt.Sprintf("%d", best), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	t, err := parse(instr)
	if err != nil {
		return nil, err
	}

	_, hits := search(t)

	return fmt.Sprintf("%d", hits), nil
}
