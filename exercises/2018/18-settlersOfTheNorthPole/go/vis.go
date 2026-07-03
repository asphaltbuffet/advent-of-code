package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Grayscale-safe palette, tiered by brightness so the three acre types read
// without relying on hue: open is darkest, trees mid, lumberyard brightest.
var visPalette = color.Palette{
	color.RGBA{20, 24, 28, 255},    // open — near-black
	color.RGBA{0, 158, 115, 255},   // trees — bluish green, mid brightness
	color.RGBA{240, 228, 66, 255},  // lumberyard — bright yellow
}

func visIndex(c byte) uint8 {
	switch c {
	case trees:
		return 1
	case lumberyard:
		return 2
	default:
		return 0
	}
}

// Vis animates the automaton: every warm-up minute plus two full passes of the
// cycle the state settles into, so the forest self-organizing and the eventual
// repeating loop are both visible.
func (e Exercise) Vis(instr, outdir string) error {
	const scale = 6

	grid := parse(instr)
	h, w := len(grid), len(grid[0])

	// Find the cycle: the minute each state first appears.
	states := [][][]byte{}
	seen := map[string]int{}
	cycleStart := -1
	for {
		key := string(bytesOf(grid))
		if first, ok := seen[key]; ok {
			cycleStart = first
			break
		}
		seen[key] = len(states)
		states = append(states, grid)
		grid = step(grid)
	}
	period := len(states) - cycleStart

	// Frames: the early minutes in full (the forest organizing fastest), then the
	// long warm-up sampled every few minutes to keep the file small, then two full
	// passes of the cycle so the repeating loop is visible.
	const fullEarly = 30 // show the first minutes frame-by-frame
	const sample = 4     // then keep every 4th warm-up minute
	frameIdx := make([]int, 0, fullEarly+cycleStart/sample+2*period)
	for i := 0; i < cycleStart; i++ {
		if i < fullEarly || i%sample == 0 {
			frameIdx = append(frameIdx, i)
		}
	}
	for pass := 0; pass < 2; pass++ {
		for i := 0; i < period; i++ {
			frameIdx = append(frameIdx, cycleStart+i)
		}
	}

	anim := &gif.GIF{}
	for _, idx := range frameIdx {
		st := states[idx]
		img := image.NewPaletted(image.Rect(0, 0, w*scale, h*scale), visPalette)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				ci := visIndex(st[y][x])
				for sy := 0; sy < scale; sy++ {
					for sx := 0; sx < scale; sx++ {
						img.SetColorIndex(x*scale+sx, y*scale+sy, ci)
					}
				}
			}
		}
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 8) // 80ms per frame
	}
	// Hold the last frame a beat before looping.
	if n := len(anim.Delay); n > 0 {
		anim.Delay[n-1] = 120
	}

	f, err := os.Create(filepath.Join(outdir, "settlers-of-the-north-pole.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}
