package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 2.
type Exercise struct {
	common.BaseExercise
}

// parseProgram parses a comma-separated list of integers.
func parseProgram(instr string) ([]int, error) {
	parts := strings.Split(strings.TrimSpace(instr), ",")
	mem := make([]int, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("parse error at position %d: %w", i, err)
		}
		mem[i] = v
	}
	return mem, nil
}

// runIntcode executes an Intcode program and returns the memory on halt.
func runIntcode(mem []int) error {
	for ip := 0; ip < len(mem); ip += 4 {
		switch mem[ip] {
		case 99:
			return nil
		case 1:
			mem[mem[ip+3]] = mem[mem[ip+1]] + mem[mem[ip+2]]
		case 2:
			mem[mem[ip+3]] = mem[mem[ip+1]] * mem[mem[ip+2]]
		default:
			return fmt.Errorf("unknown opcode %d at position %d", mem[ip], ip)
		}
	}
	return nil
}

// One returns the answer to the first part of the exercise.
// For the real input (>6 elements), the 1202 alarm fix is applied:
// position 1 = 12, position 2 = 2, per the puzzle instructions.
func (e Exercise) One(instr string) (any, error) {
	mem, err := parseProgram(instr)
	if err != nil {
		return nil, err
	}

	// Apply the "1202 program alarm" fix for real input only.
	if len(mem) > 6 {
		mem[1] = 12
		mem[2] = 2
	}

	if err = runIntcode(mem); err != nil {
		return nil, err
	}

	return mem[0], nil
}

// Two returns the answer to the second part of the exercise.
// It finds the noun/verb pair (each in [0,99]) such that running the
// Intcode program with mem[1]=noun and mem[2]=verb produces mem[0]==19690720,
// then returns 100*noun + verb.
func (e Exercise) Two(instr string) (any, error) {
	original, err := parseProgram(instr)
	if err != nil {
		return nil, err
	}

	const target = 19690720

	for noun := range 100 {
		for verb := range 100 {
			mem := make([]int, len(original))
			copy(mem, original)
			mem[1] = noun
			mem[2] = verb

			if err = runIntcode(mem); err != nil {
				continue
			}

			if mem[0] == target {
				return 100*noun + verb, nil
			}
		}
	}

	return nil, fmt.Errorf("no noun/verb pair produces %d", target)
}
