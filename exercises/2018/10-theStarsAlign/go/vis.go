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
		return errors.New("no stars to visualize")
	}

	t := converge(stars)
	minX, minY, maxX, maxY := starBounds(stars, t)

	gw, gh := maxX-minX+1, maxY-minY+1
	const (
		scale  = 10
		margin = 20
		topPad = 26
	)
	imgW := gw*scale + 2*margin
	imgH := gh*scale + 2*margin + topPad

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	bg := color.RGBA{0x0a, 0x0c, 0x12, 0xff}
	starC := color.RGBA{0xdf, 0xe8, 0xf6, 0xff}
	glow := color.RGBA{0x5a, 0x74, 0x9a, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}

	fillRect10(img, 0, 0, imgW, imgH, imgW, imgH, bg)
	drawStars10(img, stars, t, minX, minY, scale, margin, topPad, imgW, imgH, glow, starC)
	drawText10(img, fmt.Sprintf("the stars align at t=%d", t), margin, 18, white)

	f, err := os.Create(filepath.Join(outdir, "the-stars-align.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func starBounds(stars []star, t int) (int, int, int, int) {
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
	return minX, minY, maxX, maxY
}

func fillRect10(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if x >= 0 && y >= 0 && x < imgW && y < imgH {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func drawStars10(
	img *image.RGBA, stars []star, t, minX, minY, scale, margin, topPad, imgW, imgH int, glow, starC color.RGBA,
) {
	for _, s := range stars {
		gx := (s.x + s.vx*t - minX) * scale
		gy := (s.y+s.vy*t-minY)*scale + topPad
		fillRect10(img, margin+gx-1, margin+gy-1, margin+gx+scale+1, margin+gy+scale+1, imgW, imgH, glow)
		fillRect10(img, margin+gx, margin+gy, margin+gx+scale-1, margin+gy+scale-1, imgW, imgH, starC)
	}
}

func drawText10(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
