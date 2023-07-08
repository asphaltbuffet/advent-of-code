package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Action is the struct for an action.
type Action struct {
	qty int
	src int
	tgt int
}

type Stacks [][]string

// Exercise for Advent of Code 2022 day 5
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise. // answer: SHMSDGZVC
func (c Exercise) One(instr string) (any, error) {
	stacks, actions, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// loop through actions and apply to stacks
	err = processActionsSingleMove(stacks, actions)
	if err != nil {
		return nil, fmt.Errorf("processing movements: %w", err)
	}

	return topCrates(stacks), nil
}

// Two returns the answer to the second part of the exercise. // answer: VRZGHDFBQ
func (c Exercise) Two(instr string) (any, error) {
	stacks, actions, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// loop through actions and apply to stacks
	err = processActionsBulkMove(stacks, actions)
	if err != nil {
		return nil, fmt.Errorf("processing movement: %w", err)
	}

	return topCrates(stacks), nil
}

func parse(data string) (Stacks, []Action, error) {
	stackstr, actionstr, _ := strings.Cut(data, "\n\n")

	stackstr, _, _ = strings.Cut(stackstr, "\n 1")

	// read crates into appropriate stacks
	stacks, err := parseStacks(stackstr)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing stacks: %w", err)
	}

	actions, err := parseActions(actionstr)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing actions: %w", err)
	}

	return stacks, actions, nil
}

func parseActions(data string) ([]Action, error) {
	var actions []Action

	for _, line := range strings.Split(data, "\n") {
		a := Action{}

		n, err := fmt.Sscanf(line, "move %d from %d to %d", &a.qty, &a.src, &a.tgt)
		if err != nil || n != 3 {
			return nil, fmt.Errorf("parsing action line: %w", err)
		}

		// normalize stack numbers to indexes
		a.src--
		a.tgt--

		actions = append(actions, a)
	}

	return actions, nil
}

func processActionsSingleMove(stacks Stacks, actions []Action) error {
	for _, a := range actions {
		// move each crate separately
		for i := 0; i < a.qty; i++ {
			// pop
			n := len(stacks[a.src]) - 1
			if n < 0 {
				return fmt.Errorf("stack %d is empty", a.src)
			}

			c := stacks[a.src][n]
			stacks[a.src] = stacks[a.src][:n]

			// push
			stacks[a.tgt] = append(stacks[a.tgt], c)
		}
	}

	return nil
}

func parseStacks(data string) (Stacks, error) {
	const w int = 4

	stackCount := getNumberOfStacks(data)

	stacks := make(Stacks, stackCount)

	for i := 0; i < stackCount; i++ {
		stacks[i] = make([]string, 0)
	}

	for _, line := range strings.Split(data, "\n") {
		if line == "" {
			return nil, fmt.Errorf("unexpected blank line")
		}

		for i := 0; i < stackCount; i++ {
			crate := string(line[(i*w)+1])

			// don't add blank areas to stack
			if crate != " " {
				stacks[i] = append([]string{crate}, stacks[i]...)
			}
		}
	}

	return stacks, nil
}

func debugPrintStacks(stacks Stacks) {
	for i, v := range stacks {
		fmt.Printf("stack[%d]: %v\n", i+1, v)
	}
}

// getNumberOfStacks returns the number of stacks of crates.
func getNumberOfStacks(data string) int {
	i := strings.Index(data, "\n") + 1
	if i%4 != 0 {
		return 0
	}

	return i / 4
}

func topCrates(stacks Stacks) string {
	var sb strings.Builder

	for _, s := range stacks {
		sb.WriteString(s[len(s)-1])
	}

	return sb.String()
}

func processActionsBulkMove(stacks Stacks, actions []Action) error {
	for _, a := range actions {
		// pop
		n := len(stacks[a.src])
		if n < 0 {
			return fmt.Errorf("stack %d is empty", a.src)
		}

		b := n - a.qty
		c := stacks[a.src][b:]
		stacks[a.src] = stacks[a.src][:b]

		// push
		stacks[a.tgt] = append(stacks[a.tgt], c...)
	}

	return nil
}
