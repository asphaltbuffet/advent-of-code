package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 19.
type Exercise struct {
	common.BaseExercise
}

func parseN(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One returns the answer to the first part of the exercise. This is the
// Josephus problem (k=2): with n = 2^m + l, the winner is 2l + 1.
func (e Exercise) One(instr string) (any, error) {
	n := parseN(instr)
	pow := 1
	for pow*2 <= n {
		pow *= 2
	}
	l := n - pow
	return 2*l + 1, nil
}

// Two returns the answer to the second part of the exercise (steal from across
// the circle). With p = the largest power of 3 <= n:
//   - n == p          -> n
//   - n <= 2p         -> n - p
//   - otherwise       -> 2n - 3p
func (e Exercise) Two(instr string) (any, error) {
	n := parseN(instr)
	p := 1
	for p*3 <= n {
		p *= 3
	}
	switch {
	case n == p:
		return n, nil
	case n-p <= p:
		return n - p, nil
	default:
		return 2*n - 3*p, nil
	}
}
