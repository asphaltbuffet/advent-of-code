// Package aoc22_09 contains the solution for day 9 of Advent of Code 2022.
package aoc22_09 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"strconv"
	"strings"
)

const debug bool = false

// Point is an x/y tuple of ints.
type Point struct {
	X, Y int
}

// D9P1 returns the solution for 2022 day 9 part 1.
// answer: 6266
func D9P1(data []string) string {
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

	return strconv.Itoa(len(visited))
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

var minX, maxX, minY, maxY int

// D9P2 returns the solution for 2022 day 9 part 2.
// answer: 2369
func D9P2(data []string) string {
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
		if debug {
			fmt.Printf("\n== %s ==\n\n", line)
		}

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

				if debug {
					setDebugMinMax(tailLocation)
				}
			}

			if debug {
				fmt.Printf("head: %v, tail: %v\n", headLocation, tailLocation)
			}
		}

		if debug {
			PrintState(tailLocation)
		}
	}

	return strconv.Itoa(len(visited))
}

func setDebugMinMax(tailLocation []Point) {
	for _, t := range tailLocation {
		if t.X < minX {
			minX = t.X
		}

		if t.X > maxX {
			maxX = t.X
		}

		if t.Y < minY {
			minY = t.Y
		}

		if t.Y > maxY {
			maxY = t.Y
		}
	}
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
