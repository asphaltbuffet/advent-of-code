// Package aoc22_14 contains the solution for day 14 of Advent of Code 2022.
package aoc22_14 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"strconv"
)

// Day14 is the exercise environment.
type Day14 struct {
	Tiles map[Point]Tile
	MinX  int
	MaxX  int
	MaxY  int
}

// D14P1 returns the solution for 2022 day 14 part 1.
//
// https://adventofcode.com/2022/day/14
//
// answer: 672
func D14P1(data []string) string {
	day14 := Day14{
		Tiles: make(map[Point]Tile, 0),
		MinX:  0,
		MaxX:  0,
		MaxY:  0,
	}

	rocks, err := InputToPoints(data)
	if err != nil {
		fmt.Printf("ERROR: processing input data: %v", err)
	}

	// add the rocks
	err = day14.AddRocks(rocks)
	if err != nil {
		fmt.Printf("ERROR: adding rocks: %v", err)
	}

	// add the source
	day14.Tiles[root.Coord] = root

	err = day14.BuildGraph(Point{root.Coord.X, root.Coord.Y + 1})
	if err != nil {
		switch err {
		case ErrVoidPath:
			// we've reached the bottom, that's fine
			fmt.Println("simulation complete")
		default:
			fmt.Printf("ERROR: building graph: %v\n", err)
		}
	}

	day14.RenderGraph()

	count := 0

	for _, t := range day14.Tiles {
		if t.Type == Sand {
			count++
		}
	}

	return strconv.Itoa(count)
}

// ProcessInput takes raw input and populates the exercise's environment.
func ProcessInput(data []string) ([][]Point, error) {
	// parse the input by line
	points, err := InputToPoints(data)
	if err != nil {
		return nil, fmt.Errorf("converting input to points: %w", err)
	}

	return points, nil
}

// D14P2 returns the solution for 2022 day 14 part 2.
// answer: 26831
func D14P2(data []string) string {
	day14 := Day14{
		Tiles: make(map[Point]Tile, 0),
		MinX:  0,
		MaxX:  0,
		MaxY:  0,
	}

	rocks, err := InputToPoints(data)
	if err != nil {
		fmt.Printf("ERROR: processing input data: %v", err)
	}

	// add the rocks
	err = day14.AddRocks(rocks)
	if err != nil {
		fmt.Printf("ERROR: adding rocks: %v", err)
	}

	// increase max Y to account for the floor
	day14.MaxY += 2

	err = day14.BuildGraphWithFloor(root.Coord)
	if err != nil {
		switch err {
		case ErrVoidPath: // we've reached the bottom, that's fine

		default:
			fmt.Printf("ERROR: building graph: %v\n", err)
		}
	}

	day14.RenderGraph()

	count := 0

	for _, t := range day14.Tiles {
		if t.Type == Sand {
			count++
		}
	}

	return strconv.Itoa(count)
}
