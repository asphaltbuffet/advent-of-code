package exercises

import (
	"fmt"
	"slices"
	"strings"
)

type Point struct {
	X, Y int
}

//go:generate stringer -type=PipeShape
type PipeShape rune

const (
	NoPipe         PipeShape = '.'
	VerticalPipe   PipeShape = '|'
	HorizontalPipe PipeShape = '-'
	NECornerPipe   PipeShape = 'L'
	NWCornerPipe   PipeShape = 'J'
	SWCornerPipe   PipeShape = '7'
	SECornerPipe   PipeShape = 'F'
	StartPipe      PipeShape = 'S'
	InvalidPipe    PipeShape = 'X'
)

func GetShape(r rune) PipeShape {
	switch p := PipeShape(r); p {
	case NoPipe, StartPipe, VerticalPipe, HorizontalPipe, NECornerPipe, NWCornerPipe, SWCornerPipe, SECornerPipe:
		return p

	default:
		return InvalidPipe
	}
}

type Pipe struct {
	Pos   Point
	To    []Point
	Shape PipeShape
}

func parseInput(s string) (map[Point]*Pipe, Point, error) {
	lines := strings.Split(s, "\n")
	pipes := make(map[Point]*Pipe, len(lines)*len(lines[0]))

	var start Point
	for y, line := range lines {
		for x, r := range line {
			shape := GetShape(r)
			if shape == InvalidPipe {
				return nil, Point{-1, -1}, fmt.Errorf("invalid pipe shape: %c", r)
			}

			pos := Point{x, y}
			connections := getConnections(pos, shape)
			pipes[pos] = &Pipe{
				Pos:   pos,
				To:    connections,
				Shape: shape,
			}

			if shape == StartPipe {
				start = pos
			}
		}
	}

	return pipes, start, nil
}

func getConnections(pos Point, shape PipeShape) []Point {
	// ground has no connections
	if shape == NoPipe {
		return []Point{}
	}

	connections := make([]Point, 0, 2) // only start pipe can have 4 connections; all others have 2

	// make north connections
	switch shape {
	case StartPipe, VerticalPipe, NECornerPipe, NWCornerPipe:
		connections = append(connections, Point{pos.X, pos.Y - 1}) // not excluding out-of-bounds points
	}

	// make east connections
	switch shape {
	case StartPipe, HorizontalPipe, NECornerPipe, SECornerPipe:
		connections = append(connections, Point{pos.X + 1, pos.Y}) // not excluding out-of-bounds points
	}

	// make south connections
	switch shape {
	case StartPipe, VerticalPipe, SWCornerPipe, SECornerPipe:
		connections = append(connections, Point{pos.X, pos.Y + 1}) // not excluding out-of-bounds points
	}

	// make west connections
	switch shape {
	case StartPipe, HorizontalPipe, NWCornerPipe, SWCornerPipe:
		connections = append(connections, Point{pos.X - 1, pos.Y}) // not excluding out-of-bounds points
	}

	return connections
}

func findPath(m map[Point]*Pipe, start Point) ([]Point, map[Point]bool, error) {
	isOnPath := make(map[Point]bool, len(m))
	isOnPath[start] = true
	path := []Point{start}

	var next, cur *Pipe
	cur = m[start]

	for {
		np := getNextPoint(m, isOnPath, cur)
		if np == (Point{-1, -1}) {
			if slices.Contains(cur.To, start) {
				// if we've found the start again, we're done
				path = append(path, start)
				break
			}

			return nil, nil, fmt.Errorf("no next point found for %+v", cur)
		}

		isOnPath[np] = true
		path = append(path, np)
		next = m[np]

		cur = next
	}

	return path, isOnPath, nil
}

func getNextPoint(m map[Point]*Pipe, seen map[Point]bool, cur *Pipe) Point {
	for _, p := range cur.To {
		// skip if we just came from this pipe
		if seen[p] {
			continue
		}

		// make sure the pipes can connect
		if slices.Contains(m[p].To, cur.Pos) {
			return p
		}
	}

	// should only get here if we've
	return Point{-1, -1}
}

func countInside(s string) int {
	m, start, err := parseInput(s)
	if err != nil {
		return 0
	}

	path, onpath, err := findPath(m, start)
	if err != nil {
		return 0
	}

	var count int

	lines := strings.Split(s, "\n")

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[0]); x++ {
			tmp := (Point{x, y})

			switch {
			case onpath[tmp]:
				continue
			case windingNumber(tmp, path):
				count++
			}
		}
	}

	return count
}

func windingNumber(pt Point, path []Point) bool {
	wn := 0 // the winding number counter

	// loop through all edges of the polygon
	for i := 0; i < len(path)-1; i++ {
		if path[i].Y <= pt.Y {
			if path[i+1].Y > pt.Y { // an upward crossing
				if isLeft(path[i], path[i+1], pt) > 0 { // point is left of edge
					wn++ // have a valid up intersect
				}
			}
		} else {
			if path[i+1].Y <= pt.Y { // a downward crossing
				if isLeft(path[i], path[i+1], pt) < 0 { // point is right of edge
					wn-- // have a valid down intersect
				}
			}
		}
	}

	return wn != 0
}

// test if a point is left (>0) or right (<0) of an edge
func isLeft(a, b, c Point) int {
	return (b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y)
}
