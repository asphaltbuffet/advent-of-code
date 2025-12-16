package exercises

import (
	"strings"
)

type Point struct {
	X int
	Y int
}

type Floor map[Point]bool

func NewFloor(s string) (*Floor, error) {
	f := make(Floor, 0)

	for y, l := range strings.Split(s, "\n") {
		for x, v := range l {
			if v == '@' {
				p := Point{x, y}
				f[p] = true
			}
		}
	}

	return &f, nil
}

func (f Floor) CountRolls() int {
	var count int
	for k, v := range f {
		if v && f.CanAccess(k) {
			count++
		}
	}

	return count
}

func (f Floor) CanAccess(p Point) bool {
	var adj int

	for _, p := range p.Adjacent() {
		if f[p] {
			adj++
		}

		if adj >= 4 {
			return false
		}
	}
	return adj < 4
}

func (p Point) Adjacent() []Point {

	return []Point{
		Point{p.X - 1, p.Y - 1},
		Point{p.X - 1, p.Y + 1},
		Point{p.X + 1, p.Y + 1},
		Point{p.X + 1, p.Y - 1},
		Point{p.X - 1, p.Y},
		Point{p.X + 1, p.Y},
		Point{p.X, p.Y - 1},
		Point{p.X, p.Y + 1},
	}
}

func (f Floor) RemoveRolls() int {
	var removed int
	var done bool

	for !done {
		// assume done until we remove a roll
		done = true
		for p, v := range f {
			if !v || f.CanAccess(p) {
				delete(f, p)
				removed++
				done = false
			}
		}
	}

	return removed
}
