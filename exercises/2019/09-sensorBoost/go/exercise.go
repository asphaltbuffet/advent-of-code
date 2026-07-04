package exercises

import (
	"errors"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

var errNoOutput = errors.New("program produced no output")

// Exercise for Advent of Code 2019 day 9.
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

// runIntcode runs the Intcode program in m with the given inputs and returns all outputs.
//
//nolint:gocognit,funlen // Intcode VM opcode dispatch is inherently complex
func runIntcode(m mem, inputs []int) []int {
	ip := 0
	relBase := 0
	inputIdx := 0
	var outputs []int

	param := func(offset, mode int) int {
		raw := m.get(ip + offset)
		switch mode {
		case 0: // position
			return m.get(raw)
		case 1: // immediate
			return raw
		case 2: // relative
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
		case 1: // add
			m.set(dest(3, m3), param(1, m1)+param(2, m2))
			ip += 4
		case 2: // multiply
			m.set(dest(3, m3), param(1, m1)*param(2, m2))
			ip += 4
		case 3: // input
			m.set(dest(1, m1), inputs[inputIdx])
			inputIdx++
			ip += 2
		case 4: // output
			outputs = append(outputs, param(1, m1))
			ip += 2
		case 5: // jump-if-true
			if param(1, m1) != 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 6: // jump-if-false
			if param(1, m1) == 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 7: // less than
			v := 0
			if param(1, m1) < param(2, m2) {
				v = 1
			}
			m.set(dest(3, m3), v)
			ip += 4
		case 8: // equals
			v := 0
			if param(1, m1) == param(2, m2) {
				v = 1
			}
			m.set(dest(3, m3), v)
			ip += 4
		case 9: // adjust relative base
			relBase += param(1, m1)
			ip += 2
		case 99: // halt
			return outputs
		}
	}
}

// One returns the answer to the first part of the exercise.
// Runs the Intcode program with input 1 and returns the single BOOST keycode output.
func (e Exercise) One(instr string) (any, error) {
	m := parseProgram(instr)
	outputs := runIntcode(m, []int{1})
	if len(outputs) == 0 {
		return nil, errNoOutput
	}
	return outputs[len(outputs)-1], nil
}

// Two returns the answer to the second part of the exercise.
// Runs the Intcode program with input 2 and returns the distress signal coordinates.
func (e Exercise) Two(instr string) (any, error) {
	m := parseProgram(instr)
	outputs := runIntcode(m, []int{2})
	if len(outputs) == 0 {
		return nil, errNoOutput
	}
	return outputs[len(outputs)-1], nil
}
