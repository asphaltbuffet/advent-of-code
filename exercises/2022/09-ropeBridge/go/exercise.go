package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Point is an x/y tuple of ints.
type Point struct {
	X, Y int
}

// Exercise for Advent of Code 2022 day 9
type Exercise struct {
	common.BaseExercise
}

var minX, maxX, minY, maxY int

// One returns the answer to the first part of the exercise.
// answer: 6266
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")

	// initialize locations of head and tail
	visited := make(map[Point]bool)
	headLocation := Point{0, 0}
	tailLocation := Point{0, 0}
	visited[tailLocation] = true

	// iterate over the data
	for _, line := range data {
		// get the direction and distance
		direction, right, _ := strings.Cut(line, " ")

		distance, err := strconv.Atoi(right)
		if err != nil {
			fmt.Printf("invalid distance: %v\n", err)
		}

		// calculate the headMovement
		var headMovement Point

		switch direction {
		case "U":
			headMovement = Point{0, 1}
		case "D":
			headMovement = Point{0, -1}
		case "L":
			headMovement = Point{-1, 0}
		case "R":
			headMovement = Point{1, 0}
		default:
			panic("invalid direction")
		}

		// move the head
		for i := 0; i < distance; i++ {
			headLocation.X += headMovement.X
			headLocation.Y += headMovement.Y

			tailMovement := CalculateMovement(headLocation, tailLocation)

			tailLocation.X += tailMovement.X
			tailLocation.Y += tailMovement.Y

			// fmt.Printf("head: %v, tail: %v\n", headLocation, tailLocation)

			visited[tailLocation] = true
		}
	}

	return len(visited), nil
}

// Two returns the answer to the second part of the exercise.
// answer: 2369
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")

	// track the bounds of the grid (for debugging)
	minX, maxX, minY, maxY = 0, 0, 0, 0

	// initialize locations of head and tail
	visited := make(map[Point]bool)
	headLocation := Point{0, 0}

	// initialize the tail locations. Tail 0 is the head, tail 9 is the end.
	var tailLocation []Point

	for i := 0; i < 10; i++ {
		tailLocation = append(tailLocation, Point{0, 0})
	}

	visited[tailLocation[9]] = true

	// iterate over the data
	for _, line := range data {
		// get the direction and distance
		direction, right, _ := strings.Cut(line, " ")

		distance, err := strconv.Atoi(right)
		if err != nil {
			fmt.Printf("invalid distance: %v\n", err)
		}

		// calculate the headMovement
		headMovement := getMovement(direction)

		// move the head
		for i := 0; i < distance; i++ {
			headLocation.X += headMovement.X
			headLocation.Y += headMovement.Y
			tailLocation[0] = headLocation

			for j := 1; j < 10; j++ {
				tailMovement := CalculateMovement(tailLocation[j-1], tailLocation[j])

				tailLocation[j].X += tailMovement.X
				tailLocation[j].Y += tailMovement.Y

				if j == 9 {
					visited[tailLocation[j]] = true
				}
			}
		}
	}

	return len(visited), nil
}

// Abs returns the absolute value of an int.
func Abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

// CalculateMovement returns the movement of tail to follow the head. No movement is returned
// if the tail is within 1 of the head.
func CalculateMovement(h, t Point) Point {
	m := Point{h.X - t.X, h.Y - t.Y}

	if Abs(m.X) <= 1 && Abs(m.Y) <= 1 {
		return Point{0, 0}
	}

	if Abs(m.X) > 1 {
		m.X /= Abs(m.X)
	}

	if Abs(m.Y) > 1 {
		m.Y /= Abs(m.Y)
	}

	return m
}

func getMovement(direction string) Point {
	var headMovement Point

	switch direction {
	case "U":
		headMovement = Point{0, 1}
	case "D":
		headMovement = Point{0, -1}
	case "L":
		headMovement = Point{-1, 0}
	case "R":
		headMovement = Point{1, 0}
	default:
		panic("invalid direction")
	}

	return headMovement
}

// PrintState prints the state of the grid.
func PrintState(tailLocation []Point) {
	dimX := maxX - minX + 1
	dimY := maxY - minY + 1

	grid := make([][]string, dimY)

	for i := range grid {
		grid[i] = make([]string, dimX)
	}

	for i := range grid {
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	for i, t := range tailLocation {
		if grid[t.Y-minY][t.X-minX] == "." {
			grid[t.Y-minY][t.X-minX] = strconv.Itoa(i)
		}
	}

	grid[tailLocation[0].Y-minY][tailLocation[0].X-minX] = "H"

	if grid[0-minY][0-minX] == "." {
		grid[0-minY][0-minX] = "s"
	}

	for i := dimY - 1; i >= 0; i-- {
		fmt.Println(strings.Join(grid[i], ""))
	}
}
