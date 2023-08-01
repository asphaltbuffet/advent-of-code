package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	open = "."
	wall = "#"
)

type tile struct {
	position point
	content  string
	up       *tile
	down     *tile
	left     *tile
	right    *tile
}

type point struct {
	x, y int
}

type board struct {
	tiles    map[point]*tile
	width    int
	height   int
	location *tile
	facing   int
	visited  []*tile
}

var (
	up    point = point{0, -1}
	down  point = point{0, 1}
	left  point = point{-1, 0}
	right point = point{1, 0}
)

func newBoard(mm [][]string) *board {
	b := board{
		tiles:    make(map[point]*tile),
		location: nil,
		facing:   0,
		visited:  []*tile{},
		width:    0,
		height:   len(mm),
	}

	// create tiles
	for y, row := range mm {
		if len(row) > b.width {
			b.width = len(row)
		}

		for x, content := range row {
			// ignore empty spaces
			if content == " " {
				continue
			}

			p := point{x, y}

			b.tiles[p] = &tile{
				position: p,
				content:  content,
			}

			// initially link each tile to itself
			b.tiles[p].up = b.tiles[p]
			b.tiles[p].down = b.tiles[p]
			b.tiles[p].left = b.tiles[p]
			b.tiles[p].right = b.tiles[p]
		}
	}

	// link tiles together
	for _, t := range b.tiles {
		if t.content == wall {
			continue
		}

		// check for adjacent open tiles
		if n := b.nextTile(t.position, up); n != nil {
			t.up = n
		}

		if n := b.nextTile(t.position, down); n != nil {
			t.down = n
		}

		if n := b.nextTile(t.position, left); n != nil {
			t.left = n
		}

		if n := b.nextTile(t.position, right); n != nil {
			t.right = n
		}
	}

	// get start tile
	for i := 0; i < b.width; i++ {
		if t, ok := b.tiles[point{i, 0}]; ok && t.content == open {
			b.location = t
			b.visited = append(b.visited, t)

			break
		}
	}

	return &b
}

func (b *board) move(path []string) {
	// fmt.Printf("path: %v\n", path)
	// fmt.Printf("start=%v\n", b.location.position)

	for _, p := range path {
		switch p {
		case "R":
			// fmt.Println("before turning right: ", b.facing)
			b.facing = (b.facing + 1) % 4

			// fmt.Println("after turning right: ", b.facing)
		case "L":
			// fmt.Println("before turning left: ", b.facing)

			b.facing = ((b.facing - 1) + 4) % 4

			// fmt.Println("after turning left: ", b.facing)
		default:
			// move forward
			n, _ := strconv.Atoi(p)

			// fmt.Printf("moving %s, facing %d\n", p, b.facing)

			for i := 0; i < n; i++ {
				// prev := b.location

				switch b.facing {
				case 0:
					b.location = b.location.right
				case 1:
					b.location = b.location.down
				case 2:
					b.location = b.location.left
				case 3:
					b.location = b.location.up
				}

				// if b.location != prev {
				// 	fmt.Printf("location: %v\n", b.location.position)
				// }

				// mark tile as visited
				// b.visited = append(b.visited, b.location)
			}
		}
	}
}

func (b *board) nextTile(p point, d point) *tile {
	nextPoint := p.add(d)
	nextPoint.y = (nextPoint.y + (b.height)) % (b.height)
	nextPoint.x = (nextPoint.x + (b.width)) % (b.width)

	if c, ok := b.tiles[nextPoint]; !ok {
		// empty space
		return b.nextTile(nextPoint, d)
	} else if c.content == wall {
		// wall; no connection
		return nil
	} else {
		// open space; connect to it
		return b.tiles[nextPoint]
	}
}

func (p point) add(b point) point {
	return point{p.x + b.x, p.y + b.y}
}

func (b *board) debugPrint() {
	for y := 0; y < b.height; y++ {
		var sb strings.Builder

		for x := 0; x < b.width; x++ {
			p := point{x, y}

			if t, ok := b.tiles[p]; ok {
				sb.WriteString(getBoxDrawing(t))
			} else {
				sb.WriteString("░")
			}
		}

		fmt.Println(sb.String())
	}
}

func getBoxDrawing(t *tile) string {
	if t.content == wall {
		return "·"
	}

	switch {
	case t != t.up && t != t.down && t != t.left && t != t.right:
		// 0 | 0 | 0 | 0
		return "┼"
	case t != t.up && t != t.down && t != t.left && t == t.right:
		// 0 | 0 | 0 | 1
		return "┤"
	case t != t.up && t != t.down && t == t.left && t != t.right:
		// 0 | 0 | 1 | 0
		return "├"
	case t != t.up && t != t.down && t == t.left && t == t.right:
		// 0 | 0 | 1 | 1
		return "│"
	case t != t.up && t == t.down && t != t.left && t != t.right:
		// 0 | 1 | 0 | 0
		return "┴"
	case t != t.up && t == t.down && t != t.left && t == t.right:
		// 0 | 1 | 0 | 1
		return "┘"
	case t != t.up && t == t.down && t == t.left && t != t.right:
		// 0 | 1 | 1 | 0
		return "└"
	case t != t.up && t == t.down && t == t.left && t == t.right:
		// 0 | 1 | 1 | 1
		return "╵"
	case t == t.up && t != t.down && t != t.left && t != t.right:
		// 1 | 0 | 0 | 0
		return "┬"
	case t == t.up && t != t.down && t != t.left && t == t.right:
		// 1 | 0 | 0 | 1
		return "┐"
	case t == t.up && t != t.down && t == t.left && t != t.right:
		// 1 | 0 | 1 | 0
		return "┌"
	case t == t.up && t != t.down && t == t.left && t == t.right:
		// 1 | 0 | 1 | 1
		return "╷"
	case t == t.up && t == t.down && t != t.left && t != t.right:
		// 1 | 1 | 0 | 0
		return "─"
	case t == t.up && t == t.down && t != t.left && t == t.right:
		// 1 | 1 | 0 | 1
		return "╴"
	case t == t.up && t == t.down && t == t.left && t != t.right:
		// 1 | 1 | 1 | 0
		return "╶"
	case t == t.up && t == t.down && t == t.left && t == t.right:
		// 1 | 1 | 1 | 1
		return "·"
	default:
		return "X"
	}
}

func parse(input string) ([][]string, []string) {
	mm, rawPath, _ := strings.Cut(input, "\n\n")

	lines := strings.Split(mm, "\n")

	m := make([][]string, len(lines))

	for i, line := range lines {
		m[i] = strings.Split(line, "")
	}

	re := regexp.MustCompile(`(\d+|[RL])`)
	path := re.FindAllString(rawPath, -1)

	return m, path
}
