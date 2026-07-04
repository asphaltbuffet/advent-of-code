package exercises

import (
	"container/heap"
	"errors"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// robotState encodes 4 robot positions (as node indices, int8) plus keyMask for the visited map.
type robotState struct {
	robots  [4]int8
	keyMask uint32
}

// dijkState4 is a state for the 4-robot Dijkstra search.
type dijkState4 struct {
	cost    int
	robots  [4]int8
	keyMask uint32
}

// pq4 is a min-heap of dijkState4 ordered by cost.
type pq4 struct{ items []dijkState4 }

func (h *pq4) Len() int           { return len(h.items) }
func (h *pq4) Less(i, j int) bool { return h.items[i].cost < h.items[j].cost }
func (h *pq4) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *pq4) Push(x any) {
	h.items = append(h.items, x.(dijkState4)) //nolint:errcheck // heap.Interface contract
}

func (h *pq4) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}

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
// Splits the maze into 4 quadrants (transforming the map if needed) and uses
// Dijkstra with state = ([4]robot positions, keysMask) to find the minimum steps.
func (e Exercise) Two(instr string) (any, error) {
	steps, err := collectAllKeysFourRobots(instr)
	if err != nil {
		return nil, err
	}

	return steps, nil
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

// transformToFourRobots replaces the 3x3 area around the single '@' with the
// four-robot pattern required by Part Two.
func transformToFourRobots(cells []byte, center, width int) {
	// center → wall
	cells[center] = '#'
	// cardinal neighbors → walls
	cells[center-width] = '#'
	cells[center+width] = '#'
	cells[center-1] = '#'
	cells[center+1] = '#'
	// diagonal neighbors → robots
	cells[center-width-1] = '@'
	cells[center-width+1] = '@'
	cells[center+width-1] = '@'
	cells[center+width+1] = '@'
}

// parseMazeFour parses the maze and returns a mazeGrid plus 4 start positions.
// If the input has a single '@', the map is transformed; if it already has 4 '@', they are used as-is.
func parseMazeFour(instr string) (mazeGrid, [4]int, error) {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	if len(lines) == 0 {
		return mazeGrid{}, [4]int{}, errors.New("empty input")
	}

	width := len(lines[0])
	cells := make([]byte, len(lines)*width)
	mg := mazeGrid{cells: cells, width: width, startPos: -1}

	for i := range mg.keyPos {
		mg.keyPos[i] = -1
	}

	var startPositions [4]int
	numStarts := 0

	for y, line := range lines {
		for x := 0; x < len(line) && x < width; x++ {
			c := line[x]
			pos := y*width + x
			cells[pos] = c

			switch {
			case c == '@':
				if numStarts < 4 {
					startPositions[numStarts] = pos
					numStarts++
				}

				mg.startPos = pos
			case c >= 'a' && c <= 'z':
				mg.keyPos[c-'a'] = pos
				mg.numKeys++
			}
		}
	}

	// Single start: apply the 3x3 transformation.
	if numStarts == 1 {
		center := mg.startPos
		transformToFourRobots(cells, center, width)
		startPositions[0] = center - width - 1
		startPositions[1] = center - width + 1
		startPositions[2] = center + width - 1
		startPositions[3] = center + width + 1
		numStarts = 4
	}

	if numStarts != 4 {
		return mazeGrid{}, [4]int{}, errors.New("expected 1 or 4 start positions")
	}

	return mg, startPositions, nil
}

// buildGraphFour builds the key graph for Part Two: numKeys+4 nodes total
// (indices 0..numKeys-1 are keys, numKeys..numKeys+3 are the 4 robot starts).
func buildGraphFour(mg mazeGrid, activeKeys []int, compactOf [26]int, startPositions [4]int) [][]edge {
	numKeys := len(activeKeys)
	graph := make([][]edge, numKeys+4)

	for i, pos := range startPositions {
		graph[numKeys+i] = bfsEdges(pos, mg, compactOf)
	}

	for ci, ki := range activeKeys {
		graph[ci] = bfsEdges(mg.keyPos[ki], mg, compactOf)
	}

	return graph
}

// dijkstraFourRobots runs Dijkstra with 4-robot state.
//
//nolint:gocognit // 4-robot Dijkstra; complexity is inherent to the algorithm
func dijkstraFourRobots(graph [][]edge, activeKeys []int, allKeys uint32, startNodes [4]int) (int, error) {
	numKeys := len(activeKeys)

	initRobots := [4]int8{
		int8(startNodes[0]),
		int8(startNodes[1]),
		int8(startNodes[2]),
		int8(startNodes[3]),
	}

	best := make(map[robotState]int)
	h := &pq4{items: []dijkState4{{cost: 0, robots: initRobots, keyMask: 0}}}
	heap.Init(h)

	for h.Len() > 0 {
		cur := heap.Pop(h).(dijkState4) //nolint:errcheck // heap.Interface contract

		if cur.keyMask == allKeys {
			return cur.cost, nil
		}

		rs := robotState{robots: cur.robots, keyMask: cur.keyMask}
		if prev, ok := best[rs]; ok && prev <= cur.cost {
			continue
		}

		best[rs] = cur.cost

		// Try moving each robot to each reachable key.
		for r := range cur.robots {
			node := int(cur.robots[r])

			for _, e := range graph[node] {
				// Check we have all required door-keys.
				if cur.keyMask&e.reqKeys != e.reqKeys {
					continue
				}

				ki := activeKeys[e.to]

				// Skip already-collected keys.
				if cur.keyMask&(1<<uint(ki)) != 0 {
					continue
				}

				newMask := cur.keyMask | (1 << uint(ki))
				newCost := cur.cost + e.dist
				newRobots := cur.robots
				newRobots[r] = int8(e.to)

				// Validate int8 range (node indices: 0..numKeys+3).
				if e.to >= numKeys+4 {
					continue
				}

				nrs := robotState{robots: newRobots, keyMask: newMask}
				if prev, ok := best[nrs]; ok && prev <= newCost {
					continue
				}

				heap.Push(h, dijkState4{cost: newCost, robots: newRobots, keyMask: newMask})
			}
		}
	}

	return 0, errors.New("no solution found")
}

func collectAllKeysFourRobots(instr string) (int, error) {
	mg, startPositions, err := parseMazeFour(instr)
	if err != nil {
		return 0, err
	}

	activeKeys, compactOf := buildKeyIndex(mg.keyPos, mg.numKeys)
	graph := buildGraphFour(mg, activeKeys, compactOf, startPositions)

	// Start node indices: numKeys, numKeys+1, numKeys+2, numKeys+3
	numKeys := len(activeKeys)
	startNodes := [4]int{numKeys, numKeys + 1, numKeys + 2, numKeys + 3}
	allKeys := uint32((1 << mg.numKeys) - 1)

	return dijkstraFourRobots(graph, activeKeys, allKeys, startNodes)
}
