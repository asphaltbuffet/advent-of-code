package aoc22

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/cmd"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

var (
	yearCmd  *cobra.Command
	dayGroup *cobra.Group
)

// Get2022Command creates a new command for the 2022 year.
func Get2022Command() *cobra.Command {
	if yearCmd == nil {
		yearCmd = &cobra.Command{
			Use:   "2022",
			Short: "Exercises for 2022 AoC",
			Long:  "https://adventofcode.com/2022",
			Run:   RunYearCmd,
		}

		cmd.GetRootCommand().AddCommand(yearCmd)
		yearCmd.GroupID = "Years"

		dayGroup = &cobra.Group{ID: "days", Title: "Exercises"}
		yearCmd.AddGroup(dayGroup)
	}

	return yearCmd
}

// RunYearCmd is the entry point for the 2022 year.
func RunYearCmd(c *cobra.Command, args []string) {
	if ok, _ := cmd.GetRootCommand().PersistentFlags().GetBool("all"); ok {
		runAllExercises(c, args)
		return
	}

	_ = c.Help()
}

func runAllExercises(cmd *cobra.Command, args []string) {
	color.Set(color.FgYellow)

	cmd.Printf("┌──────────────────┒\n")
	cmd.Printf("│     AoC %-4s     ┃\n", cmd.Name())
	cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

	color.Unset()

	exercises := aoc.Filter(cmd.Commands(), func(c *cobra.Command) bool { return c.GroupID == "days" })

	// fmt.Printf("\tfound %d exercises for %s:\n", len(exercises), yearCmd.Name())

	for _, exercise := range exercises {
		// fmt.Printf("\t\t%s\n", exercise.Name())
		exercise.Run(exercise, args)
	}
}
