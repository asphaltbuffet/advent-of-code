package aoc22

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 5, D5P1, D5P2, Get2022Command())
}

// D5P1 returns the solution for 2022 day 5 part 1
// answer:
func D5P1(data []string) string {
	// get number of stacks
	stackCount := getNumberOfStacks(data)
	fmt.Printf("%d stacks found\n", stackCount)

	//
	stacks := make([][]string, stackCount)

	for i := 0; i < stackCount; i++ {
		stacks[i] = make([]string, 0)
	}

	return "incomplete"
}

// D5P2 returns the solution for 2022 day 5 part 2
// answer:
func D5P2(data []string) string {
	return "not implemented"
}

func getNumberOfStacks(data string) int {
	return (len(data) + 1) / 4
}
