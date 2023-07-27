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

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	for i, a := range entries {
		for _, b := range entries[i+1:] {
			if a+b == targetSum {
				return a * b, nil
			}
		}
	}

	return nil, fmt.Errorf("no answer found")
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	for i, a := range entries {
		for _, b := range entries[i+1:] {
			// skip if a+b is already greater than 2020
			if a+b >= targetSum {
				continue
			}

			for _, c := range entries[i+2:] {
				if a+b+c == targetSum {
					return a * b * c, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no answer found")
}

func parse(instr string) ([]int, error) {
	entries := make([]int, 0, strings.Count(instr, "\n"))

	for _, e := range strings.Split(instr, "\n") {
		entry, err := strconv.Atoi(e)
		if err != nil {
			return nil, fmt.Errorf("parsing entry: %w", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
