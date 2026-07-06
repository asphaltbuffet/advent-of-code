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
	color.RGBA{20, 24, 28, 255},   // open — near-black
	color.RGBA{0, 158, 115, 255},  // trees — bluish green, mid brightness
	color.RGBA{240, 228, 66, 255}, // lumberyard — bright yellow
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

func findCycle(grid [][]byte) ([][][]byte, int) {
	var states [][][]byte
	seen := map[string]int{}
	for {
		key := string(bytesOf(grid))
		if first, ok := seen[key]; ok {
			return states, first
		}
		seen[key] = len(states)
		states = append(states, grid)
		grid = step(grid)
	}
}

func renderForest(st [][]byte, w, h, scale int) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, w*scale, h*scale), visPalette)
	for y := range h {
		for x := range w {
			ci := visIndex(st[y][x])
			for sy := range scale {
				for sx := range scale {
					img.SetColorIndex(x*scale+sx, y*scale+sy, ci)
				}
			}
		}
	}
	return img
}

// Vis animates the automaton: every warm-up minute plus two full passes of the
// cycle the state settles into, so the forest self-organizing and the eventual
// repeating loop are both visible.
func (e Exercise) Vis(instr, outdir string) error {
	const scale = 6

	grid := parse(instr)
	h, w := len(grid), len(grid[0])

	states, cycleStart := findCycle(grid)
	period := len(states) - cycleStart

	const fullEarly = 30
	const sample = 4
	frameIdx := make([]int, 0, fullEarly+cycleStart/sample+2*period)
	for i := range cycleStart {
		if i < fullEarly || i%sample == 0 {
			frameIdx = append(frameIdx, i)
		}
	}
	for range 2 {
		for i := range period {
			frameIdx = append(frameIdx, cycleStart+i)
		}
	}

	anim := &gif.GIF{}
	for _, idx := range frameIdx {
		anim.Image = append(anim.Image, renderForest(states[idx], w, h, scale))
		anim.Delay = append(anim.Delay, 8)
	}
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
