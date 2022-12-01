package aoc21

import (
	"strconv"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	NewDay1Command()
}

// NewDay1Command creates a new command for the 2021 day 1 exercise.
func NewDay1Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "1",
		Short: "day 1 exercise for 2021 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			got := Day1part1(d.PartOne.Input)
			cmd.Printf("Day %s Part 1: %s\n", cmd.Name(), got)

			got = Day1part2(d.PartTwo.Input)
			cmd.Printf("Day %s Part 2: %s\n", cmd.Name(), got)
		},
	}

	Get2021Command().AddCommand(cmd)

	return cmd
}

// Day1part1 returns the solution for 2021 day 1 part 1
// answer: 1711
func Day1part1(data []string) string {
	formattedData, _ := aoc.ConvertStringSliceToIntSlice(data)
	return strconv.Itoa(increasingCount(formattedData))
}

func increasingCount(data []int) int {
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

// Day1part2 returns the solution for 2021 day 1 part 2
// answer: 1743
func Day1part2(data []string) string {
	formattedData, _ := aoc.ConvertStringSliceToIntSlice(data)

	w := []int{}

	// calculate a new array of windowed values
	for i := 2; i < len(formattedData); i++ {
		w = append(w, formattedData[i]+formattedData[i-1]+formattedData[i-2])
	}

	return strconv.Itoa(increasingCount(w))
}
