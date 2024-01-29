package exercises

import (
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 20.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	defer common.Close()

	m, err := loadInput(input)
	if err != nil {
		return nil, err
	}

	var low, high int

	for i := 0; i < 1000; i++ {
		msgQueue.Send("button", "broadcaster", Low)
		l, h := m.ProcessQueue()

		low += l
		high += h
	}

	return (low * high), nil
}

// Two returns the answer to the second part of the exercise.
// not: 2_237_580_152_160_073 (too high)
func (e Exercise) Two(input string) (any, error) {
	defer common.Close()

	cfg, err := loadInput(input)
	if err != nil {
		return nil, err
	}

	src, err := cfg.GetFeed("rx")
	if err != nil {
		return nil, err
	}

	cycles := cfg.ProcessUntilHigh(src)

	return cycles, nil
}
