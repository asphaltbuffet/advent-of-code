package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 14.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) (string, map[string]byte, error) {
	tmpl, rulesPart, ok := strings.Cut(strings.TrimSpace(instr), "\n\n")
	if !ok {
		return "", nil, fmt.Errorf("input missing blank-line separator")
	}
	tmpl = strings.TrimSpace(tmpl)

	rules := map[string]byte{}
	for _, line := range strings.Split(strings.TrimSpace(rulesPart), "\n") {
		pair, insert, ok := strings.Cut(strings.TrimSpace(line), " -> ")
		if !ok || len(pair) != 2 || len(insert) != 1 {
			return "", nil, fmt.Errorf("bad rule %q", line)
		}
		rules[pair] = insert[0]
	}

	return tmpl, rules, nil
}

// polymerize runs the pair-insertion process for the given number of steps,
// tracking counts of adjacent pairs instead of the (exponentially long) string,
// and returns the spread: most-common element count minus least-common.
func polymerize(tmpl string, rules map[string]byte, steps int) int {
	pairs := map[string]int{}
	for i := 0; i+1 < len(tmpl); i++ {
		pairs[tmpl[i:i+2]]++
	}

	for s := 0; s < steps; s++ {
		next := make(map[string]int, len(pairs))
		for pair, n := range pairs {
			if ins, ok := rules[pair]; ok {
				// AB -> C splits into AC and CB.
				next[string([]byte{pair[0], ins})] += n
				next[string([]byte{ins, pair[1]})] += n
			} else {
				next[pair] += n
			}
		}
		pairs = next
	}

	// Each element is the first char of exactly one pair; the template's last
	// character is never anyone's first char, so seed it once.
	counts := map[byte]int{tmpl[len(tmpl)-1]: 1}
	for pair, n := range pairs {
		counts[pair[0]] += n
	}

	minC, maxC := -1, 0
	for _, n := range counts {
		if minC < 0 || n < minC {
			minC = n
		}
		if n > maxC {
			maxC = n
		}
	}

	return maxC - minC
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	tmpl, rules, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", polymerize(tmpl, rules, 10)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	tmpl, rules, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", polymerize(tmpl, rules, 40)), nil
}
