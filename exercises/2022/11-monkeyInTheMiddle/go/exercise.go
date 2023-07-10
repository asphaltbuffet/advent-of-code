package exercises

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 11
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 151312
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	d := Day11{
		Monkeys: []*Monkey{},
		Product: 1,
	}

	// Read monkeys into structs
	err := d.ParseInput(data)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Process 20 rounds
	for i := 0; i < 20; i++ {
		err = d.ProcessRound()
		if err != nil {
			return nil, fmt.Errorf("processing round %d: %w", i+1, err)
		}
	}

	// find top 2 monkeys based on inspection count
	h := getHeap(d.Monkeys)

	// log.Printf("monkey heap: %+v", h)

	first, _ := heap.Pop(h).(Monkey)
	second, _ := heap.Pop(h).(Monkey)

	// log.Printf("top 2 monkeys: %v, %v", first, second)

	// multiply top 2 counts together & return as string
	return first.Count * second.Count, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 51382025916
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	d := Day11{
		Monkeys: []*Monkey{},
		Product: 1,
	}

	// Read monkeys into structs
	err := d.ParseInput(data)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// Process 10,000 rounds
	for i := 0; i < 10000; i++ {
		// log.Printf("Round %d", i)
		err = d.ProcessRoundPart2()
		if err != nil {
			return nil, fmt.Errorf("processing round %d: %w", i+1, err)
		}
	}

	// find top 2 monkeys based on inspection count
	h := getHeap(d.Monkeys)

	// log.Printf("monkey heap: %+v", h)

	first, _ := heap.Pop(h).(Monkey)
	second, _ := heap.Pop(h).(Monkey)

	// log.Printf("top 2 monkeys: %v, %v", first, second)

	// multiply top 2 counts together & return as string
	return first.Count * second.Count, nil
}

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
