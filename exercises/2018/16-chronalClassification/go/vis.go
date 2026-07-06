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
	cand := buildCandidateMatrix(samples)
	resolved := resolveOpcodeMapping(cand)

	const (
		cell   = 26
		left   = 46
		top    = 70
		margin = 10
	)
	imgW := left + 16*cell + 2*margin
	imgH := top + 16*cell + 2*margin

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	drawMatrix16(img, cand, resolved, cell, left, top, margin, imgW, imgH)

	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	dim := color.RGBA{0x8a, 0x93, 0xa2, 0xff}
	cx := left + margin
	cy := top + margin
	names := []string{"ar", "ai", "mr", "mi", "br", "bi", "or", "oi", "sr", "si", "Gi", "Gt", "Gr", "Ei", "Et", "Er"}
	for i := range 16 {
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

func buildCandidateMatrix(samples []sample) [][]bool {
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
		for op := range 16 {
			if !valid[op] {
				cand[opcode][op] = false
			}
		}
	}
	return cand
}

func singleCandidate(work [][]bool, opcode int) (int, bool) {
	n, only := 0, -1
	for op := range 16 {
		if work[opcode][op] {
			n++
			only = op
		}
	}
	return only, n == 1
}

func resolveOpcodeMapping(cand [][]bool) []int {
	resolved := make([]int, 16)
	for i := range resolved {
		resolved[i] = -1
	}
	work := make([][]bool, 16)
	for i := range work {
		work[i] = append([]bool(nil), cand[i]...)
	}
	for done := 0; done < 16; {
		for opcode := range 16 {
			if resolved[opcode] != -1 {
				continue
			}
			if only, ok := singleCandidate(work, opcode); ok {
				resolved[opcode] = only
				done++
				for other := range 16 {
					if other != opcode {
						work[other][only] = false
					}
				}
			}
		}
	}
	return resolved
}

func setPixel16(img *image.RGBA, x, y, imgW, imgH int, c color.RGBA) {
	if x >= 0 && y >= 0 && x < imgW && y < imgH {
		img.SetRGBA(x, y, c)
	}
}

func fillRect16(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			setPixel16(img, x, y, imgW, imgH, c)
		}
	}
}

func drawMatrix16(img *image.RGBA, cand [][]bool, resolved []int, cell, left, top, margin, imgW, imgH int) {
	bg := color.RGBA{0x0a, 0x0c, 0x12, 0xff}
	ruledOut := color.RGBA{0x1b, 0x20, 0x2a, 0xff}
	candidate := color.RGBA{0x6a, 0x8a, 0xb0, 0xff}
	resolvedC := color.RGBA{0xff, 0x8a, 0x00, 0xff}
	gridC := color.RGBA{0x30, 0x36, 0x42, 0xff}

	cx := left + margin
	cy := top + margin
	fillRect16(img, 0, 0, imgW, imgH, imgW, imgH, bg)
	for row := range 16 {
		for col := range 16 {
			x0 := cx + col*cell
			y0 := cy + row*cell
			c := ruledOut
			if cand[row][col] {
				c = candidate
			}
			fillRect16(img, x0+1, y0+1, x0+cell-1, y0+cell-1, imgW, imgH, c)
			if resolved[row] == col {
				drawResolvedBox16(img, x0, y0, cell, imgW, imgH, resolvedC)
			}
		}
	}
	for i := range 17 {
		fillRect16(img, cx+i*cell, cy, cx+i*cell+1, cy+16*cell, imgW, imgH, gridC)
		fillRect16(img, cx, cy+i*cell, cx+16*cell, cy+i*cell+1, imgW, imgH, gridC)
	}
}

func drawResolvedBox16(img *image.RGBA, x0, y0, cell, imgW, imgH int, c color.RGBA) {
	for t := range 2 {
		for x := x0; x < x0+cell; x++ {
			setPixel16(img, x, y0+t, imgW, imgH, c)
			setPixel16(img, x, y0+cell-1-t, imgW, imgH, c)
		}
		for y := y0; y < y0+cell; y++ {
			setPixel16(img, x0+t, y, imgW, imgH, c)
			setPixel16(img, x0+cell-1-t, y, imgW, imgH, c)
		}
	}
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
