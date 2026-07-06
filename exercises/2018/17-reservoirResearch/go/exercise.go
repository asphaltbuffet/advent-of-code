package exercises

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 17.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	g := simulate(instr)
	return g.count(true), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	g := simulate(instr)
	return g.count(false), nil
}

// Tile states in the scan.
const (
	sand    = 0 // untouched
	clay    = 1 // '#'
	flowing = 2 // '|' water that has passed through
	settled = 3 // '~' water at rest
)

type grid struct {
	tiles      [][]byte // [y][x] indexed absolute, offset by minX
	minX, maxX int
	minY, maxY int
}

// simulate parses the clay veins and runs the water flow from the spring at x=500.
func simulate(instr string) *grid {
	type vein struct{ x1, x2, y1, y2 int }

	re := regexp.MustCompile(`-?\d+`)
	var veins []vein
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		if line == "" {
			continue
		}
		nums := re.FindAllString(line, -1)
		a, _ := strconv.Atoi(nums[0])
		b, _ := strconv.Atoi(nums[1])
		c, _ := strconv.Atoi(nums[2])

		var v vein
		if line[0] == 'x' { // x=a, y=b..c
			v = vein{a, a, b, c}
		} else { // y=a, x=b..c
			v = vein{b, c, a, a}
		}
		veins = append(veins, v)

		minX = min(minX, v.x1)
		maxX = max(maxX, v.x2)
		minY = min(minY, v.y1)
		maxY = max(maxY, v.y2)
	}

	// Pad x by one on each side so water can spill past the outermost clay.
	minX--
	maxX++

	g := &grid{
		minX: minX, maxX: maxX,
		minY: minY, maxY: maxY,
	}
	g.tiles = make([][]byte, maxY+1)
	for y := range g.tiles {
		g.tiles[y] = make([]byte, maxX-minX+1)
	}
	for _, v := range veins {
		for y := v.y1; y <= v.y2; y++ {
			for x := v.x1; x <= v.x2; x++ {
				g.tiles[y][x-minX] = clay
			}
		}
	}

	g.flow(500, 0)
	return g
}

func (g *grid) at(x, y int) byte     { return g.tiles[y][x-g.minX] }
func (g *grid) set(x, y int, b byte) { g.tiles[y][x-g.minX] = b }

// isFloor reports whether the tile below (x, y) can hold water: clay or settled.
func (g *grid) isFloor(x, y int) bool {
	b := g.at(x, y)
	return b == clay || b == settled
}

// flow drops water from (x, y) downward, spreading across floors and settling
// enclosed rows. It returns true if (x, y) itself comes to rest as settled water,
// so its caller (the row above) knows this tile is now a floor and can spread over it.
func (g *grid) flow(x, y int) bool {
	g.set(x, y, flowing)

	below := y + 1
	if below > g.maxY {
		return false // falls off the bottom of the scan
	}
	if !g.isFloor(x, below) {
		if g.at(x, below) == sand {
			if !g.flow(x, below) {
				return false // fell through and kept flowing
			}
		} else {
			return false // already flowing through here: not a floor
		}
	}

	// A floor is below: spread left and right. When a side reaches an open edge,
	// spill down it; if that spill settles into a new floor, keep extending the row.
	leftX, leftWall := g.spread(x, y, -1)
	rightX, rightWall := g.spread(x, y, +1)

	if leftWall && rightWall {
		for fx := leftX; fx <= rightX; fx++ {
			g.set(fx, y, settled)
		}
		return true
	}
	return false
}

// spread walks from x in direction dir along row y, marking flowing water. It
// stops at clay (a wall) or at an open edge it cannot seal. At each open edge it
// spills water downward; if that spill settles, it extends past the edge and keeps
// going. Returns the furthest reachable x and whether it ended against a wall.
func (g *grid) spread(x, y, dir int) (int, bool) {
	for {
		if g.at(x+dir, y) == clay {
			return x, true
		}
		x += dir
		g.set(x, y, flowing)
		if !g.isFloor(x, y+1) {
			// Open edge: try to seal it by filling below.
			if g.at(x, y+1) == sand && g.flow(x, y+1) {
				continue // spill settled into a floor; extend the row
			}
			return x, false // water escapes here
		}
	}
}

// count returns reachable water tiles: all flowing+settled if includeFlowing,
// else settled only. Only rows within [minY, maxY] count.
func (g *grid) count(includeFlowing bool) int {
	n := 0
	for y := g.minY; y <= g.maxY; y++ {
		for _, t := range g.tiles[y] {
			if t == settled || (includeFlowing && t == flowing) {
				n++
			}
		}
	}
	return n
}
