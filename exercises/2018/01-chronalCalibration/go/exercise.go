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

// parse turns the input into a slice of signed frequency changes, tolerating
// blank lines and any surrounding whitespace (the elf may deliver input with or
// without a trailing newline).
func parse(instr string) ([]int, error) {
	fields := strings.Fields(instr)
	changes := make([]int, 0, len(fields))

	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		changes = append(changes, n)
	}

	return changes, nil
}

// One returns the answer to the first part of the exercise.
// answer: 416
func (c Exercise) One(instr string) (any, error) {
	changes, err := parse(instr)
	if err != nil {
		return nil, err
	}

	sum := 0
	for _, n := range changes {
		sum += n
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 56752
func (c Exercise) Two(instr string) (any, error) {
	changes, err := parse(instr)
	if err != nil {
		return nil, err
	}

	sum := 0
	seen := map[int]bool{0: true}

	for {
		for _, n := range changes {
			sum += n

			if seen[sum] {
				return sum, nil
			}

			seen[sum] = true
		}
	}
}
