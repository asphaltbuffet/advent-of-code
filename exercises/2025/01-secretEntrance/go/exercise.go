package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var zero int
	pos := 50

	for _, l := range strings.Split(instr, "\n") {
		direction := l[0]
		val, _ := strconv.Atoi(l[1:])

		if direction == 'L' {
			pos = pos + 100 - (val % 100)
		} else {
			pos = pos + val
		}

		pos %= 100

		if pos == 0 {
			zero++
		}
	}

	return fmt.Sprintf("%d", zero), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}
