package exercises

import "fmt"

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("[%d, %d]", p.x, p.y)
}