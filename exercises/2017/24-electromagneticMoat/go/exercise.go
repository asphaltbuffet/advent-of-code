package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 24.
type Exercise struct {
	common.BaseExercise
}

type component struct{ a, b int }

func parseComponents(instr string) []component {
	var comps []component
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		comps = append(comps, component{a, b})
	}
	return comps
}

// build explores every bridge via backtracking, updating the strongest bridge
// and the strongest among the longest.
func build(comps []component, used []bool, port, strength, length int,
	bestStr *int, bestLen *[2]int) {

	// Record this bridge (the empty bridge scores 0, which never wins).
	if strength > *bestStr {
		*bestStr = strength
	}
	if length > bestLen[0] || (length == bestLen[0] && strength > bestLen[1]) {
		*bestLen = [2]int{length, strength}
	}

	for i, c := range comps {
		if used[i] {
			continue
		}
		var next int
		switch {
		case c.a == port:
			next = c.b
		case c.b == port:
			next = c.a
		default:
			continue
		}
		used[i] = true
		build(comps, used, next, strength+c.a+c.b, length+1, bestStr, bestLen)
		used[i] = false
	}
}

// solve returns the strongest bridge and the strongest of the longest bridges.
func solve(instr string) (int, int) {
	comps := parseComponents(instr)
	used := make([]bool, len(comps))
	bestStr := 0
	bestLen := [2]int{0, 0}
	build(comps, used, 0, 0, 0, &bestStr, &bestLen)
	return bestStr, bestLen[1]
}

// One returns the strength of the strongest bridge.
func (e Exercise) One(instr string) (any, error) {
	strongest, _ := solve(instr)
	return strongest, nil
}

// Two returns the strength of the strongest among the longest bridges.
func (e Exercise) Two(instr string) (any, error) {
	_, longestStrongest := solve(instr)
	return longestStrongest, nil
}
