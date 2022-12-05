package common

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewDayCommand creates a new command for the given day.
func NewDayCommand(year, day int, part1, part2 func([]string) string, parentCmd *cobra.Command) *cobra.Command {
	// check if year is between 2015 and 2022
	if year < 2015 || year > time.Now().Year() {
		panic("year must be between 2015 and this year")
	}

	// check if day is between 1 and 25
	if day < 1 || day > 25 {
		panic("day must be between 1 and 25")
	}

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%d", day),
		Short: fmt.Sprintf("day %d exercise for %d AoC", day, year),
		Run: func(cmd *cobra.Command, args []string) {
			d, err := NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				return
			}

			color.Set(color.FgYellow)

			cmd.Printf("┌──────────────────┒\n")
			cmd.Printf("│      Day %-2s      ┃\n", cmd.Name())
			cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

			color.Unset()

			got := part1(d.PartOne.Input)
			cmd.Printf("Part 1: %s\n", got)

			got = part2(d.PartTwo.Input)
			cmd.Printf("Part 2: %s\n", got)
		},
	}

	parentCmd.AddCommand(cmd)
	cmd.GroupID = "days"

	return cmd
}
