package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 19.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// 260478 (too low)
func (e Exercise) One(instr string) (any, error) {
	rawWorkflows, rawParts, _ := strings.Cut(instr, "\n\n")
	wfData := strings.Split(rawWorkflows, "\n")
	pData := strings.Split(rawParts, "\n")

	// fmt.Printf("%d workflows, %d parts\n", len(wfData), len(pData))

	workflows := make(map[string]*Workflow)
	parts := make([]*Part, 0)

	for _, line := range wfData {
		w := parseWorkflow(line)
		workflows[w.Name] = w
	}

	for _, line := range pData {
		p := parsePart(line)
		parts = append(parts, p)
		// fmt.Printf("part %d: %v\n", i, p)
	}

	var sum int
	for _, p := range parts {
		if p.Process(workflows, "in") == Accepted {
			// fmt.Printf("part %d: Accepted\n", i)
			sum += p.Value
		}
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	rawWorkflows, _, _ := strings.Cut(instr, "\n\n")
	wfData := strings.Split(rawWorkflows, "\n")

	workflows := make(map[string]*Workflow)

	for _, line := range wfData {
		w := parseWorkflow(line)
		workflows[w.Name] = w
	}

	result2 := countCombinations(workflows, "in", []PartRange{{1, 4000}, {1, 4000}, {1, 4000}, {1, 4000}})

	return result2, nil
}
