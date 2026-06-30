package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 20.
type Exercise struct {
	common.BaseExercise
}

// target reads the present count from the input.
func target(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One: elf n gives 10*n presents to every multiple of n, so house h receives
// 10 * (sum of divisors of h). A sieve adds each elf's gift to its multiples;
// return the lowest house meeting the target. The bound target/10 suffices
// because house h alone (elf h) already delivers 10*h there.
func (e Exercise) One(instr string) (any, error) {
	t := target(instr)
	limit := t/10 + 1
	houses := make([]int, limit+1)

	for n := 1; n <= limit; n++ {
		for h := n; h <= limit; h += n {
			houses[h] += 10 * n
		}
	}

	for h := 1; h <= limit; h++ {
		if houses[h] >= t {
			return h, nil
		}
	}
	return -1, nil
}

// Two: each elf gives 11*n presents but visits only its first 50 houses.
func (e Exercise) Two(instr string) (any, error) {
	t := target(instr)
	limit := t/10 + 1
	houses := make([]int, limit+1)

	for n := 1; n <= limit; n++ {
		for c, h := 1, n; c <= 50 && h <= limit; c, h = c+1, h+n {
			houses[h] += 11 * n
		}
	}

	for h := 1; h <= limit; h++ {
		if houses[h] >= t {
			return h, nil
		}
	}
	return -1, nil
}
