package aoc22_08 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

// Coordinate is a type that represents a coordinate.
type Coordinate struct {
	X int
	Y int
}

// Visibility is a type that represents the visibility of a tree.
type Visibility int

// Visibility is a type that represents the visibility of a tree.
const (
	Unknown Visibility = iota // default behavior of "Unknown" visibility is by design.
	Yes
	No
)

// Tree is a type that represents a tree.
type Tree struct {
	Name     string
	Location Coordinate
	Height   int
}

// GetVisibleTreesFromNorth returns the visible trees from the north.
func GetVisibleTreesFromNorth(tm [][]int) []int {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Tree())
	_ = g.AddVertex(-1)

	for col := 1; col < dimX-1; col++ {
		max := 0
		lastVertex := -1

		for row := 0; row < dimY-1; row++ {
			switch {
			case row == 0:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)

			case tm[row][col] > tm[row-1][col] && tm[row][col] > max:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)
			}

			if max < tm[row][col] {
				max = tm[row][col]
			}
		}
	}

	// file, _ := os.Create("./north_graph.gv")
	// _ = draw.DOT(g, file)

	order, err := graph.TopologicalSort(g)
	if err != nil {
		fmt.Printf("error sorting: %v", err)
	}

	// fmt.Printf("north: %v\n)", order)

	return order
}

// GetVisibleTreesFromSouth returns the visible trees from the south.
func GetVisibleTreesFromSouth(tm [][]int) []int {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Tree())
	_ = g.AddVertex(-1)

	for col := 1; col < dimX-1; col++ {
		max := 0
		lastVertex := -1

		for row := dimY - 1; row >= 0; row-- {
			switch {
			case row == dimY-1:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)

			case tm[row][col] > tm[row+1][col] && tm[row][col] > max:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)
			}

			if max < tm[row][col] {
				max = tm[row][col]
			}
		}
	}

	// file, _ := os.Create("./north_graph.gv")
	// _ = draw.DOT(g, file)

	order, err := graph.TopologicalSort(g)
	if err != nil {
		fmt.Printf("error sorting: %v", err)
	}

	// fmt.Printf("south: %v\n)", order)

	return order
}

// GetVisibleTreesFromWest returns the visible trees from the west.
func GetVisibleTreesFromWest(tm [][]int) []int {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Tree())
	_ = g.AddVertex(-1)

	for row := 1; row < dimY-1; row++ {
		max := 0
		lastVertex := -1

		for col := 0; col < dimX; col++ {
			switch {
			case col == 0:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)

			case tm[row][col] > tm[row][col-1] && tm[row][col] > max:
				_ = g.AddVertex(col + (row * dimX))
				_ = g.AddEdge(lastVertex, col+(row*dimX))
				lastVertex = col + (row * dimX)
			}

			if max < tm[row][col] {
				max = tm[row][col]
			}
		}
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		fmt.Printf("error sorting: %v", err)
	}

	// fmt.Printf("west: %v\n)", order)

	return order
}

// GetVisibleTreesFromEast returns the visible trees from the east.
func GetVisibleTreesFromEast(tm [][]int) []int {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Tree())
	_ = g.AddVertex(-1)

	for row := 1; row < dimY-1; row++ {
		max := 0
		lastVertex := -1

		for col := dimX - 1; col >= 0; col-- {
			cur := col + (row * dimX)

			switch {
			case col == dimX-1:
				_ = g.AddVertex(cur)
				_ = g.AddEdge(-1, cur)
				lastVertex = col + (row * dimX)

			case tm[row][col] > tm[row][col+1] && tm[row][col] > max:
				_ = g.AddVertex(cur)
				_ = g.AddEdge(lastVertex, cur)
				lastVertex = col + (row * dimX)
			}

			if max < tm[row][col] {
				max = tm[row][col]
			}
		}
	}

	// file, _ := os.Create("./north_graph.gv")
	// _ = draw.DOT(g, file)

	order, err := graph.TopologicalSort(g)
	if err != nil {
		fmt.Printf("error sorting: %v", err)
	}

	// fmt.Printf("east: %v\n)", order)

	return order
}

// GetTreeMap takes the input data and creates a 2D array of trees.
func GetTreeMap(data []string) [][]int {
	t := make([][]int, dimY)

	for i, line := range data {
		t[i] = aoc.Map(strings.Split(line, ""), func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		})
	}

	return t
}
