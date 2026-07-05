package exercises

import (
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 3.
type Exercise struct {
	common.BaseExercise
}

var numRe = regexp.MustCompile(`\d+`)

// parseNums returns every integer in the input, in order.
func parseNums(instr string) []int {
	matches := numRe.FindAllString(instr, -1)
	nums := make([]int, len(matches))
	for i, s := range matches {
		nums[i], _ = strconv.Atoi(s)
	}
	return nums
}

// valid reports whether sides a, b, c form a triangle.
func valid(a, b, c int) bool {
	return a+b > c && a+c > b && b+c > a
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	nums := parseNums(instr)
	count := 0
	for i := 0; i+2 < len(nums); i += 3 {
		if valid(nums[i], nums[i+1], nums[i+2]) {
			count++
		}
	}
	return count, nil
}

// Two returns the answer to the second part of the exercise. Triangles are read
// down each column across blocks of three rows (3 numbers per row).
func (e Exercise) Two(instr string) (any, error) {
	nums := parseNums(instr)
	count := 0
	// Each block is 9 numbers (3 rows x 3 cols); columns are (0,3,6) etc.
	for i := 0; i+8 < len(nums); i += 9 {
		for col := range 3 {
			if valid(nums[i+col], nums[i+col+3], nums[i+col+6]) {
				count++
			}
		}
	}
	return count, nil
}
