package aoc22

import (
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/cmd"
)

var yearCmd *cobra.Command

// Get2021Command creates a new command for the 2021 year.
func Get2021Command() *cobra.Command {
	if yearCmd == nil {
		yearCmd = &cobra.Command{
			Use:   "2022",
			Short: "Exercises for 2022 AoC",
		}

		cmd.GetRootCommand().AddCommand(yearCmd)
	}

	return yearCmd
}