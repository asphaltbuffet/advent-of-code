package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 15.
type Exercise struct {
	common.BaseExercise
}

var moveDelta = map[byte][2]int{
	'^': {-1, 0},
	'v': {1, 0},
	'<': {0, -1},
	'>': {0, 1},
}

// warehouse is a mutable grid plus the robot position.
type warehouse struct {
	grid [][]byte
	r, c int
}

// parse splits the input into the warehouse grid and the concatenated moves.
func parse(instr string) (warehouse, string) {
	parts := strings.SplitN(strings.TrimRight(instr, "\n"), "\n\n", 2)
	var grid [][]byte
	for _, line := range strings.Split(parts[0], "\n") {
		grid = append(grid, []byte(line))
	}
	moves := strings.ReplaceAll(parts[1], "\n", "")
	w := warehouse{grid: grid}
	w.r, w.c = w.findRobot()
	return w, moves
}

func (w *warehouse) findRobot() (int, int) {
	for r := range w.grid {
		for c := range w.grid[r] {
			if w.grid[r][c] == '@' {
				return r, c
			}
		}
	}
	return -1, -1
}

// widen doubles the grid width per Part Two's expansion rules.
func (w *warehouse) widen() {
	out := make([][]byte, len(w.grid))
	for r := range w.grid {
		var row []byte
		for _, ch := range w.grid[r] {
			switch ch {
			case '#':
				row = append(row, '#', '#')
			case 'O':
				row = append(row, '[', ']')
			case '.':
				row = append(row, '.', '.')
			case '@':
				row = append(row, '@', '.')
			}
		}
		out[r] = row
	}
	w.grid = out
	w.r, w.c = w.findRobot()
}

// step applies one move, pushing boxes as needed. Handles both narrow ('O')
// and wide ('[' ']') boxes.
func (w *warehouse) step(move byte) {
	d, ok := moveDelta[move]
	if !ok {
		return
	}
	dr, dc := d[0], d[1]

	if dr == 0 {
		// Horizontal: scan a chain of boxes (O or [ ]).
		nc := w.c + dc
		for w.grid[w.r][nc] == 'O' || w.grid[w.r][nc] == '[' || w.grid[w.r][nc] == ']' {
			nc += dc
		}
		if w.grid[w.r][nc] == '#' {
			return
		}
		for x := nc; x != w.c; x -= dc {
			w.grid[w.r][x] = w.grid[w.r][x-dc]
		}
		w.grid[w.r][w.c] = '.'
		w.c += dc
		return
	}

	// Vertical.
	ahead := w.grid[w.r+dr][w.c]
	switch ahead {
	case '#':
		return
	case '.':
		w.grid[w.r][w.c] = '.'
		w.r += dr
		w.grid[w.r][w.c] = '@'
		return
	case 'O':
		// Narrow box chain straight up/down.
		nr := w.r + dr
		for w.grid[nr][w.c] == 'O' {
			nr += dr
		}
		if w.grid[nr][w.c] == '#' {
			return
		}
		w.grid[nr][w.c] = 'O'
		w.grid[w.r][w.c] = '.'
		w.r += dr
		w.grid[w.r][w.c] = '@'
		return
	}

	// Wide box: collect the connected set, abort if any blocked.
	boxes := make(map[[2]int]bool)
	if !w.collectVertical(w.r+dr, w.c, dr, boxes) {
		return
	}
	for b := range boxes {
		w.grid[b[0]][b[1]] = '.'
		w.grid[b[0]][b[1]+1] = '.'
	}
	for b := range boxes {
		w.grid[b[0]+dr][b[1]] = '['
		w.grid[b[0]+dr][b[1]+1] = ']'
	}
	w.grid[w.r][w.c] = '.'
	w.r += dr
	w.grid[w.r][w.c] = '@'
}

