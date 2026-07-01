package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 12
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// incorrect: 354
// answer: 330
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")

	_, startDist, _ := bfsFromEnd(data)
	if startDist < 0 {
		return nil, fmt.Errorf("no path from start to end")
	}

	return startDist, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 321
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")

	_, _, minFromLow := bfsFromEnd(data)
	if minFromLow < 0 {
		return nil, fmt.Errorf("no path from any low point to end")
	}

	return minFromLow, nil
}

// bfsFromEnd does a single breadth-first search outward from the end tile. The
// climb rule (a step may rise at most one level) is inverted: traveling
// backwards from the end, a move to a neighbor is legal when that neighbor is
// at most one level *below* the current tile. One pass yields the step count to
// the end for every reachable tile, so both parts read straight off it:
//   - startDist: distance from the 'S' tile (part one)
//   - minFromLow: minimum distance over all elevation-'a' tiles (part two)
// Any value is -1 when no path exists.
func bfsFromEnd(data []string) (dist [][]int, startDist, minFromLow int) {
	rows := len(data)
	cols := len(data[0])

	heights := make([][]int, rows)
	dist = make([][]int, rows)
	var endPos [2]int
	var startPos [2]int

	for r := 0; r < rows; r++ {
		heights[r] = make([]int, cols)
		dist[r] = make([]int, cols)
		for cc := 0; cc < cols; cc++ {
			ch := rune(data[r][cc])
			heights[r][cc] = GetHeight(ch)
			dist[r][cc] = -1
			switch ch {
			case end:
				endPos = [2]int{r, cc}
			case start:
				startPos = [2]int{r, cc}
			}
		}
	}

	dist[endPos[0]][endPos[1]] = 0
	queue := [][2]int{endPos}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		r, cc := cur[0], cur[1]

		for _, d := range dirs {
			nr, nc := r+d[0], cc+d[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols || dist[nr][nc] != -1 {
				continue
			}
			// Reverse climb rule: stepping from neighbor -> current is legal
			// forward when heights[cur] - heights[neighbor] <= 1.
			if heights[r][cc]-heights[nr][nc] > 1 {
				continue
			}
			dist[nr][nc] = dist[r][cc] + 1
			queue = append(queue, [2]int{nr, nc})
		}
	}

	startDist = dist[startPos[0]][startPos[1]]

	minFromLow = -1
	for r := 0; r < rows; r++ {
		for cc := 0; cc < cols; cc++ {
			if heights[r][cc] == 0 && dist[r][cc] >= 0 {
				if minFromLow < 0 || dist[r][cc] < minFromLow {
					minFromLow = dist[r][cc]
				}
			}
		}
	}

	return dist, startDist, minFromLow
}

const (
	lowest  = 'a'
	highest = 'z'
	start   = 'S'
	end     = 'E'
)

// GetHeight calculates the height of of a map location. a = 0 -> z = 25.
func GetHeight(c rune) int {
	switch c {
	case start:
		return 0 // lowest - lowest

	case end:
		return int(highest - lowest)

	default:
		return int(c - lowest)
	}
}

