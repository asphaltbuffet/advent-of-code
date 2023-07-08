package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
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

// Exercise for Advent of Code 2022 day 2
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise. // answer: 11906
func (c Exercise) One(instr string) (any, error) {
	score := 0

	for _, line := range strings.Split(instr, "\n") {
		score += scores[line]
	}

	return score, nil
}

// Two returns the answer to the second part of the exercise. // answer: 11186
func (c Exercise) Two(instr string) (any, error) {
	score := 0

	for _, line := range strings.Split(instr, "\n") {
		score += scores[plays[line]]
	}

	return score, nil
}
