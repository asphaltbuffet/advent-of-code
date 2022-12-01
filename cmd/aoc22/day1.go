package aoc22

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	NewDay1Command()
}

// NewDay1Command creates a new command for the 2021 day 1 exercise
func NewDay1Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "1",
		Short: "day 1 exercises for 2022 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			got := D1P1(d.PartOne.Input)
			cmd.Printf("Day %s Part 1: %s\n", cmd.Name(), got)

			got = D1P2(d.PartTwo.Input)
			cmd.Printf("Day %s Part 2: %s\n", cmd.Name(), got)
		},
	}

	Get2021Command().AddCommand(cmd)

	return cmd
}

// D1P1 returns the solution for 2022 day 1 part 1
// answer: 70720
func D1P1(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0])
}

// D1P2 returns the solution for 2022 day 1 part 2
// answer:
func D1P2(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0] + calories[1] + calories[2])
}
