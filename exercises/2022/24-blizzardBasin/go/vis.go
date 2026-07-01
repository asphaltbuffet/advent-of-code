package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the blizzard field the expedition must cross (GIF). Each frame is
// one minute: every blizzard is extrapolated to its position and drawn in a
// direction-coded color (→ red, ← blue, ↑ green, ↓ yellow); cells where two or
// more blizzards overlap are white, the momentary gaps being the only safe tiles
// to stand on. The entry and exit gaps in the walls are marked. The animation
// runs for the duration of the first crossing (the part-one answer).
func (e Exercise) Vis(instr, outdir string) error {
	b, err := parseInput(instr)
	if err != nil {
		return err
	}

	frames, err := startToEnd(b, 0)
	if err != nil {
		return err
	}
	if frames > 180 {
		frames = 180 // keep the GIF a sensible size
	}

	rows := b.totalRows
	cols := b.totalCols

	const cell = 6
	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 empty
		color.RGBA{0x22, 0x26, 0x33, 0xff}, // 1 wall
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 2 → east
		color.RGBA{0x40, 0x80, 0xff, 0xff}, // 3 ← west
		color.RGBA{0x3f, 0xd0, 0x7a, 0xff}, // 4 ↑ north
		color.RGBA{0xf0, 0xd0, 0x40, 0xff}, // 5 ↓ south
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // 6 overlap / portals
	}

	dirIdx := map[string]uint8{">": 2, "<": 3, "^": 4, "v": 5}

	// Grid is rows+2 tall (wall rows top and bottom) and cols+2 wide.
	gw := cols + 2
	gh := rows + 2

	anim := &gif.GIF{}
	for t := 0; t <= frames; t++ {
		img := image.NewPaletted(image.Rect(0, 0, gw*cell, gh*cell), pal)

		// Walls around the border.
		fill := func(gx, gy int, idx uint8) {
			for dy := 0; dy < cell; dy++ {
				for dx := 0; dx < cell; dx++ {
					img.SetColorIndex(gx*cell+dx, gy*cell+dy, idx)
				}
			}
		}
		for gx := 0; gx < gw; gx++ {
			fill(gx, 0, 1)
			fill(gx, gh-1, 1)
		}
		for gy := 0; gy < gh; gy++ {
			fill(0, gy, 1)
			fill(gw-1, gy, 1)
		}
		// Entry/exit portals (start at top, end at bottom).
		fill(b.start.x+1, 0, 6)
		fill(b.end.x+1, gh-1, 6)

		// Count blizzards per cell so overlaps show white.
		count := map[point]uint8{}
		dir := map[point]uint8{}
		for _, w := range b.winds {
			p := w.extrapolatePosition(t)
			count[p]++
			dir[p] = dirIdx[w.char]
		}
		for p, n := range count {
			idx := dir[p]
			if n > 1 {
				idx = 6
			}
			fill(p.x+1, p.y+1, idx)
		}

		anim.Image = append(anim.Image, img)
		delay := 8
		if t == 0 {
			delay = 120
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "blizzard-basin.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}
