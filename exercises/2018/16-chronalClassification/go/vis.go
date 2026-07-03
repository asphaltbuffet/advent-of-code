package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders the opcode-deduction as a 16×16 matrix: rows are opcode numbers,
// columns are the sixteen operations. A cell is bright where that operation is an
// initial candidate for the number (consistent with every sample), dim where it is
// ruled out. The single resolved operation per number — the mapping the solver
// pins down by elimination — is boxed in a distinct accent. Candidacy is carried
// by brightness and the resolution by a boxed marker, so the matrix reads in
// grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	samples, _ := parse(instr)

	// Initial candidate matrix: cand[opcode][op] = consistent with all samples.
	cand := make([][]bool, 16)
	for i := range cand {
		cand[i] = make([]bool, 16)
		for j := range cand[i] {
			cand[i][j] = true
		}
	}
	for _, s := range samples {
		opcode := s.ins[0]
		valid := map[int]bool{}
		for _, op := range s.matches() {
			valid[op] = true
		}
		for op := 0; op < 16; op++ {
			if !valid[op] {
				cand[opcode][op] = false
			}
		}
	}

	// Resolve the one-to-one mapping by elimination (same as part two).
	resolved := make([]int, 16)
	for i := range resolved {
		resolved[i] = -1
	}
	work := make([][]bool, 16)
	for i := range work {
		work[i] = append([]bool(nil), cand[i]...)
	}
	for done := 0; done < 16; {
		for opcode := 0; opcode < 16; opcode++ {
			if resolved[opcode] != -1 {
				continue
			}
			n, only := 0, -1
			for op := 0; op < 16; op++ {
				if work[opcode][op] {
					n++
					only = op
				}
			}
			if n == 1 {
				resolved[opcode] = only
				done++
				for other := 0; other < 16; other++ {
					if other != opcode {
						work[other][only] = false
					}
				}
			}
		}
	}

	const (
		cell   = 26
		left   = 46 // room for row labels
		top    = 70 // room for rotated-ish column labels
		margin = 10
	)
	W := left + 16*cell + 2*margin
	H := top + 16*cell + 2*margin

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	bg := color.RGBA{0x0a, 0x0c, 0x12, 0xff}
	ruledOut := color.RGBA{0x1b, 0x20, 0x2a, 0xff} // eliminated candidate
	candidate := color.RGBA{0x6a, 0x8a, 0xb0, 0xff} // still-possible candidate
	resolvedC := color.RGBA{0xff, 0x8a, 0x00, 0xff} // the pinned mapping
	grid := color.RGBA{0x30, 0x36, 0x42, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	dim := color.RGBA{0x8a, 0x93, 0xa2, 0xff}

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

	cx := left + margin
	cy := top + margin

	for row := 0; row < 16; row++ {
		for col := 0; col < 16; col++ {
			x0 := cx + col*cell
			y0 := cy + row*cell
			c := ruledOut
			if cand[row][col] {
				c = candidate
			}
			fill(x0+1, y0+1, x0+cell-1, y0+cell-1, c)
			// box the resolved mapping cell
			if resolved[row] == col {
				for t := 0; t < 2; t++ {
					for x := x0; x < x0+cell; x++ {
						set(x, y0+t, resolvedC)
						set(x, y0+cell-1-t, resolvedC)
					}
					for y := y0; y < y0+cell; y++ {
						set(x0+t, y, resolvedC)
						set(x0+cell-1-t, y, resolvedC)
					}
				}
			}
		}
	}

	// Grid lines.
	for i := 0; i <= 16; i++ {
		fill(cx+i*cell, cy, cx+i*cell+1, cy+16*cell, grid)
		fill(cx, cy+i*cell, cx+16*cell, cy+i*cell+1, grid)
	}

	// Row labels (opcode numbers) and column labels (operation names, first 3
	// letters stacked would be busy — use a short two-letter abbreviation).
	names := []string{"ar", "ai", "mr", "mi", "br", "bi", "or", "oi", "sr", "si", "Gi", "Gt", "Gr", "Ei", "Et", "Er"}
	for i := 0; i < 16; i++ {
		drawText16(img, opLabel(i), 8, cy+i*cell+cell/2+4, dim)
		drawText16(img, names[i], cx+i*cell+cell/2-6, top-6, dim)
	}

	drawText16(img, "opcode-number  x  operation   (orange = resolved mapping)", left, 20, white)
	drawText16(img, "rows: opcode #   cols: op", 8, 40, dim)

	f, err := os.Create(filepath.Join(outdir, "chronal-classification.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// opLabel formats a small right-aligned opcode number.
func opLabel(n int) string {
	if n < 10 {
		return " " + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func drawText16(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
