package exercises

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

const scale = 20

// Vis renders a PNG of the composited Space Image Format image.
// White pixels (1) are drawn as bright white; black pixels (0) as near-black.
// Each pixel cell is scale×scale display pixels, yielding a 500×120 image.
//
//nolint:gocognit // multi-layer image decoding and rendering is inherently branchy
func (e Exercise) Vis(instr, outdir string) error {
	data := bytes.TrimSpace([]byte(instr))

	// composite: start with all transparent
	img := make([]byte, pixels)
	for i := range img {
		img[i] = '2'
	}
	for i := 0; i+pixels <= len(data); i += pixels {
		layer := data[i : i+pixels]
		for j, b := range layer {
			if img[j] == '2' {
				img[j] = b
			}
		}
	}

	// render to PNG
	imgW := width * scale
	imgH := height * scale
	out := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	dark := color.RGBA{0x11, 0x11, 0x11, 0xFF}
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

	for row := range height {
		for col := range width {
			c := dark
			if img[row*width+col] == '1' {
				c = white
			}
			for dy := range scale {
				for dx := range scale {
					out.SetRGBA(col*scale+dx, row*scale+dy, c)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "vis.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, out)
}
