// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

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
	}

	return rootCmd
}

// RunRootCmd is the entry point for the CLI.
func RunRootCmd(cmd *cobra.Command, args []string) {
	cmd.Println("RunRootCmd")
}
