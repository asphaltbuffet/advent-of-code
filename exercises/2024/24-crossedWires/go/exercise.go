package exercises

import (
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 24.
type Exercise struct {
	common.BaseExercise
}

// gate is a boolean operation producing a wire.
type gate struct {
	a, op, b string
}

// parse returns the initial wire values and the gate definitions by output wire.
func parse(instr string) (map[string]int, map[string]gate) {
	parts := strings.SplitN(strings.TrimRight(instr, "\n"), "\n\n", 2)

	values := map[string]int{}
	for _, line := range strings.Split(parts[0], "\n") {
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ": ", 2)
		v := 0
		if strings.TrimSpace(kv[1]) == "1" {
			v = 1
		}
		values[kv[0]] = v
	}

	gates := map[string]gate{}
	if len(parts) > 1 {
		for _, line := range strings.Split(parts[1], "\n") {
			if line == "" {
				continue
			}
			f := strings.Fields(line) // a OP b -> out
			gates[f[4]] = gate{f[0], f[1], f[2]}
		}
	}
	return values, gates
}

// eval resolves wire w via memoized recursion over the gate DAG.
func eval(w string, values map[string]int, gates map[string]gate, memo map[string]int) int {
	if v, ok := values[w]; ok {
		return v
	}
	if v, ok := memo[w]; ok {
		return v
	}
	g := gates[w]
	a := eval(g.a, values, gates, memo)
	b := eval(g.b, values, gates, memo)
	var r int
	switch g.op {
	case "AND":
		r = a & b
	case "OR":
		r = a | b
	case "XOR":
		r = a ^ b
	}
	memo[w] = r
	return r
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	values, gates := parse(instr)
	if len(gates) == 0 {
		return 0, nil
	}
	memo := map[string]int{}

	var result int
	for w := range gates {
		if w[0] != 'z' {
			continue
		}
		bit := (int(w[1]-'0'))*10 + int(w[2]-'0')
		if eval(w, values, gates, memo) == 1 {
			result |= 1 << uint(bit)
		}
	}
	return result, nil
}

// Two finds the eight wires whose gates violate ripple-carry-adder structure.
func (e Exercise) Two(instr string) (any, error) {
	_, gates := parse(instr)
	if len(gates) == 0 {
		return "", nil
	}

	// Highest z wire is the carry-out and is allowed to be an OR.
	maxZ := ""
	for w := range gates {
		if w[0] == 'z' && w > maxZ {
			maxZ = w
		}
	}

	isXY := func(s string) bool { return s != "" && (s[0] == 'x' || s[0] == 'y') }

	wrong := map[string]bool{}
	for out, g := range gates {
		switch {
		// A z output must come from XOR, except the final carry-out bit.
		case out[0] == 'z' && out != maxZ && g.op != "XOR":
			wrong[out] = true
		// A XOR not feeding z must combine x/y inputs (the half-sum gate).
		case g.op == "XOR" && out[0] != 'z' && !isXY(g.a) && !isXY(g.b):
			wrong[out] = true
		// An AND on real x/y bits (not bit 0) must feed an OR.
		case g.op == "AND" && isXY(g.a) && isXY(g.b) && g.a[1:] != "00":
			if !feedsOp(out, gates, "OR") {
				wrong[out] = true
			}
		// A XOR of x/y bits (not bit 0) must feed another XOR.
		case g.op == "XOR" && isXY(g.a) && isXY(g.b) && g.a[1:] != "00":
			if !feedsOp(out, gates, "XOR") {
				wrong[out] = true
			}
		}
	}

	var names []string
	for w := range wrong {
		names = append(names, w)
	}
	sort.Strings(names)
	return strings.Join(names, ","), nil
}

// feedsOp reports whether wire w is an input to at least one gate of the given op.
func feedsOp(w string, gates map[string]gate, op string) bool {
	for _, g := range gates {
		if g.op == op && (g.a == w || g.b == w) {
			return true
		}
	}
	return false
}
