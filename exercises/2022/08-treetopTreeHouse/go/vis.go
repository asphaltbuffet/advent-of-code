package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the forest as a PNG. Every tree is shaded by its height; trees
// visible from outside the grid (part one) glow gold while hidden interior trees
// stay dark green, so the sparse visible perimeter-and-ridges pattern stands out.
// The single tree with the best scenic score (part two) is ringed in red.
func (c Exercise) Vis(instr, outdir string) error {
	data := strings.Split(instr, "\n")
	dimX = len(data[0])
	dimY = len(data)
	m := GetTreeMap(data)
	rows := len(m)
	cols := len(m[0])

	visible := make([][]bool, rows)
	for r := range visible {
		visible[r] = make([]bool, cols)
	}

	// Sweep each of the four directions, marking a tree visible when it is taller
	// than everything seen so far along that line.
	for r := 0; r < rows; r++ {
		hi := -1
		for cc := 0; cc < cols; cc++ { // from west
			if m[r][cc] > hi {
				visible[r][cc] = true
				hi = m[r][cc]
			}
		}
		hi = -1
		for cc := cols - 1; cc >= 0; cc-- { // from east
			if m[r][cc] > hi {
				visible[r][cc] = true
				hi = m[r][cc]
			}
		}
	}
	for cc := 0; cc < cols; cc++ {
		hi := -1
		for r := 0; r < rows; r++ { // from north
			if m[r][cc] > hi {
				visible[r][cc] = true
				hi = m[r][cc]
			}
		}
		hi = -1
		for r := rows - 1; r >= 0; r-- { // from south
			if m[r][cc] > hi {
				visible[r][cc] = true
				hi = m[r][cc]
			}
		}
	}

	// Best scenic score (interior only, matching part two).
	bestR, bestC, best := 0, 0, -1
	for r := 1; r < rows-1; r++ {
		for cc := 1; cc < cols-1; cc++ {
			h := m[r][cc]
			score := CalculateScoreUp(h, r, cc, m) * CalculateScoreDown(h, r, cc, m) *
				CalculateScoreLeft(h, r, cc, m) * CalculateScoreRight(h, r, cc, m)
			if score > best {
				best, bestR, bestC = score, r, cc
			}
		}
	}

	const scale = 6
	const pad = 4
	W := cols*scale + 2*pad
	H := rows*scale + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x08, 0x0a, 0x10, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	put := func(cx, cy int, col color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(pad+cx*scale+dx, pad+cy*scale+dy, col)
			}
		}
	}

	for r := 0; r < rows; r++ {
		for cc := 0; cc < cols; cc++ {
			h := float64(m[r][cc]) / 9 // 0..1
			var col color.RGBA
			if visible[r][cc] {
				// Gold, brighter for taller visible trees.
				col = color.RGBA{uint8(0xc0 + 0x3f*h), uint8(0x90 + 0x50*h), 0x20, 0xff}
			} else {
				// Hidden: dark green shaded by height.
				col = color.RGBA{0x10, uint8(0x30 + 0x60*h), 0x18, 0xff}
			}
			put(cc, r, col)
		}
	}

	// Ring the best scenic tree in red.
	ring := color.RGBA{0xff, 0x30, 0x40, 0xff}
	for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}} {
		put(bestC+d[1], bestR+d[0], ring)
	}
	put(bestC, bestR, color.RGBA{0xff, 0xff, 0xff, 0xff})

	f, err := os.Create(filepath.Join(outdir, "treetop-treehouse.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
