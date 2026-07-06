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

	p1x, p1y, _ := best(sat, 3)
	b2x, b2y, b2s := findBestAnySize(sat)

	const (
		scale  = 2
		margin = 16
		topPad = 24
	)
	plot := gridSize * scale
	imgW := plot + 2*margin
	imgH := plot + 2*margin + topPad

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	bg := color.RGBA{0x08, 0x0a, 0x0e, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	p1c := color.RGBA{0x56, 0xB4, 0xE9, 0xff}
	p2c := color.RGBA{0xD5, 0x5E, 0x00, 0xff}

	fillRect11(img, 0, 0, imgW, imgH, imgW, imgH, bg)
	drawHeatmap11(img, serial, scale, margin, topPad, imgW, imgH)
	drawBox11(img, b2x, b2y, b2s, scale, margin, topPad, imgW, imgH, p2c, 2)
	drawBox11(img, p1x, p1y, 3, scale, margin, topPad, imgW, imgH, p1c, 1)
	drawText11(img, fmt.Sprintf("power heatmap  |  P1 3x3 @ %d,%d  |  P2 %dx%d @ %d,%d",
		p1x, p1y, b2s, b2s, b2x, b2y), margin, 16, white)

	f, err := os.Create(filepath.Join(outdir, "chronal-charge.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func findBestAnySize(sat [][]int) (int, int, int) {
	bx, by, bs, bsum := 0, 0, 0, -1<<62
	for size := 1; size <= gridSize; size++ {
		x, y, sum := best(sat, size)
		if sum > bsum {
			bsum, bx, by, bs = sum, x, y, size
		}
	}
	return bx, by, bs
}

func setPixel11(img *image.RGBA, x, y, imgW, imgH int, c color.RGBA) {
	if x >= 0 && y >= 0 && x < imgW && y < imgH {
		img.SetRGBA(x, y, c)
	}
}

func fillRect11(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			setPixel11(img, x, y, imgW, imgH, c)
		}
	}
}

func drawHeatmap11(img *image.RGBA, serial, scale, margin, topPad, imgW, imgH int) {
	for gy := 1; gy <= gridSize; gy++ {
		for gx := 1; gx <= gridSize; gx++ {
			p := cellPower(gx, gy, serial)
			v := uint8((p + 5) * 255 / 9)
			c := color.RGBA{v, v, v, 0xff}
			px := margin + (gx-1)*scale
			py := margin + topPad + (gy-1)*scale
			fillRect11(img, px, py, px+scale, py+scale, imgW, imgH, c)
		}
	}
}

func drawRing11(img *image.RGBA, x0, y0, x1, y1, pad, imgW, imgH int, col color.RGBA) {
	for x := x0 - pad; x <= x1+pad; x++ {
		setPixel11(img, x, y0-pad, imgW, imgH, col)
		setPixel11(img, x, y1+pad, imgW, imgH, col)
	}
	for y := y0 - pad; y <= y1+pad; y++ {
		setPixel11(img, x0-pad, y, imgW, imgH, col)
		setPixel11(img, x1+pad, y, imgW, imgH, col)
	}
}

func drawBox11(img *image.RGBA, gx, gy, size, scale, margin, topPad, imgW, imgH int, c color.RGBA, thick int) {
	x0 := margin + (gx-1)*scale
	y0 := margin + topPad + (gy-1)*scale
	x1 := x0 + size*scale
	y1 := y0 + size*scale
	halo := color.RGBA{0x00, 0x00, 0x00, 0xff}
	drawRing11(img, x0, y0, x1, y1, thick+1, imgW, imgH, halo)
	for t := 0; t <= thick; t++ {
		drawRing11(img, x0, y0, x1, y1, t, imgW, imgH, c)
	}
}

func drawText11(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
