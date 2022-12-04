package aoc22

import (
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/cmd"
)

var (
	yearCmd  *cobra.Command
)

// Get2022Command creates a new command for the 2021 year.
func Get2022Command() *cobra.Command {
	if yearCmd == nil {
		yearCmd = &cobra.Command{
			Use:   "2022",
			Short: "Exercises for 2022 AoC",
			Long:  "https://adventofcode.com/2022",
		}

		yearCmd.Flags().Bool("all", false, "process all exercises")

		cmd.GetRootCommand().AddCommand(yearCmd)
		yearCmd.GroupID = "Years"
	}

	return yearCmd
}
