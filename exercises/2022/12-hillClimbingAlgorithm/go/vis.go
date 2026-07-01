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

// Vis renders the height map as a shaded-relief terrain (dark low ground to
// bright ridges) and overlays the single-source distance field computed by the
// reverse BFS: reachable cells are tinted by how many steps they sit from the
// end, so the "flood" spreading back from the summit is visible. The start is
// marked red and the end cyan.
func (c Exercise) Vis(instr, outdir string) error {
	data := strings.Split(instr, "\n")
	dist, _, _ := bfsFromEnd(data)

	rows := len(data)
	cols := len(data[0])

	maxD := 0
	for r := range dist {
		for _, v := range dist[r] {
			if v > maxD {
				maxD = v
			}
		}
	}

	const scale = 10
	const pad = 6
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
			ch := rune(data[r][cc])
			switch ch {
			case start:
				put(cc, r, color.RGBA{0xff, 0x44, 0x55, 0xff})
				continue
			case end:
				put(cc, r, color.RGBA{0x3d, 0xe0, 0xe0, 0xff})
				continue
			}

			h := GetHeight(ch)          // 0..25
			base := 0.25 + 0.75*float64(h)/25 // relief brightness

			d := dist[r][cc]
			if d < 0 {
				// unreachable: gray relief only
				g := uint8(base * 120)
				put(cc, r, color.RGBA{g, g, uint8(base * 140), 0xff})
				continue
			}
			// Reachable: hue by distance, value by elevation relief.
			hue := math.Mod(float64(d)/float64(maxD+1)*260, 360)
			put(cc, r, hsvTerrain(hue, 0.6, base))
		}
	}

	f, err := os.Create(filepath.Join(outdir, "hill-climbing.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func hsvTerrain(h, s, v float64) color.RGBA {
	cc := v * s
	hp := math.Mod(h, 360) / 60
	x := cc * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = cc, x, 0
	case hp < 2:
		r, g, b = x, cc, 0
	case hp < 3:
		r, g, b = 0, cc, x
	case hp < 4:
		r, g, b = 0, x, cc
	case hp < 5:
		r, g, b = x, 0, cc
	default:
		r, g, b = cc, 0, x
	}
	m := v - cc
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}
