package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 5.
type Exercise struct {
	common.BaseExercise
}

// react fully collapses the polymer: adjacent units of the same letter but
// opposite case annihilate, and each removal may expose a new reacting pair. A
// stack handles the cascade in one O(n) pass — two ASCII letters react exactly
// when they differ only in the case bit, i.e. x^y == 0x20. Units whose lowercase
// form equals skip are ignored entirely; skip == 0 keeps every unit.
func react(polymer string, skip byte) int {
	stack := make([]byte, 0, len(polymer))

	for i := 0; i < len(polymer); i++ {
		u := polymer[i]
		if skip != 0 && u|0x20 == skip {
			continue
		}
		if n := len(stack); n > 0 && stack[n-1]^u == 0x20 {
			stack = stack[:n-1]
		} else {
			stack = append(stack, u)
		}
	}

	return len(stack)
}

// One returns the answer to the first part of the exercise.
// answer: 11476
func (e Exercise) One(instr string) (any, error) {
	return react(strings.TrimSpace(instr), 0), nil
}

// Two returns the answer to the second part of the exercise.
// answer: 5446
func (e Exercise) Two(instr string) (any, error) {
	polymer := strings.TrimSpace(instr)

	best := len(polymer)
	for c := byte('a'); c <= 'z'; c++ {
		if n := react(polymer, c); n < best {
			best = n
		}
	}

	return best, nil
}
