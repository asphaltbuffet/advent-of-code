package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 14
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer: 672
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")
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
		default:
			fmt.Printf("ERROR: building graph: %v\n", err)
		}
	}

	// day14.RenderGraph()

	count := 0

	for _, t := range day14.Tiles {
		if t.Type == Sand {
			count++
		}
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 26831
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")
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

	// day14.RenderGraph()

	count := 0

	for _, t := range day14.Tiles {
		if t.Type == Sand {
			count++
		}
	}

	return count, nil
}

// Day14 is the exercise environment.
type Day14 struct {
	Tiles map[Point]Tile
	MinX  int
	MaxX  int
	MaxY  int
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
