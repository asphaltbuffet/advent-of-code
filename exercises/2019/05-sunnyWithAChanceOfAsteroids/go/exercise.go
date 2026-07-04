package exercises

import (
	"errors"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

var errNoOutput = errors.New("program produced no output")

// Exercise for Advent of Code 2019 day 5.
type Exercise struct {
	common.BaseExercise
}

// parseProgram parses a comma-separated Intcode program into a []int.
func parseProgram(instr string) []int {
	parts := strings.Split(strings.TrimSpace(instr), ",")
	mem := make([]int, len(parts))
	for i, p := range parts {
		n, _ := strconv.Atoi(strings.TrimSpace(p))
		mem[i] = n
	}
	return mem
}

// runIntcode runs the Intcode program with the given input queue and returns all outputs.
//
//nolint:gocognit // Intcode VM opcode dispatch is inherently complex
func runIntcode(prog []int, inputs []int) []int {
	// Copy memory so the original is not modified.
	mem := make([]int, len(prog))
	copy(mem, prog)

	ip := 0
	inIdx := 0
	var outputs []int

	param := func(offset, mode int) int {
		v := mem[ip+offset]
		if mode == 0 {
			return mem[v] // position mode
		}
		return v // immediate mode
	}

	for {
		instr := mem[ip]
		opcode := instr % 100
		modeC := (instr / 100) % 10
		modeB := (instr / 1000) % 10
		// modeA (writes) is always position mode; not needed.

		switch opcode {
		case 1: // add
			mem[mem[ip+3]] = param(1, modeC) + param(2, modeB)
			ip += 4
		case 2: // multiply
			mem[mem[ip+3]] = param(1, modeC) * param(2, modeB)
			ip += 4
		case 3: // input
			mem[mem[ip+1]] = inputs[inIdx]
			inIdx++
			ip += 2
		case 4: // output
			outputs = append(outputs, param(1, modeC))
			ip += 2
		case 5: // jump-if-true
			if param(1, modeC) != 0 {
				ip = param(2, modeB)
			} else {
				ip += 3
			}
		case 6: // jump-if-false
			if param(1, modeC) == 0 {
				ip = param(2, modeB)
			} else {
				ip += 3
			}
		case 7: // less-than
			if param(1, modeC) < param(2, modeB) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 8: // equals
			if param(1, modeC) == param(2, modeB) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 99: // halt
			return outputs
		default:
			// Unknown opcode — halt to avoid infinite loop.
			return outputs
		}
	}
}

// One returns the answer to the first part of the exercise.
// Runs the Intcode program with input=1 (AC unit ID) and returns the diagnostic code
// (last output value; all preceding outputs should be 0).
func (e Exercise) One(instr string) (any, error) {
	prog := parseProgram(instr)
	outputs := runIntcode(prog, []int{1})
	if len(outputs) == 0 {
		return nil, errNoOutput
	}
	return outputs[len(outputs)-1], nil
}

// Two returns the answer to the second part of the exercise.
// Runs the Intcode program with input=5 (thermal radiator controller ID) and returns the diagnostic code.
func (e Exercise) Two(instr string) (any, error) {
	prog := parseProgram(instr)
	outputs := runIntcode(prog, []int{5})
	if len(outputs) == 0 {
		return nil, errNoOutput
	}
	return outputs[len(outputs)-1], nil
}
