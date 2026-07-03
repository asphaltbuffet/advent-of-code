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

// Vis renders the stars at the instant they align into the message. Each star is
// a bright block on a dark field, scaled up so the letters are crisp. The letters
// are carried by brightness alone against the background, so the image reads in
// grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	stars, err := parse(instr)
	if err != nil {
		return err
	}
	if len(stars) == 0 {
		return fmt.Errorf("no stars to visualize")
	}

	t := converge(stars)

	minX, minY := 1<<62, 1<<62
	maxX, maxY := -(1 << 62), -(1 << 62)
	for _, s := range stars {
		x, y := s.x+s.vx*t, s.y+s.vy*t
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	gw, gh := maxX-minX+1, maxY-minY+1
	const (
		scale  = 10 // pixels per star cell
		margin = 20
		topPad = 26
	)
	W := gw*scale + 2*margin
	H := gh*scale + 2*margin + topPad

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	bg := color.RGBA{0x0a, 0x0c, 0x12, 0xff}
	star := color.RGBA{0xdf, 0xe8, 0xf6, 0xff} // near-white light
	glow := color.RGBA{0x5a, 0x74, 0x9a, 0xff} // dim halo for depth
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}

	fill := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if x >= 0 && y >= 0 && x < W && y < H {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}

	fill(0, 0, W, H, bg)

	for _, s := range stars {
		gx := (s.x + s.vx*t - minX) * scale
		gy := (s.y+s.vy*t-minY)*scale + topPad
		// a soft halo then a solid core, both in cool tones that stay legible in gray
		fill(margin+gx-1, margin+gy-1, margin+gx+scale+1, margin+gy+scale+1, glow)
		fill(margin+gx, margin+gy, margin+gx+scale-1, margin+gy+scale-1, star)
	}

	drawText10(img, fmt.Sprintf("the stars align at t=%d", t), margin, 18, white)

	f, err := os.Create(filepath.Join(outdir, "the-stars-align.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawText10(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
