// Package aoc22_11 contains the solution for day 11 of Advent of Code 2022.
package aoc22_11 //nolint:revive,stylecheck // I don't care about the package name

import (
	"container/heap"
	"fmt"
)

// Monkey struct represents necessary aspect of a monkey for day 11.
// Items is a slice of worry levels for each item currently held by this monkey.
// Operation is a function that modifies a worry level after inspected.
// Test determines where an item is moved after inspection based on worry level.
type Monkey struct {
	ID        int
	Items     []int
	Divisor   int
	Operator  string
	Scalar    int
	TargetOne int
	TargetTwo int
	Count     int
}

// Day11 is the exercise environment.
type Day11 struct {
	Monkeys []*Monkey
	Product int
}

// D11P1 returns the solution for 2022 day 11 part 1.
//
// https://adventofcode.com/2022/day/11
//
// answer:
func D11P1(data []string) string {
	d := Day11{
		Monkeys: []*Monkey{},
		Product: 1,
	}

	// Read monkeys into structs
	err := d.ParseInput(data)
	if err != nil {
		return fmt.Sprintf("parsing input: %v", err)
	}

	// Process 20 rounds
	for i := 0; i < 20; i++ {
		err = d.ProcessRound()
		if err != nil {
			return fmt.Sprintf("ERROR: processing round %d: %v", i+1, err)
		}
	}

	// find top 2 monkeys based on inspection count
	h := getHeap(d.Monkeys)

	// log.Printf("monkey heap: %+v", h)

	first, _ := heap.Pop(h).(Monkey)
	second, _ := heap.Pop(h).(Monkey)

	// log.Printf("top 2 monkeys: %v, %v", first, second)

	// multiply top 2 counts together & return as string
	return fmt.Sprintf("%d", first.Count*second.Count)
}

// D11P2 returns the solution for 2022 day 11 part 2.
// answer:
func D11P2(data []string) string {
	d := Day11{
		Monkeys: []*Monkey{},
		Product: 1,
	}

	// Read monkeys into structs
	err := d.ParseInput(data)
	if err != nil {
		return fmt.Sprintf("parsing input: %v", err)
	}

	// Process 10,000 rounds
	for i := 0; i < 10000; i++ {
		// log.Printf("Round %d", i)
		err = d.ProcessRoundPart2()
		if err != nil {
			return fmt.Sprintf("ERROR: processing round %d: %v", i+1, err)
		}
	}

	// find top 2 monkeys based on inspection count
	h := getHeap(d.Monkeys)

	// log.Printf("monkey heap: %+v", h)

	first, _ := heap.Pop(h).(Monkey)
	second, _ := heap.Pop(h).(Monkey)

	// log.Printf("top 2 monkeys: %v, %v", first, second)

	// multiply top 2 counts together & return as string
	return fmt.Sprintf("%d", first.Count*second.Count)
}
