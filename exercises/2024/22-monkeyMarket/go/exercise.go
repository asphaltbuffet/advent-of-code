package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 22.
type Exercise struct {
	common.BaseExercise
}

const pruneMask = 0xFFFFFF // modulo 16777216 == 2^24

// next advances the secret number one step.
func next(s int) int {
	s ^= s << 6
	s &= pruneMask
	s ^= s >> 5
	s &= pruneMask
	s ^= s << 11
	s &= pruneMask
	return s
}

func parseSecrets(instr string) []int {
	var out []int
	for _, f := range strings.Fields(instr) {
		n, _ := strconv.Atoi(f)
		out = append(out, n)
	}
	return out
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	sum := 0
	for _, s := range parseSecrets(instr) {
		for i := 0; i < 2000; i++ {
			s = next(s)
		}
		sum += s
	}
	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	// totals[key] accumulates bananas across all buyers for a 4-change window.
	// key packs four deltas (each shifted +9 into [0,18]) in base 19.
	totals := make(map[int]int)

	for _, s := range parseSecrets(instr) {
		seen := make(map[int]bool)
		prevPrice := s % 10
		key := 0
		for i := 0; i < 2000; i++ {
			s = next(s)
			price := s % 10
			delta := price - prevPrice
			prevPrice = price
			key = (key*19 + (delta + 9)) % (19 * 19 * 19 * 19)
			// After 4 deltas (i >= 3) the key is a valid 4-window.
			if i >= 3 && !seen[key] {
				seen[key] = true
				totals[key] += price
			}
		}
	}

	best := 0
	for _, v := range totals {
		if v > best {
			best = v
		}
	}
	return best, nil
}
