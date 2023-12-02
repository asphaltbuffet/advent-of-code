package exercises

import (
	"fmt"
	"strings"

	"github.com/caarlos0/log"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 1.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var sum int

	for _, line := range strings.Split(instr, "\n") {
		sum += getCalibrationValue(line)
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var sum int
	for _, line := range strings.Split(instr, "\n") {
		sum += part2(line)
	}

	return sum, nil
}

func getCalibrationValue(line string) int {
	nums := []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	left, leftOk := firstNumber(line, nums)
	right, rightOk := lastNumber(line, nums)

	if !rightOk || !leftOk {
		log.Warnf("could not find left and right numbers in %s", line)
		return 0
	}

	return (left * 10) + right
}

func part2(line string) int {
	nums := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	left, leftOk := firstNumber(line, nums)
	right, rightOk := lastNumber(line, nums)

	if !rightOk || !leftOk {
		log.Warnf("could not find left and right numbers in %s", line)
		return 0
	}

	return (left * 10) + right
}

func firstNumber(s string, numbers []string) (int, bool) {
	var (
		index  = len(s)
		number string
		found  bool
	)

	for _, num := range numbers {
		idx := strings.Index(s, num)
		if idx >= 0 && idx < index {
			index = idx
			number = num
			found = true
		}
	}

	if !found {
		return -1, false
	}

	n, err := stringToNumber(number)

	return n, err == nil
}

func lastNumber(s string, numbers []string) (int, bool) {
	var (
		index  = -1
		number string
		found  bool
	)

	for _, num := range numbers {
		idx := strings.LastIndex(s, num)
		if idx > index {
			index = idx
			number = num
			found = true
		}
	}

	if !found {
		return -1, false
	}

	n, err := stringToNumber(number)

	return n, err == nil
}

func stringToNumber(s string) (int, error) {
	nums := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	for i, num := range nums {
		if num == s {
			return i%9 + 1, nil
		}
	}

	return -1, fmt.Errorf("could not convert %s to number", s)
}
