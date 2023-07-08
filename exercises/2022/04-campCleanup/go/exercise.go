package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

type elf struct {
	low  int
	high int
}

func isAnyOverlapping(a, b *elf) bool {
	return (a.low <= b.low && a.high >= b.low) ||
		(b.low <= a.low && b.high >= a.low)
}

func isFullyOverlapping(a, b *elf) bool {
	return (a.low <= b.low && a.high >= b.high) ||
		(b.low <= a.low && b.high >= a.high)
}

func parsePair(line string) (*elf, *elf) {
	var a, b elf

	_, err := fmt.Sscanf(line, "%d-%d,%d-%d", &a.low, &a.high, &b.low, &b.high)
	if err != nil {
		fmt.Printf("error parsing line: %s", line)
	}

	return &a, &b
}

// Exercise for Advent of Code 2022 day 4
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise. // answer: 532
func (c Exercise) One(instr string) (any, error) {
	count := 0

	for _, line := range strings.Split(instr, "\n") {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isFullyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (c Exercise) Two(instr string) (any, error) {
	count := 0

	for _, line := range strings.Split(instr, "\n") {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isAnyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return count, nil
}
