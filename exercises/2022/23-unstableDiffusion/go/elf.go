package exercises

import (
	"crypto/sha1"
	"fmt"
	"math"
	"sort"
)

var cardinals = [][]point{
	// N
	{
		point{-1, -1}, // NW
		point{0, -1},  // N
		point{1, -1},  // NE
	},
	// S
	{
		point{-1, 1}, // SW
		point{0, 1},  // S
		point{1, 1},  // SE
	},
	// W
	{
		point{-1, -1}, // NW
		point{-1, 0},  // W
		point{-1, 1},  // SW
	},
	// E
	{
		point{1, -1}, // NE
		point{1, 0},  // E
		point{1, 1},  // SE
	},
}

var orthagonals = []point{
	{0, -1}, // N
	{0, 1},  // S
	{-1, 0}, // W
	{1, 0},  // E
}

func diffuse(elfLocations map[point]string, part int) (int, error) {
	startDirection := 0
	round := 1

	for (part == 1 && round <= 10) || part == 2 {
		plannedMoves, targetCounts := planElfMoves(elfLocations, startDirection)

		var moved bool
		elfLocations, moved = updateElfLocations(plannedMoves, targetCounts)
		startDirection++

		// Part two: the first round in which no elf moves is the answer. This is
		// cheaper and exact compared with hashing the whole board each round.
		if part == 2 && !moved {
			return round, nil
		}

		round++
	}

	minY, maxY := math.MaxInt16, math.MinInt16
	minX, maxX := math.MaxInt16, math.MinInt16

	for e := range elfLocations {
		if e.x < minX {
			minX = e.x
		}

		if e.x > maxX {
			maxX = e.x
		}

		if e.y < minY {
			minY = e.y
		}

		if e.y > maxY {
			maxY = e.y
		}
	}

	w := maxX - minX + 1
	h := maxY - minY + 1

	return w*h - len(elfLocations), nil
}

func hashState(elfLocations map[point]string) string {
	// ref: https://gist.github.com/anonymouse64/96435627839bbec0b87a2e518abaa9a7
	h := sha1.New() //nolint:gosec // using for hash, not crypto

	state := make([]string, len(elfLocations))
	i := 0

	for c := range elfLocations {
		state[i] = fmt.Sprintf("%v", c)
		i++
	}

	sort.Strings(state)

	for _, s := range state {
		b := sha1.Sum([]byte(fmt.Sprintf("%v", s))) //nolint:gosec // using for hash, not crypto
		h.Write(b[:])
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func planElfMoves(elfLocations map[point]string, startDirection int) (map[point]point, map[point]int) {
	plannedMoves := make(map[point]point)
	targetedCounts := make(map[point]int)

	for p, v := range elfLocations {
		if v == "#" {
			target := getMove(elfLocations, p, startDirection)
			plannedMoves[p] = target
			targetedCounts[target]++
		}
	}

	return plannedMoves, targetedCounts
}

func getMove(elfLocations map[point]string, p point, startDirection int) point {
	// if no neighbors, stay
	if !hasAdjascent(elfLocations, p) {
		return p
	}

	for i := 0; i < 4; i++ {
		directionIndex := (i + startDirection) % 4
		cardinal := cardinals[directionIndex]
		adjascentCount := 0

		for _, d := range cardinal {
			if elfLocations[p.add(d)] == "#" {
				adjascentCount++
			}
		}

		if adjascentCount == 0 {
			return p.add(orthagonals[directionIndex])
		}
	}

	// neighbors in all directions, stay
	return p
}

func hasAdjascent(elfLocations map[point]string, p point) bool {
	for _, cardinal := range cardinals {
		for _, d := range cardinal {
			if elfLocations[p.add(d)] == "#" {
				return true
			}
		}
	}

	return false
}

func updateElfLocations(plannedMoves map[point]point, targets map[point]int) (map[point]string, bool) {
	// reset map, but only if elves are not blocked...
	elfLocations := make(map[point]string)
	moved := false

	for start, end := range plannedMoves {
		if targets[end] > 1 {
			// stay
			elfLocations[start] = "#"
		} else {
			// move
			elfLocations[end] = "#"
			if end != start {
				moved = true
			}
		}
	}

	return elfLocations, moved
}
