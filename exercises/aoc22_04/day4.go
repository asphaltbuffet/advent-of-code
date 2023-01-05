// Package aoc22_04 contains the solution for day 4 of Advent of Code 2022.
package aoc22_04 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"strconv"
)

// D4P1 returns the solution for 2022 day 4 part 1
// answer: 532
func D4P1(data []string) string {
	count := 0

	for _, line := range data {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isFullyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return strconv.Itoa(count)
}

type elf struct {
	low  int
	high int
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

// D4P2 returns the solution for 2022 day 4 part 2
// answer: 854
func D4P2(data []string) string {
	count := 0

	for _, line := range data {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isAnyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return strconv.Itoa(count)
}

func isAnyOverlapping(a, b *elf) bool {
	return (a.low <= b.low && a.high >= b.low) ||
		(b.low <= a.low && b.high >= a.low)
}
