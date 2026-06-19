package exercises

import (
	"math/bits"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 17.
type Exercise struct {
	common.BaseExercise
}

// parse reads the container sizes. The target volume is not in the input, so
// it is inferred: the 5-container AoC example fills 25 liters, the real 150.
func parse(instr string) ([]int, int) {
	var sizes []int
	for _, line := range strings.Fields(instr) {
		n, _ := strconv.Atoi(line)
		sizes = append(sizes, n)
	}

	target := 150
	if len(sizes) <= 5 {
		target = 25
	}

	return sizes, target
}

// sizeCounts returns, for every container count k, how many subsets of the
// containers sum to exactly target. Index k holds the number of valid subsets
// using k containers. With <= 20 containers the 2^n bitmask sweep is trivial.
func sizeCounts(sizes []int, target int) []int {
	counts := make([]int, len(sizes)+1)

	for mask := 0; mask < (1 << len(sizes)); mask++ {
		sum := 0
		for i, s := range sizes {
			if mask&(1<<i) != 0 {
				sum += s
			}
		}
		if sum == target {
			counts[bits.OnesCount(uint(mask))]++
		}
	}

	return counts
}

// One returns the number of container combinations that total the target.
func (e Exercise) One(instr string) (any, error) {
	sizes, target := parse(instr)

	total := 0
	for _, c := range sizeCounts(sizes, target) {
		total += c
	}

	return total, nil
}

// Two returns how many combinations use the fewest containers possible.
func (e Exercise) Two(instr string) (any, error) {
	sizes, target := parse(instr)

	for _, c := range sizeCounts(sizes, target) {
		if c > 0 {
			return c, nil // first non-zero size bucket is the minimum
		}
	}

	return 0, nil
}
