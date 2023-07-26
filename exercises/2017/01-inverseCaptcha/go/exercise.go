package exercises

import (
	"container/ring"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	r := parse(instr)

	sum := 0

	for i := 0; i < r.Len(); i++ {
		v, _ := r.Value.(int)
		n, _ := r.Next().Value.(int)

		// fmt.Printf("%d & %d\n", v, n)

		if v == n {
			sum += v
		}

		r = r.Next()
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	return nil, nil
}

func parse(instr string) *ring.Ring {
	r := ring.New(len(instr))

	for _, c := range instr {
		r.Value, _ = strconv.Atoi(string(c))
		r = r.Next()
	}

	return r
}
