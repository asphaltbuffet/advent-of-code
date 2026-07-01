package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 3.
type Exercise struct {
	common.BaseExercise
}

func parse(in string) []string {
	var lines []string

	for _, line := range strings.Split(in, "\n") {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		lines = append(lines, l)
	}

	return lines
}

// onesMajority reports whether column c holds at least as many 1s as 0s across
// the given numbers. The "at least" (ties count as 1) matches both the gamma
// rule and the oxygen-rating tie-break.
func onesMajority(nums []string, c int) bool {
	ones := 0
	for _, n := range nums {
		if n[c] == '1' {
			ones++
		}
	}

	return ones*2 >= len(nums)
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) == 0 {
		return nil, fmt.Errorf("no diagnostic numbers in input")
	}

	width := len(nums[0])

	gamma := 0
	for c := 0; c < width; c++ {
		gamma <<= 1
		if onesMajority(nums, c) {
			gamma |= 1
		}
	}

	// epsilon is the bitwise complement of gamma over the fixed width.
	epsilon := gamma ^ ((1 << width) - 1)

	return fmt.Sprintf("%d", gamma*epsilon), nil
}

// filter narrows the list one column at a time, keeping the numbers whose bit in
// that column matches the wanted value; keepOnes selects which value is wanted
// (most-common for oxygen, least-common for CO2). It stops when one number remains.
func filter(nums []string, keepMostCommon bool) string {
	width := len(nums[0])

	for c := 0; c < width && len(nums) > 1; c++ {
		majorityOne := onesMajority(nums, c)

		var want byte = '0'
		if majorityOne == keepMostCommon {
			want = '1'
		}

		var kept []string
		for _, n := range nums {
			if n[c] == want {
				kept = append(kept, n)
			}
		}

		nums = kept
	}

	return nums[0]
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) == 0 {
		return nil, fmt.Errorf("no diagnostic numbers in input")
	}

	oxygen, err := strconv.ParseInt(filter(nums, true), 2, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing oxygen rating: %w", err)
	}

	co2, err := strconv.ParseInt(filter(nums, false), 2, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing CO2 rating: %w", err)
	}

	return fmt.Sprintf("%d", oxygen*co2), nil
}
