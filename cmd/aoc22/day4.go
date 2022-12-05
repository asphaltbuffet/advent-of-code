package aoc22

import (
	"fmt"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 4, D4P1, D4P2, Get2022Command())
}

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
