package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 12.
type Exercise struct {
	common.BaseExercise
}

// parseGraph reads the pipe listing into an adjacency map.
func parseGraph(instr string) map[int][]int {
	graph := map[int][]int{}
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "<->", 2)
		id, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		for p := range strings.SplitSeq(parts[1], ",") {
			peer, _ := strconv.Atoi(strings.TrimSpace(p))
			graph[id] = append(graph[id], peer)
		}
	}
	return graph
}

// component returns every program reachable from start, marking them in seen.
func component(graph map[int][]int, start int, seen map[int]bool) int {
	size := 0
	stack := []int{start}
	seen[start] = true
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		size++
		for _, next := range graph[cur] {
			if !seen[next] {
				seen[next] = true
				stack = append(stack, next)
			}
		}
	}
	return size
}

// One returns the size of the group containing program 0.
func (e Exercise) One(instr string) (any, error) {
	graph := parseGraph(instr)
	return component(graph, 0, map[int]bool{}), nil
}

// Two returns the number of distinct groups.
func (e Exercise) Two(instr string) (any, error) {
	graph := parseGraph(instr)
	seen := map[int]bool{}
	groups := 0
	for id := range graph {
		if !seen[id] {
			component(graph, id, seen)
			groups++
		}
	}
	return groups, nil
}

// componentNodes returns the sorted list of nodes reachable from start.
type vec struct{ x, y float64 }

// layoutComponent runs a small Fruchterman-Reingold force simulation on the
// given nodes and returns their positions (centred on the origin) plus the
// component's bounding radius.

// mapToSlice converts an undirected adjacency-set map into the []int adjacency
// form componentNodes expects.
// hslHex renders an HSL colour as a hex string.
