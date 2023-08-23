package exercises

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ParseInput puts monkey information into struct
func (d *Day11) ParseInput(instr string) error {
	// Read monkeys into structs
	for _, m := range strings.Split(instr, "\n\n") {
		monkey, err := ParseMonkey(m)
		if err != nil {
			return fmt.Errorf("parsing monkey: %w", err)
		}

		d.Monkeys = append(d.Monkeys, monkey)

		d.Product *= monkey.Divisor
	}

	log.Printf("product: %d", &d.Product)

	return nil
}

// ParseMonkey parses a single monkey.
func ParseMonkey(mstr string) (*Monkey, error) {
	input := strings.Split(mstr, "\n")

	var id int
	_, err := fmt.Sscanf(input[0], "Monkey %d:", &id)
	if err != nil {
		return nil, fmt.Errorf("parsing monkey number: %w", err)
	}

	// starting items
	_, tokens, _ := strings.Cut(input[1], ": ")

	var items []int
	for _, n := range strings.Split(tokens, ", ") {
		i, itemErr := strconv.Atoi(n)
		if itemErr != nil {
			return nil, fmt.Errorf("parsing starting item: %w", itemErr)
		}

		items = append(items, i)
	}

	// operation
	var operator string
	var scalar int
	_, err = fmt.Sscanf(input[2], "  Operation: new = old %s %d", &operator, &scalar)
	if err != nil {
		operator = "^"
		scalar = 2
	}

	// test
	var divisor int
	_, _ = fmt.Sscanf(input[3], "  Test: divisible by %d", &divisor)

	// monkey targets
	var t1, t2 int
	_, _ = fmt.Sscanf(input[4], "    If true: throw to monkey %d", &t1)
	_, _ = fmt.Sscanf(input[5], "    If false: throw to monkey %d", &t2)

	return &Monkey{
		ID:        id,
		Items:     items,
		Divisor:   divisor,
		Operator:  operator,
		Scalar:    scalar,
		TargetOne: t1,
		TargetTwo: t2,
	}, nil
}
