package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 6.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	parent := make(map[string]string)
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ")", 2)
		parent[parts[1]] = parts[0]
	}

	// Memoized depth computation.
	depth := make(map[string]int)
	var getDepth func(node string) int
	getDepth = func(node string) int {
		if d, ok := depth[node]; ok {
			return d
		}
		p, ok := parent[node]
		if !ok {
			// This is COM (root).
			depth[node] = 0
			return 0
		}
		d := 1 + getDepth(p)
		depth[node] = d
		return d
	}

	total := 0
	for node := range parent {
		total += getDepth(node)
	}
	return total, nil
}

// Two returns the answer to the second part of the exercise.
// It finds the minimum orbital transfers between the objects YOU and SAN orbit.
func (e Exercise) Two(instr string) (any, error) {
	parent := make(map[string]string)
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ")", 2)
		parent[parts[1]] = parts[0]
	}

	// Collect ancestors of YOU with their distances from parent(YOU).
	youAncestors := make(map[string]int)
	dist := 0
	cur := parent["YOU"]
	for cur != "" {
		youAncestors[cur] = dist
		dist++
		cur = parent[cur]
	}

	// Walk ancestors of SAN upward; first common ancestor gives the answer.
	dist = 0
	cur = parent["SAN"]
	for cur != "" {
		if d, ok := youAncestors[cur]; ok {
			return d + dist, nil
		}
		dist++
		cur = parent[cur]
	}

	return -1, nil
}
