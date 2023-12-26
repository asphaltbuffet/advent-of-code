package exercises

import (
	"fmt"
	"strings"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Point util.Point2D

type Contraption struct {
	Points    map[Point]Tile
	Energized map[Point]map[Direction]bool
	Width     int
	Height    int
}

type Tile rune

const (
	Empty           Tile = '.'
	SplitVertical   Tile = '|'
	SplitHorizontal Tile = '-'
	LeftSlant       Tile = '/'
	RightSlant      Tile = '\\'
)

// var energized map[Point]map[Direction]bool

func parseInput(s string) (*Contraption, error) {
	lines := strings.Split(s, "\n")
	m := make(map[Point]Tile, len(lines)*len(lines[0]))
	energized := make(map[Point]map[Direction]bool)

	for y, line := range lines {
		for x, r := range line {
			m[Point{x, y}] = Tile(r)
			energized[Point{x, y}] = make(map[Direction]bool)
		}
	}

	return &Contraption{
		Points:    m,
		Energized: energized,
		Width:     len(lines[0]),
		Height:    len(lines),
	}, nil
}

func getTyleType(m map[Point]rune, p Point) Tile {
	switch m[p] {
	case '.':
		return Empty

	case '|':
		return SplitVertical

	case '-':
		return SplitHorizontal

	case '/':
		return LeftSlant

	case '\\':
		return RightSlant
	}
	return '.'
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func (c *Contraption) Power(pt Point, d Direction) {
	p := pt

	for {
		if _, ok := c.Points[p]; !ok {
			// in a point that doesn't exist
			return
		}

		if c.Energized[p][d] {
			// we've already been here
			return
		}

		c.Energized[p][d] = true

		switch c.Points[p] {
		case Empty:
			p = add(p, d)

		case SplitVertical:
			switch d {
			case Left, Right:
				np := add(p, Up)
				c.Power(np, Up)

				p = add(p, Down)
				c.Power(p, Down)
				return
			case Up, Down:
				p = add(p, d)
			default:
				fmt.Println("invalid direction: ", d)
				return
			}

		case SplitHorizontal:
			switch d {
			case Up, Down:
				np := add(p, Right)
				c.Power(np, Right)

				p = add(p, Left)
				c.Power(p, Left)
				return
			case Left, Right:
				p = add(p, d)
			default:
				fmt.Println("invalid direction: ", d)
				return
			}
		case LeftSlant:
			switch d {
			case Up:
				p = add(p, Right)
				d = Right

			case Right:
				p = add(p, Up)
				d = Up

			case Down:
				p = add(p, Left)
				d = Left

			case Left:
				p = add(p, Down)
				d = Down
			}

		case RightSlant:

			switch d {
			case Up:
				p = add(p, Left)
				d = Left

			case Right:
				p = add(p, Down)
				d = Down

			case Down:
				p = add(p, Right)
				d = Right

			case Left:
				p = add(p, Up)
				d = Up
			}
		}
	}
}

func add(p Point, d Direction) Point {
	switch d {
	case Up:
		return Point{p.X, p.Y - 1}

	case Down:
		return Point{p.X, p.Y + 1}

	case Right:
		return Point{p.X + 1, p.Y}

	case Left:
		return Point{p.X - 1, p.Y}
	}

	fmt.Println("invalid direction: ", d)
	return Point{-1, -1}
}

func (c *Contraption) CountEnergized() int {
	count := 0

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			v, ok := c.Energized[Point{x, y}]

			if ok && (v[Up] || v[Down] || v[Right] || v[Left]) {
				count++
			}
		}
	}

	return count
}

func (c *Contraption) DebugPrintEnergized() {
	for y := 0; y < c.Height; y++ {
		sb := strings.Builder{}

		for x := 0; x < c.Width; x++ {
			v, ok := c.Energized[Point{x, y}]
			switch {
			case !ok:
				sb.WriteString("?")

			case v[Up], v[Down], v[Right], v[Left]:
				sb.WriteString("#")

			default:
				sb.WriteString(".")
			}
		}

		fmt.Println(sb.String())
	}
}

func (c *Contraption) DebugPrintContraption() {
	fmt.Println("Width: ", c.Width)
	fmt.Println("Height: ", c.Height)

	fmt.Println("Points:")
	for y := 0; y < c.Height; y++ {
		sb := strings.Builder{}

		for x := 0; x < c.Width; x++ {
			t := c.Points[Point{x, y}]
			sb.WriteRune(rune(t))
		}

		fmt.Println(sb.String())
	}

	fmt.Println()
}

func (c *Contraption) Clone() *Contraption {
	m := make(map[Point]Tile, len(c.Points))
	energized := make(map[Point]map[Direction]bool)

	// copy points
	for k, v := range c.Points {
		m[k] = v
		energized[k] = make(map[Direction]bool)
	}

	// do not copy energized map

	return &Contraption{
		Points:    m,
		Energized: energized,
		Width:     c.Width,
		Height:    c.Height,
	}
}
