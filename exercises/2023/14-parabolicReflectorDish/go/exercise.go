package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 14.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// not: 166856 (too high)
func (e Exercise) One(input string) (any, error) {
	data, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	data.tiltNorth()

	return data.calcLoad(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	const cycles = 1_000_000_000

	memos := make(map[string]int)

	data, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	data.transpose()

	for i := 0; i < cycles; i++ {
		data.spin()

		hash := data.hash()
		if v, ok := memos[hash]; ok {
			// fmt.Printf("hash: %s seen before at i=%d and again at i=%d\n", hash, v, i)
			// fmt.Printf("cycle length: %d\n", i-v)

			// fmt.Printf("skipping %d cycles\n", (cycles - i) / (i - v))
			i += ((cycles - i) / (i - v)) * (i - v)

			// reset memos
			memos = map[string]int{}
		}

		memos[hash] = i
	}

	data.transpose()

	return data.calcLoad(), nil
}
