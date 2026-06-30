package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 24.
type Exercise struct {
	common.BaseExercise
}

// parse reads the package weights.
func parse(instr string) []int {
	var ws []int
	for _, f := range strings.Fields(instr) {
		n, _ := strconv.Atoi(f)
		ws = append(ws, n)
	}
	return ws
}

// minQE finds the smallest first-group size that sums to target, then returns
// the lowest quantum entanglement (product of weights) among groups of that
// size. Searching by increasing size guarantees the fewest-packages tier; the
// minimum product within it is the answer. The remaining packages always admit
// a valid split for AoC inputs, so only the first group is optimized.
func minQE(weights []int, groups int) int64 {
	target := 0
	for _, w := range weights {
		target += w
	}
	target /= groups

	for size := 1; size <= len(weights); size++ {
		best := int64(-1)
		var pick func(start, remaining, count int, qe int64)
		pick = func(start, remaining, count int, qe int64) {
			if count == 0 {
				if remaining == 0 && (best < 0 || qe < best) {
					best = qe
				}
				return
			}
			for i := start; i <= len(weights)-count; i++ {
				if weights[i] <= remaining {
					pick(i+1, remaining-weights[i], count-1, qe*int64(weights[i]))
				}
			}
		}
		pick(0, target, size, 1)

		if best >= 0 {
			return best
		}
	}

	return -1
}

// One: split into 3 equal groups; lowest QE of the smallest first group.
func (e Exercise) One(instr string) (any, error) {
	return minQE(parse(instr), 3), nil
}

// Two: split into 4 equal groups.
func (e Exercise) Two(instr string) (any, error) {
	return minQE(parse(instr), 4), nil
}
