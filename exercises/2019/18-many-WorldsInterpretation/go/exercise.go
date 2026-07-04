package exercises

import (
	"container/heap"
	"errors"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 18.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// Minimum steps to collect all keys in the maze.
func (e Exercise) One(instr string) (any, error) {
	steps, err := collectAllKeys(instr)
	if err != nil {
		return nil, err
	}

	return steps, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(_ string) (any, error) {
	return nil, errors.New("part 2 not implemented")
}

// edge represents a connection between two nodes in the key graph.
type edge struct {
	to      int    // compact index of destination key
	dist    int    // steps
	reqKeys uint32 // bitmask of door-keys required (bit i = key 'a'+i needed)
}

// dijkState is a state for the outer Dijkstra search.
type dijkState struct {
	cost    int
	node    int    // compact node index (numKeys = start position)
	keyMask uint32 // collected keys bitmask
}

// pq is a min-heap of dijkState ordered by cost.
type pq struct{ items []dijkState }

func (h *pq) Len() int           { return len(h.items) }
func (h *pq) Less(i, j int) bool { return h.items[i].cost < h.items[j].cost }
func (h *pq) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *pq) Push(x any) {
	h.items = append(h.items, x.(dijkState)) //nolint:errcheck // heap.Interface contract
}
func (h *pq) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}

// mazeGrid holds the parsed maze.
type mazeGrid struct {
	cells    []byte
	width    int
	startPos int
	keyPos   [26]int
	numKeys  int
}

// bfsCell is a BFS frontier cell for key-graph precomputation.
type bfsCell struct {
	pos     int
	reqKeys uint32 // doors encountered on the path from the source
}

func parseMaze(instr string) (mazeGrid, error) {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	if len(lines) == 0 {
		return mazeGrid{}, errors.New("empty input")
	}

	width := len(lines[0])
	cells := make([]byte, len(lines)*width)
	mg := mazeGrid{cells: cells, width: width, startPos: -1}

	for i := range mg.keyPos {
		mg.keyPos[i] = -1
	}

	for y, line := range lines {
		for x := 0; x < len(line) && x < width; x++ {
			c := line[x]
			pos := y*width + x
			cells[pos] = c

			switch {
			case c == '@':
				mg.startPos = pos
			case c >= 'a' && c <= 'z':
				mg.keyPos[c-'a'] = pos
				mg.numKeys++
			}
		}
	}

	if mg.startPos == -1 {
		return mazeGrid{}, errors.New("no start position found")
	}

	return mg, nil
}

// bfsEdges does a plain BFS from startPos (visiting each cell once) and returns
// all reachable keys with their distance and the door-keys required on the shortest path.
//
//nolint:gocognit // BFS with 4-direction expansion; complexity is inherent to the algorithm
func bfsEdges(startPos int, mg mazeGrid, compactOf [26]int) []edge {
	dirs := [4]int{-mg.width, mg.width, -1, 1}
	visited := make([]bool, len(mg.cells))
	visited[startPos] = true
	queue := []bfsCell{{pos: startPos}}
	var result []edge

	for dist := 1; len(queue) > 0; dist++ {
		next := make([]bfsCell, 0, len(queue)*2)

		for _, s := range queue {
			x := s.pos % mg.width

			for di, off := range dirs {
				if (di == 2 && x == 0) || (di == 3 && x == mg.width-1) {
					continue
				}

				np := s.pos + off
				if np < 0 || np >= len(mg.cells) || visited[np] {
					continue
				}

				c := mg.cells[np]
				if c == '#' {
					continue
				}

				visited[np] = true
				req := s.reqKeys

				if c >= 'A' && c <= 'Z' {
					req |= 1 << uint(c-'A')
				}

				next = append(next, bfsCell{pos: np, reqKeys: req})

				if c >= 'a' && c <= 'z' {
					if ci := compactOf[c-'a']; ci >= 0 {
						result = append(result, edge{to: ci, dist: dist, reqKeys: req})
					}
				}
			}
		}

		queue = next
	}

	return result
}

func collectAllKeys(instr string) (int, error) {
	mg, err := parseMaze(instr)
	if err != nil {
		return 0, err
	}

	activeKeys, compactOf := buildKeyIndex(mg.keyPos, mg.numKeys)
	graph := buildGraph(mg, activeKeys, compactOf)

	return dijkstraKeys(graph, activeKeys, uint32((1<<mg.numKeys)-1), mg.numKeys)
}

func buildKeyIndex(keyPos [26]int, numKeys int) ([]int, [26]int) {
	activeKeys := make([]int, 0, numKeys)

	for i, p := range keyPos {
		if p >= 0 {
			activeKeys = append(activeKeys, i)
		}
	}

	var compactOf [26]int
	for i := range compactOf {
		compactOf[i] = -1
	}

	for ci, ki := range activeKeys {
		compactOf[ki] = ci
	}

	return activeKeys, compactOf
}

func buildGraph(mg mazeGrid, activeKeys []int, compactOf [26]int) [][]edge {
	graph := make([][]edge, mg.numKeys+1)
	graph[mg.numKeys] = bfsEdges(mg.startPos, mg, compactOf) // start node = numKeys

	for ci, ki := range activeKeys {
		graph[ci] = bfsEdges(mg.keyPos[ki], mg, compactOf)
	}

	return graph
}

func dijkstraKeys(graph [][]edge, activeKeys []int, allKeys uint32, startNode int) (int, error) {
	type visitKey struct {
		node    int
		keyMask uint32
	}

	best := make(map[visitKey]int)
	h := &pq{items: []dijkState{{cost: 0, node: startNode, keyMask: 0}}}
	heap.Init(h)

	for h.Len() > 0 {
		cur := heap.Pop(h).(dijkState) //nolint:errcheck // heap.Interface contract

		if cur.keyMask == allKeys {
			return cur.cost, nil
		}

		vk := visitKey{node: cur.node, keyMask: cur.keyMask}
		if prev, ok := best[vk]; ok && prev <= cur.cost {
			continue
		}

		best[vk] = cur.cost

		for _, e := range graph[cur.node] {
			if cur.keyMask&e.reqKeys != e.reqKeys {
				continue
			}

			ki := activeKeys[e.to]

			if cur.keyMask&(1<<uint(ki)) != 0 {
				continue
			}

			newMask := cur.keyMask | (1 << uint(ki))
			newCost := cur.cost + e.dist
			nvk := visitKey{node: e.to, keyMask: newMask}

			if prev, ok := best[nvk]; ok && prev <= newCost {
				continue
			}

			heap.Push(h, dijkState{cost: newCost, node: e.to, keyMask: newMask})
		}
	}

	return 0, errors.New("no solution found")
}
