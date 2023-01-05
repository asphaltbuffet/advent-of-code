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
	// Read monkeys into structs
	for i := 0; i < len(data); i += 7 {
		m := data[i : i+6]

		monkey, err := ParseMonkey(m)
		if err != nil {
			return fmt.Errorf("parsing monkey: %w", err)
		}

		d.Monkeys = append(d.Monkeys, monkey)

		d.Product *= monkey.Divisor
	}

	log.Printf("product: %d", d.Product)

	return nil
}

// ParseMonkey parses a single monkey.
func ParseMonkey(input []string) (*Monkey, error) {
	if len(input) != 6 {
		return nil, fmt.Errorf("invalid number of input lines: %d", len(input))
	}

	m := &Monkey{}

	_, n, ok := strings.Cut(input[0], " ")
	if !ok {
		return nil, fmt.Errorf("' ' not found in '%s'", input[0])
	}

	num, err := strconv.Atoi(strings.Trim(n, " :"))
	if err != nil {
		return nil, fmt.Errorf("parsing monkey number: %w", err)
	}

	m.ID = num

	// log.Printf("monkey number: %d", num)

	// starting items
	trimmed := strings.Trim(input[1], " ")
	tokens := strings.Split(trimmed, " ")

	for _, token := range tokens[2:] {
		t := strings.Trim(token, ", ")

		i, itemErr := strconv.Atoi(t)
		if itemErr != nil {
			return nil, fmt.Errorf("parsing starting item: %w", itemErr)
		}

		m.Items = append(m.Items, i)
	}

	// log.Printf("monkey items: %v", m.Items)

	// operation
	trimmed = strings.Trim(input[2], " ")
	tokens = strings.Split(trimmed, " ")

	if tokens[3] == tokens[5] {
		m.Operator = "^"
		m.Scalar = 2
	} else {
		m.Operator = tokens[4]

		m.Scalar, err = strconv.Atoi(tokens[5])
		if err != nil {
			return nil, fmt.Errorf("parsing scalar: %w", err)
		}
	}

	// log.Printf("monkey operator: %s", m.Operator)
	// log.Printf("monkey scalar: %d", m.Scalar)

	// test
	trimmed = strings.Trim(input[3], " ")
	tokens = strings.Split(trimmed, " ")

	m.Divisor, err = strconv.Atoi(tokens[3])
	if err != nil {
		return nil, fmt.Errorf("parsing divisor: %w", err)
	}

	// log.Printf("monkey divisor: %d", m.Divisor)

	// monkey targets
	trimmed = strings.Trim(input[4], " ")
	tokens = strings.Split(trimmed, " ")

	m.TargetOne, err = strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("parsing target one: %w", err)
	}

	trimmed = strings.Trim(input[5], " ")
	tokens = strings.Split(trimmed, " ")

	m.TargetTwo, err = strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("parsing target two: %w", err)
	}

	// log.Printf("monkey target one: %d", m.TargetOne)
	// log.Printf("monkey target two: %d", m.TargetTwo)

	return m, nil
}
