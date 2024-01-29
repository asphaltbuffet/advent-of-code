package exercises

import (
	"fmt"
	"strings"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Point util.Point2D

type GardenMap struct {
	Tiles    map[Point]*Tile
	Start    *Point
	Boundary Point
	Steps    int
}

type Tile struct {
	Point Point
	Type  TileType
}

type TileType rune

const (
	Garden TileType = '.'
	Rock   TileType = '#'
	Start  TileType = 'S'
)

type Direction int

const (
	Up Direction = iota + 1
	Right
	Down
	Left
)

var Moves = map[Direction]Point{
	Up:    {X: 0, Y: -1},
	Right: {X: 1, Y: 0},
	Down:  {X: 0, Y: 1},
	Left:  {X: -1, Y: 0},
}

func parseInput(in string) (map[Point]TileType, Point, int, error) {
	lines := strings.Split(in, "\n")
	if len(lines) != len(lines[0]) {
		return nil, Point{}, 0, fmt.Errorf("input is not square")
	}

	var (
		tiles = map[Point]TileType{}
		dim   = len(lines) // this assumes the input is square

		start Point
	)

	for y, line := range lines {
		for x, r := range line {
			switch t := TileType(r); t {
			case Start:
				start = Point{x, y}
				fallthrough

			case Garden, Rock:
				tiles[Point{x, y}] = t

			default:
				return nil, Point{}, 0, fmt.Errorf("unknown tile type: %c", r)
			}
		}
	}

	return tiles, start, dim, nil
}

func manhattanDistance(p1, p2 *Point) int {
	return util.AbsInt(p1.X-p2.X) + util.AbsInt(p1.Y-p2.Y)
}

// wrap converts a point to it's equivalent point within the map bounds.
func (p *Point) wrap(max Point) Point {
	nx := (p.X%max.X + max.X) % max.X
	ny := (p.Y%max.Y + max.Y) % max.Y

	return Point{X: nx, Y: ny}
}
