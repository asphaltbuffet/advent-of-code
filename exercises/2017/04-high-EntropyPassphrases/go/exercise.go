package exercises

import (
	"slices"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 4.
type Exercise struct {
	common.BaseExercise
}

// countValid returns how many passphrases are valid, where key normalises each
// word so that repeated keys mark an invalid passphrase.
func countValid(instr string, key func(string) string) int {
	valid := 0
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}
		seen := make(map[string]bool, len(words))
		ok := true
		for _, w := range words {
			k := key(w)
			if seen[k] {
				ok = false
				break
			}
			seen[k] = true
		}
		if ok {
			valid++
		}
	}
	return valid
}

// One counts passphrases with no duplicate words.
func (e Exercise) One(instr string) (any, error) {
	return countValid(instr, func(w string) string { return w }), nil
}

// Two counts passphrases with no two words that are anagrams of each other.
// Sorting a word's letters gives anagrams a shared key.
func (e Exercise) Two(instr string) (any, error) {
	return countValid(instr, func(w string) string {
		b := []byte(w)
		slices.Sort(b)
		return string(b)
	}), nil
}
