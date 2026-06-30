package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 19.
type Exercise struct {
	common.BaseExercise
}

// parse returns the available towel patterns and the list of desired designs.
func parse(instr string) (patterns, designs []string) {
	parts := strings.SplitN(strings.TrimSpace(instr), "\n\n", 2)
	for _, p := range strings.Split(parts[0], ",") {
		patterns = append(patterns, strings.TrimSpace(p))
	}
	for _, d := range strings.Split(strings.TrimSpace(parts[1]), "\n") {
		if d = strings.TrimSpace(d); d != "" {
			designs = append(designs, d)
		}
	}
	return patterns, designs
}

// arrangements returns the number of distinct ways to form design by
// concatenating patterns (0 if impossible). ways[i] counts ways to form the
// prefix of length i.
func arrangements(design string, patterns []string) int {
	n := len(design)
	ways := make([]int, n+1)
	ways[0] = 1
	for i := 1; i <= n; i++ {
		for _, p := range patterns {
			if len(p) <= i && design[i-len(p):i] == p {
				ways[i] += ways[i-len(p)]
			}
		}
	}
	return ways[n]
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	patterns, designs := parse(instr)
	count := 0
	for _, d := range designs {
		if arrangements(d, patterns) > 0 {
			count++
		}
	}
	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	patterns, designs := parse(instr)
	total := 0
	for _, d := range designs {
		total += arrangements(d, patterns)
	}
	return total, nil
}
