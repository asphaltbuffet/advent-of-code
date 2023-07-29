package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 20.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// wrong: 9476 (too high)
// answer: 4066
func (e Exercise) One(instr string) (any, error) {
	f, err := parse(instr)
	if err != nil {
		return nil, err
	}

	f.decrypted.Do(func(v digit) {
		println(v.value)
	})

	err = f.decrypt()
	if err != nil {
		return nil, err
	}

	// fmt.Println(f.decryptedToString())

	c1, c2, c3 := f.getCoordinates()

	fmt.Printf("coordinates = %d %d %d\n", c1, c2, c3)

	return c1 + c2 + c3, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	return nil, nil
}
