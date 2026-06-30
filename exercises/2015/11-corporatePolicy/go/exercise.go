package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 11.
type Exercise struct {
	common.BaseExercise
}

// increment advances the password by one, treating it as a base-26 odometer
// over 'a'..'z' with carry from the rightmost character.
func increment(p []byte) {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == 'z' {
			p[i] = 'a'
			continue
		}
		p[i]++
		break
	}
}

// valid reports whether p satisfies all three corporate policy rules.
func valid(p []byte) bool {
	// Rule 2: no i, o, or l.
	for _, c := range p {
		if c == 'i' || c == 'o' || c == 'l' {
			return false
		}
	}

	// Rule 1: one increasing straight of at least three letters.
	straight := false
	for i := 0; i+2 < len(p); i++ {
		if p[i+1] == p[i]+1 && p[i+2] == p[i]+2 {
			straight = true
			break
		}
	}
	if !straight {
		return false
	}

	// Rule 3: at least two different, non-overlapping pairs.
	pairs := map[byte]bool{}
	for i := 0; i+1 < len(p); i++ {
		if p[i] == p[i+1] {
			pairs[p[i]] = true
			i++ // non-overlapping
		}
	}
	return len(pairs) >= 2
}

// next returns the next valid password strictly after p.
func next(p []byte) []byte {
	increment(p)
	for !valid(p) {
		increment(p)
	}
	return p
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	p := []byte(strings.TrimSpace(instr))
	return string(next(p)), nil
}

// Two returns the answer to the second part of the exercise: the next valid
// password after part one's result.
func (e Exercise) Two(instr string) (any, error) {
	p := []byte(strings.TrimSpace(instr))
	return string(next(next(p))), nil
}
