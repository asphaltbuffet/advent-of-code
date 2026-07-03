package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 19.
type Exercise struct {
	common.BaseExercise
}

type instruction struct {
	op      string
	a, b, c int
}

type program struct {
	ipReg int
	instr []instruction
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	p := parse(instr)
	regs := run(p, [6]int{}, -1)
	return regs[0], nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	p := parse(instr)

	// With register 0 = 1 the program builds a large number in one register, then
	// sums its divisors with an O(n²) loop — far too slow to simulate. Run just
	// long enough for that number to be assembled, then compute the divisor sum
	// directly. The setup finishes before the ip returns to the small main loop.
	regs := run(p, [6]int{1, 0, 0, 0, 0, 0}, 1000)
	n := max6(regs)
	return sumDivisors(n), nil
}

// run executes the program from the given starting registers. If maxSteps >= 0 it
// stops after that many instructions; otherwise it runs to completion (ip out of
// range).
func run(p program, regs [6]int, maxSteps int) [6]int {
	ip := 0
	for step := 0; ip >= 0 && ip < len(p.instr); step++ {
		if maxSteps >= 0 && step >= maxSteps {
			break
		}
		regs[p.ipReg] = ip
		regs = apply(p.instr[ip], regs)
		ip = regs[p.ipReg] + 1
	}
	return regs
}

// apply executes a single instruction against the registers.
func apply(in instruction, r [6]int) [6]int {
	a, b, c := in.a, in.b, in.c
	switch in.op {
	case "addr":
		r[c] = r[a] + r[b]
	case "addi":
		r[c] = r[a] + b
	case "mulr":
		r[c] = r[a] * r[b]
	case "muli":
		r[c] = r[a] * b
	case "banr":
		r[c] = r[a] & r[b]
	case "bani":
		r[c] = r[a] & b
	case "borr":
		r[c] = r[a] | r[b]
	case "bori":
		r[c] = r[a] | b
	case "setr":
		r[c] = r[a]
	case "seti":
		r[c] = a
	case "gtir":
		r[c] = boolToInt(a > r[b])
	case "gtri":
		r[c] = boolToInt(r[a] > b)
	case "gtrr":
		r[c] = boolToInt(r[a] > r[b])
	case "eqir":
		r[c] = boolToInt(a == r[b])
	case "eqri":
		r[c] = boolToInt(r[a] == b)
	case "eqrr":
		r[c] = boolToInt(r[a] == r[b])
	}
	return r
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// sumDivisors returns the sum of all positive divisors of n (what the assembly
// laboriously computes).
func sumDivisors(n int) int {
	sum := 0
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			sum += d
			if d != n/d {
				sum += n / d
			}
		}
	}
	return sum
}

func max6(r [6]int) int {
	m := r[0]
	for _, v := range r[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// parse reads the "#ip" directive and the instruction list.
func parse(instr string) program {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	p := program{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "#ip" {
			p.ipReg, _ = strconv.Atoi(fields[1])
			continue
		}
		a, _ := strconv.Atoi(fields[1])
		b, _ := strconv.Atoi(fields[2])
		c, _ := strconv.Atoi(fields[3])
		p.instr = append(p.instr, instruction{fields[0], a, b, c})
	}
	return p
}
