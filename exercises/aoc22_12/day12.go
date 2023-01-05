// Package aoc22_12 contains the solution for day 12 of Advent of Code 2022.
package aoc22_12 //nolint:revive,stylecheck // I don't care about the package name

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/fatih/color"
)

const debug12 bool = false

const (
	lowest  = 'a'
	highest = 'z'
	start   = 'S'
	end     = 'E'
)

// Point is an X, Y coordinate.
type Point struct {
	X int
	Y int
}

// Location is a set of coordinates on the map with a height.
type Location struct {
	Coord  Point
	Height int
}

// D12P1 returns the solution for 2022 day 12 part 1.
//
// https://adventofcode.com/2022/day/12
//
// incorrect: 354
// answer: 330
func D12P1(data []string) string {
	g, startPoint, endPoint, err := PopulateFromInput(data)
	if err != nil {
		return fmt.Sprintf("error: %v\n", err)
	}

	// calculate the shortest path
	path, _ := graph.ShortestPath(g, startPoint, endPoint)

	if debug12 {
		fmt.Printf("path from %v to %v: %v\n", startPoint, endPoint, path)
	}

	// PrintPath(path, data)
	//
	// file, _ := os.Create("./mapGraph1.gv")
	// err = draw.DOT(g, file)
	// if err != nil {
	// 	return fmt.Sprintf("error: %v\n", err)
	// }

	return strconv.Itoa(len(path) - 1) // number of steps is one less than the number of locations in the path.
}

func locationHash(l Location) Point {
	return l.Coord
}

// PopulateFromInput populates the graph with the input data and sets up possible travel between points.
func PopulateFromInput(data []string) (graph.Graph[Point, Location], Point, Point, error) {
	var sPoint, ePoint Point

	dimY := len(data) - 1

	g := graph.New(locationHash, graph.Directed(), graph.Weighted())

	// parse the input
	for row, line := range data {
		// parse the line
		for col, c := range line {
			cur := Location{Point{X: col, Y: dimY - row}, GetHeight(c)}

			switch c {
			case start:
				sPoint = cur.Coord

				if debug12 {
					fmt.Printf("start point: %v\n", sPoint)
				}

			case end:
				ePoint = cur.Coord

				if debug12 {
					fmt.Printf("end point: %v\n", ePoint)
				}
			}

			err := g.AddVertex(cur, graph.VertexAttribute("label", string(c)))
			if err != nil {
				return nil, Point{}, Point{}, fmt.Errorf("adding location vertex: %w", err)
			}

			// add travel connectivity to the vertex above
			tgt, err := g.Vertex(Point{X: col, Y: dimY - row + 1})

			switch {
			case errors.Is(err, graph.ErrVertexNotFound):
				// no vertex above
				if debug12 {
					fmt.Printf("no vertex at %v\n", Point{X: col, Y: dimY - row + 1})
				}

			case err != nil:
				return nil, Point{}, Point{}, fmt.Errorf("getting vertex at %v: %w", Point{X: col, Y: dimY - row + 1}, err)

			default:
				setTravelEdges(g, cur, tgt)
			}

			// add travel connectivity to the vertex left
			tgt, err = g.Vertex(Point{X: col - 1, Y: dimY - row})

			switch {
			case errors.Is(err, graph.ErrVertexNotFound):
				// no vertex above
				if debug12 {
					fmt.Printf("no vertex at %v\n", Point{X: col - 1, Y: dimY - row})
				}

			case err != nil:
				return nil, Point{}, Point{}, fmt.Errorf("getting vertex at %v: %w", Point{X: col - 1, Y: dimY - row}, err)

			default:
				setTravelEdges(g, cur, tgt)
			}
		}
	}

	return g, sPoint, ePoint, nil
}

func setTravelEdges(g graph.Graph[Point, Location], srcLoc, tgtLoc Location) {
	diff := srcLoc.Height - tgtLoc.Height

	switch {
	case diff > 1: // steep downward; can only travel down to target
		_ = g.AddEdge(srcLoc.Coord, tgtLoc.Coord, graph.EdgeWeight(1))

	case diff < -1: // steep upward; can only travel up to target
		_ = g.AddEdge(tgtLoc.Coord, srcLoc.Coord, graph.EdgeWeight(1))

	default: // +/- 1 change or level; can travel both ways
		_ = g.AddEdge(srcLoc.Coord, tgtLoc.Coord, graph.EdgeWeight(1))
		_ = g.AddEdge(tgtLoc.Coord, srcLoc.Coord, graph.EdgeWeight(1))
	}
}

// GetHeight calculates the height of of a map location. a = 0 -> z = 25.
func GetHeight(c rune) int {
	switch c {
	case start:
		return 0 // lowest - lowest

	case end:
		return int(highest - lowest)

	default:
		return int(c - lowest)
	}
}

// PrintPath prints the path on the map.
func PrintPath(path []Point, data []string) {
	pMap := make(map[Point]bool)

	for _, p := range path {
		pMap[p] = true
	}

	fmt.Println("Map with path:")

	// red := color.New(color.FgRed).SprintFunc()

	for j, line := range data {
		for i, s := range strings.Split(line, "") {
			if pMap[Point{X: i, Y: j}] {
				color.Set(color.FgHiRed, color.Bold)
				fmt.Printf("%s", s)
				color.Unset()
			} else {
				fmt.Printf("%s", s)
			}
		}

		fmt.Println()
	}
}

// D12P2 returns the solution for 2022 day 12 part 2.
// answer: 321
func D12P2(data []string) string {
	g, _, endPoint, err := PopulateFromInput(data)
	if err != nil {
		return fmt.Sprintf("error: %v\n", err)
	}

	minY := len(data) - 1
	minPath := len(data) * len(data[0])

	for j := 0; j < len(data); j++ {
		for i := 0; i < len(data[0]); i++ {
			v, err := g.Vertex(Point{X: i, Y: minY - j})
			if err != nil {
				return fmt.Sprintf("error getting vertex at %v: %v\n", Point{X: i, Y: minY - j}, err)
			}

			if v.Height == 0 {
				path, err := graph.ShortestPath(g, v.Coord, endPoint)
				if err != nil {
					if debug12 {
						fmt.Printf("path from %v to %v: %v\n", v.Coord, endPoint, path)
					}

					continue // no path from here
				}

				if len(path)-1 < minPath {
					minPath = len(path) - 1
				}
			}
		}
	}

	return strconv.Itoa(minPath)
}
