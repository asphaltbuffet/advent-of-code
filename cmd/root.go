// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "advent-of-code [sub-commands] [flags]",
		Short: "advent-of-code is a collection of AoC solutions",
		Long:  `advent-of-code is a collection of AoC solutions`,
		// Args:  cobra.MinimumNArgs(1),
		Run: RunRootCmd,
	}

	rootCmd.Flags().Bool("svg", false, "output SVG")

	rootCmd.AddCommand(new2021Command())

	cobra.CheckErr(rootCmd.Execute())
}

// RunRootCmd is the entry point for the CLI.
func RunRootCmd(cmd *cobra.Command, args []string) {
	cmd.Println("RunRootCmd")
}
