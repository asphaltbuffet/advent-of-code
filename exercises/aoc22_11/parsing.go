package aoc22_11 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ParseInput puts monkey information into struct
// Read monkeys into structs
// - starting items
// - operation
// - test + monkey targets.
func (d *Day11) ParseInput(data []string) error {
	return nil
}

// ParseMonkey parses a single monkey.
func ParseMonkey(m []string) error {
	if len(m) != 6 {
		return fmt.Errorf("invalid number of input lines: %d", len(m))
	}

	_, n, ok := strings.Cut(m[0], " ")
	if !ok {
		return fmt.Errorf("' ' not found in '%s'", m[0])
	}

	num, err := strconv.Atoi(strings.Trim(n, " :"))
	if err != nil {
		return fmt.Errorf("parsing monkey number: %w", err)
	}

	log.Printf("monkey number: %d", num)

	return nil
}
