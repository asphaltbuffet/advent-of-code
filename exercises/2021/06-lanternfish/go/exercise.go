package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 6.
type Exercise struct {
	common.BaseExercise
}

// parse reads the comma-separated timers into a histogram of counts per timer
// value (0..8). Tracking counts rather than individual fish keeps the work
// constant regardless of how large the population grows.
func parse(instr string) ([9]int, error) {
	var counts [9]int

	for _, f := range strings.Split(strings.TrimSpace(instr), ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}

		n, err := strconv.Atoi(f)
		if err != nil {
			return counts, fmt.Errorf("parsing timer %q: %w", f, err)
		}
		if n < 0 || n > 8 {
			return counts, fmt.Errorf("timer %d out of range 0..8", n)
		}

		counts[n]++
	}

	return counts, nil
}

// grow advances the population by the given number of days. Each day the fish at
// timer 0 both reset to 6 and spawn newborns at 8; every other bucket shifts
// down by one.
func grow(counts [9]int, days int) int {
	for d := 0; d < days; d++ {
		spawning := counts[0]
		for i := 0; i < 8; i++ {
			counts[i] = counts[i+1]
		}
		counts[6] += spawning // resetting parents
		counts[8] = spawning  // newborns
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	return total
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	counts, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", grow(counts, 80)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	counts, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", grow(counts, 256)), nil
}
