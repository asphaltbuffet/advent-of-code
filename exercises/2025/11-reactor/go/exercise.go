package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 11.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
		}
	}()

	memo := make(map[string]int)
	dd := GetDevices(instr)

	return dd.Trace(memo, "you", "out"), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
		}
	}()

	var path []string
	memo := make(map[string]int)
	dd := GetDevices(instr)

	if dd.PathExists("dac", "fft") {
		path = []string{"dac", "fft", "out"}
	} else {
		path = []string{"fft", "dac", "out"}
	}

	return dd.Trace(memo, "svr", path...), nil
}
