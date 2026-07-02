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

// Vis renders the fabric as a coverage heatmap. Every claim marks the cells it
// covers; a cell's brightness encodes how many claims cover it, so uncontested
// single-claim fabric is dim and the contested overlaps — the Part One answer —
// glow brightly. The single intact claim, whose cells are all covered exactly
// once (the Part Two answer), is framed in vermilion and labeled. Coverage reads
// as brightness and the answer reads by outline and label, so the heatmap
// survives grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	claims, err := parse(instr)
	if err != nil {
		return err
	}
	if len(claims) == 0 {
		return fmt.Errorf("no claims to visualize")
	}

	// Fabric extent (with a small margin) and the max coverage seen.
	maxX, maxY := 0, 0
	for _, c := range claims {
		if r := c.left + c.width; r > maxX {
			maxX = r
		}
		if b := c.top + c.height; b > maxY {
			maxY = b
		}
	}
	const margin = 8
	W, H := maxX+margin, maxY+margin

	cover := coverage(claims)
	maxCov := 1
	for _, n := range cover {
		if n > maxCov {
			maxCov = n
		}
	}

	// Find the intact claim (Part Two) for the outline/label.
	var intact *claim
	for i := range claims {
		c := &claims[i]
		ok := true
		for x := c.left; x < c.left+c.width && ok; x++ {
			for y := c.top; y < c.top+c.height; y++ {
				if cover[[2]int{x, y}] > 1 {
					ok = false
					break
				}
			}
		}
		if ok {
			intact = c
			break
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for i := 0; i < W*H; i++ {
		img.Pix[i*4], img.Pix[i*4+1], img.Pix[i*4+2], img.Pix[i*4+3] = bg.R, bg.G, bg.B, bg.A
	}

	// Sequential dark->light ramp: single claim is dim blue, overlaps brighten
	// toward near-white. Brightness alone carries coverage, so grayscale works.
	shade := func(n int) color.RGBA {
		if n <= 0 {
			return bg
		}
		if n == 1 {
			return color.RGBA{0x1c, 0x3a, 0x5e, 0xff} // dim blue: uncontested
		}
		// map 2..maxCov onto a brightening ramp
		denom := float64(maxCov - 1)
		if denom <= 0 {
			denom = 1
		}
		t := float64(n-1) / denom
		if t > 1 {
			t = 1
		}
		// from mid blue (#2166AC-ish) to bright yellow-white for hottest overlap
		r := uint8(0x33 + t*(0xF0-0x33))
		g := uint8(0x66 + t*(0xE4-0x66))
		b := uint8(0xAC + t*(0x42-0xAC))
		return color.RGBA{r, g, b, 0xff}
	}

	for pos, n := range cover {
		x, y := pos[0], pos[1]
		if x < 0 || y < 0 || x >= W || y >= H {
			continue
		}
		img.SetRGBA(x, y, shade(n))
	}

	// Frame the intact claim in vermilion and label it.
	if intact != nil {
		vermilion := color.RGBA{0xD5, 0x5E, 0x00, 0xff}
		x0, y0 := intact.left, intact.top
		x1, y1 := intact.left+intact.width-1, intact.top+intact.height-1
		// draw a 2px frame slightly outside the claim so its own fill stays visible
		drawRect(img, x0-2, y0-2, x1+2, y1+2, vermilion)
		drawRect(img, x0-1, y0-1, x1+1, y1+1, vermilion)

		label := fmt.Sprintf("intact #%d", intact.id)
		lx := x0
		ly := y0 - 6
		if ly < 12 {
			ly = y1 + 16
		}
		drawLabel(img, label, lx, ly, vermilion)
	}

	f, err := os.Create(filepath.Join(outdir, "no-matter-how-you-slice-it.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// drawRect strokes a 1px rectangle border in col.
func drawRect(img *image.RGBA, x0, y0, x1, y1 int, col color.RGBA) {
	for x := x0; x <= x1; x++ {
		img.SetRGBA(x, y0, col)
		img.SetRGBA(x, y1, col)
	}
	for y := y0; y <= y1; y++ {
		img.SetRGBA(x0, y, col)
		img.SetRGBA(x1, y, col)
	}
}

// drawLabel writes text at (x, baseline y) in col using the built-in bitmap font.
func drawLabel(img *image.RGBA, s string, x, y int, col color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}
