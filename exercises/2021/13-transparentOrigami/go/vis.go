package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the folding (GIF). The first frame is the full sheet of dots;
// each subsequent frame applies one more fold, so the paper visibly halves again
// and again until the dots line up into the eight-letter message (part two). Dots
// are drawn bright on a dark sheet; every frame is scaled to a fixed canvas so
// the collapse is easy to follow. Brightness alone encodes dot vs empty, so it
// reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	dots, folds, err := parse(instr)
	if err != nil {
		return err
	}

	// Snapshot the dot set before and after each fold.
	frames := []map[point]bool{cloneDots(dots)}
	cur := dots
	for _, f := range folds {
		cur = apply(cur, f)
		frames = append(frames, cloneDots(cur))
	}

	// Fixed output canvas; each frame's bounds are scaled to fit inside it.
	const canvas = 620
	const border = 10

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 background
		color.RGBA{0x1a, 0x20, 0x2c, 0xff}, // 1 sheet area
		color.RGBA{0x56, 0xB4, 0xE9, 0xff}, // 2 dot
		color.RGBA{0xF0, 0xE4, 0x42, 0xff}, // 3 final dots (message)
	}

	anim := &gif.GIF{}
	for fi, fr := range frames {
		maxX, maxY := 0, 0
		for p := range fr {
			if p.x > maxX {
				maxX = p.x
			}
			if p.y > maxY {
				maxY = p.y
			}
		}
		w := maxX + 1
		h := maxY + 1

		// Scale so the sheet fits the canvas while keeping square pixels.
		inner := canvas - 2*border
		scale := inner / max(w, 1)
		if s := inner / max(h, 1); s < scale {
			scale = s
		}
		if scale < 1 {
			scale = 1
		}

		img := image.NewPaletted(image.Rect(0, 0, canvas, canvas), pal)
		// Sheet backdrop.
		offX := (canvas - w*scale) / 2
		offY := (canvas - h*scale) / 2
		for y := 0; y < h*scale; y++ {
			for x := 0; x < w*scale; x++ {
				img.SetColorIndex(offX+x, offY+y, 1)
			}
		}

		dotIdx := uint8(2)
		if fi == len(frames)-1 {
			dotIdx = 3 // highlight the final message
		}
		for p := range fr {
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetColorIndex(offX+p.x*scale+dx, offY+p.y*scale+dy, dotIdx)
				}
			}
		}

		anim.Image = append(anim.Image, img)
		delay := 90
		if fi == 0 {
			delay = 150
		}
		if fi == len(frames)-1 {
			delay = 400
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "transparent-origami.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func cloneDots(d map[point]bool) map[point]bool {
	out := make(map[point]bool, len(d))
	for p := range d {
		out[p] = true
	}
	return out
}
