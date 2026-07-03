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
		return fmt.Errorf("no coordinates to visualize")
	}

	threshold := 10000
	if len(pts) <= 10 {
		threshold = 32
	}

	minX, minY, maxX, maxY := bounds(pts)
	const border = 12
	x0, y0 := minX-border, minY-border
	W := (maxX - minX) + 2*border + 1
	H := (maxY - minY) + 2*border + 1

	// Owner index per cell (-1 = tie/unowned), plus totals for the safe region.
	owner := make([]int, W*H)
	safe := make([]bool, W*H)
	area := make([]int, len(pts))
	infinite := make([]bool, len(pts))

	for gy := 0; gy < H; gy++ {
		for gx := 0; gx < W; gx++ {
			wx, wy := x0+gx, y0+gy
			best, bestIdx, tie, total := 1<<30, -1, false, 0
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
			idx := gy*W + gx
			if tie {
				owner[idx] = -1
			} else {
				owner[idx] = bestIdx
				// only count area within the true bounding box (matches Part One)
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

	// Largest finite territory (Part One).
	largestIdx, largest := -1, -1
	for i, a := range area {
		if !infinite[i] && a > largest {
			largest, largestIdx = a, i
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	// The per-owner partition is drawn as a muted backdrop: a single dim blue
	// fill with thin dark boundaries between territories, so the Voronoi shape is
	// legible without a rainbow that would collapse in grayscale. The two answers
	// carry the brightness: the largest finite territory is filled bright, the
	// tie cells stay darkest, and the safe-region contour (below) is brightest.
	fillDim := color.RGBA{0x22, 0x33, 0x47, 0xff}      // ordinary territory
	fillWin := color.RGBA{0x86, 0xb8, 0xe0, 0xff}      // largest finite (Part One)
	boundary := color.RGBA{0x11, 0x16, 0x1d, 0xff}     // between territories
	tie := color.RGBA{0x08, 0x0a, 0x0e, 0xff}          // contested (darkest)

	ownerAt := func(gx, gy int) int {
		if gx < 0 || gy < 0 || gx >= W || gy >= H {
			return -2
		}
		return owner[gy*W+gx]
	}
	for gy := 0; gy < H; gy++ {
		for gx := 0; gx < W; gx++ {
			o := owner[gy*W+gx]
			var c color.RGBA
			switch {
			case o < 0:
				c = tie
			case o == largestIdx:
				c = fillWin
			default:
				c = fillDim
			}
			// darken cells on an owner boundary to trace the Voronoi edges
			if o >= 0 && (ownerAt(gx+1, gy) != o || ownerAt(gx, gy+1) != o) {
				c = boundary
			}
			img.SetRGBA(gx, gy, c)
		}
	}

	// Outline the safe region (Part Two): draw a bright contour where a safe cell
	// borders a non-safe one.
	contour := color.RGBA{0xF0, 0xE4, 0x42, 0xff} // yellow
	at := func(gx, gy int) bool {
		if gx < 0 || gy < 0 || gx >= W || gy >= H {
			return false
		}
		return safe[gy*W+gx]
	}
	for gy := 0; gy < H; gy++ {
		for gx := 0; gx < W; gx++ {
			if !at(gx, gy) {
				continue
			}
			if !at(gx-1, gy) || !at(gx+1, gy) || !at(gx, gy-1) || !at(gx, gy+1) {
				img.SetRGBA(gx, gy, contour)
			}
		}
	}

	// Mark coordinate seeds with a small white cross.
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
		if w := 7 * len(lbl); lx+w > W-2 {
			lx = W - 2 - w
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
	drawText(img, fmt.Sprintf("safe region (total<%d)", threshold), 4, H-6, contour)

	f, err := os.Create(filepath.Join(outdir, "chronal-coordinates.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawText(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
