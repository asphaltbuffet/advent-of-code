package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 23.
type Exercise struct {
	common.BaseExercise
}

func parseProgram(instr string) [][]string {
	var prog [][]string
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			prog = append(prog, strings.Fields(line))
		}
	}
	return prog
}

// One runs the program literally with a=0 and counts mul invocations.
func (e Exercise) One(instr string) (any, error) {
	prog := parseProgram(instr)
	regs := map[string]int{}
	val := func(x string) int {
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
		return regs[x]
	}

	muls := 0
	for ip := 0; ip >= 0 && ip < len(prog); {
		op := prog[ip]
		switch op[0] {
		case "set":
			regs[op[1]] = val(op[2])
		case "sub":
			regs[op[1]] -= val(op[2])
		case "mul":
			regs[op[1]] *= val(op[2])
			muls++
		case "jnz":
			if val(op[1]) != 0 {
				ip += val(op[2])
				continue
			}
		}
		ip++
	}
	return muls, nil
}

// Two decompiles the program: with a=1 it counts the composite numbers in the
// range the setup computes, stepping by the loop's increment. Rather than run
// billions of instructions, we run only the setup to recover b and c, then
// count non-primes directly.
func (e Exercise) Two(instr string) (any, error) {
	prog := parseProgram(instr)

	// Execute the program with a=1 until the first big back-jump (the main
	// loop, "jnz 1 -N"), so registers b and c hold the range bounds.
	regs := map[string]int{"a": 1}
	val := func(x string) int {
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
		return regs[x]
	}
	for ip := 0; ip >= 0 && ip < len(prog); {
		op := prog[ip]
		// The linear setup (which fixes b and c) contains only forward jumps;
		// the first backward jump marks the start of the counting loops.
		if op[0] == "jnz" {
			if n, err := strconv.Atoi(op[2]); err == nil && n < 0 {
				break
			}
		}
		switch op[0] {
		case "set":
			regs[op[1]] = val(op[2])
		case "sub":
			regs[op[1]] -= val(op[2])
		case "mul":
			regs[op[1]] *= val(op[2])
		case "jnz":
			if val(op[1]) != 0 {
				ip += val(op[2])
				continue
			}
		}
		ip++
	}

	b, c := regs["b"], regs["c"]
	step := loopStep(prog)

	h := 0
	for n := b; n <= c; n += step {
		if !isPrime(n) {
			h++
		}
	}
	return h, nil
}

// loopStep reads the increment applied to b at the end of each outer loop: the
// last "sub b -N" instruction (an earlier one scales b during setup), defaulting
// to 17.
func loopStep(prog [][]string) int {
	step := 17
	for _, op := range prog {
		if op[0] == "sub" && op[1] == "b" {
			if n, err := strconv.Atoi(op[2]); err == nil && n < 0 {
				step = -n
			}
		}
	}
	return step
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}
