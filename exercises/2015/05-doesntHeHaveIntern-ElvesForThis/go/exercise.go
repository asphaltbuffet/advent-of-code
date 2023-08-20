package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 5.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	count := 0

	for _, line := range strings.Split(instr, "\n") {
		h := hasVowels(line)
		d := hasDoubles(line)
		b := hasBad(line)

		if h && d && !b {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	count := 0

	for _, line := range strings.Split(instr, "\n") {
		p := hasPair(line)
		s := hasSeparated(line)

		if p && s {
			count++
		}
	}

	return count, nil
}

func hasVowels(s string) bool {
	v := 0

	for _, c := range []string{"a", "e", "i", "o", "u"} {
		v += strings.Count(s, c)
	}

	return v >= 3
}

func hasDoubles(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}

	return false
}

func hasBad(s string) bool {
	for _, c := range []string{"ab", "cd", "pq", "xy"} {
		if strings.Contains(s, c) {
			return true
		}
	}

	return false
}

func hasPair(s string) bool {
	for i := 0; i < len(s)-2; i++ {
		if strings.Count(s, s[i:i+2]) > 1 {
			return true
		}
	}

	return false
}

func hasSeparated(s string) bool {
	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+2] {
			return true
		}
	}

	return false
}
