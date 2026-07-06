package exercises

import (
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 15.
type Exercise struct {
	common.BaseExercise
}

// unit is an elf or goblin combatant.
type unit struct {
	x, y  int
	kind  byte // 'E' or 'G'
	hp    int
	alive bool
}

// reading order iterates neighbors up, left, right, down so a breadth-first search
// naturally discovers the reading-order-first shortest path.
var readingDirs = [4][2]int{{0, -1}, {-1, 0}, {1, 0}, {0, 1}}

// parseCave reads the map, lifting units off it (their squares become open floor).
func parseCave(instr string) ([][]byte, []*unit) {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	grid := make([][]byte, len(lines))
	var units []*unit

	for y, line := range lines {
		grid[y] = []byte(line)
		for x := range len(grid[y]) {
			if c := grid[y][x]; c == 'E' || c == 'G' {
				units = append(units, &unit{x: x, y: y, kind: c, hp: 200, alive: true})
				grid[y][x] = '.'
			}
		}
	}

	return grid, units
}

// point packs a coordinate for map keys.
type point struct{ x, y int }

// occupied maps every live unit's square to that unit.
func occupied(units []*unit) map[point]*unit {
	occ := make(map[point]*unit, len(units))
	for _, u := range units {
		if u.alive {
			occ[point{u.x, u.y}] = u
		}
	}
	return occ
}

// combat runs the battle. elfAP is the elves' attack power (goblins always hit for
// 3). When stopOnElfDeath is set, the fight aborts (returning -1) the moment an elf
// dies, which drives the part-two search. Otherwise it returns the outcome:
// completed rounds × remaining total HP.
func combat(instr string, elfAP int, stopOnElfDeath bool) int {
	grid, units := parseCave(instr)
	rounds := 0
	for {
		sortUnits(units)
		if done, result := combatRound(grid, units, elfAP, stopOnElfDeath, rounds); done {
			return result
		}
		rounds++
	}
}

// adjacentEnemy returns the weakest enemy adjacent to (x, y), preferring lower HP
// then reading order, or nil if none is adjacent.
func adjacentEnemy(units []*unit, x, y int, kind byte) *unit {
	var best *unit
	for _, dir := range readingDirs {
		nx, ny := x+dir[0], y+dir[1]
		for _, e := range units {
			if !e.alive || e.kind == kind || e.x != nx || e.y != ny {
				continue
			}
			if best == nil || e.hp < best.hp {
				best = e
			}
		}
	}
	return best
}

// stepToward finds the first step of the shortest path to the nearest square in
// range of an enemy, with all ties broken by reading order. It returns the step
// and whether a move is possible.
func stepToward(grid [][]byte, occ map[point]*unit, units []*unit, u *unit) (int, int, bool) {
	inRange := reachableTargets(grid, occ, units, u)
	if len(inRange) == 0 {
		return 0, 0, false
	}

	start := point{u.x, u.y}
	dist := map[point]int{start: 0}
	firstStep := map[point]point{}
	queue := []point{start}

	var chosen point
	found := false
	bestDist := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if found && dist[cur] > bestDist {
			break
		}
		if inRange[cur] && cur != start {
			if !found {
				found, bestDist, chosen = true, dist[cur], cur
			}
		}

		queue = expandNeighbors(queue, cur, start, grid, occ, dist, firstStep)
	}

	if !found {
		return 0, 0, false
	}
	fs := firstStep[chosen]
	return fs.x, fs.y, true
}

// open reports whether (x, y) is floor and unoccupied.
func open(grid [][]byte, occ map[point]*unit, x, y int) bool {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return false
	}
	if grid[y][x] != '.' {
		return false
	}
	_, taken := occ[point{x, y}]
	return !taken
}

func anyEnemyAlive(units []*unit, kind byte) bool {
	for _, e := range units {
		if e.alive && e.kind != kind {
			return true
		}
	}
	return false
}

func sumHP(units []*unit) int {
	total := 0
	for _, e := range units {
		if e.alive {
			total += e.hp
		}
	}
	return total
}

func reachableTargets(grid [][]byte, occ map[point]*unit, units []*unit, u *unit) map[point]bool {
	inRange := map[point]bool{}
	for _, e := range units {
		if !e.alive || e.kind == u.kind {
			continue
		}
		for _, dir := range readingDirs {
			nx, ny := e.x+dir[0], e.y+dir[1]
			if open(grid, occ, nx, ny) {
				inRange[point{nx, ny}] = true
			}
		}
	}
	return inRange
}

func expandNeighbors(
	queue []point, cur, start point, grid [][]byte, occ map[point]*unit, dist map[point]int, firstStep map[point]point,
) []point {
	for _, dir := range readingDirs {
		nx, ny := cur.x+dir[0], cur.y+dir[1]
		np := point{nx, ny}
		if !open(grid, occ, nx, ny) {
			continue
		}
		if _, seen := dist[np]; seen {
			continue
		}
		dist[np] = dist[cur] + 1
		if cur == start {
			firstStep[np] = np
		} else {
			firstStep[np] = firstStep[cur]
		}
		queue = append(queue, np)
	}
	return queue
}

func sortUnits(units []*unit) {
	sort.Slice(units, func(i, j int) bool {
		if units[i].y != units[j].y {
			return units[i].y < units[j].y
		}
		return units[i].x < units[j].x
	})
}

func combatRound(grid [][]byte, units []*unit, elfAP int, stopOnElfDeath bool, rounds int) (bool, int) {
	for _, u := range units {
		if !u.alive {
			continue
		}
		if !anyEnemyAlive(units, u.kind) {
			return true, rounds * sumHP(units)
		}
		occ := occupied(units)
		if adjacentEnemy(units, u.x, u.y, u.kind) == nil {
			if fx, fy, ok := stepToward(grid, occ, units, u); ok {
				u.x, u.y = fx, fy
			}
		}
		if done, result := attackTarget(units, u, elfAP, stopOnElfDeath); done {
			return true, result
		}
	}
	return false, 0
}

func attackTarget(units []*unit, u *unit, elfAP int, stopOnElfDeath bool) (bool, int) {
	target := adjacentEnemy(units, u.x, u.y, u.kind)
	if target == nil {
		return false, 0
	}
	ap := 3
	if u.kind == 'E' {
		ap = elfAP
	}
	target.hp -= ap
	if target.hp <= 0 {
		target.alive = false
		if stopOnElfDeath && target.kind == 'E' {
			return true, -1
		}
	}
	return false, 0
}

// One returns the answer to the first part of the exercise.
// Answer: 263327
func (e Exercise) One(instr string) (any, error) {
	return combat(instr, 3, false), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 77872
func (e Exercise) Two(instr string) (any, error) {
	// Find the lowest elf attack power (at least 4) that wins with no elf deaths.
	for ap := 4; ; ap++ {
		if outcome := combat(instr, ap, true); outcome != -1 {
			return outcome, nil
		}
	}
}
