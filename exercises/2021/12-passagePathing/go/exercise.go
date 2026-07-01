package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 12.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) (map[string][]string, error) {
	adj := map[string][]string{}

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		a, b, ok := strings.Cut(line, "-")
		if !ok {
			return nil, fmt.Errorf("edge missing `-` separator: %q", line)
		}

		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	return adj, nil
}

// isSmall reports whether a cave is small (lowercase), i.e. visit-limited.
func isSmall(cave string) bool {
	return cave == strings.ToLower(cave)
}

// countPaths does a DFS from cur to "end". visited tracks small caves already on
// the current path; canDouble is true while a single small cave may still be
// visited a second time (part two). Big caves are never restricted.
func countPaths(adj map[string][]string, cur string, visited map[string]bool, canDouble bool) int {
	if cur == "end" {
		return 1
	}

	total := 0
	for _, next := range adj[cur] {
		if next == "start" {
			continue // never return to the start
		}

		if isSmall(next) && visited[next] {
			if !canDouble {
				continue
			}
			// Spend the one allowed double-visit on this small cave.
			total += countPaths(adj, next, visited, false)
			continue
		}

		if isSmall(next) {
			visited[next] = true
			total += countPaths(adj, next, visited, canDouble)
			visited[next] = false
		} else {
			total += countPaths(adj, next, visited, canDouble)
		}
	}

	return total
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	adj, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	n := countPaths(adj, "start", map[string]bool{}, false)

	return fmt.Sprintf("%d", n), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	adj, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	n := countPaths(adj, "start", map[string]bool{}, true)

	return fmt.Sprintf("%d", n), nil
}
