package aoc22

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
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

var dimX, dimY int

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 8, D8P1, D8P2, Get2022Command())
}

// D8P1 returns the solution for 2022 day 8 part 1.
// answer: 1805
func D8P1(data []string) string {
	dimX = len(data[0])
	dimY = len(data)

	m := GetTreeMap(data)

	n := GetVisibleTreesFromNorth(m)
	s := GetVisibleTreesFromSouth(m)
	e := GetVisibleTreesFromEast(m)
	w := GetVisibleTreesFromWest(m)

	// fmt.Printf("trees: %+v\n", trees)
	// file, _ := os.Create("./mygraph.gv")
	// _ = draw.DOT(tGraph, file)
	out := []int{0, dimX - 1, dimY * (dimX - 1), (dimX * dimY) - 1} // corners
	out = append(out, n...)
	out = append(out, s...)
	out = append(out, e...)
	out = append(out, w...)
	out = aoc.Unique(out)

	sort.Ints(out)
	// fmt.Printf("visible trees: %v\n", out)

	return strconv.Itoa(len(out) - 1) // -1 for the "edge" tree
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

// D8P2 returns the solution for 2022 day 8 part 2.
// answer: 444528
func D8P2(data []string) string {
	dimX = len(data[0])
	dimY = len(data)

	m := GetTreeMap(data)
	maxScenicScore := 0

	for r := 1; r < dimY-1; r++ {
		for c := 1; c < dimX-1; c++ {
			h := m[r][c]
			up := CalculateScoreUp(h, r, c, m)
			down := CalculateScoreDown(h, r, c, m)
			left := CalculateScoreLeft(h, r, c, m)
			right := CalculateScoreRight(h, r, c, m)

			score := up * down * left * right

			// fmt.Printf("[%d,%d]=%d, score: %d * %d * %d * %d = %d\n", r, c, h, up, down, left, right, score)

			if score > maxScenicScore {
				maxScenicScore = score
			}
		}
	}

	return strconv.Itoa(maxScenicScore)
}

func CalculateScoreUp(h, r, c int, m [][]int) int {
	switch {
	case r == 0:
		return 0
	case h > m[r-1][c]:
		return 1 + CalculateScoreUp(h, r-1, c, m)
	default:
		return 1
	}
}

func CalculateScoreDown(h, r, c int, m [][]int) int {
	switch {
	case r == dimY-1:
		return 0
	case h > m[r+1][c]:
		return 1 + CalculateScoreDown(h, r+1, c, m)
	default:
		return 1
	}
}

func CalculateScoreLeft(h, r, c int, m [][]int) int {
	switch {
	case c == 0:
		return 0
	case h > m[r][c-1]:
		return 1 + CalculateScoreLeft(h, r, c-1, m)
	default:
		return 1
	}
}

func CalculateScoreRight(h, r, c int, m [][]int) int {
	switch {
	case c == dimX-1:
		return 0
	case h > m[r][c+1]:
		return 1 + CalculateScoreRight(h, r, c+1, m)
	default:
		return 1
	}
}
