package exercises

import (
	"slices"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 7.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	lines := strings.Split(instr, "\n")
	hands := make([]Hand, 0, len(lines))

	for _, line := range lines {
		hands = append(hands, parseHand(line, false))
	}

	slices.SortFunc(hands, handSort)

	var total int
	for i, hand := range hands {
		total += hand.Bid * (i + 1)

		// fmt.Printf("pHand[i]: %q v=%d, s=%d, win=%d\n", hand.Cards, hand.Value, hand.Strength, total)
	}

	return total, nil
}

// Two returns the answer to the second part of the exercise.
// not: 246948152
func (e Exercise) Two(instr string) (any, error) {
	lines := strings.Split(instr, "\n")
	hands := make([]Hand, 0, len(lines))

	for _, line := range lines {
		hands = append(hands, parseHand(line, true))
	}

	slices.SortFunc(hands, handSort)

	var total int
	for i, hand := range hands {
		total += hand.Bid * (i + 1)

		// fmt.Printf("pHand[i]: %q v=%d, s=%d, win=%d\n", hand.Cards, hand.Value, hand.Strength, total)
	}

	return total, nil
}
