package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

const (
	targetSum int = 2020
)

// Exercise for Advent of Code 2020 day 1.
type Exercise struct {
	common.BaseExercise
}

// One finds the two entries summing to 2020 and returns their product.
func (c Exercise) One(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Complement lookup: for each entry, check whether 2020-entry was seen.
	seen := make(map[int]bool, len(entries))
	for _, a := range entries {
		if seen[targetSum-a] {
			return a * (targetSum - a), nil
		}
		seen[a] = true
	}

	return nil, fmt.Errorf("no answer found")
}

// Two finds the three entries summing to 2020 and returns their product.
func (c Exercise) Two(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Fix the first entry, then use a complement set over the remainder so the
	// inner search is linear: overall O(n^2).
	for i, a := range entries {
		need := targetSum - a
		seen := make(map[int]bool, len(entries)-i)
		for _, b := range entries[i+1:] {
			if b < need && seen[need-b] {
				return a * b * (need - b), nil
			}
			seen[b] = true
		}
	}

	return nil, fmt.Errorf("no answer found")
}

func parse(instr string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	entries := make([]int, 0, len(lines))

	for _, e := range lines {
		entry, err := strconv.Atoi(strings.TrimSpace(e))
		if err != nil {
			return nil, fmt.Errorf("parsing entry %q: %w", e, err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
