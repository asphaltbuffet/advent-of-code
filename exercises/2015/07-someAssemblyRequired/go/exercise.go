package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 7.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	circuit, solved := parse(instr)

	answer := solve(circuit, solved, "a")

	if answer < 0 {
		answer += (1 << 16)
	}

	return answer, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}

type wire struct {
	eval   func(...int) int
	params []string
}

func parse(instr string) (map[string]wire, map[string]int) {
	lines := strings.Split(instr, "\n")

	circuit := make(map[string]wire, len(lines))
	solved := make(map[string]int, len(lines))

	for _, line := range lines {
		rawSignal, name, _ := strings.Cut(line, " -> ")

		n, err := strconv.Atoi(rawSignal)
		if err != nil {
			circuit[name] = newWire(rawSignal)
		} else {
			solved[name] = n
		}
	}

	return circuit, solved
}

func newWire(rawSignal string) wire {
	tokens := strings.Fields(rawSignal)

	if len(tokens) == 1 {
		return wire{eq, []string{tokens[0]}}
	}

	switch {
	case tokens[0] == "NOT":
		return wire{not, []string{tokens[1]}}
	case tokens[1] == "AND":
		return wire{and, []string{tokens[0], tokens[2]}}
	case tokens[1] == "OR":
		return wire{or, []string{tokens[0], tokens[2]}}
	case tokens[1] == "LSHIFT":
		return wire{lshift, []string{tokens[0], tokens[2]}}
	case tokens[1] == "RSHIFT":
		return wire{rshift, []string{tokens[0], tokens[2]}}
	}

	panic("unknown wire type")
}

func eq(n ...int) int {
	if len(n) != 1 {
		panic("eq requires 1 argument")
	}

	return n[0]
}

func or(n ...int) int {
	if len(n) != 2 {
		panic("or requires 2 arguments")
	}

	return n[0] | n[1]
}

func and(n ...int) int {
	if len(n) != 2 {
		panic("and requires 2 arguments")
	}

	return n[0] & n[1]
}

func lshift(n ...int) int {
	if len(n) != 2 {
		panic("lshift requires 2 arguments")
	}

	return n[0] << n[1]
}

func rshift(n ...int) int {
	if len(n) != 2 {
		panic("rshift requires 2 arguments")
	}

	return n[0] >> n[1]
}

func not(n ...int) int {
	if len(n) != 1 {
		panic("not requires 1 arguments")
	}

	return ^n[0]
}

func solve(circuit map[string]wire, solved map[string]int, name string) int {
	// fmt.Printf("solving [%s]: %s\n", name, circuit[name])
	if n, ok := solved[name]; ok {
		return n
	}

	w := circuit[name]
	params := make([]int, len(w.params))

	for i, p := range w.params {
		if n, ok := solved[p]; ok {
			params[i] = n
			continue
		}

		n, err := strconv.Atoi(p)
		if err == nil {
			params[i] = n
		} else {
			params[i] = solve(circuit, solved, p)
			solved[p] = params[i]
		}
	}

	// fmt.Printf("solved %s -> %d\n", name, w.eval(params...))

	return w.eval(params...)
}
