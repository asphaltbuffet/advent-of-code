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

// Exercise for Advent of Code 2015 day 18.
type Exercise struct {
	common.BaseExercise
}

// grid is a square light grid; true is on.
type grid [][]bool

// parse reads the grid. The step count is not in the input and the AoC example
// uses different counts per part (4 for part 1, 5 for part 2 with stuck
// corners); the real 100x100 input runs 100 for both.
func parse(instr string) grid {
	var g grid
	for _, line := range strings.Fields(instr) {
		row := make([]bool, len(line))
		for i, c := range line {
			row[i] = c == '#'
		}
		g = append(g, row)
	}
	return g
}

// steps returns the animation length: the small AoC example runs 4 steps for
// part 1 and 5 for part 2; the real input runs 100 for both.
func (g grid) steps(stuck bool) int {
	if len(g) > 6 {
		return 100
	}
	if stuck {
		return 5
	}
	return 4
}

// neighborsOn counts the lit cells among the eight neighbours of (r, c).
func (g grid) neighborsOn(r, c int) int {
	n := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := r+dr, c+dc
			if nr >= 0 && nr < len(g) && nc >= 0 && nc < len(g[nr]) && g[nr][nc] {
				n++
			}
		}
	}
	return n
}

// step applies one Game of Life generation, returning a new grid.
func (g grid) step() grid {
	next := make(grid, len(g))
	for r := range g {
		next[r] = make([]bool, len(g[r]))
		for c := range g[r] {
			on := g.neighborsOn(r, c)
			if g[r][c] {
				next[r][c] = on == 2 || on == 3
			} else {
				next[r][c] = on == 3
			}
		}
	}
	return next
}

// stickCorners forces the four corner lights on (part two).
func (g grid) stickCorners() {
	last := len(g) - 1
	g[0][0], g[0][last] = true, true
	g[last][0], g[last][last] = true, true
}

func (g grid) countOn() int {
	n := 0
	for _, row := range g {
		for _, on := range row {
			if on {
				n++
			}
		}
	}
	return n
}

// run animates the grid for the given number of steps; if stuck is set the four
// corners are forced on initially and after every step.
func run(g grid, steps int, stuck bool) int {
	if stuck {
		g.stickCorners()
	}
	for i := 0; i < steps; i++ {
		g = g.step()
		if stuck {
			g.stickCorners()
		}
	}
	return g.countOn()
}

// One returns the number of lights on after animating the grid.
func (e Exercise) One(instr string) (any, error) {
	g := parse(instr)
	return run(g, g.steps(false), false), nil
}

// Two returns the count with the four corners stuck on.
func (e Exercise) Two(instr string) (any, error) {
	g := parse(instr)
	return run(g, g.steps(true), true), nil
}

// Visualization constants: each light is rendered as a cellPx square block and
// every frame is held for frameDelay hundredths of a second (0.5s).
const (
	cellPx     = 4
	frameDelay = 50 // 1/100s units -> 0.5s per frame
	gifName    = "yard.gif"
)

// frame renders the grid as a paletted image: lit cells bright, dark otherwise.
func (g grid) frame(palette color.Palette) *image.Paletted {
	h, w := len(g), len(g[0])
	img := image.NewPaletted(image.Rect(0, 0, w*cellPx, h*cellPx), palette)

	for r := 0; r < h; r++ {
		for c := 0; c < len(g[r]); c++ {
			idx := uint8(0)
			if g[r][c] {
				idx = 1
			}
			for py := 0; py < cellPx; py++ {
				for px := 0; px < cellPx; px++ {
					img.SetColorIndex(c*cellPx+px, r*cellPx+py, idx)
				}
			}
		}
	}

	return img
}

// Vis animates the (corner-stuck) yard and writes an animated GIF to outdir,
// one frame per step held for 0.5s.
func (e Exercise) Vis(instr string, outdir string) error {
	g := parse(instr)
	steps := g.steps(true)

	palette := color.Palette{
		color.RGBA{R: 0x0a, G: 0x0a, B: 0x1e, A: 0xff}, // off: dark navy
		color.RGBA{R: 0xff, G: 0xd1, B: 0x66, A: 0xff}, // on: warm yellow
	}

	anim := &gif.GIF{}
	appendFrame := func() {
		anim.Image = append(anim.Image, g.frame(palette))
		anim.Delay = append(anim.Delay, frameDelay)
	}

	// Capture the initial state (corners already stuck), then each generation.
	g.stickCorners()
	appendFrame()
	for i := 0; i < steps; i++ {
		g = g.step()
		g.stickCorners()
		appendFrame()
	}

	f, err := os.Create(filepath.Join(outdir, gifName))
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, anim)
}
