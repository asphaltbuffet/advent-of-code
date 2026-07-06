package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 13.
type Exercise struct {
	common.BaseExercise
}

// cart is a mine cart: its position, heading, and how many intersections it has
// crossed (which selects left / straight / right in a repeating cycle).
type cart struct {
	x, y   int
	dx, dy int
	turns  int
	dead   bool
}

// parse reads the track grid and lifts the carts off it, replacing each cart glyph
// with the straight track it sits on.
func parse(instr string) ([][]byte, []*cart) {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	grid := make([][]byte, len(lines))
	var carts []*cart

	for y, line := range lines {
		grid[y] = []byte(line)
		for x := range len(grid[y]) {
			var dx, dy int
			var under byte
			switch grid[y][x] {
			case '<':
				dx, dy, under = -1, 0, '-'
			case '>':
				dx, dy, under = 1, 0, '-'
			case '^':
				dx, dy, under = 0, -1, '|'
			case 'v':
				dx, dy, under = 0, 1, '|'
			default:
				continue
			}
			carts = append(carts, &cart{x: x, y: y, dx: dx, dy: dy})
			grid[y][x] = under
		}
	}

	return grid, carts
}

// at returns the track glyph at (x, y), treating out-of-range as empty space.
func at(grid [][]byte, x, y int) byte {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return ' '
	}
	return grid[y][x]
}

// advance moves one cart a single step, turning it according to the track it lands
// on (curves reflect the heading; intersections cycle left / straight / right).
func (c *cart) advance(grid [][]byte) {
	c.x += c.dx
	c.y += c.dy

	switch at(grid, c.x, c.y) {
	case '/':
		c.dx, c.dy = -c.dy, -c.dx
	case '\\':
		c.dx, c.dy = c.dy, c.dx
	case '+':
		switch c.turns % 3 {
		case 0: // left
			c.dx, c.dy = c.dy, -c.dx
		case 2: // right
			c.dx, c.dy = -c.dy, c.dx
		}
		c.turns++
	}
}

// simulate runs the carts until the first crash (part one) or until one cart
// remains (part two), returning that location as "x,y".
func simulate(grid [][]byte, carts []*cart, lastStanding bool) string {
	for {
		sortCarts(carts)
		for _, c := range carts {
			if c.dead {
				continue
			}
			c.advance(grid)
			if crash := findCollision(carts, c); crash != nil {
				if !lastStanding {
					return fmt.Sprintf("%d,%d", c.x, c.y)
				}
				c.dead, crash.dead = true, true
			}
		}
		if lastStanding {
			carts = aliveCarts(carts)
			if len(carts) == 1 {
				return fmt.Sprintf("%d,%d", carts[0].x, carts[0].y)
			}
		}
	}
}

func findCollision(carts []*cart, c *cart) *cart {
	for _, other := range carts {
		if other != c && !other.dead && other.x == c.x && other.y == c.y {
			return other
		}
	}
	return nil
}

func aliveCarts(carts []*cart) []*cart {
	live := carts[:0]
	for _, c := range carts {
		if !c.dead {
			live = append(live, c)
		}
	}
	return live
}

func sortCarts(carts []*cart) {
	sort.Slice(carts, func(i, j int) bool {
		if carts[i].y != carts[j].y {
			return carts[i].y < carts[j].y
		}
		return carts[i].x < carts[j].x
	})
}

// One returns the answer to the first part of the exercise.
// Answer: 8,9
func (e Exercise) One(instr string) (any, error) {
	grid, carts := parse(instr)
	return simulate(grid, carts, false), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 73,33
func (e Exercise) Two(instr string) (any, error) {
	grid, carts := parse(instr)
	return simulate(grid, carts, true), nil
}
