package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 7.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) ([]int, error) {
	var positions []int

	for _, f := range strings.Split(strings.TrimSpace(instr), ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}

		n, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("parsing position %q: %w", f, err)
		}

		positions = append(positions, n)
	}

	return positions, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// cheapestAlignment evaluates the total fuel to gather every crab at each
// candidate position across the full range and returns the minimum. cost maps a
// distance to its fuel, so the same sweep serves both parts.
func cheapestAlignment(positions []int, cost func(d int) int) int {
	lo, hi := positions[0], positions[0]
	for _, p := range positions {
		if p < lo {
			lo = p
		}
		if p > hi {
			hi = p
		}
	}

	best := -1
	for target := lo; target <= hi; target++ {
		total := 0
		for _, p := range positions {
			total += cost(abs(p - target))
		}
		if best < 0 || total < best {
			best = total
		}
	}

	return best
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	positions, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Linear fuel: one unit per step.
	fuel := cheapestAlignment(positions, func(d int) int { return d })

	return fmt.Sprintf("%d", fuel), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	positions, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Triangular fuel: step k costs k, so distance d costs d(d+1)/2.
	fuel := cheapestAlignment(positions, func(d int) int { return d * (d + 1) / 2 })

	return fmt.Sprintf("%d", fuel), nil
}
