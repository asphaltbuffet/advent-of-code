package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 2.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	reports, err := parse(instr)
	if err != nil {
		return nil, err
	}

	var count int
	for _, r := range reports {
		if ok := r.isSafe(); ok {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	reports, err := parse(instr)
	if err != nil {
		return nil, err
	}

	var count int
	for _, r := range reports {
		r := r
		if ok := r.isSafe(); ok {
			count++
		} else {
			for i := 0; i < len(r.values); i++ {
				tmp := Report{values: removeAt(r.values, i)}

				if altOk := tmp.isSafe(); altOk {
					count++
					break
				}
			}
		}
	}

	return count, nil
}
