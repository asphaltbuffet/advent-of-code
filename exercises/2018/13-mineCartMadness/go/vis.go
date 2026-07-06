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
	cols := gridWidth(grid)

	const (
		cell = 4
		pad  = 6
	)
	imgW := cols*cell + 2*pad
	imgH := rows*cell + 2*pad

	pal := cartPalette()
	crashes := map[[2]int]bool{}

	anim := &gif.GIF{}
	appendFrame(anim, renderCartFrame(grid, carts, crashes, pal, rows, imgW, imgH), 40)

	const sampleEvery = 60
	tick := 0
	for {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y != carts[j].y {
				return carts[i].y < carts[j].y
			}
			return carts[i].x < carts[j].x
		})

		crashedThisTick := tickCarts(carts, grid, crashes)

		live := carts[:0]
		for _, c := range carts {
			if !c.dead {
				live = append(live, c)
			}
		}
		carts = live
		tick++

		if crashedThisTick {
			appendFrame(anim, renderCartFrame(grid, carts, crashes, pal, rows, imgW, imgH), 120)
		} else if tick%sampleEvery == 0 {
			appendFrame(anim, renderCartFrame(grid, carts, crashes, pal, rows, imgW, imgH), 6)
		}

		if len(carts) <= 1 {
			appendFrame(anim, renderCartFrame(grid, carts, crashes, pal, rows, imgW, imgH), 400)
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

func gridWidth(grid [][]byte) int {
	cols := 0
	for _, r := range grid {
		if len(r) > cols {
			cols = len(r)
		}
	}
	return cols
}

func cartPalette() color.Palette {
	return color.Palette{
		color.RGBA{0x0a, 0x0c, 0x12, 0xff}, // 0 background
		color.RGBA{0x26, 0x2c, 0x38, 0xff}, // 1 track
		color.RGBA{0x8f, 0xd4, 0xff, 0xff}, // 2 cart
		color.RGBA{0xff, 0x8a, 0x00, 0xff}, // 3 crash
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // 4 halo
	}
}

func appendFrame(anim *gif.GIF, img *image.Paletted, delay int) {
	anim.Image = append(anim.Image, img)
	anim.Delay = append(anim.Delay, delay)
}

func cartBlock(img *image.Paletted, x, y, cell, pad int, idx uint8) {
	for dy := range cell {
		for dx := range cell {
			img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
		}
	}
}

func cartMarker(img *image.Paletted, x, y, cell, pad, imgW, imgH int, idx uint8) {
	const haloIdx = 4
	for dy := -1; dy < cell+1; dy++ {
		for dx := -1; dx < cell+1; dx++ {
			px, py := pad+x*cell+dx, pad+y*cell+dy
			if px < 0 || py < 0 || px >= imgW || py >= imgH {
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

func renderCartFrame(
	grid [][]byte, carts []*cart, crashes map[[2]int]bool, pal color.Palette, rows, imgW, imgH int,
) *image.Paletted {
	const (
		bgIdx    = 0
		trackIdx = 1
		cartIdx  = 2
		crashIdx = 3
		cell     = 4
		pad      = 6
	)
	img := image.NewPaletted(image.Rect(0, 0, imgW, imgH), pal)
	for i := range img.Pix {
		img.Pix[i] = bgIdx
	}
	for y := range rows {
		for x := range len(grid[y]) {
			if grid[y][x] != ' ' && grid[y][x] != 0 {
				cartBlock(img, x, y, cell, pad, trackIdx)
			}
		}
	}
	for p := range crashes {
		cartMarker(img, p[0], p[1], cell, pad, imgW, imgH, crashIdx)
	}
	for _, c := range carts {
		if !c.dead {
			cartMarker(img, c.x, c.y, cell, pad, imgW, imgH, cartIdx)
		}
	}
	return img
}

func tickCarts(carts []*cart, grid [][]byte, crashes map[[2]int]bool) bool {
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
	return crashedThisTick
}
