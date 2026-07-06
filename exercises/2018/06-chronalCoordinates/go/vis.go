package exercises

import (
	"errors"
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

// Vis renders the Manhattan-distance Voronoi map. Each cell is tinted by the
// coordinate that owns it (contested ties stay dark). The largest finite
// territory — the Part One answer — is brightened, and the safe region where the
// summed distance to every coordinate stays under the threshold (the Part Two
// answer) is outlined as a bright contour. Coordinate seeds are marked. Owner
// tints vary in brightness and the two answers are carried by brightness and a
// contour, so the map reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	pts, err := parse(instr)
	if err != nil {
		return err
	}
	if len(pts) == 0 {
		return errors.New("no coordinates to visualize")
	}

	threshold := 10000
	if len(pts) <= 10 {
		threshold = 32
	}

	minX, minY, maxX, maxY := bounds(pts)
	const border = 12
	x0, y0 := minX-border, minY-border
	imgW := (maxX - minX) + 2*border + 1
	imgH := (maxY - minY) + 2*border + 1

	owner, safe, largestIdx, largest := buildVoronoi(pts, imgW, imgH, x0, y0, minX, minY, maxX, maxY, threshold)

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	drawVoronoiMap(img, owner, imgW, imgH, largestIdx)

	contour := color.RGBA{0xF0, 0xE4, 0x42, 0xff}
	drawSafeContour(img, safe, imgW, imgH, contour)

	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	for _, p := range pts {
		cx, cy := p.x-x0, p.y-y0
		for d := -1; d <= 1; d++ {
			img.SetRGBA(cx+d, cy, white)
			img.SetRGBA(cx, cy+d, white)
		}
	}

	// Label the largest finite territory near its seed, kept inside the image.
	if largestIdx >= 0 {
		p := pts[largestIdx]
		lbl := fmt.Sprintf("largest finite = %d", largest)
		lx := p.x - x0 + 4
		if w := 7 * len(lbl); lx+w > imgW-2 {
			lx = imgW - 2 - w
		}
		if lx < 2 {
			lx = 2
		}
		ly := p.y - y0 - 4
		if ly < 10 {
			ly = p.y - y0 + 14
		}
		drawText(img, lbl, lx, ly, white)
	}
	drawText(img, fmt.Sprintf("safe region (total<%d)", threshold), 4, imgH-6, contour)

	f, err := os.Create(filepath.Join(outdir, "chronal-coordinates.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func classifyCell(pts []point, wx, wy int) (int, int, bool) {
	best, bestIdx, total := 1<<30, -1, 0
	tie := false
	for i, p := range pts {
		d := abs(p.x-wx) + abs(p.y-wy)
		total += d
		switch {
		case d < best:
			best, bestIdx, tie = d, i, false
		case d == best:
			tie = true
		}
	}
	return bestIdx, total, tie
}

func buildVoronoi(pts []point, imgW, imgH, x0, y0, minX, minY, maxX, maxY, threshold int) ([]int, []bool, int, int) {
	owner := make([]int, imgW*imgH)
	safe := make([]bool, imgW*imgH)
	area := make([]int, len(pts))
	infinite := make([]bool, len(pts))

	for gy := range imgH {
		for gx := range imgW {
			wx, wy := x0+gx, y0+gy
			bestIdx, total, tie := classifyCell(pts, wx, wy)
			idx := gy*imgW + gx
			if tie {
				owner[idx] = -1
			} else {
				owner[idx] = bestIdx
				if wx >= minX && wx <= maxX && wy >= minY && wy <= maxY {
					area[bestIdx]++
					if wx == minX || wx == maxX || wy == minY || wy == maxY {
						infinite[bestIdx] = true
					}
				}
			}
			safe[idx] = total < threshold
		}
	}

	largestIdx, largest := -1, -1
	for i, a := range area {
		if !infinite[i] && a > largest {
			largest, largestIdx = a, i
		}
	}
	return owner, safe, largestIdx, largest
}

func drawVoronoiMap(img *image.RGBA, owner []int, imgW, imgH, largestIdx int) {
	fillDim := color.RGBA{0x22, 0x33, 0x47, 0xff}
	fillWin := color.RGBA{0x86, 0xb8, 0xe0, 0xff}
	boundary := color.RGBA{0x11, 0x16, 0x1d, 0xff}
	tieC := color.RGBA{0x08, 0x0a, 0x0e, 0xff}

	ownerAt := func(gx, gy int) int {
		if gx < 0 || gy < 0 || gx >= imgW || gy >= imgH {
			return -2
		}
		return owner[gy*imgW+gx]
	}
	for gy := range imgH {
		for gx := range imgW {
			o := owner[gy*imgW+gx]
			var c color.RGBA
			switch {
			case o < 0:
				c = tieC
			case o == largestIdx:
				c = fillWin
			default:
				c = fillDim
			}
			if o >= 0 && (ownerAt(gx+1, gy) != o || ownerAt(gx, gy+1) != o) {
				c = boundary
			}
			img.SetRGBA(gx, gy, c)
		}
	}
}

func drawSafeContour(img *image.RGBA, safe []bool, imgW, imgH int, contour color.RGBA) {
	safeAt := func(gx, gy int) bool {
		if gx < 0 || gy < 0 || gx >= imgW || gy >= imgH {
			return false
		}
		return safe[gy*imgW+gx]
	}
	for gy := range imgH {
		for gx := range imgW {
			if !safeAt(gx, gy) {
				continue
			}
			if !safeAt(gx-1, gy) || !safeAt(gx+1, gy) || !safeAt(gx, gy-1) || !safeAt(gx, gy+1) {
				img.SetRGBA(gx, gy, contour)
			}
		}
	}
}

func drawText(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
