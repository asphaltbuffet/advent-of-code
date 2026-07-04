package exercises

import (
	"maps"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 15.
type Exercise struct {
	common.BaseExercise
}

// mem is a map-backed unbounded memory for Intcode.
type mem map[int]int

func (m mem) get(addr int) int  { return m[addr] }
func (m mem) set(addr, val int) { m[addr] = val }

// parseProgram parses a comma-separated Intcode program into a mem map.
func parseProgram(input string) mem {
	m := make(mem)
	for i, s := range strings.Split(strings.TrimSpace(input), ",") {
		v, _ := strconv.Atoi(s)
		m[i] = v
	}

	return m
}

// vm holds the full Intcode VM state so it can be driven synchronously.
type vm struct {
	m       mem
	ip      int
	relBase int
}

func newVM(prog mem) *vm {
	m := make(mem, len(prog))
	maps.Copy(m, prog)

	return &vm{m: m}
}

func (v *vm) param(offset, mode int) int {
	raw := v.m.get(v.ip + offset)

	switch mode {
	case 0:
		return v.m.get(raw)
	case 1:
		return raw
	case 2:
		return v.m.get(v.relBase + raw)
	}

	panic("unknown mode")
}

func (v *vm) dest(offset, mode int) int {
	raw := v.m.get(v.ip + offset)

	switch mode {
	case 0:
		return raw
	case 2:
		return v.relBase + raw
	}

	panic("unknown dest mode")
}

// step runs the VM until it produces output.
// send is the value to provide when input is requested.
// Returns the output value.
func (v *vm) step(send int) int { //nolint:cyclop // Intcode VM opcode dispatch is inherently complex
	for {
		instr := v.m.get(v.ip)
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10
		m3 := (instr / 10000) % 10

		switch op {
		case 1:
			v.m.set(v.dest(3, m3), v.param(1, m1)+v.param(2, m2))
			v.ip += 4
		case 2:
			v.m.set(v.dest(3, m3), v.param(1, m1)*v.param(2, m2))
			v.ip += 4
		case 3:
			v.m.set(v.dest(1, m1), send)
			v.ip += 2
		case 4:
			out := v.param(1, m1)
			v.ip += 2

			return out
		case 5:
			if v.param(1, m1) != 0 {
				v.ip = v.param(2, m2)
			} else {
				v.ip += 3
			}
		case 6:
			if v.param(1, m1) == 0 {
				v.ip = v.param(2, m2)
			} else {
				v.ip += 3
			}
		case 7:
			val := 0
			if v.param(1, m1) < v.param(2, m2) {
				val = 1
			}

			v.m.set(v.dest(3, m3), val)
			v.ip += 4
		case 8:
			val := 0
			if v.param(1, m1) == v.param(2, m2) {
				val = 1
			}

			v.m.set(v.dest(3, m3), val)
			v.ip += 4
		case 9:
			v.relBase += v.param(1, m1)
			v.ip += 2
		case 99:
			return -1
		}
	}
}

type point struct{ x, y int }

// directions: 1=N, 2=S, 3=W, 4=E
var dirs = [4]int{1, 2, 3, 4}

var delta = map[int]point{
	1: {0, -1},
	2: {0, 1},
	3: {-1, 0},
	4: {1, 0},
}

var opposite = map[int]int{1: 2, 2: 1, 3: 4, 4: 3}

const (
	cellWall   = 0
	cellOpen   = 1
	cellOxygen = 2
)

// exploreMaze uses DFS with backtracking to map the entire reachable area.
// Returns the grid map and the position of the oxygen system.
func exploreMaze(prog mem) (map[point]int, point) {
	grid := make(map[point]int)
	grid[point{0, 0}] = cellOpen

	machine := newVM(prog)
	cur := point{0, 0}
	oxygenPos := point{}

	var dfs func()
	dfs = func() {
		for _, dir := range dirs {
			d := delta[dir]
			next := point{cur.x + d.x, cur.y + d.y}

			if _, seen := grid[next]; seen {
				continue
			}

			// Try to move there.
			resp := machine.step(dir)
			grid[next] = resp

			if resp == cellWall {
				// Didn't move; wall recorded.
				continue
			}

			if resp == cellOxygen {
				oxygenPos = next
			}

			// Moved successfully; recurse.
			prev := cur
			cur = next
			dfs()
			cur = prev

			// Backtrack.
			machine.step(opposite[dir])
		}
	}

	dfs()

	return grid, oxygenPos
}

// bfsDistance returns the shortest path distance from start to target in grid.
func bfsDistance(grid map[point]int, start, target point) int {
	type state struct {
		pos  point
		dist int
	}

	visited := map[point]bool{start: true}
	queue := []state{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.pos == target {
			return cur.dist
		}

		for _, dir := range dirs {
			d := delta[dir]
			next := point{cur.pos.x + d.x, cur.pos.y + d.y}

			if visited[next] {
				continue
			}

			cell, known := grid[next]
			if !known || cell == cellWall {
				continue
			}

			visited[next] = true
			queue = append(queue, state{next, cur.dist + 1})
		}
	}

	return -1 // not reachable
}

// One returns the answer to the first part of the exercise.
// Explores the maze via DFS+backtracking, then BFS to find shortest path to oxygen.
func (e Exercise) One(instr string) (any, error) {
	grid, oxygenPos := exploreMaze(parseProgram(instr))
	dist := bfsDistance(grid, point{0, 0}, oxygenPos)

	return dist, nil
}

// Two returns the answer to the second part of the exercise.
// Explores the maze via DFS+backtracking, then BFS flood-fill from the oxygen
// system to find the maximum distance to any reachable open cell.
func (e Exercise) Two(instr string) (any, error) {
	grid, oxygenPos := exploreMaze(parseProgram(instr))

	type state struct {
		pos  point
		dist int
	}

	visited := map[point]bool{oxygenPos: true}
	queue := []state{{oxygenPos, 0}}
	maxDist := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.dist > maxDist {
			maxDist = cur.dist
		}

		for _, dir := range dirs {
			d := delta[dir]
			next := point{cur.pos.x + d.x, cur.pos.y + d.y}

			if visited[next] {
				continue
			}

			cell, known := grid[next]
			if !known || cell == cellWall {
				continue
			}

			visited[next] = true
			queue = append(queue, state{next, cur.dist + 1})
		}
	}

	return maxDist, nil
}
