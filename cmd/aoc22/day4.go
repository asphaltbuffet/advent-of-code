package aoc22

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	newDay4Command()
}

func newDay4Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "4",
		Short: "day 4 exercise for 2022 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				return
			}

			cmd.Printf("┌──────────────────┒\n")
			cmd.Printf("│      Day %-2s      ┃\n", cmd.Name())
			cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

			got := D4P1(d.PartOne.Input)
			cmd.Printf("Part 1: %s\n", got)

			got = D4P2(d.PartTwo.Input)
			cmd.Printf("Part 2: %s\n", got)
		},
	}

	Get2022Command().AddCommand(cmd)

	return cmd
}

// D4P1 returns the solution for 2022 day 4 part 1
// answer: 532
func D4P1(data []string) string {
	count := 0

	for _, line := range data {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isFullyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return strconv.Itoa(count)
}

type elf struct {
	low  int
	high int
}

func isFullyOverlapping(a, b *elf) bool {
	return (a.low <= b.low && a.high >= b.high) ||
		(b.low <= a.low && b.high >= a.high)
}

func parsePair(line string) (*elf, *elf) {
	var a, b elf

	_, err := fmt.Sscanf(line, "%d-%d,%d-%d", &a.low, &a.high, &b.low, &b.high)
	if err != nil {
		fmt.Printf("error parsing line: %s", line)
	}

	return &a, &b
}

// D4P2 returns the solution for 2022 day 4 part 2
// answer: 854
func D4P2(data []string) string {
	count := 0

	for _, line := range data {
		elfOne, elfTwo := parsePair(line)
		// log.Infof("elfOne: %v, elfTwo: %v\n", elfOne, elfTwo)

		if isAnyOverlapping(elfOne, elfTwo) {
			count++
		}
	}

	return strconv.Itoa(count)
}

func isAnyOverlapping(a, b *elf) bool {
	return (a.low <= b.low && a.high >= b.low) ||
		(b.low <= a.low && b.high >= a.low)
}