// collectVertical gathers every wide box (keyed by its '[' cell) that must move
// when pushing into (r,c) vertically; returns false if any hits a wall.
func (w *warehouse) collectVertical(r, c, dr int, boxes map[[2]int]bool) bool {
	if w.grid[r][c] == ']' {
		c--
	}
	if w.grid[r][c] != '[' {
		return true
	}
	if boxes[[2]int{r, c}] {
		return true
	}
	boxes[[2]int{r, c}] = true
	for _, cc := range []int{c, c + 1} {
		switch w.grid[r+dr][cc] {
		case '#':
			return false
		case '[', ']':
			if !w.collectVertical(r+dr, cc, dr, boxes) {
				return false
			}
		}
	}
	return true
}

// gps sums the GPS coordinates of all boxes (O or the '[' of wide boxes).
func (w *warehouse) gps() int {
	sum := 0
	for r := range w.grid {
		for c := range w.grid[r] {
			if w.grid[r][c] == 'O' || w.grid[r][c] == '[' {
				sum += 100*r + c
			}
		}
	}
	return sum
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	w, moves := parse(instr)
	for i := 0; i < len(moves); i++ {
		w.step(moves[i])
	}
	return w.gps(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	w, moves := parse(instr)
	w.widen()
	for i := 0; i < len(moves); i++ {
		w.step(moves[i])
	}
	return w.gps(), nil
}

// --- Visualization (animated GIF) ---

// visPalette maps warehouse cell bytes to colours.
var visPalette = color.Palette{
	color.RGBA{0x0a, 0x0a, 0x12, 0xff}, // 0 floor / background
	color.RGBA{0x3a, 0x3a, 0x55, 0xff}, // 1 wall
	color.RGBA{0xe0, 0x9b, 0x3e, 0xff}, // 2 box
	color.RGBA{0x4d, 0xe0, 0x8a, 0xff}, // 3 robot
}

func paletteIndex(ch byte) uint8 {
	switch ch {
	case '#':
		return 1
	case 'O', '[', ']':
		return 2
	case '@':
		return 3
	default:
		return 0
	}
}

// renderFrame draws the current warehouse into a paletted image.
func (w *warehouse) renderFrame(scale int) *image.Paletted {
	rows := len(w.grid)
	cols := 0
	for _, row := range w.grid {
		if len(row) > cols {
			cols = len(row)
		}
	}
	img := image.NewPaletted(image.Rect(0, 0, cols*scale, rows*scale), visPalette)
	for r := 0; r < rows; r++ {
		for c := 0; c < len(w.grid[r]); c++ {
			idx := paletteIndex(w.grid[r][c])
			x0, y0 := c*scale, r*scale
			for y := y0; y < y0+scale; y++ {
				for x := x0; x < x0+scale; x++ {
					img.SetColorIndex(x, y, idx)
				}
			}
		}
	}
	return img
}

// animate runs the moves, capturing a frame every `sample` moves, and writes a
// GIF to outdir/name.
func animate(w warehouse, moves string, scale, sample int, path string) error {
	anim := &gif.GIF{}
	add := func() {
		anim.Image = append(anim.Image, w.renderFrame(scale))
		anim.Delay = append(anim.Delay, 4) // 40ms per frame
	}
	add() // initial state
	for i := 0; i < len(moves); i++ {
		w.step(moves[i])
		if i%sample == 0 {
			add()
		}
	}
	add() // final state
	anim.Delay[len(anim.Delay)-1] = 300 // hold the final frame

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

// Vis emits two animated GIFs: the narrow (Part 1) and wide (Part 2) warehouses.
func (e Exercise) Vis(instr string, outdir string) error {
	// Sample so even ~20k-move inputs produce a manageable GIF (<~400 frames).
	w1, moves := parse(instr)
	sample := len(moves)/300 + 1

	if err := animate(w1, moves, 8, sample, filepath.Join(outdir, "warehouse-p1.gif")); err != nil {
		return err
	}

	w2, _ := parse(instr)
	w2.widen()
	return animate(w2, moves, 6, sample, filepath.Join(outdir, "warehouse-p2.gif"))
}
