package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 10.
type Exercise struct {
	common.BaseExercise
}

var pairs = map[rune]rune{'(': ')', '[': ']', '{': '}', '<': '>'}

// scanLine runs one line through a bracket stack. If a closing bracket does not
// match the most recent opener the line is corrupted and that bracket is
// returned. Otherwise the line is incomplete and the leftover open brackets
// (top of stack first) are returned as the completion needed to close it.
func scanLine(line string) (corrupt rune, completion []rune) {
	var stack []rune

	for _, c := range line {
		if close, ok := pairs[c]; ok {
			stack = append(stack, close) // remember the close we now expect
			continue
		}

		// c is a closing bracket; it must match the expected one.
		if len(stack) == 0 || stack[len(stack)-1] != c {
			return c, nil
		}
		stack = stack[:len(stack)-1]
	}

	// Remaining expected closers, innermost first, complete the line.
	for i := len(stack) - 1; i >= 0; i-- {
		completion = append(completion, stack[i])
	}

	return 0, completion
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	errorScore := map[rune]int{')': 3, ']': 57, '}': 1197, '>': 25137}

	total := 0
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if corrupt, _ := scanLine(line); corrupt != 0 {
			total += errorScore[corrupt]
		}
	}

	return fmt.Sprintf("%d", total), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	completeScore := map[rune]int{')': 1, ']': 2, '}': 3, '>': 4}

	var scores []int
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		corrupt, completion := scanLine(line)
		if corrupt != 0 || len(completion) == 0 {
			continue // skip corrupted (and any already-complete) lines
		}

		score := 0
		for _, c := range completion {
			score = score*5 + completeScore[c]
		}
		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return nil, fmt.Errorf("no incomplete lines to score")
	}

	sort.Ints(scores)

	return fmt.Sprintf("%d", scores[len(scores)/2]), nil
}
