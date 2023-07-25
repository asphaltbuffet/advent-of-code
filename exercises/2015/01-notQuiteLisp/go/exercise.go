package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	up := strings.Count(instr, "(")
	down := strings.Count(instr, ")")

	// fmt.Printf("%d - %d\n", up, down)

	return up - down, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	floor := 0

	for i, c := range instr {
		switch c {
		case '(':
			// up
			floor++
		case ')':
			// down
			floor--
		}

		if floor < 0 {
			return i + 1, nil
		}
	}

	return 0, fmt.Errorf("santa never enters the basement")
}
