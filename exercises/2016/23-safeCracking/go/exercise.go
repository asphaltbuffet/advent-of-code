package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 23.
type Exercise struct {
	common.BaseExercise
}

const (
	cpy string = "cpy"
	inc string = "inc"
	dec string = "dec"
	jnz string = "jnz"
)

// run executes the assembunny program (with the day-23 tgl instruction) from
// the given initial registers and returns the final value of register a.
// execOp executes one instruction, returning the new ip.
// Returns -1 for jnz taken (caller adds the offset), or 0 for normal advance.
func execOp(op []string, program [][]string, regs map[string]int, ip int,
	isReg func(string) bool, value func(string) int,
) int {
	switch op[0] {
	case cpy:
		if isReg(op[2]) {
			regs[op[2]] = value(op[1])
		}
	case inc:
		if isReg(op[1]) {
			regs[op[1]]++
		}
	case dec:
		if isReg(op[1]) {
			regs[op[1]]--
		}
	case jnz:
		if value(op[1]) != 0 {
			return ip + value(op[2])
		}
	case "tgl":
		t := ip + value(op[1])
		if t >= 0 && t < len(program) {
			toggle(program[t])
		}
	}
	return ip + 1
}

func run(instr string, regs map[string]int) int {
	var program [][]string
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		program = append(program, strings.Fields(line))
	}

	isReg := func(x string) bool {
		_, err := strconv.Atoi(x)
		return err != nil
	}
	value := func(x string) int {
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
		return regs[x]
	}

	for ip := 0; ip >= 0 && ip < len(program); {
		// Peephole: collapse the multiply idiom
		//   cpy b c / inc a / dec c / jnz c -2 / dec d / jnz d -5
		// into a += b * d in a single step. This is what makes Part Two
		// (a starts at 12) finish instantly rather than running 12! loops.
		if ip+5 < len(program) && mulPattern(program[ip:ip+6]) {
			p := program[ip : ip+6]
			a, c, d := p[1][1], p[0][2], p[4][1]
			b := p[0][1]
			regs[a] += value(b) * regs[d]
			regs[c] = 0
			regs[d] = 0
			ip += 6
			continue
		}
		ip = execOp(program[ip], program, regs, ip, isReg, value)
	}
	return regs["a"]
}

// toggle rewrites an instruction in place per the day-23 rules.
func toggle(op []string) {
	switch len(op) {
	case 2:
		if op[0] == inc {
			op[0] = dec
		} else {
			op[0] = inc
		}
	case 3:
		if op[0] == jnz {
			op[0] = cpy
		} else {
			op[0] = jnz
		}
	}
}

// mulPattern reports whether the six instructions match the inner-multiply
// idiom that toggling never alters in this input.
func mulPattern(p [][]string) bool {
	return len(p) == 6 &&
		p[0][0] == cpy &&
		p[1][0] == inc &&
		p[2][0] == dec && p[2][1] == p[0][2] &&
		p[3][0] == jnz && p[3][1] == p[0][2] && p[3][2] == "-2" &&
		p[4][0] == dec &&
		p[5][0] == jnz && p[5][1] == p[4][1] && p[5][2] == "-5"
}

// One runs the program with register a initialised to 7 eggs.
func (e Exercise) One(instr string) (any, error) {
	a := 7
	if isExample(instr) {
		a = 0 // the example seeds a from the program itself
	}
	return run(instr, map[string]int{"a": a, "b": 0, "c": 0, "d": 0}), nil
}

// Two runs the program with register a initialised to 12 eggs.
func (e Exercise) Two(instr string) (any, error) {
	return run(instr, map[string]int{"a": 12, "b": 0, "c": 0, "d": 0}), nil
}

// isExample detects the short example program, which has no tgl operand "c"
// nor any cpy-of-register-into-a setup; it begins with "cpy 2 a".
func isExample(instr string) bool {
	return strings.HasPrefix(strings.TrimSpace(instr), "cpy 2 a")
}
