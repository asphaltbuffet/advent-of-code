package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

// Point is an x/y tuple of ints.
type Point struct {
	X, Y int
}

// Content is the type of a tile.
type Content string

// Tile types.
const (
	Unknown Content = ""
	Air     Content = "·"
	Rock    Content = "*"
	Source  Content = "+"
	Sand    Content = "o"
)

// Tile is a space in the environment.
type Tile struct {
	Coord Point
	Type  Content
}

// InputToPoints converts raw input into a slice of slices of points.
func InputToPoints(data []string) ([][]Point, error) {
	points := make([][]Point, len(data))

	for i, line := range data {
		points[i] = []Point{}
		tokens := strings.Split(line, "->")

		if len(tokens) < 2 {
			return nil, fmt.Errorf("tokenizing input line %d: %s", i, line)
		}

		for _, token := range tokens {
			point, err := ParseToken(token)
			if err != nil {
				return nil, fmt.Errorf("parsing token: %w", err)
			}

			points[i] = append(points[i], point)
		}
	}

	return points, nil
}

// GetPointsBetween creates an inclusive slice of tiles between two points lying on the same X or Y axis.
func GetPointsBetween(left, right Point) ([]Point, error) {
	modifier := Point{0, 0}

	switch {
	case left.X == right.X && left.Y == right.Y:
		return []Point{left}, fmt.Errorf("left and right are the same point")

	case left.X != right.X && left.Y != right.Y:
		return []Point{}, fmt.Errorf("invalid diagonal line: %v -> %v", left, right)

	case left.X > right.X: // look at using the sgn function
		modifier = Point{-1, 0}

	case left.X < right.X:
		modifier = Point{1, 0}

	case left.Y > right.Y:
		modifier = Point{0, -1}

	case left.Y < right.Y:
		modifier = Point{0, 1}
	}

	p := make([]Point, 0)

	tmp := left

	for {
		p = append(p, tmp)

		tmp.X += modifier.X
		tmp.Y += modifier.Y

		if tmp.X == right.X && tmp.Y == right.Y {
			break
		}
	}

	p = append(p, right)

	return p, nil
}

// ParseToken parses a string of format ' into a Point.
func ParseToken(token string) (Point, error) {
	left, right, ok := strings.Cut(strings.Trim(token, " "), ",")
	if !ok {
		return Point{}, fmt.Errorf("trimming and cutting token: %s", token)
	}

	x, err := strconv.Atoi(left)
	if err != nil {
		return Point{}, fmt.Errorf("parsing x coordinate: %w", err)
	}

	y, err := strconv.Atoi(right)
	if err != nil {
		return Point{}, fmt.Errorf("parsing y coordinate: %w", err)
	}

	return Point{X: x, Y: y}, nil
}
