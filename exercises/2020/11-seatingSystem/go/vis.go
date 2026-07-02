package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the seats settling under the Part Two line-of-sight rules, which
// let a seat's influence travel across aisles so the pattern settles less locally
// than Part One. Floor is dark, empty seats are mid blue, and occupied seats are
// bright yellow, three tones with a wide brightness gap so the states stay
// distinct in grayscale. The final settled frame holds so the answer is readable.
func (e Exercise) Vis(instr, outdir string) error {
	g := parseGrid(instr)

	const scale = 5
	W, H := g.w*scale, g.h*scale

	// Palette index 0 = floor, 1 = empty seat, 2 = occupied seat.
	pal := color.Palette{
		color.RGBA{0x14, 0x18, 0x1e, 0xff}, // floor (dark)
		color.RGBA{0x00, 0x72, 0xB2, 0xff}, // empty seat (blue)
		color.RGBA{0xF0, 0xE4, 0x42, 0xff}, // occupied seat (bright yellow)
	}

	cur := make([]byte, len(g.cells))
	copy(cur, g.cells)
	next := make([]byte, len(g.cells))
	work := grid{cur, g.w, g.h}

	var frames []*image.Paletted
	var delays []int
	capture := func() {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		for y := 0; y < g.h; y++ {
			for x := 0; x < g.w; x++ {
				var idx uint8
				switch work.at(x, y) {
				case empty:
					idx = 1
				case occupied:
					idx = 2
				}
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.SetColorIndex(x*scale+dx, y*scale+dy, idx)
					}
				}
			}
		}
		frames = append(frames, img)
		delays = append(delays, 12)
	}

	capture()
	iter := 0
	for {
		changed := false
		for y := 0; y < g.h; y++ {
			for x := 0; x < g.w; x++ {
				c := work.at(x, y)
				nc := c
				switch c {
				case empty:
					if visibleOccupied(work, x, y) == 0 {
						nc = occupied
					}
				case occupied:
					if visibleOccupied(work, x, y) >= 5 {
						nc = empty
					}
				}
				next[y*g.w+x] = nc
				if nc != c {
					changed = true
				}
			}
		}
		copy(cur, next)
		iter++
		// Capture every frame while the layout churns; sample the settling tail.
		if !changed || iter <= 12 || iter%2 == 0 {
			capture()
		}
		if !changed {
			break
		}
	}
	delays[len(delays)-1] = 400

	f, err := os.Create(filepath.Join(outdir, "seating-system.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, &gif.GIF{Image: frames, Delay: delays})
}
