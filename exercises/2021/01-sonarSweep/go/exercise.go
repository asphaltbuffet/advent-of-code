package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func parse(in string) ([]int, error) {
	var lines []int

	for _, line := range strings.Split(in, "\n") {
		if line == "" {
			continue
		}

		l := strings.TrimSpace(line)

		n, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		lines = append(lines, n)
	}

	return lines, nil
}

type Exercise struct {
	common.BaseExercise
}

func (c Exercise) One(instr string) (any, error) {
	data, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return increasingCount(data), nil
}

func (c Exercise) Two(instr string) (any, error) {
	data, err := parse(instr)
	if err != nil {
		return nil, err
	}

	w := []int{}

	// calculate a new array of windowed values
	for i := 2; i < len(data); i++ {
		w = append(w, data[i]+data[i-1]+data[i-2])
	}

	return increasingCount(w), nil
}

func increasingCount(data []int) int {
	count := 0
	prev := data[0]

	for i := 1; i < len(data); i++ {
		curr := data[i]
		if curr > prev {
			count++
		}

		prev = curr
	}

	return count
}
