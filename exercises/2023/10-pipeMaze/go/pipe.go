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

func parseInput(s string) (map[Point]Pipe, Point, error) {
	lines := strings.Split(s, "\n")
	pipes := make(map[Point]Pipe, len(lines)*len(lines[0]))

	var start Point
	for y, line := range lines {
		for x, r := range line {
			shape := GetShape(r)
			if shape == InvalidPipe {
				return nil, Point{-1, -1}, fmt.Errorf("invalid pipe shape: %c", r)
			} else if shape == NoPipe {
				continue
			}
			pos := Point{x, y}
			connections := getConnections(pos, shape)
			pipes[pos] = Pipe{Pos: pos, To: connections, Shape: shape}

			if shape == StartPipe {
				start = pos
			}
		}
	}

	return pipes, start, nil
}

func getConnections(pos Point, shape PipeShape) []Point {
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

func findPathLength(m map[Point]Pipe, start Point) (int, error) {
	path := []Pipe{m[start]}

	// fmt.Printf("starting at %v\n", start)

	var next, last, cur Pipe
	last = m[start]
	cur = m[start]

	for {
		np := getNextPoint(m, last, cur)
		if np == (Point{-1, -1}) {
			return 0, fmt.Errorf("no next point found")
		}

		// if we've found the start again, we're done
		if np == start {
			break
		}

		next = m[np]

		path = append(path, next)

		// fmt.Printf("current path: %+v\n", path)

		last = cur
		cur = next
	}

	return len(path), nil
}

func getNextPoint(m map[Point]Pipe, last, cur Pipe) Point {
	for _, p := range cur.To {
		// skip if we just came from this pipe
		if last.Pos == p {
			continue
		}

		if slices.Contains(m[p].To, cur.Pos) {
			// fmt.Printf("found next pipe at %v\n", p)
			return p
		}
	}

	// if we get here, something has gone very wrong
	fmt.Printf("no next pipe found for %+v\n", cur)
	return Point{-1, -1}
}
