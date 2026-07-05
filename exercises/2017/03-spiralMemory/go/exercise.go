package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 3.
type Exercise struct {
	common.BaseExercise
}

type coord struct{ x, y int }

// spiralWalk yields the grid coordinates of the Ulam-style spiral in order,
// starting at the origin (square 1) and turning counterclockwise. It calls
// yield for each square; yield returns false to stop.
func spiralWalk(yield func(coord) bool) {
	// Direction order: right, up, left, down (counterclockwise).
	dirs := []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	pos := coord{0, 0}
	if !yield(pos) {
		return
	}
	dir, runLen, sinceTurn, legsAtLen := 0, 1, 0, 0
	for {
		d := dirs[dir]
		pos.x += d.x
		pos.y += d.y
		if !yield(pos) {
			return
		}
		sinceTurn++
		if sinceTurn == runLen {
			sinceTurn = 0
			dir = (dir + 1) % 4
			// Run length grows every two legs: 1,1,2,2,3,3,...
			legsAtLen++
			if legsAtLen == 2 {
				legsAtLen = 0
				runLen++
			}
		}
	}
}

func parseTarget(instr string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(instr))
	return n
}

// One returns the Manhattan distance from square N back to square 1.
func (e Exercise) One(instr string) (any, error) {
	target := parseTarget(instr)
	dist := 0
	i := 1
	spiralWalk(func(c coord) bool {
		if i == target {
			dist = abs(c.x) + abs(c.y)
			return false
		}
		i++
		return true
	})
	return dist, nil
}

// Two returns the first value written that exceeds the puzzle input, where each
// square stores the sum of its already-filled (including diagonal) neighbours.
func (e Exercise) Two(instr string) (any, error) {
	target := parseTarget(instr)
	vals := map[coord]int{}
	result := 0
	first := true
	spiralWalk(func(c coord) bool {
		if first {
			vals[c] = 1
			first = false
			if 1 > target {
				result = 1
				return false
			}
			return true
		}
		sum := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				sum += vals[coord{c.x + dx, c.y + dy}]
			}
		}
		vals[c] = sum
		if sum > target {
			result = sum
			return false
		}
		return true
	})
	return result, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
