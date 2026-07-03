package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 25.
type Exercise struct {
	common.BaseExercise
}

// emitsClock reports whether the assembunny program, run with register a set to
// the given seed, transmits an unending 0,1,0,1,... clock signal.
//
// "Unending" is decided by cycle detection: if the full machine state (ip plus
// registers plus expected next-bit) recurs while the signal has stayed correct,
// the loop will repeat that valid output forever.
func emitsClock(program [][]string, seed int) bool {
	regs := map[string]int{"a": seed, "b": 0, "c": 0, "d": 0}

	value := func(x string) int {
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
		return regs[x]
	}

	type state struct {
		ip, a, b, c, d, want int
	}
	seen := map[state]bool{}
	want := 0 // next expected output bit

	for ip := 0; ip >= 0 && ip < len(program); {
		op := program[ip]
		if op[0] == "out" {
			bit := value(op[1])
			if bit != want {
				return false
			}
			want ^= 1
			s := state{ip, regs["a"], regs["b"], regs["c"], regs["d"], want}
			if seen[s] {
				return true // state repeats with a valid signal so far
			}
			seen[s] = true
			ip++
			continue
		}

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
	return false
}

// One finds the lowest positive seed for register a that yields the clock.
func (e Exercise) One(instr string) (any, error) {
	var program [][]string
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		program = append(program, strings.Fields(line))
	}
	for seed := 1; ; seed++ {
		if emitsClock(program, seed) {
			return seed, nil
		}
	}
}

// Two has no puzzle: day 25 completes the year once every other star is earned.
func (e Exercise) Two(instr string) (any, error) {
	return "Merry Christmas!", nil
}
