// Package cmd contains all CLI commands used by the application.
package cmd

import (
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
			Use:     "aoc [command]",
			Aliases: []string{"advent-of-code", "advent"},
			Version: "0.0.1",
			Short:   "aoc is a collection of AoC solutions",
			Long:    `aoc is a collection of AoC solutions`,
			Run:     RunRootCmd,
		}

		rootCmd.PersistentFlags().Bool("all", false, "run all solutions")

		yearGroup = &cobra.Group{ID: "Years", Title: "Advent of Code events"}
		rootCmd.AddGroup(yearGroup)
	}

	return rootCmd
}

// RunRootCmd is the entry point for the CLI.
func RunRootCmd(cmd *cobra.Command, args []string) {
	ok, err := rootCmd.Flags().GetBool("all")
	if err != nil {
		cmd.PrintErrf("error: %s", err)
	}

	if ok {
		yearCommands := aoc.Filter(rootCmd.Commands(), func(c *cobra.Command) bool { return c.GroupID == "Years" })
		// fmt.Printf("found %d years:\n", len(yearCommands))

		for _, yearCmd := range yearCommands {
			yearCmd.Run(yearCmd, args)
		}

		return
	}

	_ = cmd.Help()
}
