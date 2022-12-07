package aoc21

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 5, D2P1, D2P1, Get2021Command())
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
