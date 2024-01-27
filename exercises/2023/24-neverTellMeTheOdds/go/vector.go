package exercises

import "fmt"

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

type Vector2D struct {
	X float64
	Y float64
}

func (v Vector2D) Cross(v2 Vector2D) float64 {
	return (v.X * v2.Y) - (v.Y * v2.X)
}

func (v Vector3D) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}
