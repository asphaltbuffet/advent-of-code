package aoc22

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	newDay3Command()
}

func newDay3Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "3",
		Short: "day 3 exercise for 2022 AoC",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := aoc.NewExercise(cmd.Parent().Name(), cmd.Name())
			if err != nil {
				fmt.Printf("creating new exercise: %v\n", err)
				return
			}

			color.Set(color.FgYellow)

			cmd.Printf("┌──────────────────┒\n")
			cmd.Printf("│      Day %-2s      ┃\n", cmd.Name())
			cmd.Printf("┕━━━━━━━━━━━━━━━━━━┛\n")

			color.Unset()

			got := D3P1(d.PartOne.Input)
			cmd.Printf("Part 1: %s\n", got)

			got = D3P2(d.PartTwo.Input)
			cmd.Printf("Part 2: %s\n", got)
		},
	}

	Get2022Command().AddCommand(cmd)
	cmd.GroupID = "days"

	return cmd
}

var priorityValue = map[byte]int{
	'a': 1, 'b': 2, 'c': 3, 'd': 4, 'e': 5, 'f': 6, 'g': 7, 'h': 8, 'i': 9, 'j': 10,
	'k': 11, 'l': 12, 'm': 13, 'n': 14, 'o': 15, 'p': 16, 'q': 17, 'r': 18, 's': 19, 't': 20,
	'u': 21, 'v': 22, 'w': 23, 'x': 24, 'y': 25, 'z': 26,

	'A': 27, 'B': 28, 'C': 29, 'D': 30, 'E': 31, 'F': 32, 'G': 33, 'H': 34, 'I': 35, 'J': 36,
	'K': 37, 'L': 38, 'M': 39, 'N': 40, 'O': 41, 'P': 42, 'Q': 43, 'R': 44, 'S': 45, 'T': 46,
	'U': 47, 'V': 48, 'W': 49, 'X': 50, 'Y': 51, 'Z': 52,
}

// D3P1 returns the solution for 2022 day 3 part 1
// answer:
func D3P1(data []string) string {
	score := 0

	for _, line := range data {
		score += scoreMispacked(line)
	}

	return strconv.Itoa(score)
}

// D3P2 returns the solution for 2022 day 3 part 2
// answer:
func D3P2(data []string) string {
	score := 0

	for i := 0; i < len(data); i += 3 {
		score += scoreBadges(data[i], data[i+1], data[i+2])
	}

	return strconv.Itoa(score)
}

func scoreMispacked(line string) int {
	compartmentOne := map[byte]bool{}

	for i := 0; i < len(line)/2; i++ {
		compartmentOne[line[i]] = true
	}

	// fmt.Printf("compartmentOne: %s\n", line[:len(line)/2])
	// fmt.Printf("compartmentTwo: %s\n", line[len(line)/2:])

	priority := 0

	for i := len(line) / 2; i < len(line); i++ {
		if _, ok := compartmentOne[line[i]]; ok {
			// fmt.Printf("found match: %q\n", line[i])
			priority += priorityValue[line[i]]

			// only count the first match
			delete(compartmentOne, line[i])
		}
	}

	return priority
}

func scoreBadges(a, b, c string) int {
	sharedItems := aoc.Unique([]byte(a))

	sharedItems = aoc.Filter(sharedItems, func(item byte) bool {
		return strings.Contains(b, string(item)) && strings.Contains(c, string(item))
	})

	// fmt.Printf("DEBUG: shared items: %s\n", sharedItems)

	return priorityValue[sharedItems[0]]
}
