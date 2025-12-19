package exercises

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Line struct {
	P1 Point
	P2 Point
}

type Rect []Point

type Floor struct {
	Min Point
	Max Point

	points   []Point
	segments []Line
	Horiz    []*Line
	Vert     []*Line
}

func NewFloor(s string) (*Floor, error) {
	lines := strings.Split(s, "\n")

	pp := make([]Point, len(lines))
	minP, maxP := Point{math.MaxInt, math.MaxInt}, Point{0, 0}

	for i, l := range lines {
		tokens := strings.Split(l, ",")
		x, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, err
		}

		pp[i] = Point{X: x, Y: y}
		minP.X = min(minP.X, x)
		minP.Y = min(minP.Y, y)
		maxP.X = max(maxP.X, x)
		maxP.Y = max(maxP.Y, y)
	}

	return &Floor{
		points:   pp,
		segments: getLines(pp),
		Min:      minP,
		Max:      maxP,
	}, nil
}

func getLines(pp []Point) []Line {
	lines := []Line{}

	for i := 1; i <= len(pp); i++ {
		from := pp[i-1]
		to := pp[i%len(pp)]
		lines = append(lines, Line{from, to})
	}

	return lines
}

func abs(a int) int {
	if a < 0 {
		a = -a
	}

	return a
}

func (f Floor) BiggestRect() int {
	bfr := 0

	for i, p1 := range f.points {
		for j, p2 := range slices.Backward(f.points) {
			if i == j {
				break
			}

			area := (abs(p1.X-p2.X) + 1) * (abs(p1.Y-p2.Y) + 1)
			bfr = max(bfr, area)
		}
	}

	return bfr
}

func (f *Floor) BoundedRect() int {
	bfr := 0

	for i, p1 := range f.points {
		for j, p2 := range slices.Backward(f.points) {
			if i == j {
				break
			}

			o1, o2 := Point{p1.X, p2.Y}, Point{p2.X, p1.Y}
			r := Rect{p1, o1, p2, o2}
			area := (abs(p1.X-p2.X) + 1) * (abs(p1.Y-p2.Y) + 1)
			if area > bfr && f.IsValidRectangle(r) {
				bfr = area
			}
		}
	}

	return bfr
}

func (f *Floor) IsValidRectangle(r Rect) bool {
	// check corners first
	for _, c := range r {
		if !f.OnEdge(c) && !f.InsidePoly(c) {
			return false
		}
	}

	r1 := Point{X: min(r[0].X, r[2].X), Y: min(r[0].Y, r[2].Y)}
	r2 := Point{X: max(r[0].X, r[2].X), Y: max(r[0].Y, r[2].Y)}

	a := f.points[len(f.points)-1]
	for _, b := range f.points {
		if edgeIntersectsRectangle(a, b, r1, r2) {
			return false
		}

		a = b
	}

	return true
}

func edgeIntersectsRectangle(e1, e2, r1, r2 Point) bool {
	// horiz
	if e1.Y == e2.Y {
		if r1.Y < e1.Y && e1.Y < r2.Y {
			if max(e1.X, e2.X) > r1.X && min(e1.X, e2.X) < r2.X {
				return true
			}
		}
	} else { // vert
		if r1.X < e1.X && e1.X < r2.X {
			if max(e1.Y, e2.Y) > r1.Y && min(e1.Y, e2.Y) < r2.Y {
				return true
			}
		}
	}

	return false
}

func (f Floor) OnEdge(p Point) bool {
	a := f.points[len(f.points)-1]
	for _, b := range f.points {
		if onSegment(a, b, p) {
			return true
		}

		a = b
	}

	return false
}

func onSegment(a, b, p Point) bool {
	if a.Y == b.Y && p.Y == a.Y {
		if a.X > b.X {
			a.X, b.X = b.X, a.X
		}
		return p.X >= a.X && p.X <= b.X
	}

	if a.X == b.X && p.X == a.X {
		if a.Y > b.Y {
			a.Y, b.Y = b.Y, a.Y
		}
		return p.Y >= a.Y && p.Y <= b.Y
	}

	return false
}

func (f Floor) InsidePoly(p Point) bool {
	a := f.points[0]
	inside := intersects(p, f.points[len(f.points)-1], a)
	for _, b := range f.points[1:] {
		if intersects(p, a, b) {
			inside = !inside
		}
		a = b
	}

	return inside
}

func intersects(p, a, b Point) bool {
	return (a.Y > p.Y) != (b.Y > p.Y) &&
		p.X < (b.X-a.X)*(p.Y-a.Y)/(b.Y-a.Y)+a.X
}
