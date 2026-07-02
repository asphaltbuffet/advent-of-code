package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 8.
type Exercise struct {
	common.BaseExercise
}

// instruction is one operation of the handheld's boot code.
type instruction struct {
	op  string
	arg int
}

func parse(instr string) ([]instruction, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	prog := make([]instruction, 0, len(lines))
	for _, line := range lines {
		op, arg, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok {
			return nil, fmt.Errorf("bad instruction %q", line)
		}
		n, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("bad argument in %q: %w", line, err)
		}
		prog = append(prog, instruction{op, n})
	}
	return prog, nil
}

// run executes the program until it either loops (an instruction is about to run
// a second time) or the instruction pointer runs off the end. It returns the
// accumulator and whether the program terminated normally.
func run(prog []instruction) (acc int, terminated bool) {
	seen := make([]bool, len(prog))
	pc := 0
	for pc < len(prog) {
		if seen[pc] {
			return acc, false
		}
		seen[pc] = true
		in := prog[pc]
		switch in.op {
		case "acc":
			acc += in.arg
			pc++
		case "jmp":
			pc += in.arg
		default: // nop
			pc++
		}
	}
	return acc, true
}

// One returns the accumulator value just before any instruction runs twice.
func (e Exercise) One(instr string) (any, error) {
	prog, err := parse(instr)
	if err != nil {
		return nil, err
	}
	acc, _ := run(prog)
	return fmt.Sprintf("%d", acc), nil
}

// Two flips one jmp<->nop so the program terminates and returns its accumulator.
func (e Exercise) Two(instr string) (any, error) {
	prog, err := parse(instr)
	if err != nil {
		return nil, err
	}

	for i := range prog {
		orig := prog[i].op
		switch orig {
		case "jmp":
			prog[i].op = "nop"
		case "nop":
			prog[i].op = "jmp"
		default:
			continue
		}
		if acc, ok := run(prog); ok {
			return fmt.Sprintf("%d", acc), nil
		}
		prog[i].op = orig // revert and try the next candidate
	}

	return nil, fmt.Errorf("no single flip terminates the program")
}
