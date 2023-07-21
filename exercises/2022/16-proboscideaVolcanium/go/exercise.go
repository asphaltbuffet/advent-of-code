package exercises

import (
	"fmt"
	"os"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

const MaxDistance = 999999

var (
	valves    map[string]valve
	distances map[string]map[string]int
	nonZero   []string
)

type Sequence struct {
	flow    int
	visited map[string]bool
}

type valve struct {
	flow        int
	connections []string
}

// Exercise for Advent of Code 2022 day 16.
type Exercise struct {
	common.BaseExercise
}

func parse(input string) (map[string]valve, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	valves = map[string]valve{}

	for _, line := range strings.Split(input, "\n") {
		var name string

		parts := strings.Split(line, "; ")
		v := valve{}

		_, err := fmt.Sscanf(parts[0], "Valve %s has flow rate=%d", &name, &v.flow)
		if err != nil {
			return nil, fmt.Errorf("parsing valve name and flow rate: %w", err)
		}

		connections := strings.Split(parts[1], ", ")

		if v.flow > 0 {
			nonZero = append(nonZero, name)
		}

		// update first entry to remove leading string
		connections[0] = connections[0][len(connections[0])-2:]
		v.connections = connections
		valves[name] = v
	}

	return valves, nil
}

// Vis is the visualization of the exercise.
func (e Exercise) Vis(instr string, outdir string) error {
	sb := strings.Builder{}

	rooms, err := parse(instr)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	connections := map[string]bool{}

	sb.WriteString("flowchart TD\n")

	for k, v := range rooms {
		sb.WriteString(fmt.Sprintf("\t%s[%s: %d]\n", k, k, v.flow))

		for _, c := range v.connections {
			if _, ok := connections[fmt.Sprintf("%s --- %s", c, k)]; !ok {
				connections[fmt.Sprintf("%s --- %s", k, c)] = true
			}
		}
	}

	for k := range connections {
		sb.WriteString(fmt.Sprintf("\t%s\n", k))
	}

	err = os.WriteFile(outdir+"/vis.mmd", []byte(sb.String()), 0o600)
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
func (e Exercise) One(instr string) (any, error) {
	var err error

	valves, err = parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	distances = floydWarshall(valves)

	sequences := getSequences("AA", 30, Sequence{0, make(map[string]bool)})

	max := 0

	// look at all possible sequences and find the best flow
	for i := 0; i < len(sequences); i++ {
		seq := sequences[i]

		if seq.flow > max {
			max = seq.flow
		}
	}

	return max, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	var err error

	valves, err = parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	distances = floydWarshall(valves)

	sequences := getSequences("AA", 26, Sequence{0, make(map[string]bool)})

	max := 0

	for i := 0; i < len(sequences)-1; i++ {
		me := sequences[i]
		if len(me.visited) == 0 {
			continue
		}

		for j := i + 1; j < len(sequences); j++ {
			elephant := sequences[j]
			if len(elephant.visited) == 0 {
				continue
			}

			combinedFlow := me.flow + elephant.flow

			if combinedFlow > max && hasNoOverlap(me.visited, elephant.visited) {
				max = combinedFlow
			}
		}
	}

	return max, nil
}

func hasNoOverlap(m1, m2 map[string]bool) bool {
	for key := range m1 {
		if m2[key] {
			return false
		}
	}
	return true
}

func (s Sequence) copy() Sequence {
	return Sequence{s.flow, copyMap(s.visited)}
}

func getSequences(curValve string, time int, curSeqence Sequence) []Sequence {
	sequences := []Sequence{curSeqence}

	// add all non-zero valves that haven't been visited
	for i := 0; i < len(nonZero); i++ {
		next := nonZero[i]
		newTime := time - distances[curValve][next] - 1

		if curSeqence.visited[next] || newTime <= 0 {
			continue
		}

		newSeq := Sequence{
			curSeqence.flow + newTime*valves[next].flow,
			copyMap(curSeqence.visited),
		}

		newSeq.visited[next] = true

		// get all sequences branching from this valve
		// add all the results to sequence list
		sequences = append(sequences, getSequences(next, newTime, newSeq)...)
	}

	return sequences
}

func floydWarshall(valves map[string]valve) map[string]map[string]int {
	dist := make(map[string]map[string]int, len(valves))

	// set initial distances
	for i := range valves {
		dist[i] = make(map[string]int)

		for j := range valves {
			switch {
			case i == j:
				dist[i][j] = 0
			case contains(valves[i].connections, j):
				dist[i][j] = 1
			default:
				dist[i][j] = MaxDistance
			}
		}
	}

	// run floyd-warshall
	for k := range valves {
		for i := range valves {
			for j := range valves {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}

func copyMap(m map[string]bool) map[string]bool {
	mcopy := make(map[string]bool)
	for k, v := range m {
		mcopy[k] = v
	}

	return mcopy
}
