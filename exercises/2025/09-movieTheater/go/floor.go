package exercises

import (
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

func ParsePoints(s string) []Point {
	lines := strings.Split(s, "\n")

	points := make([]Point, len(lines))
	for i, l := range lines {
		tokens := strings.Split(l, ",")
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])

		points[i] = Point{X: x, Y: y}
	}

	return points
}

func abs(a int) int {
	if a < 0 {
		a = -a
	}

	return a
}

func BiggestRect(points []Point) int {
	bfr := 0

	for i, p1 := range points {
		for j, p2 := range slices.Backward(points) {
			if i == j {
				break
			}

			area := (abs(p1.X-p2.X) + 1) * (abs(p1.Y-p2.Y) + 1)
			// fmt.Printf("%v-%v=%d\n", p1, p2, area)
			bfr = max(bfr, area)
		}
	}

	return bfr
}
