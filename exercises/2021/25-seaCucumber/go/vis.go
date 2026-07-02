package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the sea cucumber herds settling into gridlock (GIF). The east herd
// and south herd are drawn in two well-separated colorblind-safe colors that also
// differ strongly in brightness, so the two flows stay distinct in grayscale. The
// animation samples steps so the file stays small while still showing the herds
// sliding, congesting, and finally freezing.
func (e Exercise) Vis(instr, outdir string) error {
	g := parseGrid(instr)

	const (
		scale     = 3
		maxFrames = 80
	)

	// Okabe-Ito with a wide luminance gap so the herds separate in grayscale:
	// yellow (bright) for east, blue (dark) for south, over a near-black
	// seafloor. Palette index 0 = floor, 1 = east, 2 = south.
	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff},
		color.RGBA{0xF0, 0xE4, 0x42, 0xff}, // '>' east, bright
		color.RGBA{0x00, 0x72, 0xB2, 0xff}, // 'v' south, dark
	}

	W, H := g.w*scale, g.h*scale
	var frames []*image.Paletted
	var delays []int

	// Estimate total steps to pick a sampling stride, then re-run capturing.
	total := 1
	tmp := g
	tmp.cells = append([]byte(nil), g.cells...)
	for step(&tmp) {
		total++
	}
	stride := (total + maxFrames - 1) / maxFrames
	if stride < 1 {
		stride = 1
	}

	capture := func(cur grid) {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		for y := 0; y < cur.h; y++ {
			for x := 0; x < cur.w; x++ {
				var idx uint8
				switch cur.at(x, y) {
				case '>':
					idx = 1
				case 'v':
					idx = 2
				}
				if idx == 0 {
					continue
				}
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.SetColorIndex(x*scale+dx, y*scale+dy, idx)
					}
				}
			}
		}
		frames = append(frames, img)
		delays = append(delays, 8)
	}

	// Replay, capturing every `stride` steps plus the final frozen state.
	run := g
	run.cells = append([]byte(nil), g.cells...)
	i := 0
	capture(run)
	for step(&run) {
		i++
		if i%stride == 0 {
			capture(run)
		}
	}
	capture(run) // final frozen frame
	delays[len(delays)-1] = 300

	f, err := os.Create(filepath.Join(outdir, "sea-cucumber.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, &gif.GIF{Image: frames, Delay: delays})
}
