package exercises

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 22.
type Exercise struct {
	common.BaseExercise
}

// Region types (also erosion % 3): rocky=0, wet=1, narrow=2.
// Tools: neither=0, torch=1, climbing gear=2. A region forbids the tool equal to
// its type, so the two allowed tools in region r are all tools except r.

type cave struct {
	depth        int
	tx, ty       int
	erosionCache map[[2]int]int
}

func (c *cave) erosion(x, y int) int {
	if v, ok := c.erosionCache[[2]int{x, y}]; ok {
		return v
	}
	var geo int
	switch {
	case x == 0 && y == 0, x == c.tx && y == c.ty:
		geo = 0
	case y == 0:
		geo = x * 16807
	case x == 0:
		geo = y * 48271
	default:
		geo = c.erosion(x-1, y) * c.erosion(x, y-1)
	}
	e := (geo + c.depth) % 20183
	c.erosionCache[[2]int{x, y}] = e
	return e
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	c := parse(instr)
	risk := 0
	for y := 0; y <= c.ty; y++ {
		for x := 0; x <= c.tx; x++ {
			risk += c.erosion(x, y) % 3
		}
	}
	return risk, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	c := parse(instr)
	return fastestRescue(c), nil
}

const torch = 1 // start and finish equipped with the torch

// margin is how far past the target the search grid extends. A shortest path only
// detours a little beyond the target (a 7-minute tool switch bounds any useful
// detour), so this is comfortably sufficient.
const margin = 50

// fastestRescue is a Dijkstra over (x, y, tool) on a region grid padded past the
// target. Moving to an adjacent region costs 1 minute; switching to the region's
// other allowed tool costs 7. A tool is allowed in a region iff it differs from the
// region's type. Start at (0,0) with the torch; reach the target with the torch.
func fastestRescue(c *cave) int {
	w, h := c.tx+margin, c.ty+margin

	// Precompute region types over the padded grid.
	region := make([][]int, h+1)
	for y := 0; y <= h; y++ {
		region[y] = make([]int, w+1)
		for x := 0; x <= w; x++ {
			region[y][x] = c.erosion(x, y) % 3
		}
	}

	// Flat distance array indexed by (y*(w+1)+x)*3 + tool.
	idx := func(x, y, tool int) int { return (y*(w+1)+x)*3 + tool }
	dist := make([]int, (w+1)*(h+1)*3)
	for i := range dist {
		dist[i] = 1 << 30
	}
	dist[idx(0, 0, torch)] = 0
	goal := idx(c.tx, c.ty, torch)

	pq := &priorityQueue{{x: 0, y: 0, tool: torch, cost: 0}}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		if cur.cost > dist[idx(cur.x, cur.y, cur.tool)] {
			continue
		}
		if idx(cur.x, cur.y, cur.tool) == goal {
			return cur.cost
		}

		// Switch to the other tool allowed in this region (7 minutes).
		for tool := 0; tool < 3; tool++ {
			if tool != region[cur.y][cur.x] && tool != cur.tool {
				relax(dist, pq, idx, cur.x, cur.y, tool, cur.cost+7)
			}
		}
		// Move to an adjacent region keeping the current tool (1 minute).
		for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 0 || ny < 0 || nx > w || ny > h {
				continue
			}
			if region[ny][nx] != cur.tool {
				relax(dist, pq, idx, nx, ny, cur.tool, cur.cost+1)
			}
		}
	}
	return -1
}

func relax(dist []int, pq *priorityQueue, idx func(x, y, tool int) int, x, y, tool, nd int) {
	if nd < dist[idx(x, y, tool)] {
		dist[idx(x, y, tool)] = nd
		heap.Push(pq, item{x: x, y: y, tool: tool, cost: nd})
	}
}

// parse reads the depth and target coordinates.
func parse(instr string) *cave {
	var depth, tx, ty int
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if strings.HasPrefix(line, "depth:") {
			fmt.Sscanf(line, "depth: %d", &depth)
		} else if strings.HasPrefix(line, "target:") {
			fmt.Sscanf(line, "target: %d,%d", &tx, &ty)
		}
	}
	return &cave{depth: depth, tx: tx, ty: ty, erosionCache: map[[2]int]int{}}
}

// item and priorityQueue implement a min-heap of states ordered by cost.
type item struct {
	x, y, tool int
	cost       int
}

type priorityQueue []item

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x any)        { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}
