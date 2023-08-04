package exercises

import (
	"image"
	"math"

	"github.com/asphaltbuffet/advent-of-code/pkg/pq"
	"github.com/asphaltbuffet/advent-of-code/pkg/set"
)

const (
	rightD direction = iota
	downD
	leftD
	upD
)

type direction int

type adjBlock struct {
	block *block
	dir   direction
}

type tileType rune

const (
	empty tileType = ' '
	open  tileType = '.'
	wall  tileType = '#'
)

type cmd int

const (
	move cmd = iota
	rotate
)

type action struct {
	command cmd
	value   int // if rotate, then -1 for L and +1 for R
}

func (dir direction) invert() direction {
	return (dir + 2) % 4
}

type stateD struct {
	row, col int
	face     direction
}

func followPath(b board, is3D bool) stateD {
	cur := b.start
	// start looking right
	dir := rightD
	st := stateD{cur[0], cur[1], dir}

	if is3D {
		// 3d needs the faces of the cube
		b.blocks = cubeFold(b)
	}

	// get movement instructions
	for _, inst := range b.actions {
		if inst.command == rotate {
			dir += direction(inst.value)
			if dir == -1 {
				dir = 3
			} else if dir == 4 {
				dir = 0
			}
		} else {
			if !is3D {
				// move by an amount
				cur = b.move2D(cur, dir, inst.value)
			} else {
				// move by an amount
				cur, dir = b.move3D(cur, dir, inst.value)
			}
		}

		st = stateD{cur[0], cur[1], dir}
	}

	return st
}

// turn the 2d map into 6 cube faces of {size}x{size}
func cubeFold(b board) map[image.Rectangle]*block {
	blocks := map[image.Rectangle]*block{}

	for row := 0; row <= b.height; row += b.blockSize {
		for col := 0; col <= b.width; col += b.blockSize {
			cell := b.get(row, col)
			if cell == open || cell == wall {
				min := image.Point{row, col}
				max := image.Point{row + b.blockSize, col + b.blockSize}
				rect := image.Rectangle{min, max}
				blocks[rect] = &block{
					dims:           rect,
					adjascentZones: set.Set[image.Rectangle]{},
				}
			}
		}
	}

	// right, down, left, up (like direction type)
	adjascent := [4]image.Point{
		{0, b.blockSize},
		{b.blockSize, 0},
		{0, -b.blockSize},
		{-b.blockSize, 0},
	}

	type state struct {
		block          *block
		fromDir, toDir direction
		currentRect    image.Rectangle
		priority       int
	}

	// BFS for adjascent blocks on edges
	queue := make(pq.PriorityQueue[state], len(blocks)*len(adjascent))

	index := 0
	// populate priority queue
	for rect, block := range blocks {
		for i, vector := range adjascent {
			dir := direction(i)
			cur := rect.Add(vector)

			queue.NewItem(&state{
				block:       block,
				fromDir:     dir,
				toDir:       dir,
				currentRect: cur,
				priority:    0,
			}, 0, index)
			index++
		}
	}

	queue.Init()

	adjascentCount := 0
	// each face has 4 adjascent blocks; stop when we've found them all
	maxAdjascentCount := 6 * 4

	for adjascentCount != maxAdjascentCount {
		cur := queue.Get()

		// don't proceed if all adjascent blocks are found for this block
		if len(cur.block.adjascentZones) == 4 {
			continue
		}

		block, isNonEmpty := blocks[cur.currentRect]

		if block == cur.block {
			// thou shalt not be neighbors with thyself
			continue
		}

		if isNonEmpty {
			if cur.block.addAdjascent(block, cur.fromDir, cur.toDir) {
				// added both
				adjascentCount += 2
			}

			continue
		}

		// else, continue path-finding
		for i, vec := range adjascent {
			dir := direction(i)
			nextRect := cur.currentRect.Add(vec)

			queue.PushValue(&state{
				block:       cur.block,
				fromDir:     cur.fromDir,
				toDir:       dir,
				currentRect: nextRect,
				priority:    cur.priority + 1,
			}, cur.priority+1)
		}
	}

	return blocks
}

func (b board) move3D(cur [2]int, dir direction, steps int) ([2]int, direction) {
	for i := 0; i < steps; i++ {
		next, newDir := b.getNext3D(cur[0], cur[1], dir)

		if next == cur {
			break
		}

		cur = next
		dir = newDir
	}

	return cur, dir
}

func moveOne(pos *[2]int, dir direction) {
	p := image.Point{pos[1], pos[0]}
	switch dir {
	case rightD:
		p.X++
	case leftD:
		p.X--
	case upD:
		p.Y--
	case downD:
		p.Y++
	}
	pos[0], pos[1] = p.Y, p.X
}

func (b *board) getNext3D(r, c int, dir direction) ([2]int, direction) {
	orig := [2]int{r, c}

	// assume we'll keep going this direction
	newDir := dir

	next := orig
	moveOne(&next, dir)

	// for loop because the `empty` case may hit a wall
	for {
		var tile tileType
		// check if out of bounds
		if next[0] < 0 || next[1] < 0 || next[0] >= b.height || next[1] >= b.width {
			// next cell is empty
			tile = empty
		} else {
			// we can get next cell
			tile = b.get(next[0], next[1])
		}

		switch tile {
		case open:
			return next, newDir
		case wall:
			// return original (didn't move)
			return orig, dir
		case empty:
			// move to new cube face in another direction
			next, newDir = b.rotateTile(orig, dir)
		}
	}
}

func (b board) rotateTile(pos [2]int, dir direction) ([2]int, direction) {
	curPoint := image.Point{pos[0], pos[1]}

	for rect, block := range b.blocks {
		if curPoint.In(rect) {
			var a *adjBlock

			switch dir {
			case upD:
				a = block.up
			case downD:
				a = block.down
			case leftD:
				a = block.left
			case rightD:
				a = block.right
			}

			newDir := a.dir

			// rotate current position clockwise until current dir lines up with new dir
			rotations := int(newDir - dir)

			size := rect.Size()

			// rotate a point around an origin algorithm
			px, py := float64(curPoint.X), float64(curPoint.Y)
			origin := rect.Min.Add(size.Div(2))
			ox, oy := float64(origin.X)-0.5, float64(origin.Y)-0.5

			// angle in radians; invert to go clockwise
			theta := -(math.Pi / 2) * float64(rotations)

			// need to round; `int()` only floors
			rx := math.Round(math.Cos(theta)*(px-ox) - math.Sin(theta)*(py-oy) + ox)
			ry := math.Round(math.Sin(theta)*(px-ox) + math.Cos(theta)*(py-oy) + oy)

			next := [2]int{int(rx), int(ry)}

			// move in new direction by one
			moveOne(&next, newDir)

			// annoying transitions here
			curPoint = image.Point{next[0], next[1]}

			// mod to next cube face (mod is virtually teleporting)
			curPoint = curPoint.Mod(a.block.dims)

			// set next for outer for loop to check for walls
			next = [2]int{curPoint.X, curPoint.Y}

			return next, newDir
		}
	}

	return [2]int{}, 0
}
