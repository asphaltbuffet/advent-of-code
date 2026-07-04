package exercises

import (
	"errors"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 20.
type Exercise struct {
	common.BaseExercise
}

type pos struct{ x, y int }

var dirs = []pos{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

func buildGrid(instr string) []string {
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	maxW := 0
	for _, row := range grid {
		if len(row) > maxW {
			maxW = len(row)
		}
	}
	for i, row := range grid {
		if len(row) < maxW {
			grid[i] = row + strings.Repeat(" ", maxW-len(row))
		}
	}
	return grid
}

func gridCell(grid []string, p pos) byte {
	if p.y < 0 || p.y >= len(grid) || p.x < 0 || p.x >= len(grid[p.y]) {
		return ' '
	}
	return grid[p.y][p.x]
}

func isLetter(b byte) bool { return b >= 'A' && b <= 'Z' }

// findPortalTile returns the '.' tile adjacent to a letter pair at (x,y)-(x2,y2).
// candidates are the two positions on either side of the pair.
func adjacentDot(grid []string, candidates []pos) (pos, bool) {
	for _, c := range candidates {
		if gridCell(grid, c) == '.' {
			return c, true
		}
	}
	return pos{}, false
}

func scanLabels(grid []string) map[string][]pos {
	labels := map[string][]pos{}
	for y := range grid {
		for x := range grid[y] {
			p := pos{x, y}
			b := gridCell(grid, p)
			if !isLetter(b) {
				continue
			}
			// Horizontal pair: current + right.
			if r := (pos{x + 1, y}); isLetter(gridCell(grid, r)) {
				label := string([]byte{b, gridCell(grid, r)})
				if tile, ok := adjacentDot(grid, []pos{{x - 1, y}, {x + 2, y}}); ok {
					labels[label] = append(labels[label], tile)
				}
			}
			// Vertical pair: current + below.
			if d := (pos{x, y + 1}); isLetter(gridCell(grid, d)) {
				label := string([]byte{b, gridCell(grid, d)})
				if tile, ok := adjacentDot(grid, []pos{{x, y - 1}, {x, y + 2}}); ok {
					labels[label] = append(labels[label], tile)
				}
			}
		}
	}
	return labels
}

func parseMaze(instr string) ([]string, map[pos]pos, pos, pos) {
	grid := buildGrid(instr)
	labels := scanLabels(grid)
	portals := map[pos]pos{}
	var start, end pos
	for label, tiles := range labels {
		switch label {
		case "AA":
			start = tiles[0]
		case "ZZ":
			end = tiles[0]
		default:
			if len(tiles) == 2 {
				portals[tiles[0]] = tiles[1]
				portals[tiles[1]] = tiles[0]
			}
		}
	}
	return grid, portals, start, end
}

func bfsMaze(grid []string, portals map[pos]pos, start, end pos) int {
	visited := map[pos]bool{start: true}
	queue := []pos{start}
	steps := 0
	for len(queue) > 0 {
		next := make([]pos, 0, len(queue))
		for _, cur := range queue {
			if cur == end {
				return steps
			}
			for _, d := range dirs {
				nb := pos{cur.x + d.x, cur.y + d.y}
				if gridCell(grid, nb) == '.' && !visited[nb] {
					visited[nb] = true
					next = append(next, nb)
				}
			}
			if partner, ok := portals[cur]; ok && !visited[partner] {
				visited[partner] = true
				next = append(next, partner)
			}
		}
		queue = next
		steps++
	}
	return -1
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	grid, portals, start, end := parseMaze(instr)
	result := bfsMaze(grid, portals, start, end)
	if result < 0 {
		return nil, errors.New("no path found")
	}
	return result, nil
}

type p2BFSState struct {
	p     pos
	level int
}

// mazeBounds returns the min/max x,y of all '.' and '#' tiles in grid.
func mazeBounds(grid []string) (int, int, int, int) {
	minX, minY := 1<<30, 1<<30
	maxX, maxY := 0, 0
	for y := range grid {
		for x := range grid[y] {
			c := grid[y][x]
			if c != '.' && c != '#' {
				continue
			}
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	return minX, minY, maxX, maxY
}

// outerChecker returns a function that reports whether a portal tile is on the outer ring.
func outerChecker(minX, minY, maxX, maxY int) func(pos) bool {
	return func(p pos) bool {
		return p.x <= minX+3 || p.x >= maxX-3 || p.y <= minY+3 || p.y >= maxY-3
	}
}

// portalJumpTarget computes the next state after a portal jump, or returns the
// zero state if the jump is blocked.
func portalJumpTarget(portals map[pos]pos, cur p2BFSState, isOuter func(pos) bool) p2BFSState {
	partner, ok := portals[cur.p]
	if !ok {
		return p2BFSState{}
	}
	const maxLevel = 50
	if isOuter(cur.p) {
		if cur.level > 0 {
			return p2BFSState{partner, cur.level - 1}
		}
		return p2BFSState{} // outer portal at level 0 is blocked
	}
	if cur.level < maxLevel {
		return p2BFSState{partner, cur.level + 1}
	}
	return p2BFSState{}
}

// bfsPart2 performs the recursive BFS for Part Two.
func bfsPart2(grid []string, portals map[pos]pos, start, end pos) int {
	minX, minY, maxX, maxY := mazeBounds(grid)
	isOuter := outerChecker(minX, minY, maxX, maxY)

	startState := p2BFSState{start, 0}
	endState := p2BFSState{end, 0}

	visited := map[p2BFSState]bool{startState: true}
	queue := []p2BFSState{startState}
	steps := 0

	for len(queue) > 0 {
		next := make([]p2BFSState, 0, len(queue))
		for _, cur := range queue {
			if cur == endState {
				return steps
			}
			next = expandP2Neighbors(grid, portals, cur, isOuter, visited, next)
		}
		queue = next
		steps++
	}
	return -1
}

// expandP2Neighbors appends all unvisited neighbors (cardinal + portal) of cur.
func expandP2Neighbors(
	grid []string,
	portals map[pos]pos,
	cur p2BFSState,
	isOuter func(pos) bool,
	visited map[p2BFSState]bool,
	next []p2BFSState,
) []p2BFSState {
	for _, d := range dirs {
		nb := pos{cur.p.x + d.x, cur.p.y + d.y}
		if gridCell(grid, nb) != '.' {
			continue
		}
		ns := p2BFSState{nb, cur.level}
		if !visited[ns] {
			visited[ns] = true
			next = append(next, ns)
		}
	}
	if ns := portalJumpTarget(portals, cur, isOuter); ns != (p2BFSState{}) && !visited[ns] {
		visited[ns] = true
		next = append(next, ns)
	}
	return next
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	grid, portals, start, end := parseMaze(instr)
	result := bfsPart2(grid, portals, start, end)
	if result < 0 {
		return nil, errors.New("no path found")
	}
	return result, nil
}
