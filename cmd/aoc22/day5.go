package aoc22

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 5, D5P1, D5P2, Get2022Command())
}

// D5P1 returns the solution for 2022 day 5 part 1
// answer:
func D5P1(data []string) string {
	// get number of stacks
	stackCount := getNumberOfStacks(data[0])
	fmt.Printf("%d stacks found\n", stackCount)

	stacks := make([][]string, stackCount)

	for i := 0; i < stackCount; i++ {
		stacks[i] = parseStack(data, i)
	}

	for i, v := range stacks {
		fmt.Printf("stack[%d]: %v\n", i+1, v)
	}

	var (
		qty int
		src int
		tgt int
	)

	for i := getNumberOfStacks(data[0]) + 2; i < len(data); i++ {
		fmt.Printf("line %d: %s\n", i, data[i])

		_, err := fmt.Sscanf(data[i], "move %d from %d to %d", &qty, &src, &tgt)
		if err != nil {
			return fmt.Sprintf("failed to parse line \"%s\": %v", data[i], err)
		}

		// normalize stack numbers to indexes
		// move each crate separately
		for i := 0; i < qty; i++ {

			// pop
			n := len(stacks[src-1]) - 1
			c := stacks[src-1][n]
			stacks[src-1] = stacks[src-1][:n]

			// push
			stacks[tgt-1] = append(stacks[tgt-1], c)

			fmt.Printf("moved %s; %d -> %d\n", c, src-1, tgt-1)
		}
	}

	return printStacks(stacks)
}

// D5P2 returns the solution for 2022 day 5 part 2
// answer:
func D5P2(data []string) string {
	return "not implemented"
}

func getNumberOfStacks(data string) int {
	return (len(data) + 1) / 4
}

func getStackHeight(data []string) int {
	h := 0

	for i := 0; i < len(data); i++ {
		if data[i][:2] == " 1" {
			h = i + 1
			break
		}
	}

	// fmt.Printf("stack height is: %d\n", h)

	return h
}

func parseStack(data []string, stackNumber int) []string {
	height := getStackHeight(data)
	stack := make([]string, 0)

	const w int = 4

	crate := ""

	for i := 0; i < height-1; i++ {
		crate = data[i][stackNumber*w+1 : stackNumber*w+2]
		// fmt.Printf("[%d][%d]: '%s'", stackNumber, i, crate)

		// don't add blank areas to stack
		if crate != " " {
			stack = append([]string{crate}, stack...)
		}
	}
	// fmt.Println()

	return stack
}

func printStacks(stacks [][]string) string {
	var sb strings.Builder

	for _, s := range stacks {
		sb.WriteString(s[len(s)-1])
	}

	return sb.String()
}
