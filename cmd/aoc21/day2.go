package aoc21

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	newD2P1Command()
}

func newD2P1Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2",
		Short: "day 2 exercise for 2021 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			cmd.Printf("┌──────────────────┒\n")
			cmd.Printf("│      Day %s       ┃\n", cmd.Name())
			cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

			got := D2P1(d.PartOne.Input)
			cmd.Printf("Part 1: %s\n", got)

			got = D2P2(d.PartTwo.Input)
			cmd.Printf("Part 2: %s\n", got)
		},
	}

	Get2021Command().AddCommand(cmd)

	return cmd
}

type command struct {
	operation string
	value     int
}

// D2P1 returns the solution for 2021 day 2 part 1
// answer: 1714950
func D2P1(data []string) string {
	commands := aoc.Map(data,
		func(s string) command {
			op, val, _ := strings.Cut(s, " ")

			v, _ := strconv.Atoi(val)

			return command{
				operation: op,
				value:     v,
			}
		})

	horiz, depth := 0, 0

	for _, c := range commands {
		switch c.operation {
		case "forward":
			horiz += c.value
		case "down":
			depth += c.value
		case "up":
			depth -= c.value

			if depth < 0 {
				depth = 0
			}
		}
	}

	return fmt.Sprintf("%d", horiz*depth)
}

// D2P2 returns the solution for 2021 day 2 part 2
// answer: 1281977850
func D2P2(data []string) string {
	commands := aoc.Map(data,
		func(s string) command {
			op, val, _ := strings.Cut(s, " ")

			v, _ := strconv.Atoi(val)

			return command{
				operation: op,
				value:     v,
			}
		})

	var horiz, depth, aim int

	for _, c := range commands {
		switch c.operation {
		case "forward":
			horiz += c.value
			depth += c.value * aim

			if depth < 0 {
				depth = 0
			}
		case "down":
			aim += c.value
		case "up":
			aim -= c.value
		}
	}

	return fmt.Sprintf("%d", horiz*depth)
}
