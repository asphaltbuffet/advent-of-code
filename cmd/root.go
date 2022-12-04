// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

var (
	rootCmd   *cobra.Command
	yearGroup *cobra.Group
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(GetRootCommand().Execute())
}

// GetRootCommand returns the root command for the CLI.
func GetRootCommand() *cobra.Command {
	if rootCmd == nil {
		rootCmd = &cobra.Command{
			Use:     "aoc",
			Aliases: []string{"advent-of-code", "advent"},
			Version: "0.0.1",
			Short:   "aoc is a collection of AoC solutions",
			Long:    `aoc is a collection of AoC solutions`,
			Run:     RunRootCmd,
		}

		rootCmd.Flags().Bool("all", false, "run all solutions")

		yearGroup = &cobra.Group{ID: "Years", Title: "Advent of Code events"}
		rootCmd.AddGroup(yearGroup)
	}

	return rootCmd
}

// RunRootCmd is the entry point for the CLI.
func RunRootCmd(cmd *cobra.Command, args []string) {
	if ok, _ := rootCmd.Flags().GetBool("all"); ok {
		runAllYears()
		return
	}

	_ = cmd.Help()
}

func runAllYears() {
	// fmt.Printf("groups: %+v", rootCmd.Groups())
	yearCommands := aoc.Filter(rootCmd.Commands(), func(c *cobra.Command) bool { return c.GroupID == "Years" })
	// yearCommands := rootCmd.Commands()
	fmt.Printf("found %d years:\n", len(yearCommands))

	for _, yearCmd := range yearCommands {
		exercises := yearCmd.Commands()

		// fmt.Printf("\tfound %d exercises for %s:\n", len(exercises), yearCmd.Name())

		for _, exercise := range exercises {
			// fmt.Printf("\t\t%s\n", exercise.Name())
			exercise.Run(exercise, []string{})
		}
	}
}
