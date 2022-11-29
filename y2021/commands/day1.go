package y2021

import (
	"reflect"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func newDay1Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "1",
		Short: "day 1 exercise for 2021 AoC",
		Run:   runDay1Cmd,
	}

	return cmd
}

func runDay1Cmd(cmd *cobra.Command, args []string) {
	d, err := aoc.ReadTestData(cmd.Parent().Name(), cmd.Name()) // TODO: pass in a function to convert the data by line
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	got := part1(d.Test.Input) // TODO: use unit testing for this validation

	cmd.Printf("Day %s, Part 1 Test: %t\n", cmd.Name(), reflect.DeepEqual(got, d.Test.Output))

	got = part1(d.Input)
	// fmt.Printf("inputs(%d): %+v\n", len(d.Input), d.Input)
	cmd.Printf("Day %s Part 1: %d\n", cmd.Name(), got)

	got = part2(d.Test.Input) // TODO: use unit testing for this validation

	cmd.Printf("Day %s, Part 2 Test: %t\n", cmd.Name(), got == 5) // TODO: example structure should have both parts stored in it

	got = part2(d.Input)
	// fmt.Printf("inputs(%d): %+v\n", len(d.Input), d.Input)
	cmd.Printf("Day %s Part 2: %d\n", cmd.Name(), got)
}

// answer: 1711
func part1(data []int) int {
	count := 0
	prev := data[0]

	for i := 1; i < len(data); i++ {
		curr := data[i]
		if curr > prev {
			count++
		}

		prev = curr
	}

	return count
}

// answer: 1743
func part2(data []int) int {
	w := []int{}

	// calculate a new array of windowed values
	for i := 2; i < len(data); i++ {
		w = append(w, data[i]+data[i-1]+data[i-2])
	}

	return part1(w)
}
