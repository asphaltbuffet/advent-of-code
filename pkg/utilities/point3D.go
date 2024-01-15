package utilities

import "fmt"

type Point3D struct {
	X, Y, Z int
}

func (p Point3D) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}

func (p Point3D) Add(m Point3D) Point3D {
	return Point3D{p.X + m.X, p.Y + m.Y, p.Z + m.Z}
}

func (p Point3D) Sub(m Point3D) Point3D {
	return Point3D{p.X - m.X, p.Y - m.Y, p.Z - m.Z}
}
