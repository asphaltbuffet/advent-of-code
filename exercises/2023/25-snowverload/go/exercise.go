package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 25.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var visitCount map[Wire]int

	c := parseComponents(instr)
	max := make([]Wire, 3)

	for i := 0; i < 3; i++ {
		visitCount = countVisits(c)
		max[i] = removeMax(visitCount)

		removeWire(c, max[i])
	}

	size1 := countComponents(c, max[0].Src)
	size2 := len(c) - countComponents(c, max[0].Src)

	return size1 * size2, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(_ string) (any, error) {
	return "", nil
}
