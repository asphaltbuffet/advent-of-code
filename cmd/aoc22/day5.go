package aoc22

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 5, D5P1, D5P2, Get2022Command())
}

// Action is the struct for an action.
type Action struct {
	qty int
	src int
	tgt int
}

// Day5 is the struct for day 5.
type Day5 struct {
	Stacks  [][]string
	Actions []Action
}

// D5P1 returns the solution for 2022 day 5 part 1
// answer: SHMSDGZVC
func D5P1(data []string) string {
	day := Day5{
		Stacks:  [][]string{},
		Actions: []Action{},
	}

	// read crates into appropriate stacks
	day.populateStacks(data)

	// debugPrintStacks(day)

	// parse actions and store in day struct
	err := day.populateActions(data)
	if err != nil {
		return fmt.Sprintf("failed to parse stacks: %v", err)
	}

	// fmt.Printf("actions: %+v\n", day.Actions)

	// loop through actions and apply to stacks
	err = day.processActionsSingleMove()
	if err != nil {
		return fmt.Sprintf("failed to parse actions: %v", err)
	}

	return day.topCrates()
}

func (day *Day5) populateActions(data []string) error {
	for i := GetMovementSectionLine(data); i < len(data); i++ {
		a := Action{}

		// fmt.Printf("line %d: %s\n", i, data[i])

		_, err := fmt.Sscanf(data[i], "move %d from %d to %d", &a.qty, &a.src, &a.tgt)
		if err != nil {
			return fmt.Errorf("failed to parse line \"%s\": %w", data[i], err)
		}

		// normalize stack numbers to indexes
		a.src--
		a.tgt--

		day.Actions = append(day.Actions, a)
	}

	return nil
}

func (day *Day5) processActionsSingleMove() error {
	for _, a := range day.Actions {
		// move each crate separately
		for i := 0; i < a.qty; i++ {
			// pop
			n := len(day.Stacks[a.src]) - 1
			if n < 0 {
				return fmt.Errorf("stack %d is empty", a.src)
			}

			c := day.Stacks[a.src][n]
			day.Stacks[a.src] = day.Stacks[a.src][:n]

			// fmt.Printf("moved %s; %d -> %d\n", c, a.src, a.tgt)

			// push
			day.Stacks[a.tgt] = append(day.Stacks[a.tgt], c)
		}
	}

	return nil
}

func (day *Day5) populateStacks(data []string) {
	stackCount := GetNumberOfStacks(data[0])
	// fmt.Printf("%d stacks found\n", stackCount)

	day.Stacks = make([][]string, stackCount)

	for i := 0; i < stackCount; i++ {
		day.Stacks[i] = ParseStack(data, i)
	}
}

// func debugPrintStacks(day Day5) {
// 	for i, v := range day.Stacks {
// 		fmt.Printf("stack[%d]: %v\n", i+1, v)
// 	}
// }

// GetNumberOfStacks returns the size of the stack section. This does not include stack numbers or the blank line.
func GetNumberOfStacks(data string) int {
	return (len(data) + 1) / 4
}

// GetMovementSectionLine returns the line number where the movement section starts.
func GetMovementSectionLine(data []string) (h int) {
	for h = 0; h < len(data); h++ {
		if data[h] != "" && strings.HasPrefix(data[h], "move") {
			return
		}
	}

	fmt.Println("failed to find movement section")

	return -1
}

// ParseStack returns a slice of strings representing the stack at the given index.
func ParseStack(data []string, stackIndex int) []string {
	if stackIndex < 0 || stackIndex >= GetNumberOfStacks(data[0]) {
		fmt.Printf("invalid stack index: %d\n", stackIndex)
		return nil
	}

	var stack []string

	const w int = 4

	crate := ""

	for i := 0; !strings.HasPrefix(data[i], " 1"); i++ {
		crate = string(data[i][(stackIndex*w)+1])
		// fmt.Printf("[%d][%d]: '%s'", stackNumber, i, crate)

		// don't add blank areas to stack
		if crate != " " {
			stack = append([]string{crate}, stack...)
		}
	}
	// fmt.Println()

	return stack
}

func (day *Day5) topCrates() string {
	var sb strings.Builder

	for _, s := range day.Stacks {
		sb.WriteString(s[len(s)-1])
	}

	return sb.String()
}

// D5P2 returns the solution for 2022 day 5 part 2
// answer: VRZGHDFBQ
func D5P2(data []string) string {
	day := Day5{
		Stacks:  [][]string{},
		Actions: []Action{},
	}

	// read crates into appropriate stacks
	day.populateStacks(data)

	// debugPrintStacks(day)

	// parse actions and store in day struct
	err := day.populateActions(data)
	if err != nil {
		return fmt.Sprintf("failed to parse stacks: %v", err)
	}

	// fmt.Printf("actions: %+v\n", day.Actions)

	// loop through actions and apply to stacks
	err = day.processActionsBulkMove()
	if err != nil {
		return fmt.Sprintf("failed to parse actions: %v", err)
	}

	return day.topCrates()
}

func (day *Day5) processActionsBulkMove() error {
	for _, a := range day.Actions {
		// pop
		n := len(day.Stacks[a.src])
		if n < 0 {
			return fmt.Errorf("stack %d is empty", a.src)
		}

		b := n - a.qty
		c := day.Stacks[a.src][b:]
		day.Stacks[a.src] = day.Stacks[a.src][:b]

		// fmt.Printf("moved %s; %d -> %d\n", c, a.src, a.tgt)

		// push
		day.Stacks[a.tgt] = append(day.Stacks[a.tgt], c...)
	}

	return nil
}
