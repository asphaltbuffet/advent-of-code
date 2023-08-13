package exercises

import (
	"image"
	"regexp"
	"strconv"
	"strings"
)

type board struct {
	width     int
	height    int
	grid      map[int]map[int]tileType
	start     [2]int
	actions   []action
	blockSize int
	blocks    map[image.Rectangle]*block
}

func parse(instr string, squareSize int) board {
	rawMap, rawPath, _ := strings.Cut(instr, "\n\n")
	width := 0
	lines := strings.Split(rawMap, "\n")

	b := board{
		grid:      map[int]map[int]tileType{},
		height:    len(lines),
		blockSize: squareSize,
	}

	var start *[2]int

	for row, line := range lines {
		if len(line) > width {
			width = len(line)
		}

		for col, c := range line {
			t := tileType(c)
			switch t {
			case empty:
				continue
			case wall, open:
				b.set(row, col, t)

				if start == nil && t == open {
					start = &[2]int{row, col}
				}
			}
		}
	}

	b.width = width
	b.start = *start

	b.parsePath(rawPath)

	return b
}

func (b *board) parsePath(rawPath string) {
	actions := []action{}

	re := regexp.MustCompile(`(\d+|[RL])`)
	tokens := re.FindAllString(rawPath, -1)

	for _, t := range tokens {
		switch t {
		case "L":
			actions = append(actions, action{
				command: rotate,
				value:   -1,
			})
		case "R":
			actions = append(actions, action{
				command: rotate,
				value:   1,
			})
		default:
			value, _ := strconv.Atoi(t)
			actions = append(actions, action{
				command: move,
				value:   value,
			})
		}
	}

	b.actions = actions
}

func (b *board) set(r, c int, val tileType) {
	grid := b.grid
	if grid[r] == nil {
		grid[r] = map[int]tileType{}
	}

	grid[r][c] = val
}

func (b *board) get(r, c int) (val tileType) {
	t, ok := b.grid[r][c]

	if !ok {
		return empty
	}

	return t
}

func (b *board) getNext2D(r, c int, dir direction) ([2]int, bool) {
	cur := [2]int{r, c}

	for {
		next := cur
		moveOne(&next, dir)

		// negative numbers
		if next[0] < 0 {
			next[0] = b.height + next[0]
		} else if next[1] < 0 {
			next[1] = b.width + next[1]
		}

		next[0] %= b.height
		next[1] %= b.width

		cell := b.get(next[0], next[1])

		// we ignored empty here and just kept looping
		if cell == open {
			return next, false
		}

		if cell == wall {
			// return original
			return [2]int{r, c}, true
		}

		cur = next
	}
}

func (b board) move2D(cur [2]int, dir direction, steps int) [2]int {
	for i := 0; i < steps; i++ {
		next, hitWall := b.getNext2D(cur[0], cur[1], dir)

		if hitWall {
			break
		}

		cur = next
	}

	return cur
}
