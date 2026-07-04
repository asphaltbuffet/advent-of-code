package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 21.
type Exercise struct {
	common.BaseExercise
}

// mem is a map-backed unbounded memory for Intcode.
type mem map[int]int

func (m mem) get(addr int) int  { return m[addr] }
func (m mem) set(addr, val int) { m[addr] = val }

// parseProgram parses a comma-separated Intcode program into a mem map.
func parseProgram(input string) mem {
	m := make(mem)
	for i, s := range strings.Split(strings.TrimSpace(input), ",") {
		v, _ := strconv.Atoi(s)
		m[i] = v
	}

	return m
}

// intcode runs an Intcode computer in a goroutine, sending all output to the out channel.
//
//nolint:gocognit,funlen // Intcode VM opcode dispatch is inherently complex
func intcode(m mem, in <-chan int, out chan<- int) {
	ip := 0
	relBase := 0

	param := func(offset, mode int) int {
		raw := m.get(ip + offset)
		switch mode {
		case 0:
			return m.get(raw)
		case 1:
			return raw
		case 2:
			return m.get(relBase + raw)
		}
		panic("unknown mode")
	}

	dest := func(offset, mode int) int {
		raw := m.get(ip + offset)
		switch mode {
		case 0:
			return raw
		case 2:
			return relBase + raw
		}
		panic("unknown dest mode")
	}

	for {
		instr := m.get(ip)
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10
		m3 := (instr / 10000) % 10

		switch op {
		case 1:
			m.set(dest(3, m3), param(1, m1)+param(2, m2))
			ip += 4
		case 2:
			m.set(dest(3, m3), param(1, m1)*param(2, m2))
			ip += 4
		case 3:
			m.set(dest(1, m1), <-in)
			ip += 2
		case 4:
			out <- param(1, m1)
			ip += 2
		case 5:
			if param(1, m1) != 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 6:
			if param(1, m1) == 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 7:
			v := 0
			if param(1, m1) < param(2, m2) {
				v = 1
			}

			m.set(dest(3, m3), v)
			ip += 4
		case 8:
			v := 0
			if param(1, m1) == param(2, m2) {
				v = 1
			}

			m.set(dest(3, m3), v)
			ip += 4
		case 9:
			relBase += param(1, m1)
			ip += 2
		case 99:
			close(out)
			return
		}
	}
}

// runSpringscript sends an ASCII springscript program to the Intcode droid and
// returns the last output value (the hull damage reading, or the error code if
// the droid fell).
func runSpringscript(prog mem, script string) int {
	in := make(chan int, len(script))
	out := make(chan int, 4096)

	for _, ch := range script {
		in <- int(ch)
	}

	go intcode(prog, in, out)

	var last int

	for v := range out {
		last = v
	}

	return last
}

// One returns the answer to the first part of the exercise.
// Strategy: jump if there is a hole in A, B, or C (unsafe to walk through)
// AND D is solid (safe to land on).
// Logic: J = (NOT A OR NOT B OR NOT C) AND D
//
// Springscript:
//
//	NOT A J   -- J = hole at A
//	NOT B T   -- T = hole at B
//	OR  T J   -- J = hole at A or B
//	NOT C T   -- T = hole at C
//	OR  T J   -- J = hole at A, B, or C
//	AND D J   -- J = (hole ahead) AND D is solid
//	WALK
func (e Exercise) One(instr string) (any, error) {
	script := "NOT A J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J\nWALK\n"
	prog := parseProgram(instr)
	result := runSpringscript(prog, script)

	return result, nil
}

// Two returns the answer to the second part of the exercise.
// Strategy: jump if there's a hole ahead (A, B, or C) AND D is solid to land on,
// AND after landing we can either walk (E solid) OR jump again safely (H solid).
// Logic: J = (NOT A OR NOT B OR NOT C) AND D AND (E OR H)
//
// Springscript:
//
//	NOT A J   -- J = hole at A
//	NOT B T   -- T = hole at B
//	OR  T J   -- J = hole at A or B
//	NOT C T   -- T = hole at C
//	OR  T J   -- J = hole at A, B, or C
//	AND D J   -- J = (hole ahead) AND D solid
//	NOT E T   -- T = NOT E
//	NOT T T   -- T = E
//	OR  H T   -- T = E OR H
//	AND T J   -- J = J AND (E OR H)
//	RUN
func (e Exercise) Two(instr string) (any, error) {
	script := "NOT A J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J\nNOT E T\nNOT T T\nOR H T\nAND T J\nRUN\n"
	prog := parseProgram(instr)
	result := runSpringscript(prog, script)

	return result, nil
}
