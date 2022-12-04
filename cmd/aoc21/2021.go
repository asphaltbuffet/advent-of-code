package aoc21

import (
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/cmd"
)

var yearCmd *cobra.Command

// Get2021Command creates a new command for the 2021 year.
func Get2021Command() *cobra.Command {
	if yearCmd == nil {
		yearCmd = &cobra.Command{
			Use:   "2021",
			Short: "Exercises for 2021 AoC",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println("2021 Command")
				_ = cmd.Help()
			},
		}

		cmd.GetRootCommand().AddCommand(yearCmd)
		yearCmd.GroupID = "Years"
	}

	return yearCmd
}
