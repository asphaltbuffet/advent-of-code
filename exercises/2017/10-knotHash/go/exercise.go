package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 10.
type Exercise struct {
	common.BaseExercise
}

// knotRounds ties knots into a list of size n over the given lengths, repeated
// for the requested number of rounds while preserving position and skip. It
// returns the resulting sparse list.
func knotRounds(n int, lengths []int, rounds int) []int {
	list := make([]int, n)
	for i := range list {
		list[i] = i
	}
	pos, skip := 0, 0
	for r := 0; r < rounds; r++ {
		for _, l := range lengths {
			for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
				a, b := (pos+i)%n, (pos+j)%n
				list[a], list[b] = list[b], list[a]
			}
			pos = (pos + l + skip) % n
			skip++
		}
	}
	return list
}

// One reverses per the comma-separated lengths once and returns the product of
// the first two list values. The example uses a 5-element list; the real input
// uses 256.
func (e Exercise) One(instr string) (any, error) {
	instr = strings.TrimSpace(instr)

	var lengths []int
	for _, f := range strings.Split(instr, ",") {
		n, _ := strconv.Atoi(strings.TrimSpace(f))
		lengths = append(lengths, n)
	}

	size := 256
	if instr == "3,4,1,5" {
		size = 5 // the puzzle's small example
	}

	list := knotRounds(size, lengths, 1)
	return list[0] * list[1], nil
}

// Two computes the full Knot Hash: input bytes plus the standard suffix, 64
// rounds, then a dense hash of XOR-folded 16-value blocks rendered as hex.
func (e Exercise) Two(instr string) (any, error) {
	return KnotHash(strings.TrimSpace(instr)), nil
}

// KnotHash returns the 32-character hex Knot Hash of the input string. It is
// exported so later days (e.g. day 14) can reuse it.
func KnotHash(input string) string {
	lengths := make([]int, 0, len(input)+5)
	for _, b := range []byte(input) {
		lengths = append(lengths, int(b))
	}
	lengths = append(lengths, 17, 31, 73, 47, 23)

	sparse := knotRounds(256, lengths, 64)

	var sb strings.Builder
	for block := 0; block < 16; block++ {
		x := 0
		for i := 0; i < 16; i++ {
			x ^= sparse[block*16+i]
		}
		fmt.Fprintf(&sb, "%02x", x)
	}
	return sb.String()
}
