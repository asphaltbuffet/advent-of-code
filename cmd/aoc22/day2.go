package aoc22

import (
	"strconv"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	newDay2Command()
}

func newDay2Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2",
		Short: "day 2 exercise for 2022 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			cmd.Printf("┌──────────────────┒\n")
			cmd.Printf("│      Day %-2s      ┃\n", cmd.Name())
			cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

			got := D2P1(d.PartOne.Input)
			cmd.Printf("Part 1: %s\n", got)

			got = D2P2(d.PartTwo.Input)
			cmd.Printf("Part 2: %s\n", got)
		},
	}

	Get2022Command().AddCommand(cmd)

	return cmd
}

// X = 1, Y = 2, Z = 3
// Loss = 0, Draw = 3, Win = 6
var scores = map[string]int{
	"A X": 4, // RR: Draw (1 + 3)
	"A Y": 8, // RP: Win (2 + 6)
	"A Z": 3, // RS: Loss (3 + 0)

	"B X": 1, // PR: Loss (1 + 0)
	"B Y": 5, // PP: Draw (2 + 3)
	"B Z": 9, // PS: Win (3 + 6)

	"C X": 7, // SR: Win (1 + 6)
	"C Y": 2, // SP: Loss (2 + 0)
	"C Z": 6, // SS: Draw (3 + 3)
}

type shape string

const (
	rock     shape = "X"
	paper    shape = "Y"
	scissors shape = "Z"
)

// X = Loss, Y = Draw, Z = Win
var plays = map[string]string{
	"A X": "A Z",
	"A Y": "A X",
	"A Z": "A Y",

	"B X": "B X",
	"B Y": "B Y",
	"B Z": "B Z",

	"C X": "C Y",
	"C Y": "C Z",
	"C Z": "C X",
}

// D2P1 returns the solution for 2021 day 2 part 1
// answer:
func D2P1(data []string) string {
	score := 0
	for _, line := range data {
		score += scores[line]
	}

	return strconv.Itoa(score)
}

// D2P2 returns the solution for 2021 day 2 part 2
// answer: 1281977850
func D2P2(data []string) string {
	score := 0

	for _, line := range data {
		score += scores[plays[line]]
	}

	return strconv.Itoa(score)
}
