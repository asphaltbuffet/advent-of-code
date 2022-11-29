// Package y2021 contains all solutions for the 2021 year.
package y2021

import (
	"github.com/spf13/cobra"
)

// New2021Command creates a new command for the 2021 year.
func New2021Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2021",
		Short: "Exercises for 2021 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("2021 Command")
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(newDay1Command())
	// cmd.AddCommand(NewDay2Command())
	// cmd.AddCommand(NewDay3Command())
	// cmd.AddCommand(NewDay4Command())
	// cmd.AddCommand(NewDay5Command())
	// cmd.AddCommand(NewDay6Command())
	// cmd.AddCommand(NewDay7Command())
	// cmd.AddCommand(NewDay8Command())
	// cmd.AddCommand(NewDay9Command())
	// cmd.AddCommand(NewDay10Command())
	// cmd.AddCommand(NewDay11Command())
	// cmd.AddCommand(NewDay12Command())
	// cmd.AddCommand(NewDay13Command())
	// cmd.AddCommand(NewDay14Command())
	// cmd.AddCommand(NewDay15Command())
	// cmd.AddCommand(NewDay16Command())
	// cmd.AddCommand(NewDay17Command())
	// cmd.AddCommand(NewDay18Command())
	// cmd.AddCommand(NewDay19Command())
	// cmd.AddCommand(NewDay20Command())
	// cmd.AddCommand(NewDay21Command())
	// cmd.AddCommand(NewDay22Command())
	// cmd.AddCommand(NewDay23Command())
	// cmd.AddCommand(NewDay24Command())
	// cmd.AddCommand(NewDay25Command())

	return cmd
}
