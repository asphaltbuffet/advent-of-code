package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 11.
type Exercise struct {
	common.BaseExercise
}

// digits returns the number of decimal digits in n (n >= 0).
func digits(n int) int {
	d := 1
	for n >= 10 {
		n /= 10
		d++
	}
	return d
}

// blink applies one transformation to a stone, returning one or two stones.
// When only one results, second is -1.
func blink(n int) (int, int) {
	if n == 0 {
		return 1, -1
	}
	if d := digits(n); d%2 == 0 {
		half := 1
		for i := 0; i < d/2; i++ {
			half *= 10
		}
		return n / half, n % half
	}
	return n * 2024, -1
}

// countAfter returns the number of stones after the given number of blinks,
// tracking a multiset of stone value -> count (order is irrelevant).
func countAfter(instr string, blinks int) int {
	stones := make(map[int]int)
	for _, f := range strings.Fields(instr) {
		n, _ := strconv.Atoi(f)
		stones[n]++
	}

	for i := 0; i < blinks; i++ {
		next := make(map[int]int, len(stones)*2)
		for n, cnt := range stones {
			a, b := blink(n)
			next[a] += cnt
			if b != -1 {
				next[b] += cnt
			}
		}
		stones = next
	}

	total := 0
	for _, cnt := range stones {
		total += cnt
	}
	return total
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return countAfter(instr, 25), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return countAfter(instr, 75), nil
}
