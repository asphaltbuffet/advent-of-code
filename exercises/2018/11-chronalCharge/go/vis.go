package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders the 300×300 fuel-cell grid as a power heatmap: each cell's
// brightness rises with its power level (from -5 up to 4). The Part One 3×3 square
// and the Part Two best-of-any-size square are outlined and labeled in two
// colorblind-safe accents that also differ in brightness, so both boxes and the
// underlying texture stay legible in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	serial, err := parseSerial(instr)
	if err != nil {
		return err
	}

	sat := summedAreaTable(serial)

	// Part One winner (3×3) and Part Two winner (best of any size).
	p1x, p1y, _ := best(sat, 3)
	b2x, b2y, b2s, b2sum := 0, 0, 0, -1<<62
	for size := 1; size <= gridSize; size++ {
		x, y, sum := best(sat, size)
		if sum > b2sum {
			b2sum, b2x, b2y, b2s = sum, x, y, size
		}
	}

	const (
		scale  = 2 // pixels per cell
		margin = 16
		topPad = 24
	)
	plot := gridSize * scale
	W := plot + 2*margin
	H := plot + 2*margin + topPad

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	bg := color.RGBA{0x08, 0x0a, 0x0e, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	// Okabe-Ito accents: sky blue (Part One), vermillion (Part Two).
	p1c := color.RGBA{0x56, 0xB4, 0xE9, 0xff}
	p2c := color.RGBA{0xD5, 0x5E, 0x00, 0xff}

	set := func(x, y int, c color.RGBA) {
		if x >= 0 && y >= 0 && x < W && y < H {
			img.SetRGBA(x, y, c)
		}
	}
	fill := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				set(x, y, c)
			}
		}
	}

	fill(0, 0, W, H, bg)

	// Heatmap: power ranges -5..4, mapped to a 0..255 brightness ramp.
	for gy := 1; gy <= gridSize; gy++ {
		for gx := 1; gx <= gridSize; gx++ {
			p := cellPower(gx, gy, serial)
			v := uint8((p + 5) * 255 / 9)
			c := color.RGBA{v, v, v, 0xff}
			px := margin + (gx-1)*scale
			py := margin + topPad + (gy-1)*scale
			fill(px, py, px+scale, py+scale, c)
		}
	}

	// Outline a square with a dark halo so the accent stays visible over any
	// brightness in the heatmap, then the colored border on top. `thick` widens
	// the border so the two boxes are also distinguishable by line weight in
	// grayscale, not hue alone.
	halo := color.RGBA{0x00, 0x00, 0x00, 0xff}
	box := func(gx, gy, size int, c color.RGBA, thick int) {
		x0 := margin + (gx-1)*scale
		y0 := margin + topPad + (gy-1)*scale
		x1 := x0 + size*scale
		y1 := y0 + size*scale
		ring := func(pad int, col color.RGBA) {
			for x := x0 - pad; x <= x1+pad; x++ {
				set(x, y0-pad, col)
				set(x, y1+pad, col)
			}
			for y := y0 - pad; y <= y1+pad; y++ {
				set(x0-pad, y, col)
				set(x1+pad, y, col)
			}
		}
		ring(thick+1, halo) // dark outer halo
		for t := 0; t <= thick; t++ {
			ring(t, c) // colored border of the given weight
		}
	}

	box(b2x, b2y, b2s, p2c, 2) // Part Two: thicker border, drawn first
	box(p1x, p1y, 3, p1c, 1)   // Part One: thinner border

	drawText11(img, fmt.Sprintf("power heatmap  |  P1 3x3 @ %d,%d  |  P2 %dx%d @ %d,%d",
		p1x, p1y, b2s, b2s, b2x, b2y), margin, 16, white)

	f, err := os.Create(filepath.Join(outdir, "chronal-charge.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawText11(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
