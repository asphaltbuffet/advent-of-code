package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 19.
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

// query runs the drone program for a single (x, y) coordinate and returns 1 if in beam, 0 if not.
func query(prog mem, x, y int) int {
	// Copy the program for a fresh run.
	m := make(mem, len(prog))
	for k, v := range prog {
		m[k] = v
	}

	in := make(chan int, 2)
	out := make(chan int, 1)

	in <- x
	in <- y

	go intcode(m, in, out)

	return <-out
}

// One returns the answer to the first part of the exercise.
// Counts how many points in the 50x50 grid are affected by the tractor beam.
func (e Exercise) One(instr string) (any, error) {
	prog := parseProgram(instr)

	count := 0

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			count += query(prog, x, y)
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	prog := parseProgram(instr)

	y := 100
	x := 0
	for {
		// Find left edge of beam at row y (beam edges only move right).
		for query(prog, x, y) == 0 {
			x++
		}
		// Check if a 100×100 square fits: top-right corner must be in beam.
		if query(prog, x+99, y-99) == 1 {
			return x*10000 + (y - 99), nil
		}
		y++
	}
}
