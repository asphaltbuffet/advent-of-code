package utilities

import "fmt"

type Point2D struct {
	X, Y int
}

type Point3D struct {
	X, Y, Z int
}

func (p Point2D) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point3D) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}
