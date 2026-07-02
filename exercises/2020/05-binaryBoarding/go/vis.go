package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Vis renders the plane as a seat grid where seat ID = row*8 + col, laid out with
// the 128 rows running left-to-right and the 8 columns top-to-bottom (a wide
// plane). Occupied seats are filled, empty seats are dark. The highest ID (Part
// One) is outlined in orange and your seat — the one empty seat with both
// neighbors occupied (Part Two) — is marked bright yellow. The front and back
// rows are entirely empty, which is why the search brackets to the occupied
// range. Highlights use a wide brightness gap and distinct marks, so they read in
// grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	ids := parse(instr)
	present := map[int]bool{}
	maxID, lo, hi := 0, 1<<10, 0
	for _, id := range ids {
		present[id] = true
		if id > maxID {
			maxID = id
		}
		if id < lo {
			lo = id
		}
		if id > hi {
			hi = id
		}
	}
	mine := -1
	for id := lo + 1; id < hi; id++ {
		if !present[id] && present[id-1] && present[id+1] {
			mine = id
			break
		}
	}

	const (
		rows  = 128
		cols  = 8
		cell  = 10
		gap   = 1
		mLeft = 6
		mTop  = 6
	)
	// Rows run horizontally (x), columns vertically (y): a wide plane.
	W := mLeft*2 + rows*(cell+gap)
	H := mTop*2 + cols*(cell+gap)

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	occupied := color.RGBA{0x00, 0x72, 0xB2, 0xff} // blue
	empty := color.RGBA{0x20, 0x26, 0x2e, 0xff}    // dark
	mineCol := color.RGBA{0xF0, 0xE4, 0x42, 0xff}  // bright yellow: your seat
	maxCol := color.RGBA{0xE6, 0x9F, 0x00, 0xff}   // orange: highest ID

	fill := func(px, py int, c color.RGBA) {
		for dy := 0; dy < cell; dy++ {
			for dx := 0; dx < cell; dx++ {
				img.SetRGBA(px+dx, py+dy, c)
			}
		}
	}
	frame := func(px, py int, c color.RGBA) {
		for d := 0; d < cell; d++ {
			img.SetRGBA(px+d, py, c)
			img.SetRGBA(px+d, py+cell-1, c)
			img.SetRGBA(px, py+d, c)
			img.SetRGBA(px+cell-1, py+d, c)
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			id := r*cols + c
			px := mLeft + r*(cell+gap)
			py := mTop + c*(cell+gap)
			switch {
			case id == mine:
				fill(px, py, mineCol)
			case present[id]:
				fill(px, py, occupied)
			default:
				fill(px, py, empty)
			}
			if id == maxID {
				frame(px, py, maxCol)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "binary-boarding.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
