package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 5.
type Exercise struct {
	common.BaseExercise
}

func parseOffsets(instr string) []int {
	fields := strings.Fields(instr)
	offsets := make([]int, len(fields))
	for i, f := range fields {
		offsets[i], _ = strconv.Atoi(f)
	}
	return offsets
}

// steps runs the jump program, applying update to each offset after it is used,
// and returns the number of steps taken to exit the list.
func steps(offsets []int, update func(int) int) int {
	jumps := make([]int, len(offsets))
	copy(jumps, offsets)

	count, ip := 0, 0
	for ip >= 0 && ip < len(jumps) {
		off := jumps[ip]
		jumps[ip] = update(off)
		ip += off
		count++
	}
	return count
}

// One counts steps to escape when each used offset increases by 1.
func (e Exercise) One(instr string) (any, error) {
	return steps(parseOffsets(instr), func(off int) int { return off + 1 }), nil
}

// Two counts steps to escape when an offset of 3 or more decreases by 1 after
// use, and any smaller offset increases by 1.
func (e Exercise) Two(instr string) (any, error) {
	return steps(parseOffsets(instr), func(off int) int {
		if off >= 3 {
			return off - 1
		}
		return off + 1
	}), nil
}
