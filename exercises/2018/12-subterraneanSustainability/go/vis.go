package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders the automaton as a space-time diagram: time runs top to bottom, one
// row per generation, with planted pots drawn bright. Early rows churn; then the
// pattern locks into a fixed shape that drifts steadily sideways — the diagonal
// streak that makes the 50-billion-generation extrapolation possible. The
// lock-in generation is marked. Plant state is carried by brightness alone, so
// the diagram reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	pots, rules, err := parse(instr)
	if err != nil {
		return err
	}

	const gens = 220

	rows, lockGen, minX, maxX := simulatePots12(pots, rules, gens)
	width := maxX - minX + 1

	const (
		cell   = 3
		margin = 14
		topPad = 24
		lblW   = 44
	)
	imgW := lblW + width*cell + 2*margin
	imgH := topPad + (gens+1)*cell + 2*margin

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	bg := color.RGBA{0x0a, 0x0c, 0x12, 0xff}
	plant := color.RGBA{0xdf, 0xe8, 0xf6, 0xff}
	mark := color.RGBA{0xE6, 0x9F, 0x00, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	dim := color.RGBA{0x6a, 0x74, 0x86, 0xff}

	fillRect12(img, 0, 0, imgW, imgH, imgW, imgH, bg)
	drawPotRows12(img, rows, minX, cell, margin, topPad, lblW, imgW, imgH, plant, dim)

	if lockGen >= 0 {
		x0 := lblW + margin
		y := topPad + margin + lockGen*cell
		fillRect12(img, lblW+margin-6, y, lblW+margin-2, y+cell, imgW, imgH, mark)
		drawText12(img, fmt.Sprintf("locks in @ gen %d", lockGen), x0+4, y+cell+2, mark)
	}

	drawText12(img, "space-time diagram (time flows down)", 6, 16, white)

	f, err := os.Create(filepath.Join(outdir, "subterranean-sustainability.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func simulatePots12(pots state, rules map[string]bool, gens int) ([]state, int, int, int) {
	rows := make([]state, gens+1)
	rows[0] = pots
	seen := map[string]int{}
	lockGen := -1
	for g := range gens {
		if lockGen < 0 {
			sh, _ := shape(rows[g])
			if prev, ok := seen[sh]; ok {
				lockGen = prev
			} else {
				seen[sh] = g
			}
		}
		rows[g+1] = step(rows[g], rules)
	}
	minX, maxX := 1<<62, -(1 << 62)
	for _, r := range rows {
		lo, hi := bounds(r)
		if lo < minX {
			minX = lo
		}
		if hi > maxX {
			maxX = hi
		}
	}
	return rows, lockGen, minX, maxX
}

func fillRect12(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if x >= 0 && y >= 0 && x < imgW && y < imgH {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func drawPotRows12(
	img *image.RGBA, rows []state, minX, cell, margin, topPad, lblW, imgW, imgH int, plant, dim color.RGBA,
) {
	x0 := lblW + margin
	for g, r := range rows {
		y := topPad + margin + g*cell
		for i := range r {
			px := x0 + (i-minX)*cell
			fillRect12(img, px, y, px+cell, y+cell, imgW, imgH, plant)
		}
		if g%40 == 0 {
			drawText12(img, strconv.Itoa(g), 6, y+cell+4, dim)
		}
	}
}

func drawText12(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
