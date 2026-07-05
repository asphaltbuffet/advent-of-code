package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func hsv(h, s, v float64) color.RGBA {
	c := v * s
	hp := h / 60
	x := c * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = c, x, 0
	case hp < 2:
		r, g, b = x, c, 0
	case hp < 3:
		r, g, b = 0, c, x
	case hp < 4:
		r, g, b = 0, x, c
	case hp < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	m := v - c
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}

// Vis renders the 128x128 disk grid, colouring each connected region a distinct
// hue (golden-angle spread for visual separation). Free squares are a dark
// near-black; used squares show their region colour.
func (e Exercise) Vis(instr, outdir string) error {
	grid := buildGrid(strings.TrimSpace(instr))
	label, regions := labelRegions(grid)

	// Precompute a colour per region; a golden-angle hue spread keeps adjacent
	// region ids visually distinct.
	regionColor := make([]color.RGBA, regions+1)
	for id := 1; id <= regions; id++ {
		hue := math.Mod(float64(id)*137.508, 360)
		regionColor[id] = hsv(hue, 0.62, 0.95)
	}

	const cell = 6
	const pad = 8
	size := 128*cell + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	free := color.RGBA{0x18, 0x1c, 0x2a, 0xff}
	for y := range size {
		for x := range size {
			img.SetRGBA(x, y, bg)
		}
	}

	for r := range 128 {
		for c := range 128 {
			col := free
			if id := label[r][c]; id != 0 {
				col = regionColor[id]
			}
			x0, y0 := pad+c*cell, pad+r*cell
			for yy := y0; yy < y0+cell; yy++ {
				for xx := x0; xx < x0+cell; xx++ {
					img.SetRGBA(xx, yy, col)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "disk-defragmentation.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
