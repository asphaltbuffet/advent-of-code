// Package aoc22_02 contains the solution for day 2 of Advent of Code 2022.
package aoc22_02 //nolint:revive,stylecheck // I don't care about the package name

import (
	"strconv"
)

// X = 1, Y = 2, Z = 3
// Loss = 0, Draw = 3, Win = 6.
var scores = map[string]int{
	"A X": 4, // RR: Draw (1 + 3)
	"A Y": 8, // RP: Win (2 + 6)
	"A Z": 3, // RS: Loss (3 + 0)

	"B X": 1, // PR: Loss (1 + 0)
	"B Y": 5, // PP: Draw (2 + 3)
	"B Z": 9, // PS: Win (3 + 6)

	"C X": 7, // SR: Win (1 + 6)
	"C Y": 2, // SP: Loss (2 + 0)
	"C Z": 6, // SS: Draw (3 + 3)
}

// X = Loss, Y = Draw, Z = Win.
var plays = map[string]string{
	"A X": "A Z",
	"A Y": "A X",
	"A Z": "A Y",

	"B X": "B X",
	"B Y": "B Y",
	"B Z": "B Z",

	"C X": "C Y",
	"C Y": "C Z",
	"C Z": "C X",
}

// D2P1 returns the solution for 2021 day 2 part 1
// answer: 11906
func D2P1(data []string) string {
	score := 0
	for _, line := range data {
		score += scores[line]
	}

	return strconv.Itoa(score)
}

// D2P2 returns the solution for 2021 day 2 part 2
// answer: 11186
func D2P2(data []string) string {
	score := 0

	for _, line := range data {
		score += scores[plays[line]]
	}

	return strconv.Itoa(score)
}
