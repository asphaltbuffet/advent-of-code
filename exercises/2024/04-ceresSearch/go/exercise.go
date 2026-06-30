package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Exercise for Advent of Code 2024 day 4.
type Exercise struct {
	common.BaseExercise
}

var xmasDirs = [8][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// xmasHit records one XMAS occurrence by its start cell (the X) and direction.
type xmasHit struct {
	r, c, dr, dc int
}

// findXMAS returns every XMAS occurrence in the grid (all 8 directions).
func findXMAS(grid []string) []xmasHit {
	const word = "XMAS"
	rows := len(grid)
	var hits []xmasHit
	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] != 'X' {
				continue
			}
			for _, d := range xmasDirs {
				ok := true
				for k := 1; k < len(word); k++ {
					nr, nc := r+d[0]*k, c+d[1]*k
					if nr < 0 || nr >= rows || nc < 0 || nc >= len(grid[nr]) || grid[nr][nc] != word[k] {
						ok = false
						break
					}
				}
				if ok {
					hits = append(hits, xmasHit{r, c, d[0], d[1]})
				}
			}
		}
	}
	return hits
}

// findXMAS2 returns the center (A) cell of every X-MAS shape in the grid.
func findXMAS2(grid []string) [][2]int {
	rows := len(grid)
	isMAS := func(a, b byte) bool {
		return (a == 'M' && b == 'S') || (a == 'S' && b == 'M')
	}
	var centers [][2]int
	for r := 1; r < rows-1; r++ {
		for c := 1; c < len(grid[r])-1; c++ {
			if grid[r][c] != 'A' {
				continue
			}
			if isMAS(grid[r-1][c-1], grid[r+1][c+1]) && isMAS(grid[r-1][c+1], grid[r+1][c-1]) {
				centers = append(centers, [2]int{r, c})
			}
		}
	}
	return centers
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return len(findXMAS(strings.Fields(instr))), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return len(findXMAS2(strings.Fields(instr))), nil
}

// Visualization tuning.
const (
	visCell    = 16 // pixel size of each grid cell
	visPadding = 24 // outer margin in pixels
)

var (
	visBG     = color.RGBA{R: 0x0a, G: 0x0a, B: 0x1e, A: 0xff} // dark navy
	visLetter = color.RGBA{R: 0xc8, G: 0xc8, B: 0xd8, A: 0xff} // bright grey letters
	visRed    = color.RGBA{R: 0xff, G: 0x3b, B: 0x3b, A: 0xff} // part 1: XMAS
	visGreen  = color.RGBA{R: 0x3b, G: 0xff, B: 0x6a, A: 0xff} // part 2: X-MAS
)

// visLineAlpha is the per-stamp opacity of overlay lines (0..1). Low so the
// dense overlays let the letters show through and overlaps glow brighter.
const visLineAlpha = 0.45

// Vis renders the word search with part-1 (red) and part-2 (green) overlays.
func (e Exercise) Vis(instr string, outdir string) error {
	grid := strings.Fields(instr)
	rows := len(grid)
	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	w := cols*visCell + 2*visPadding
	h := rows*visCell + 2*visPadding
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range img.Pix {
		img.Pix[i] = 0 // start transparent, then fill bg
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, visBG)
		}
	}

	// center of cell (r,c) in pixel space
	center := func(r, c int) (float64, float64) {
		x := float64(visPadding) + (float64(c)+0.5)*visCell
		y := float64(visPadding) + (float64(r)+0.5)*visCell
		return x, y
	}

	// Draw the letters.
	face := basicfont.Face7x13
	drawer := &font.Drawer{Dst: img, Src: image.NewUniform(visLetter), Face: face}
	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			ch := grid[r][c]
			cx, cy := center(r, c)
			// basicfont glyph is 7 wide, 13 tall; nudge to roughly center it.
			drawer.Dot = fixed.P(int(cx)-3, int(cy)+4)
			drawer.DrawString(string(ch))
		}
	}

	// Part 1: red line along each XMAS, from the X through the S.
	for _, hit := range findXMAS(grid) {
		x0, y0 := center(hit.r, hit.c)
		x1, y1 := center(hit.r+hit.dr*3, hit.c+hit.dc*3)
		drawLine(img, x0, y0, x1, y1, visRed, 1)
	}

	// Part 2: green cross over each X-MAS (both diagonals of the 3x3 block).
	for _, ctr := range findXMAS2(grid) {
		r, c := ctr[0], ctr[1]
		ax, ay := center(r-1, c-1)
		bx, by := center(r+1, c+1)
		drawLine(img, ax, ay, bx, by, visGreen, 1)
		ax, ay = center(r-1, c+1)
		bx, by = center(r+1, c-1)
		drawLine(img, ax, ay, bx, by, visGreen, 1)
	}

	f, err := os.Create(filepath.Join(outdir, "ceres-search.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// drawLine alpha-blends a thick line over the image. width is the half-extent
// of each stamp in pixels. Each pixel the line covers is blended exactly once
// (visLineAlpha), so a single line is uniform while crossings of separate
// lines accumulate and glow brighter.
func drawLine(img *image.RGBA, x0, y0, x1, y1 float64, col color.RGBA, width int) {
	dx := x1 - x0
	dy := y1 - y0
	steps := int(math.Max(math.Abs(dx), math.Abs(dy))) * 2
	if steps == 0 {
		steps = 1
	}
	covered := make(map[int]struct{})
	w := img.Bounds().Dx()
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		px := int(math.Round(x0 + dx*t))
		py := int(math.Round(y0 + dy*t))
		for oy := -width; oy <= width; oy++ {
			for ox := -width; ox <= width; ox++ {
				x, y := px+ox, py+oy
				key := y*w + x
				if _, seen := covered[key]; seen {
					continue
				}
				covered[key] = struct{}{}
				blendPixel(img, x, y, col, visLineAlpha)
			}
		}
	}
}

// blendPixel composites col over the existing pixel at (x,y) with opacity a.
func blendPixel(img *image.RGBA, x, y int, col color.RGBA, a float64) {
	if !(image.Pt(x, y).In(img.Bounds())) {
		return
	}
	dst := img.RGBAAt(x, y)
	blend := func(s, d uint8) uint8 {
		return uint8(float64(s)*a + float64(d)*(1-a))
	}
	img.SetRGBA(x, y, color.RGBA{
		R: blend(col.R, dst.R),
		G: blend(col.G, dst.G),
		B: blend(col.B, dst.B),
		A: 0xff,
	})
}
