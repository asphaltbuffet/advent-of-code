package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 7.
type Exercise struct {
	common.BaseExercise
}

type program struct {
	weight   int
	children []string
}

// parseTower reads the program listing into a name->program map.
func parseTower(instr string) map[string]program {
	tower := map[string]program{}
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// name (weight) [-> c1, c2, ...]
		name := strings.Fields(line)[0]
		open := strings.IndexByte(line, '(')
		closeIdx := strings.IndexByte(line, ')')
		weight, _ := strconv.Atoi(line[open+1 : closeIdx])

		var children []string
		if _, after, ok := strings.Cut(line, "->"); ok {
			for c := range strings.SplitSeq(after, ",") {
				children = append(children, strings.TrimSpace(c))
			}
		}
		tower[name] = program{weight, children}
	}
	return tower
}

// root returns the one program that is never listed as anyone's child.
func root(tower map[string]program) string {
	isChild := map[string]bool{}
	for _, p := range tower {
		for _, c := range p.children {
			isChild[c] = true
		}
	}
	for name := range tower {
		if !isChild[name] {
			return name
		}
	}
	return ""
}

// One returns the name of the bottom program.
func (e Exercise) One(instr string) (any, error) {
	return root(parseTower(instr)), nil
}

// totalWeight returns the combined weight of a program's sub-tower.
func totalWeight(tower map[string]program, name string) int {
	p := tower[name]
	sum := p.weight
	for _, c := range p.children {
		sum += totalWeight(tower, c)
	}
	return sum
}

// balance descends from name looking for the single wrong-weight program. want
// is the total weight this sub-tower is expected to have. It returns the
// corrected own-weight for the offending program and true once found.
func balance(tower map[string]program, name string, want int) (int, bool) {
	p := tower[name]

	// Group children by their total weight.
	totals := make([]int, len(p.children))
	counts := map[int]int{}
	for i, c := range p.children {
		totals[i] = totalWeight(tower, c)
		counts[totals[i]]++
	}

	// If a child's total is the minority, the imbalance lives inside it.
	if len(counts) > 1 {
		var oddTotal, goodTotal int
		for t, n := range counts {
			if n == 1 {
				oddTotal = t
			} else {
				goodTotal = t
			}
		}
		for i, c := range p.children {
			if totals[i] == oddTotal {
				// The odd child should weigh goodTotal in total; recurse in
				// case its own children are further unbalanced.
				return balance(tower, c, goodTotal)
			}
		}
	}

	// Children balance (or there are none): this program is the culprit. Its
	// corrected own-weight makes the sub-tower hit want.
	childSum := 0
	for _, t := range totals {
		childSum += t
	}
	return want - childSum, true
}

// Two returns the corrected weight for the single unbalanced program.
func (e Exercise) Two(instr string) (any, error) {
	tower := parseTower(instr)
	fixed, _ := balance(tower, root(tower), 0)
	return fixed, nil
}

// culpritPath returns the chain of names from the root down to the single
// unbalanced program (inclusive).
func culpritPath(tower map[string]program, name string) []string {
	p := tower[name]
	totals := make([]int, len(p.children))
	counts := map[int]int{}
	for i, c := range p.children {
		totals[i] = totalWeight(tower, c)
		counts[totals[i]]++
	}
	if len(counts) > 1 {
		var oddTotal int
		for t, n := range counts {
			if n == 1 {
				oddTotal = t
			}
		}
		for i, c := range p.children {
			if totals[i] == oddTotal {
				return append([]string{name}, culpritPath(tower, c)...)
			}
		}
	}
	return []string{name} // balanced children: this node is the culprit
}
