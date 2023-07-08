package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

var priorityValue = map[byte]int{
	'a': 1, 'b': 2, 'c': 3, 'd': 4, 'e': 5, 'f': 6, 'g': 7, 'h': 8, 'i': 9, 'j': 10,
	'k': 11, 'l': 12, 'm': 13, 'n': 14, 'o': 15, 'p': 16, 'q': 17, 'r': 18, 's': 19, 't': 20,
	'u': 21, 'v': 22, 'w': 23, 'x': 24, 'y': 25, 'z': 26,

	'A': 27, 'B': 28, 'C': 29, 'D': 30, 'E': 31, 'F': 32, 'G': 33, 'H': 34, 'I': 35, 'J': 36,
	'K': 37, 'L': 38, 'M': 39, 'N': 40, 'O': 41, 'P': 42, 'Q': 43, 'R': 44, 'S': 45, 'T': 46,
	'U': 47, 'V': 48, 'W': 49, 'X': 50, 'Y': 51, 'Z': 52,
}

// Exercise for Advent of Code 2022 day 3
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (c Exercise) One(instr string) (any, error) {
	score := 0

	for _, line := range strings.Split(instr, "\n") {
		score += scoreMispacked(line)
	}

	return score, nil
}

// Two returns the answer to the second part of the exercise.
func (c Exercise) Two(instr string) (any, error) {
	var data []string

	data = append(data, strings.Split(instr, "\n")...)

	score := 0

	for i := 0; i < len(data); i += 3 {
		score += scoreBadges(data[i], data[i+1], data[i+2])
	}

	return score, nil
}

func scoreMispacked(line string) int {
	compartmentOne := map[byte]bool{}

	for i := 0; i < len(line)/2; i++ {
		compartmentOne[line[i]] = true
	}

	priority := 0

	for i := len(line) / 2; i < len(line); i++ {
		if _, ok := compartmentOne[line[i]]; ok {
			// fmt.Printf("found match: %q\n", line[i])
			priority += priorityValue[line[i]]

			// only count the first match
			delete(compartmentOne, line[i])
		}
	}

	return priority
}

func scoreBadges(a, b, c string) int {
	sharedItems := aoc.Unique([]byte(a))

	sharedItems = aoc.Filter(sharedItems, func(item byte) bool {
		return strings.Contains(b, string(item)) && strings.Contains(c, string(item))
	})

	// fmt.Printf("DEBUG: shared items: %s\n", sharedItems)

	return priorityValue[sharedItems[0]]
}
