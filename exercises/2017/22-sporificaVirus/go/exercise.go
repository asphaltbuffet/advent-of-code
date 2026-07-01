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

// Exercise for Advent of Code 2017 day 22.
type Exercise struct {
	common.BaseExercise
}

type pt struct{ x, y int }

// Node states.
const (
	clean = iota
	weakened
	infected
	flagged
)

// Directions, clockwise from up: up, right, down, left.
var dirs = [4]pt{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

// parseGrid loads the starting infected nodes into a state map, centred so the
// carrier begins at the origin. It returns the map and the carrier's start.
func parseGrid(instr string) map[pt]int {
	lines := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	rows := len(lines)
	cols := len(lines[0])
	grid := map[pt]int{}
	for r, line := range lines {
		for c, ch := range line {
			if ch == '#' {
				grid[pt{c - cols/2, r - rows/2}] = infected
			}
		}
	}
	return grid
}

// One runs 10000 bursts of the two-state virus and returns how many caused an
// infection.
func (e Exercise) One(instr string) (any, error) {
	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0 // up
	infections := 0

	for i := 0; i < 10000; i++ {
		if grid[pos] == infected {
			dir = (dir + 1) & 3 // turn right
			delete(grid, pos)   // becomes clean
		} else {
			dir = (dir + 3) & 3 // turn left
			grid[pos] = infected
			infections++
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}
	}
	return infections, nil
}

// Vis animates the four-state virus spreading (GIF). It runs the burst sequence
// once, snapshotting the grid at even intervals; each frame colours every node
// by state (weakened / infected / flagged) over clean space, sizing the canvas
// to the region the infection ultimately reaches.
func (e Exercise) Vis(instr, outdir string) error {
	const bursts = 400_000 // enough to grow a rich bloom without a huge canvas
	const wantFrames = 300

	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0

	type snap struct {
		state map[pt]int
		pos   pt
	}
	var snaps []snap
	interval := bursts / wantFrames

	minX, maxX, minY, maxY := 0, 0, 0, 0
	for i := 0; i < bursts; i++ {
		switch grid[pos] {
		case clean:
			dir = (dir + 3) & 3
			grid[pos] = weakened
		case weakened:
			grid[pos] = infected
		case infected:
			dir = (dir + 1) & 3
			grid[pos] = flagged
		case flagged:
			dir = (dir + 2) & 3
			delete(grid, pos)
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}

		if pos.x < minX {
			minX = pos.x
		}
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.y > maxY {
			maxY = pos.y
		}

		if (i+1)%interval == 0 {
			cp := make(map[pt]int, len(grid))
			for k, v := range grid {
				cp[k] = v
			}
			snaps = append(snaps, snap{cp, pos})
		}
	}

	w := maxX - minX + 1
	h := maxY - minY + 1
	// Scale so the longer side is roughly 720px, at an integer cell size.
	cell := 720 / max(w, h)
	if cell < 1 {
		cell = 1
	}
	W, H := w*cell, h*cell

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 clean
		color.RGBA{0x3a, 0x6a, 0x9a, 0xff}, // 1 weakened
		color.RGBA{0x2f, 0xd0, 0x9a, 0xff}, // 2 infected
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 3 flagged
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 4 carrier
	}

	anim := &gif.GIF{}
	for si, s := range snaps {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		fillCell := func(gx, gy, idx int) {
			x0, y0 := (gx-minX)*cell, (gy-minY)*cell
			for dy := 0; dy < cell; dy++ {
				for dx := 0; dx < cell; dx++ {
					img.SetColorIndex(x0+dx, y0+dy, uint8(idx))
				}
			}
		}
		for p, st := range s.state {
			fillCell(p.x, p.y, st) // weakened=1, infected=2, flagged=3
		}
		fillCell(s.pos.x, s.pos.y, 4) // carrier on top
		anim.Image = append(anim.Image, img)
		delay := 4
		if si == len(snaps)-1 {
			delay = 200
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "sporifica-virus.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Two runs 10 million bursts of the four-state virus and returns how many
// bursts turned a node infected.
func (e Exercise) Two(instr string) (any, error) {
	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0
	infections := 0

	for i := 0; i < 10_000_000; i++ {
		switch grid[pos] {
		case clean:
			dir = (dir + 3) & 3 // turn left
			grid[pos] = weakened
		case weakened:
			grid[pos] = infected // no turn
			infections++
		case infected:
			dir = (dir + 1) & 3 // turn right
			grid[pos] = flagged
		case flagged:
			dir = (dir + 2) & 3 // reverse
			delete(grid, pos)   // back to clean
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}
	}
	return infections, nil
}
