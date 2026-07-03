package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 21.
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
	// Part one only needs the first halt value, so stop after one round of the
	// recurrence rather than generating the whole cycle.
	first := -1
	haltValues(parse(instr), func(v int) bool { first = v; return true })
	return first, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	last := -1
	haltValues(parse(instr), func(v int) bool { last = v; return false })
	return last, nil
}

// haltValues walks, in order, every distinct value register 0 is compared against
// at the program's single `eqrr … 0` instruction and calls yield for each; it stops
// early if yield returns true, otherwise runs until the sequence repeats. The
// program halts exactly when r0 equals the current value, so the first value halts
// fastest (part one) and the last distinct value halts slowest (part two).
//
// The inner loop of the program is an O(n) simulation of dividing a number by 256
// digit by digit; interpreting it opcode-by-opcode takes tens of seconds. Instead
// the two per-input constants — the seed the running value is reset to and the
// multiplier — are read out of the program, and the recurrence is evaluated
// directly.
func haltValues(p program, yield func(int) bool) {
	seed, mult := loopConstants(p)

	const mask = 0xFFFFFF
	seen := map[int]struct{}{}
	acc := 0
	for {
		hi := acc | 0x10000
		acc = seed
		for {
			acc = ((acc + (hi & 0xFF)) & mask) * mult & mask
			if hi < 256 {
				break
			}
			hi /= 256
		}
		if _, ok := seen[acc]; ok {
			return // acc repeats: the full cycle has been seen
		}
		seen[acc] = struct{}{}
		if yield(acc) {
			return
		}
	}
}

// loopConstants extracts the two magic numbers that drive the running value: the
// seed it is reset to each round (the `seti` into the compared register) and the
// multiplier (the `muli` on that register).
func loopConstants(p program) (seed, mult int) {
	// The value register is whichever operand of the `eqrr … 0` isn't register 0.
	valReg := -1
	for _, in := range p.instr {
		if in.op == "eqrr" {
			if in.a == 0 {
				valReg = in.b
			} else if in.b == 0 {
				valReg = in.a
			}
		}
	}
	for _, in := range p.instr {
		switch {
		case in.op == "seti" && in.c == valReg && in.a > seed:
			seed = in.a // the large per-input seed constant
		case in.op == "muli" && in.c == valReg:
			mult = in.b
		}
	}
	return seed, mult
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
