package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	sum := 0

	for _, f := range strings.Split(instr, "\n") {
		n, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		sum += n
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	sum := 0
	seen := map[int]bool{0: true}

	for {
		for _, f := range strings.Split(instr, "\n") {
			n, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("parsing input: %w", err)
			}

			sum += n

			if seen[sum] {
				return sum, nil
			}

			seen[sum] = true
		}
	}
}
