package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates all 18 enhancement passes (GIF) on a canvas sized to the final
// grid. Each intermediate grid is centred and scaled to whole pixels; the first
// and last frames are held longer for readability.
func (e Exercise) Vis(instr, outdir string) error {
	const passes = 18

	rules := parseRules(instr)
	g := grid{".#.", "..#", "###"}

	frames := []grid{g}
	for range passes {
		g = enhance(g, rules)
		frames = append(frames, g)
	}
	canvas := len(frames[len(frames)-1]) // 2187: fits the final grid at 1px/cell

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 off
		color.RGBA{0x2f, 0xd0, 0x9a, 0xff}, // 1 on
	}

	anim := &gif.GIF{}
	for fi, fr := range frames {
		img := renderFractalFrame(fr, canvas, pal)
		anim.Image = append(anim.Image, img)
		delay := 40
		if fi == 0 || fi == len(frames)-1 {
			delay = 200
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "fractal-art.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func renderFractalFrame(fr grid, canvas int, pal color.Palette) *image.Paletted {
	n := len(fr)
	cell := max(canvas/n, 1)
	span := n * cell
	off := (canvas - span) / 2 // centre the scaled grid

	img := image.NewPaletted(image.Rect(0, 0, canvas, canvas), pal)
	for r := range n {
		row := fr[r]
		for c := range n {
			if row[c] != '#' {
				continue
			}
			x0, y0 := off+c*cell, off+r*cell
			for dy := range cell {
				for dx := range cell {
					img.SetColorIndex(x0+dx, y0+dy, 1)
				}
			}
		}
	}
	return img
}
