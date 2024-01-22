package exercises

import (
	"fmt"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Point util.Point2D

func (p Point) add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%2d, %2d)", p.X, p.Y)
}
