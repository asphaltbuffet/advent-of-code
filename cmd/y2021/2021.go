package y2021

import (
	"github.com/spf13/cobra"
)

func new2021Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2021",
		Short: "Exercises for 2021 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("2021 Command")
		},
	}

	return cmd
}
