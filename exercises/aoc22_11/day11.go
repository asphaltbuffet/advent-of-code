// Package aoc22_11 contains the solution for day 11 of Advent of Code 2022.
package aoc22_11 //nolint:revive,stylecheck // I don't care about the package name

import "fmt"

// Monkey struct represents necessary aspect of a monkey for day 11.
// Items is a slice of worry levels for each item currently held by this monkey.
// Operation is a function that modifies a worry level after inspected.
// Test determines where an item is moved after inspection based on worry level.
type Monkey struct {
	Items     []int
	Operation func(int) int
	Test      func(int) int
	Count     int
}

// Day11 is the exercise environment.
type Day11 struct {
	Monkeys []Monkey
}

// D11P1 returns the solution for 2022 day 11 part 1.
//
// https://adventofcode.com/2022/day/11
//
// answer:
func D11P1(data []string) string {
	d := Day11{
		Monkeys: []Monkey{},
	}

	// Read monkeys into structs
	err := d.ParseInput(data)
	if err != nil {
		return fmt.Sprintf("parsing input: %v", err)
	}

	// Process a round
	// - monkey checks each item: increase worry, divide by 3, move item
	// - increase count for items inspected by that monkey
	// - next monkey

	// find top 2 monkeys based on inspection counnt

	// multiply top 2 counts together & return as string
	return ""
}

// D11P2 returns the solution for 2022 day 11 part 2.
// answer:
func D11P2(data []string) string {
	return ""
}
