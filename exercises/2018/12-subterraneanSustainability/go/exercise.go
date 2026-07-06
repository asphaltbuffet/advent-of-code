package exercises

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 12.
type Exercise struct {
	common.BaseExercise
}

// state is the set of pots that currently contain a plant, keyed by index.
type state map[int]bool

// parse reads the initial state and the set of neighborhoods that grow a plant.
func parse(instr string) (state, map[string]bool, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	if len(lines) == 0 {
		return nil, nil, errors.New("empty input")
	}

	pots := state{}
	initial := strings.TrimSpace(strings.TrimPrefix(lines[0], "initial state:"))
	for i, c := range initial {
		if c == '#' {
			pots[i] = true
		}
	}

	rules := map[string]bool{}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " => ")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("bad rule %q", line)
		}
		if parts[1] == "#" {
			rules[parts[0]] = true
		}
	}

	return pots, rules, nil
}

// bounds returns the lowest and highest planted pot indices.
func bounds(pots state) (int, int) {
	first := true
	lo, hi := 0, 0
	for i := range pots {
		if first || i < lo {
			lo = i
		}
		if first || i > hi {
			hi = i
		}
		first = false
	}
	return lo, hi
}

// step advances one generation. Each pot's next value depends on the five-pot
// neighborhood centered on it; only pots within two of a live pot can change.
func step(pots state, rules map[string]bool) state {
	lo, hi := bounds(pots)
	next := state{}

	for i := lo - 2; i <= hi+2; i++ {
		var nb strings.Builder
		for j := i - 2; j <= i+2; j++ {
			if pots[j] {
				nb.WriteByte('#')
			} else {
				nb.WriteByte('.')
			}
		}
		if rules[nb.String()] {
			next[i] = true
		}
	}

	return next
}

// sumIndices totals the indices of all planted pots.
func sumIndices(pots state) int {
	sum := 0
	for i := range pots {
		sum += i
	}
	return sum
}

// shape returns the planted pattern normalized to start at index 0, so two
// generations with the same shape differ only by a horizontal shift.
func shape(pots state) (string, int) {
	lo, hi := bounds(pots)
	var b strings.Builder
	for i := lo; i <= hi; i++ {
		if pots[i] {
			b.WriteByte('#')
		} else {
			b.WriteByte('.')
		}
	}
	return b.String(), lo
}

// One returns the answer to the first part of the exercise.
// Answer: 2281
func (e Exercise) One(instr string) (any, error) {
	pots, rules, err := parse(instr)
	if err != nil {
		return nil, err
	}

	for range 20 {
		pots = step(pots, rules)
	}

	return sumIndices(pots), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 2250000000120
func (e Exercise) Two(instr string) (any, error) {
	pots, rules, err := parse(instr)
	if err != nil {
		return nil, err
	}

	const target = 50_000_000_000

	// The pattern eventually settles into a fixed shape that just drifts sideways
	// by a constant offset each generation. Once a shape repeats, the planted
	// count is fixed, so the index sum grows linearly and we can extrapolate.
	seen := map[string]struct {
		gen, lo int
	}{}

	for gen := range target {
		sh, lo := shape(pots)
		if prev, ok := seen[sh]; ok {
			drift := lo - prev.lo    // indices shifted per (gen - prev.gen)
			period := gen - prev.gen // generations for that shift
			remaining := target - gen
			count := len(pots)
			return sumIndices(pots) + remaining*count*drift/period, nil
		}
		seen[sh] = struct{ gen, lo int }{gen, lo}
		pots = step(pots, rules)
	}

	return sumIndices(pots), nil
}
