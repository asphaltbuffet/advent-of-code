package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 7.
type Exercise struct {
	common.BaseExercise
}

// equation is a target test value and its list of operands.
type equation struct {
	target  int
	numbers []int
}

func parse(instr string) []equation {
	var eqs []equation
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		colon := strings.IndexByte(line, ':')
		target, _ := strconv.Atoi(line[:colon])
		var nums []int
		for _, f := range strings.Fields(line[colon+1:]) {
			n, _ := strconv.Atoi(f)
			nums = append(nums, n)
		}
		eqs = append(eqs, equation{target, nums})
	}
	return eqs
}

// concat returns a || b, e.g. 12 || 345 = 12345.
func concat(a, b int) int {
	pow := 1
	for p := b; p > 0; p /= 10 {
		pow *= 10
	}
	return a*pow + b
}

// solvable reports whether some operator placement makes acc..nums equal target.
// Branches are pruned once acc exceeds target (operators never decrease).
func solvable(target, acc int, nums []int, concatOp bool) bool {
	if len(nums) == 0 {
		return acc == target
	}
	if acc > target {
		return false
	}
	n := nums[0]
	rest := nums[1:]
	if solvable(target, acc+n, rest, concatOp) {
		return true
	}
	if solvable(target, acc*n, rest, concatOp) {
		return true
	}
	if concatOp && solvable(target, concat(acc, n), rest, concatOp) {
		return true
	}
	return false
}

func calibrate(instr string, concatOp bool) int {
	sum := 0
	for _, eq := range parse(instr) {
		if solvable(eq.target, eq.numbers[0], eq.numbers[1:], concatOp) {
			sum += eq.target
		}
	}
	return sum
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return calibrate(instr, false), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return calibrate(instr, true), nil
}
