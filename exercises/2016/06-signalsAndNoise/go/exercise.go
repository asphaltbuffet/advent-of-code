package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 6.
type Exercise struct {
	common.BaseExercise
}

// decode recovers the message by picking, per column, the most common letter
// (mostCommon true) or the least common letter (false).
func decode(instr string, mostCommon bool) string {
	lines := strings.Fields(instr)
	if len(lines) == 0 {
		return ""
	}
	width := len(lines[0])
	counts := make([]map[byte]int, width)
	for c := range counts {
		counts[c] = map[byte]int{}
	}
	for _, line := range lines {
		for c := 0; c < width && c < len(line); c++ {
			counts[c][line[c]]++
		}
	}

	var msg strings.Builder
	for c := range width {
		var best byte
		bestN := -1
		for ch, n := range counts[c] {
			if bestN == -1 || (mostCommon && n > bestN) || (!mostCommon && n < bestN) {
				best, bestN = ch, n
			}
		}
		msg.WriteByte(best)
	}
	return msg.String()
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return decode(instr, true), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return decode(instr, false), nil
}
