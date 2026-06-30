package exercises

import (
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 25.
type Exercise struct {
	common.BaseExercise
}

const (
	first   = 20151125
	mult    = 252533
	modulus = 33554393
)

// parseRowCol pulls the row and column from the puzzle's prose by extracting
// the two integers it contains.
func parseRowCol(instr string) (row, col int) {
	var nums []int
	cur := ""
	for _, r := range instr {
		if r >= '0' && r <= '9' {
			cur += string(r)
		} else if cur != "" {
			n, _ := strconv.Atoi(cur)
			nums = append(nums, n)
			cur = ""
		}
	}
	if cur != "" {
		n, _ := strconv.Atoi(cur)
		nums = append(nums, n)
	}
	return nums[0], nums[1]
}

// modPow computes base^exp mod m by fast exponentiation.
func modPow(base, exp, m int64) int64 {
	result := int64(1)
	base %= m
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % m
		}
		base = base * base % m
		exp >>= 1
	}
	return result
}

// One returns the code at the requested grid cell. The cell (r, c) is the n-th
// code filled in diagonal order, where n = T(r+c-2) + c, and the n-th code is
// first * mult^(n-1) mod modulus.
func (e Exercise) One(instr string) (any, error) {
	row, col := parseRowCol(instr)
	diag := row + col - 2
	n := int64(diag*(diag+1)/2 + col)
	return first * modPow(mult, n-1, modulus) % modulus, nil
}

// Two has no puzzle on day 25 — the final star is earned by completing the rest.
func (e Exercise) Two(instr string) (any, error) {
	_ = instr
	return "Merry Christmas!", nil
}
