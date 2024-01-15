package utilities

import (
	"fmt"
	"math"
)

type Point2D struct {
	X, Y int
}

func (p Point2D) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point2D) Add(m Point2D) Point2D {
	return Point2D{p.X + m.X, p.Y + m.Y}
}

func (p Point2D) Sub(m Point2D) Point2D {
	return Point2D{p.X - m.X, p.Y - m.Y}
}

func (p Point2D) ManhattanDistance(m Point2D) int {
	return AbsInt(p.X-m.X) + AbsInt(p.Y-m.Y)
}

func (p Point2D) EuclideanDistance(m Point2D) float64 {
	dx := float64(p.X - m.X)
	dy := float64(p.Y - m.Y)

	return math.Sqrt(dx*dx + dy*dy)
}
