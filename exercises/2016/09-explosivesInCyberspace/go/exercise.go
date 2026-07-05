package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 9.
type Exercise struct {
	common.BaseExercise
}

// decompressedLen returns the length of the decompressed string. When recursive
// is true, markers inside a repeated section are expanded as well (part two);
// otherwise repeated data is counted literally (part one).
func decompressedLen(s string, recursive bool) int {
	total := 0
	i := 0
	for i < len(s) {
		if s[i] != '(' {
			total++
			i++
			continue
		}
		// Parse marker (AxB).
		closeIdx := strings.IndexByte(s[i:], ')') + i
		marker := s[i+1 : closeIdx]
		x := strings.SplitN(marker, "x", 2)
		take, _ := strconv.Atoi(x[0])
		reps, _ := strconv.Atoi(x[1])

		section := s[closeIdx+1 : closeIdx+1+take]
		if recursive {
			total += reps * decompressedLen(section, true)
		} else {
			total += reps * take
		}
		i = closeIdx + 1 + take
	}
	return total
}

// clean strips whitespace, which the format ignores.
func clean(instr string) string {
	return strings.Join(strings.Fields(instr), "")
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return decompressedLen(clean(instr), false), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return decompressedLen(clean(instr), true), nil
}
