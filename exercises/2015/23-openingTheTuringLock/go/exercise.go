package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 23.
type Exercise struct {
	common.BaseExercise
}

// instr23 is one decoded instruction: opcode, an optional register (0=a, 1=b),
// and an optional jump offset.
type instr23 struct {
	op     string
	reg    int
	offset int
}

// parse decodes the program. jmp uses only offset; jie/jio use reg and offset;
// hlf/tpl/inc use only reg.
func parse(input string) []instr23 {
	var prog []instr23

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		fields := strings.Fields(strings.ReplaceAll(line, ",", ""))
		ins := instr23{op: fields[0]}
		switch fields[0] {
		case "jmp":
			ins.offset, _ = strconv.Atoi(fields[1])
		case "jie", "jio":
			ins.reg = regIndex(fields[1])
			ins.offset, _ = strconv.Atoi(fields[2])
		default: // hlf, tpl, inc
			ins.reg = regIndex(fields[1])
		}
		prog = append(prog, ins)
	}

	return prog
}

func regIndex(s string) int {
	if s == "b" {
		return 1
	}
	return 0
}

// run executes the program with the given starting value of register a and
// returns the final value of register b. A program counter outside the program
// halts execution.
func run(prog []instr23, startA int) int {
	regs := [2]int{startA, 0}
	pc := 0

	for pc >= 0 && pc < len(prog) {
		ins := prog[pc]
		switch ins.op {
		case "hlf":
			regs[ins.reg] /= 2
			pc++
		case "tpl":
			regs[ins.reg] *= 3
			pc++
		case "inc":
			regs[ins.reg]++
			pc++
		case "jmp":
			pc += ins.offset
		case "jie": // jump if even
			if regs[ins.reg]%2 == 0 {
				pc += ins.offset
			} else {
				pc++
			}
		case "jio": // jump if one (not "if odd")
			if regs[ins.reg] == 1 {
				pc += ins.offset
			} else {
				pc++
			}
		}
	}

	return regs[1]
}

// One returns register b after running with a starting at 0.
func (e Exercise) One(instr string) (any, error) {
	return run(parse(instr), 0), nil
}

// Two returns register b after running with a starting at 1.
func (e Exercise) Two(instr string) (any, error) {
	return run(parse(instr), 1), nil
}
