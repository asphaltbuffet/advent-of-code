package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the image enhancement (GIF), one frame per step. The picture
// grows by a ring each step and gains detail; when the algorithm lights empty
// space (algo[0] = '#'), the infinite background flickers on for odd steps, shown
// as the frame's fill color. Lit pixels are bright on the dark (or flickered)
// background. Every frame is scaled to a fixed canvas sized to the final image,
// so the growth reads as the picture filling out from the center.
func (e Exercise) Vis(instr, outdir string) error {
	algo, im, err := parse(instr)
	if err != nil {
		return err
	}

	const steps = 50
	frames := []image2img{snapshot(im)}
	for i := 0; i < steps; i++ {
		im = enhance(algo, im)
		frames = append(frames, snapshot(im))
	}

	finalRows := len(frames[len(frames)-1].pix)
	finalCols := len(frames[len(frames)-1].pix[0])

	const canvas = 500
	scale := canvas / finalCols
	if s := canvas / finalRows; s < scale {
		scale = s
	}
	if scale < 1 {
		scale = 1
	}
	W := finalCols * scale
	H := finalRows * scale

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 dark background / unlit
		color.RGBA{0xF0, 0xE4, 0x42, 0xff}, // 1 lit pixel
		color.RGBA{0x2a, 0x24, 0x10, 0xff}, // 2 flickered (lit) background
	}

	anim := &gif.GIF{}
	for _, fr := range frames {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		// Background fill: index 2 when the infinite void is currently lit.
		bgIdx := uint8(0)
		if fr.bg {
			bgIdx = 2
		}
		for i := range img.Pix {
			img.Pix[i] = bgIdx
		}

		// Center this frame's (smaller) image in the fixed canvas.
		fr2rows := len(fr.pix)
		fr2cols := len(fr.pix[0])
		offR := (finalRows - fr2rows) / 2
		offC := (finalCols - fr2cols) / 2
		for r := 0; r < fr2rows; r++ {
			for c := 0; c < fr2cols; c++ {
				idx := bgIdx
				if fr.pix[r][c] {
					idx = 1
				}
				gr, gc := r+offR, c+offC
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.SetColorIndex((gc)*scale+dx, (gr)*scale+dy, idx)
					}
				}
			}
		}

		anim.Image = append(anim.Image, img)
		delay := 14
		if len(anim.Delay) == 0 {
			delay = 120
		}
		anim.Delay = append(anim.Delay, delay)
	}
	if len(anim.Delay) > 0 {
		anim.Delay[len(anim.Delay)-1] = 400
	}

	f, err := os.Create(filepath.Join(outdir, "trench-map.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

type image2img struct {
	pix [][]bool
	bg  bool
}

func snapshot(im enhImage) image2img {
	out := make([][]bool, len(im.pixels))
	for r := range im.pixels {
		out[r] = make([]bool, len(im.pixels[r]))
		copy(out[r], im.pixels[r])
	}
	return image2img{pix: out, bg: im.bg}
}
