package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise is the exercise for 2021 day 2
type Exercise struct {
	common.BaseExercise
}

type command struct {
	operation string
	value     int
}

func parse(in string) ([]command, error) {
	var commands []command

	for _, line := range strings.Split(in, "\n") {
		if line == "" {
			continue
		}

		l := strings.TrimSpace(line)

		op, val, _ := strings.Cut(l, " ")

		v, _ := strconv.Atoi(val)

		commands = append(commands, command{operation: op, value: v})
	}

	return commands, nil
}

// One returns the solution for 2021 day 2 part 1
func (c Exercise) One(instr string) (any, error) {
	commands, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commands from input: %w", err)
	}

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

	return fmt.Sprintf("%d", horiz*depth), nil
}

// Two returns the solution for 2021 day 2 part 2
func (c Exercise) Two(instr string) (any, error) {
	commands, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commands from input: %w", err)
	}

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

	return fmt.Sprintf("%d", horiz*depth), nil
}
