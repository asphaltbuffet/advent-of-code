package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 16.
type Exercise struct {
	common.BaseExercise
}

// regs holds the four registers of the little device.
type regs [4]int

// instruction is opcode, two inputs, and the output register.
type instruction [4]int

// operation names and their register effects. All sixteen produce the value to
// store in the output register from the current registers and the two inputs.
var operations = []struct {
	name string
	fn   func(r regs, a, b int) int
}{
	{"addr", func(r regs, a, b int) int { return r[a] + r[b] }},
	{"addi", func(r regs, a, b int) int { return r[a] + b }},
	{"mulr", func(r regs, a, b int) int { return r[a] * r[b] }},
	{"muli", func(r regs, a, b int) int { return r[a] * b }},
	{"banr", func(r regs, a, b int) int { return r[a] & r[b] }},
	{"bani", func(r regs, a, b int) int { return r[a] & b }},
	{"borr", func(r regs, a, b int) int { return r[a] | r[b] }},
	{"bori", func(r regs, a, b int) int { return r[a] | b }},
	{"setr", func(r regs, a, _ int) int { return r[a] }},
	{"seti", func(_ regs, a, _ int) int { return a }},
	{"gtir", func(r regs, a, b int) int { return boolToInt(a > r[b]) }},
	{"gtri", func(r regs, a, b int) int { return boolToInt(r[a] > b) }},
	{"gtrr", func(r regs, a, b int) int { return boolToInt(r[a] > r[b]) }},
	{"eqir", func(r regs, a, b int) int { return boolToInt(a == r[b]) }},
	{"eqri", func(r regs, a, b int) int { return boolToInt(r[a] == b) }},
	{"eqrr", func(r regs, a, b int) int { return boolToInt(r[a] == r[b]) }},
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// apply runs one operation and returns the resulting registers.
func apply(fn func(regs, int, int) int, r regs, ins instruction) regs {
	r[ins[3]] = fn(r, ins[1], ins[2])
	return r
}

func buildCandidates(samples []sample) []map[int]bool {
	candidates := make([]map[int]bool, 16)
	for i := range candidates {
		candidates[i] = map[int]bool{}
		for op := range operations {
			candidates[i][op] = true
		}
	}
	for _, s := range samples {
		opcode := s.ins[0]
		valid := map[int]bool{}
		for _, op := range s.matches() {
			valid[op] = true
		}
		for op := range candidates[opcode] {
			if !valid[op] {
				delete(candidates[opcode], op)
			}
		}
	}
	return candidates
}

func resolveOpcodes(samples []sample) []int {
	candidates := buildCandidates(samples)
	opFor := make([]int, 16)
	assigned := make([]bool, 16)
	for done := 0; done < 16; {
		for opcode := range 16 {
			if assigned[opcode] || len(candidates[opcode]) != 1 {
				continue
			}
			var op int
			for o := range candidates[opcode] {
				op = o
			}
			opFor[opcode] = op
			assigned[opcode] = true
			done++
			for other := range 16 {
				if other != opcode {
					delete(candidates[other], op)
				}
			}
		}
	}
	return opFor
}

// sample is one Before/instruction/After observation.
type sample struct {
	before regs
	ins    instruction
	after  regs
}

var numRe = regexp.MustCompile(`-?\d+`)

func nums(s string) []int {
	found := numRe.FindAllString(s, -1)
	out := make([]int, len(found))
	for i, f := range found {
		out[i], _ = strconv.Atoi(f)
	}
	return out
}

// parse splits the input into the observed samples and the test program.
func parse(instr string) ([]sample, []instruction) {
	blocks := strings.SplitN(strings.ReplaceAll(instr, "\r\n", "\n"), "\n\n\n", 2)

	var samples []sample
	for blk := range strings.SplitSeq(strings.TrimSpace(blocks[0]), "\n\n") {
		lines := strings.Split(strings.TrimSpace(blk), "\n")
		if len(lines) < 3 {
			continue
		}
		var s sample
		copy(s.before[:], nums(lines[0]))
		copy(s.ins[:], nums(lines[1]))
		copy(s.after[:], nums(lines[2]))
		samples = append(samples, s)
	}

	var program []instruction
	if len(blocks) > 1 {
		for line := range strings.SplitSeq(strings.TrimSpace(blocks[1]), "\n") {
			if line = strings.TrimSpace(line); line == "" {
				continue
			}
			var ins instruction
			copy(ins[:], nums(line))
			program = append(program, ins)
		}
	}

	return samples, program
}

// matches returns the set of operation indices whose effect reproduces the
// sample's After registers.
func (s sample) matches() []int {
	var ok []int
	for i, op := range operations {
		if apply(op.fn, s.before, s.ins) == s.after {
			ok = append(ok, i)
		}
	}
	return ok
}

// One returns the answer to the first part of the exercise.
// Answer: 547
func (e Exercise) One(instr string) (any, error) {
	samples, _ := parse(instr)

	count := 0
	for _, s := range samples {
		if len(s.matches()) >= 3 {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 582
func (e Exercise) Two(instr string) (any, error) {
	samples, program := parse(instr)
	opFor := resolveOpcodes(samples)

	var r regs
	for _, ins := range program {
		r = apply(operations[opFor[ins[0]]].fn, r, ins)
	}

	return r[0], nil
}
