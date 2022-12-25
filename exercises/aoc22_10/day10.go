// Package aoc22_10 contains the solution for day 10 of Advent of Code 2022.
package aoc22_10 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Instruction is string type for instructions.
type Instruction string

// Instruction constants.
const (
	Add  Instruction = "addx"
	Noop Instruction = "noop"
)

// Cycle constants.
const (
	Cycle20  int = 20
	Cycle60  int = 60
	Cycle100 int = 100
	Cycle140 int = 140
	Cycle180 int = 180
	Cycle220 int = 220
)

// Command is a single instruction with a value.
type Command struct {
	Instruction Instruction
	Value       int
}

// Day10 is the exercise environment.
type Day10 struct {
	Cycle    int
	X        map[int]int
	Commands []Command
	Cycle20  int
	Cycle60  int
	Cycle100 int
	Cycle140 int
	Cycle180 int
	Cycle220 int
}

// D10P1 returns the solution for 2022 day 10 part 1.
//
// https://adventofcode.com/2022/day/10
//
// answer: 17380
func D10P1(data []string) string {
	part1 := Day10{
		Cycle:    1,
		X:        map[int]int{0: 0, 1: 1},
		Commands: []Command{},
		Cycle20:  1,
		Cycle60:  1,
		Cycle100: 1,
		Cycle140: 1,
		Cycle180: 1,
		Cycle220: 1,
	}

	err := part1.Parse(data)
	if err != nil {
		return fmt.Sprintf("parsing input: %v", err)
	}

	err = part1.Process()
	if err != nil {
		return fmt.Sprintf("processing input: %v", err)
	}

	c := part1.Calculate()

	return strconv.Itoa(c)
}

// Parse converts the input data into a tokenized form.
func (d *Day10) Parse(data []string) error {
	d.Commands = []Command{}

	for i, line := range data {
		left, right, _ := strings.Cut(line, " ")
		switch left {
		case string(Add):
			val, err := strconv.Atoi(right)
			if err != nil {
				return fmt.Errorf("parsing value on line %d: %w", i, err)
			}

			d.Commands = append(d.Commands, Command{Add, val})
		case string(Noop):
			d.Commands = append(d.Commands, Command{Noop, 0})
		default:
			return fmt.Errorf("unknown instruction on line %d: %s", i, left)
		}
	}

	// log.Printf("commands: %v", d.Commands)

	return nil
}

// Process executes the commands.
func (d *Day10) Process() error {
	for i, cmd := range d.Commands {
		switch {
		case d.Cycle+2 <= Cycle20:
			d.Cycle20 += cmd.Value
			fallthrough
		case d.Cycle+2 <= Cycle60:
			d.Cycle60 += cmd.Value
			fallthrough
		case d.Cycle+2 <= Cycle100:
			d.Cycle100 += cmd.Value
			fallthrough
		case d.Cycle+2 <= Cycle140:
			d.Cycle140 += cmd.Value
			fallthrough
		case d.Cycle+2 <= Cycle180:
			d.Cycle180 += cmd.Value
			fallthrough
		case d.Cycle+2 <= Cycle220:
			d.Cycle220 += cmd.Value
		}

		switch cmd.Instruction {
		case Add:
			d.X[d.Cycle+2] = cmd.Value
			d.Cycle += 2

		case Noop:
			d.Cycle++

		default:
			return fmt.Errorf("processing command %d: %+v", i, cmd)
		}
	}

	// log.Printf("temp X: %+v", d.X)

	// for i := 0; i < d.Cycle; i++ {
	// 	d.X[i+1] += d.X[i]
	// }

	// log.Printf("final X: %v", d.X)
	log.Printf("cycle 20: %d", d.Cycle20)
	log.Printf("cycle 60: %d", d.Cycle60)
	log.Printf("cycle 100: %d", d.Cycle100)
	log.Printf("cycle 140: %d", d.Cycle140)
	log.Printf("cycle 180: %d", d.Cycle180)
	log.Printf("cycle 220: %d", d.Cycle220)

	return nil
}

// Calculate returns the result of the processing.
func (d *Day10) Calculate() int {
	log.Printf("cycles: %d, %d, %d, %d, %d, %d",
		(Cycle20 * d.Cycle20),
		(Cycle60 * d.Cycle60),
		(Cycle100 * d.Cycle100),
		(Cycle140 * d.Cycle140),
		(Cycle180 * d.Cycle180),
		(Cycle220 * d.Cycle220))

	return (Cycle20 * d.Cycle20) + (Cycle60 * d.Cycle60) + (Cycle100 * d.Cycle100) + (Cycle140 * d.Cycle140) + (Cycle180 * d.Cycle180) + (Cycle220 * d.Cycle220)
}

// D10P2 returns the solution for 2022 day 10 part 2.
// answer: FGCUZREC
func D10P2(data []string) string {
	part1 := Day10{
		Cycle:    1,
		X:        map[int]int{0: 0, 1: 1},
		Commands: []Command{},
		Cycle20:  1,
		Cycle60:  1,
		Cycle100: 1,
		Cycle140: 1,
		Cycle180: 1,
		Cycle220: 1,
	}

	err := part1.Parse(data)
	if err != nil {
		return fmt.Sprintf("parsing input: %v", err)
	}

	err = part1.Process()
	if err != nil {
		return fmt.Sprintf("processing input: %v", err)
	}

	return part1.Draw()
}

// Draw returns a string showing which pixels are lit.
func (d *Day10) Draw() string {
	const (
		Lit  string = "#"
		Dark string = " "
	)

	sb := strings.Builder{}
	sb.WriteString("\n") // start on a new line

	x := 1

	for c := 1; c < 240; c++ {
		hpos := c % 40
		x += d.X[c]

		if abs(x-hpos) <= 1 {
			sb.WriteString(Lit)
		} else {
			sb.WriteString(Dark)
		}

		if hpos == 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
