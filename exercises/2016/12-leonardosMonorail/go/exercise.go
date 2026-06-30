package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 12.
type Exercise struct {
	common.BaseExercise
}

// run executes the assembunny program with the given initial registers and
// returns the final value of register a.
func run(instr string, regs map[string]int) int {
	var program [][]string
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		program = append(program, strings.Fields(line))
	}

	// value resolves an operand: a literal integer or a register's value.
	value := func(x string) int {
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
		return regs[x]
	}

	for ip := 0; ip >= 0 && ip < len(program); {
		op := program[ip]
		switch op[0] {
		case "cpy":
			regs[op[2]] = value(op[1])
		case "inc":
			regs[op[1]]++
		case "dec":
			regs[op[1]]--
		case "jnz":
			if value(op[1]) != 0 {
				ip += value(op[2])
				continue
			}
		}
		ip++
	}
	return regs["a"]
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return run(instr, map[string]int{"a": 0, "b": 0, "c": 0, "d": 0}), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return run(instr, map[string]int{"a": 0, "b": 0, "c": 1, "d": 0}), nil
}
