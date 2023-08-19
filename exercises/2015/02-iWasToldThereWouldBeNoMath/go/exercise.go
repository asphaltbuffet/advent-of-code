package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 2.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	dims := parse(instr)

	var total int

	for _, dim := range dims {
		a := dim[0] * dim[1]
		b := dim[1] * dim[2]
		c := dim[2] * dim[0]

		min := min(a, b, c)

		total += 2*(a+b+c) + min
	}

	return total, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	return nil, nil
}

func parse(instr string) [][3]int {
	var dims [][3]int

	for _, line := range strings.Split(instr, "\n") {
		var (
			a int
			b int
			c int
		)

		_, err := fmt.Sscanf(line, "%dx%dx%d", &a, &b, &c)
		if err != nil {
			panic(err)
		}

		dims = append(dims, [3]int{a, b, c})
	}

	return dims
}
