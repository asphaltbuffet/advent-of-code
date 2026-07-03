package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 20.
type Exercise struct {
	common.BaseExercise
}

type point struct{ x, y int }

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	dist := walk(instr)
	furthest := 0
	for _, d := range dist {
		if d > furthest {
			furthest = d
		}
	}
	return furthest, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	dist := walk(instr)
	count := 0
	for _, d := range dist {
		if d >= 1000 {
			count++
		}
	}
	return count, nil
}

// walk follows the route regex and returns the shortest door-distance to every
// room reachable from the origin. Because every step through a door adds exactly
// one, tracking the running distance and keeping the minimum on revisits yields
// the shortest distances directly — no separate BFS needed.
func walk(instr string) map[point]int {
	route := strings.TrimSpace(instr)
	dist := map[point]int{{0, 0}: 0}

	var stack []point // saved positions at each open paren
	pos := point{0, 0}

	for i := 0; i < len(route); i++ {
		switch route[i] {
		case '(':
			stack = append(stack, pos)
		case '|':
			// Reset to the branch's start for the next alternative.
			pos = stack[len(stack)-1]
		case ')':
			pos = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case 'N', 'S', 'E', 'W':
			prev := pos
			switch route[i] {
			case 'N':
				pos.y--
			case 'S':
				pos.y++
			case 'E':
				pos.x++
			case 'W':
				pos.x--
			}
			nd := dist[prev] + 1
			if d, ok := dist[pos]; !ok || nd < d {
				dist[pos] = nd
			}
		}
	}
	return dist
}
