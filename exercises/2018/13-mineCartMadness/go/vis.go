package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"sort"
)

// Vis animates the carts running the track (GIF). The track is drawn dim, live
// carts bright, and every crash flashes a bright marker that lingers for the rest
// of the run. Because there are thousands of ticks, frames are sampled — but any
// tick with a crash is always captured, so the four "expensive crashes" and the
// final lone survivor all appear. Carts and crashes are carried by brightness (not
// hue alone) against the dim track, so the animation reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	grid, carts := parse(instr)

	rows := len(grid)
	cols := 0
	for _, r := range grid {
		if len(r) > cols {
			cols = len(r)
		}
	}

	// Palette indices: 0 background, 1 track (very dim), 2 cart, 3 crash, 4 halo.
	pal := color.Palette{
		color.RGBA{0x0a, 0x0c, 0x12, 0xff}, // background
		color.RGBA{0x26, 0x2c, 0x38, 0xff}, // track (very dim)
		color.RGBA{0x8f, 0xd4, 0xff, 0xff}, // cart (bright sky blue)
		color.RGBA{0xff, 0x8a, 0x00, 0xff}, // crash (bright orange)
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // halo (black, to lift markers off track)
	}
	const (
		bgIdx    = 0
		trackIdx = 1
		cartIdx  = 2
		crashIdx = 3
		haloIdx  = 4
	)

	const (
		cell = 4
		pad  = 6
	)
	W := cols*cell + 2*pad
	H := rows*cell + 2*pad

	// Persistent crash markers accumulate across the whole run.
	crashes := map[[2]int]bool{}

	block := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := 0; dy < cell; dy++ {
			for dx := 0; dx < cell; dx++ {
				img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
			}
		}
	}

	// marker draws a bright block ringed by a dark halo so events stay legible
	// against the busy track, in color and in grayscale alike.
	marker := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := -1; dy < cell+1; dy++ {
			for dx := -1; dx < cell+1; dx++ {
				px, py := pad+x*cell+dx, pad+y*cell+dy
				if px < 0 || py < 0 || px >= W || py >= H {
					continue
				}
				if dy < 0 || dy >= cell || dx < 0 || dx >= cell {
					img.SetColorIndex(px, py, haloIdx)
				} else {
					img.SetColorIndex(px, py, idx)
				}
			}
		}
	}

	render := func() *image.Paletted {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		// background
		for i := range img.Pix {
			img.Pix[i] = bgIdx
		}
		// track
		for y := 0; y < rows; y++ {
			for x := 0; x < len(grid[y]); x++ {
				if grid[y][x] != ' ' && grid[y][x] != 0 {
					block(img, x, y, trackIdx)
				}
			}
		}
		// crash markers (haloed)
		for p := range crashes {
			marker(img, p[0], p[1], crashIdx)
		}
		// live carts on top (haloed)
		for _, c := range carts {
			if !c.dead {
				marker(img, c.x, c.y, cartIdx)
			}
		}
		return img
	}

	anim := &gif.GIF{}
	push := func(delay int) {
		anim.Image = append(anim.Image, render())
		anim.Delay = append(anim.Delay, delay)
	}

	push(40) // opening frame

	const sampleEvery = 60
	tick := 0
	for {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y != carts[j].y {
				return carts[i].y < carts[j].y
			}
			return carts[i].x < carts[j].x
		})

		crashedThisTick := false
		for _, c := range carts {
			if c.dead {
				continue
			}
			c.advance(grid)
			for _, other := range carts {
				if other == c || other.dead || other.x != c.x || other.y != c.y {
					continue
				}
				c.dead, other.dead = true, true
				crashes[[2]int{c.x, c.y}] = true
				crashedThisTick = true
				break
			}
		}

		// Drop dead carts.
		live := carts[:0]
		for _, c := range carts {
			if !c.dead {
				live = append(live, c)
			}
		}
		carts = live
		tick++

		// Always capture crash ticks; otherwise sample.
		if crashedThisTick {
			push(120) // linger on a crash
		} else if tick%sampleEvery == 0 {
			push(6)
		}

		if len(carts) <= 1 {
			push(400) // hold on the final survivor
			break
		}
	}

	f, err := os.Create(filepath.Join(outdir, "mine-cart-madness.gif"))
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, anim)
}
