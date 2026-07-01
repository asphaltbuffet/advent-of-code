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

// Vis renders the board — the unfolded cube net the monkeys hand you — as a PNG.
// Each of the six faces (blockSize x blockSize regions) gets its own hue so the
// cross-shaped fold pattern is obvious; open tiles are bright, walls are darkened
// within their face's color, and the starting tile is marked white. This is the
// geometry part two treats as a cube.
func (e Exercise) Vis(instr, outdir string) error {
	bs := blockSize
	if strings.Count(instr, "\n") > 50 {
		bs = 50
	}
	b := parse(instr, bs)

	// Assign a color index to each occupied face region.
	faceHue := map[[2]int]float64{}
	nextFace := 0
	faceOf := func(r, c int) float64 {
		key := [2]int{r / bs, c / bs}
		h, ok := faceHue[key]
		if !ok {
			h = math.Mod(float64(nextFace)*60, 360)
			faceHue[key] = h
			nextFace++
		}
		return h
	}

	const scale = 4
	const pad = 6
	W := b.width*scale + 2*pad
	H := b.height*scale + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	put := func(r, c int, col color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(pad+c*scale+dx, pad+r*scale+dy, col)
			}
		}
	}

	for r, rowMap := range b.grid {
		for c, t := range rowMap {
			if t == empty {
				continue
			}
			hue := faceOf(r, c)
			val := 0.9
			sat := 0.55
			if t == wall {
				val = 0.35 // darken walls within the face color
				sat = 0.8
			}
			put(r, c, hsv22(hue, sat, val))
		}
	}

	// Mark the start tile.
	put(b.start[0], b.start[1], color.RGBA{0xff, 0xff, 0xff, 0xff})

	f, err := os.Create(filepath.Join(outdir, "monkey-map.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func hsv22(h, s, v float64) color.RGBA {
	c := v * s
	hp := math.Mod(h, 360) / 60
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
