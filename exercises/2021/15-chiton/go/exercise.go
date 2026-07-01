package exercises

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 15.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) ([][]int, error) {
	var grid [][]int

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		row := make([]int, len(line))
		for i, c := range line {
			if c < '0' || c > '9' {
				return nil, fmt.Errorf("non-digit %q in risk grid", c)
			}
			row[i] = int(c - '0')
		}
		grid = append(grid, row)
	}

	return grid, nil
}

// expand tiles the grid 5x5; a cell in tile (tr,tc) has risk raised by tr+tc,
// wrapping from 9 back to 1.
func expand(grid [][]int, factor int) [][]int {
	rows, cols := len(grid), len(grid[0])
	big := make([][]int, rows*factor)
	for r := range big {
		big[r] = make([]int, cols*factor)
		for c := range big[r] {
			base := grid[r%rows][c%cols]
			add := r/rows + c/cols
			big[r][c] = (base+add-1)%9 + 1
		}
	}
	return big
}

type node struct {
	risk, r, c int
}

type pqueue []node

func (q pqueue) Len() int           { return len(q) }
func (q pqueue) Less(i, j int) bool { return q[i].risk < q[j].risk }
func (q pqueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *pqueue) Push(x any)        { *q = append(*q, x.(node)) }
func (q *pqueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}

var moves = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// lowestRisk runs Dijkstra from the top-left to the bottom-right, summing the
// risk of every entered cell (the start cell's risk is not counted).
func lowestRisk(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	dist := make([][]int, rows)
	for r := range dist {
		dist[r] = make([]int, cols)
		for c := range dist[r] {
			dist[r][c] = -1
		}
	}

	pq := &pqueue{}
	heap.Push(pq, node{0, 0, 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(node)
		if dist[cur.r][cur.c] != -1 {
			continue // already finalized with a lower risk
		}
		dist[cur.r][cur.c] = cur.risk

		if cur.r == rows-1 && cur.c == cols-1 {
			return cur.risk
		}

		for _, m := range moves {
			nr, nc := cur.r+m[0], cur.c+m[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols || dist[nr][nc] != -1 {
				continue
			}
			heap.Push(pq, node{cur.risk + grid[nr][nc], nr, nc})
		}
	}

	return -1
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	grid, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", lowestRisk(grid)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	grid, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	return fmt.Sprintf("%d", lowestRisk(expand(grid, 5))), nil
}
