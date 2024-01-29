package exercises

import (
	"fmt"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

type City struct {
	Blocks map[Point]*Block
	Width  int
	Height int
}

var (
	up    = Point{0, -1}
	down  = Point{0, 1}
	left  = Point{-1, 0}
	right = Point{1, 0}
)

func LoadCity(input string) (*City, error) {
	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}

	c := &City{}

	lines := strings.Split(input, "\n")

	c.Width = len(lines[0])
	c.Height = len(lines)
	c.Blocks = make(map[Point]*Block)

	for y, line := range lines {
		for x, h := range line {
			p := Point{x, y}
			c.Blocks[p] = &Block{
				HeatLoss: int(h - '0'),
				Position: p,
				City:     c,
			}
		}
	}

	return c, nil
}

func (c City) Start() *Block {
	return c.Blocks[Point{0, 0}]
}

func (c City) End() *Block {
	return c.Blocks[Point{c.Width - 1, c.Height - 1}]
}
