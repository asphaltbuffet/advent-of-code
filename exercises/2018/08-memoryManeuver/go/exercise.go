package exercises

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 8.
type Exercise struct {
	common.BaseExercise
}

var intRe = regexp.MustCompile(`\d+`)

// parse reads the flat list of numbers describing the tree.
func parse(instr string) []int {
	toks := intRe.FindAllString(instr, -1)
	nums := make([]int, len(toks))
	for i, t := range toks {
		nums[i], _ = strconv.Atoi(t)
	}
	return nums
}

// metaSum reads one node starting at nums[i], returning the summed metadata of
// that node and its whole subtree along with the index just past the node.
func metaSum(nums []int, i int) (sum, next int) {
	children, meta := nums[i], nums[i+1]
	i += 2

	for c := 0; c < children; c++ {
		var s int
		s, i = metaSum(nums, i)
		sum += s
	}

	for m := 0; m < meta; m++ {
		sum += nums[i]
		i++
	}

	return sum, i
}

// nodeValue reads one node starting at nums[i], returning its value and the
// index just past the node. A leaf's value is the sum of its metadata; a node
// with children sums the values of the children its metadata entries reference
// (1-based, out-of-range references skipped, repeats counted).
func nodeValue(nums []int, i int) (value, next int) {
	children, meta := nums[i], nums[i+1]
	i += 2

	childValues := make([]int, children)
	for c := 0; c < children; c++ {
		childValues[c], i = nodeValue(nums, i)
	}

	for m := 0; m < meta; m++ {
		ref := nums[i]
		i++

		if children == 0 {
			value += ref
		} else if ref >= 1 && ref <= children {
			value += childValues[ref-1]
		}
	}

	return value, i
}

// One returns the answer to the first part of the exercise.
// Answer: 45865
func (e Exercise) One(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) < 2 {
		return nil, fmt.Errorf("input too short")
	}

	sum, _ := metaSum(nums, 0)

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 22608
func (e Exercise) Two(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) < 2 {
		return nil, fmt.Errorf("input too short")
	}

	value, _ := nodeValue(nums, 0)

	return value, nil
}
