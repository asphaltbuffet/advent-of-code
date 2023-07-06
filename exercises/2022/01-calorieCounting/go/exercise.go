package exercises

import (
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func parse(data string) ([]int, error) {
	sum := 0
	calories := sort.IntSlice{}

	for _, l := range strings.Split(data, "\n") {
		line := strings.TrimSpace(l)

		if line == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}

			sum += n
		}
	}

	sort.Slice(calories,
		func(i, j int) bool {
			return calories[i] > calories[j]
		})

	return calories, nil
}

// Exercise for Advent of Code 2022 day 1
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (c Exercise) One(instr string) (any, error) {
	cal, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return cal[0], nil
}

// Two returns the answer to the second part of the exercise.
func (c Exercise) Two(instr string) (any, error) {
	cal, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return cal[0] + cal[1] + cal[2], nil
}
