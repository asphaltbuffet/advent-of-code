package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 4.
type Exercise struct {
	common.BaseExercise
}

func parseRange(instr string) (int, int, error) {
	parts := strings.Split(strings.TrimSpace(instr), "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range: %q", instr)
	}
	lo, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	hi, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return lo, hi, nil
}

// toNum converts a 6-digit slice (each 0–9) to an integer.
func toNum(d [6]int) int {
	return d[0]*100000 + d[1]*10000 + d[2]*1000 + d[3]*100 + d[4]*10 + d[5]
}

// hasAdjacentPair returns true if any two adjacent digits are equal.
func hasAdjacentPair(d [6]int) bool {
	for i := range 5 {
		if d[i] == d[i+1] {
			return true
		}
	}
	return false
}

// countPasswords enumerates all non-decreasing 6-digit sequences within [lo, hi]
// and applies the extraFilter predicate. This is C(14,6)=3003 sequences total.
//
//nolint:gocognit // digit-rule validation branches are inherently nested
func countPasswords(lo, hi int, extraFilter func([6]int) bool) int {
	count := 0
	var d [6]int
	for d[0] = 1; d[0] <= 9; d[0]++ {
		for d[1] = d[0]; d[1] <= 9; d[1]++ {
			for d[2] = d[1]; d[2] <= 9; d[2]++ {
				for d[3] = d[2]; d[3] <= 9; d[3]++ {
					for d[4] = d[3]; d[4] <= 9; d[4]++ {
						for d[5] = d[4]; d[5] <= 9; d[5]++ {
							n := toNum(d)
							if n < lo || n > hi {
								continue
							}
							if !hasAdjacentPair(d) {
								continue
							}
							if extraFilter != nil && !extraFilter(d) {
								continue
							}
							count++
						}
					}
				}
			}
		}
	}
	return count
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	lo, hi, err := parseRange(instr)
	if err != nil {
		return nil, err
	}
	return countPasswords(lo, hi, nil), nil
}

// hasExactPair returns true if at least one run of identical adjacent digits
// has length exactly 2. A run of 3+ doesn't satisfy this on its own.
func hasExactPair(d [6]int) bool {
	i := 0
	for i < 6 {
		j := i + 1
		for j < 6 && d[j] == d[i] {
			j++
		}
		if j-i == 2 {
			return true
		}
		i = j
	}
	return false
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	lo, hi, err := parseRange(instr)
	if err != nil {
		return nil, err
	}
	return countPasswords(lo, hi, hasExactPair), nil
}
