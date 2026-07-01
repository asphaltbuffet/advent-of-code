package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the CRT output (part two) as a clean PNG. The little register
// program lights a 40x6 screen whose lit pixels spell out eight capital letters
// — the actual puzzle answer. Lit pixels are drawn on a warm glow, dark pixels
// as faint scanline cells, so the letters read clearly at a glance.
func (c Exercise) Vis(instr, outdir string) error {
	data := strings.Split(instr, "\n")
	d := Day10{
		Cycle:    1,
		X:        map[int]int{0: 0, 1: 1},
		Commands: []Command{},
		Cycle20:  1,
		Cycle60:  1,
		Cycle100: 1,
		Cycle140: 1,
		Cycle180: 1,
		Cycle220: 1,
	}

	if err := d.Parse(data); err != nil {
		return err
	}
	if err := d.Process(); err != nil {
		return err
	}

	rows := strings.Split(d.Draw(), "\n")
	gh := len(rows)
	gw := 0
	for _, r := range rows {
		if len([]rune(r)) > gw {
			gw = len([]rune(r))
		}
	}

	const scale = 16
	const pad = 12
	W := gw*scale + 2*pad
	H := gh*scale + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	dark := color.RGBA{0x18, 0x1d, 0x2c, 0xff}
	lit := color.RGBA{0xff, 0xc8, 0x4a, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	put := func(cx, cy int, col color.RGBA, inset int) {
		for dy := inset; dy < scale-inset; dy++ {
			for dx := inset; dx < scale-inset; dx++ {
				img.SetRGBA(pad+cx*scale+dx, pad+cy*scale+dy, col)
			}
		}
	}

	for ry, row := range rows {
		for rx, ch := range []rune(row) {
			if ch == '█' {
				put(rx, ry, lit, 1)
			} else {
				put(rx, ry, dark, 4)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "cathode-ray-tube.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
