package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Vis renders the tree map with the Part One toboggan run (right 3, down 1)
// traced through it. The map repeats to the right, so the path's column wraps; it
// is drawn on the natural 31-wide grid. Cells are: open ground (dark), a tree the
// run misses (muted), a tree the run hits (bright vermilion), and open ground the
// run passes over (a light ring). The count of hit trees is the answer. Colors
// are chosen with a wide brightness gap and hits are the brightest cells, so the
// run reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	grid := parse(instr)
	h := len(grid)
	w := len(grid[0])

	// Mark which cell each step of the run visits.
	visited := make(map[[2]int]bool, h)
	x := 0
	for y := 0; y < h; y++ {
		visited[[2]int{x % w, y}] = true
		x += 3
	}

	const scale = 6
	W, H := w*scale, h*scale
	img := image.NewRGBA(image.Rect(0, 0, W, H))

	ground := color.RGBA{0x14, 0x18, 0x1e, 0xff}
	tree := color.RGBA{0x2f, 0x6d, 0x50, 0xff}     // muted green tree
	hit := color.RGBA{0xF0, 0xE4, 0x42, 0xff}      // bright yellow: tree struck
	hitCore := color.RGBA{0xD5, 0x5E, 0x00, 0xff}  // vermilion core inside a hit
	pathRing := color.RGBA{0x8a, 0x94, 0xa2, 0xff} // gray ring: open cell on the path

	fill := func(cx, cy int, c color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(cx*scale+dx, cy*scale+dy, c)
			}
		}
	}
	inset := scale / 4
	if inset < 1 {
		inset = 1
	}
	ring := func(cx, cy int, c color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				if dx < inset || dx >= scale-inset || dy < inset || dy >= scale-inset {
					img.SetRGBA(cx*scale+dx, cy*scale+dy, c)
				}
			}
		}
	}

	for y := 0; y < h; y++ {
		for cx := 0; cx < w; cx++ {
			isTree := grid[y][cx] == '#'
			onPath := visited[[2]int{cx, y}]
			switch {
			case onPath && isTree:
				fill(cx, y, hit)
				ring(cx, y, hitCore) // vermilion frame inside the bright cell
			case onPath:
				fill(cx, y, ground)
				ring(cx, y, pathRing)
			case isTree:
				fill(cx, y, tree)
			default:
				fill(cx, y, ground)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "toboggan-trajectory.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
