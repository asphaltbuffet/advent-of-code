package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Vis drops the part-one sequence of rocks and renders the bottom of the
// resulting tower as a PNG. The seven-wide chamber makes a very tall, thin
// stack, so only the lowest rows are shown — enough to see the interlocking
// shapes and the trapped air pockets that make the height non-trivial to
// predict. Rock is shaded on a warm gradient by row; trapped empty cells inside
// the settled region are left dark.
func (c Exercise) Vis(instr, outdir string) error {
	s := parse(instr)
	for i := 0; i < partOneRocks; i++ {
		s.moveRock()
	}

	top := s.settledHeight
	const showRows = 240
	bottom := 0
	if top-showRows+1 > 0 {
		bottom = top - showRows + 1
	}
	rows := top - bottom + 1
	const cols = 7

	// Lay the tower on its side so it displays as a wide strip: the horizontal
	// axis is height (bottom of the tower at the left), the vertical axis is the
	// seven-wide chamber.
	const scale = 6
	const pad = 6
	W := rows*scale + 2*pad
	H := cols*scale + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	wall := color.RGBA{0x22, 0x26, 0x33, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// gh: height index (0 = bottom, left edge); gc: chamber column 0..6.
	put := func(gh, gc int, col color.RGBA) {
		sx := pad + gh*scale
		sy := pad + gc*scale
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(sx+dx, sy+dy, col)
			}
		}
	}

	for y := bottom; y <= top; y++ {
		for x := 0; x < cols; x++ {
			if s.chamber[y][x] == rock {
				f := float64(y) / float64(top+1) // 0 bottom .. 1 top
				col := color.RGBA{
					uint8(0xff),
					uint8(0x40 + 0xa0*f),
					uint8(0x20 + 0x40*f),
					0xff,
				}
				put(y-bottom, x, col)
			}
		}
	}

	// Chamber walls (top and bottom edges of the sideways strip).
	for gh := 0; gh < rows; gh++ {
		sx := pad + gh*scale
		for dx := 0; dx < scale; dx++ {
			img.SetRGBA(sx+dx, pad-1, wall)
			img.SetRGBA(sx+dx, H-pad, wall)
		}
	}

	f, err := os.Create(filepath.Join(outdir, "pyroclastic-flow.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
