package exercises

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

var (
	tracking map[string]int
	open     map[string]bool
)

type valve struct {
	flowRate    int
	connections []string
}

func parse(input string) map[string]valve {
	valves := map[string]valve{}

	for _, line := range strings.Split(input, "\n") {
		var name string

		parts := strings.Split(line, "; ")
		v := valve{}

		_, err := fmt.Sscanf(parts[0], "Valve %s has flow rate=%d", &name, &v.flowRate)
		if err != nil {
			panic("parsing valve name and flow rate" + err.Error())
		}

		connections := strings.Split(parts[1], ", ")

		if v.flowRate > 0 {
			nonZero = append(nonZero, name)
		}

		// update first entry to remove leading string
		connections[0] = connections[0][len(connections[0])-2:]
		v.connections = connections
		valves[name] = v
	}

	return valves
}

// Exercise for Advent of Code 2022 day 16.
type Exercise struct {
	common.BaseExercise
}

// Vis is the visualization of the exercise.
func (c Exercise) Vis(instr string, outdir string) error {
	sb := strings.Builder{}

	rooms := parse(instr)
	connections := map[string]bool{}

	sb.WriteString("flowchart TD\n")

	for k, v := range rooms {
		sb.WriteString(fmt.Sprintf("\t%s[%s: %d]\n", k, k, v.flowRate))

		for _, c := range v.connections {
			if _, ok := connections[fmt.Sprintf("%s --- %s", c, k)]; !ok {
				connections[fmt.Sprintf("%s --- %s", k, c)] = true
			}
		}
	}

	for k := range connections {
		sb.WriteString(fmt.Sprintf("\t%s\n", k))
	}

	err := os.WriteFile(outdir+"/vis.txt", []byte(sb.String()), 0o600)
	if err != nil {
		return fmt.Errorf("writing visualization to file: %w", err)
	}

	return nil
}

// One returns the answer to the first part of the exercise.
// wrong: 1575 (too low)
// wrong: 1566 (too low)
// wrong: 1589 (too low)
// answer: 1647
func (c Exercise) One(instr string) (any, error) {
	valves := parse(instr)

	// track best flow for each state
	tracking = map[string]int{}

	// track open valves separately
	open = map[string]bool{}

	// we can skip checking effect of opening a valve if the flow is zero
	for name, v := range valves {
		if v.flowRate == 0 {
			open[name] = true
		}
	}

	// start at AA and recurse through all paths, returning the best
	return bfs(valves, "AA", 30, 0), nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (c Exercise) Two(instr string) (any, error) {
	return nil, nil
}

func bfs(valves map[string]valve, curRoom string, timeLeft, curFlow int) int {
	if timeLeft == 0 {
		return 0
	}

	h := hash(curRoom, timeLeft, curFlow)

	// if we've already seen this state, return the best flow
	if v, ok := tracking[h]; ok {
		return v
	}

	maxFlow := 0

	// if the valve is closed, test result of opening it
	if !open[curRoom] {
		open[curRoom] = true

		newFlow := curFlow + valves[curRoom].flowRate

		maxFlow = curFlow + bfs(valves, curRoom, timeLeft-1, newFlow)

		// close it for other testing (hash will track open state history)
		open[curRoom] = false
	}

	// check adjacent valves (without opening this valve)
	for _, v := range valves[curRoom].connections {
		testFlow := curFlow + bfs(valves, v, timeLeft-1, curFlow)

		if testFlow > maxFlow {
			maxFlow = testFlow
		}
	}

	tracking[h] = maxFlow

	return maxFlow
}

func hash(valve string, timeLeft int, flow int) string {
	names := make([]string, 0, len(open))

	for k, isOpen := range open {
		if isOpen {
			names = append(names, k)
		}
	}

	sort.Strings(names)

	// switch to string builder to avoid creating intermediate strings
	var sb strings.Builder

	sb.WriteString(valve)
	sb.WriteString(fmt.Sprint(timeLeft))
	sb.WriteString(strings.Join(names, ""))
	sb.WriteString(fmt.Sprint(flow))

	// return fmt.Sprint(valve, timeLeft, names, flow)
	return sb.String()
}
