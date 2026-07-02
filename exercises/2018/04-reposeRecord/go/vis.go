package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders the guard sleep tally as a heatmap: one row per guard, columns for
// minutes 0..59, each cell shaded by how many days that guard was asleep at that
// minute. Strategy 1 (Part One) picks the guard with the most total asleep
// minutes — its row is outlined. Strategy 2 (Part Two) picks the single brightest
// cell across all guards — it is boxed and labeled. Counts read as brightness and
// the two answers by outline and label, so the heatmap reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	asleep, err := sleepByGuard(instr)
	if err != nil {
		return err
	}
	if len(asleep) == 0 {
		return fmt.Errorf("no sleep records to visualize")
	}

	// Stable row order: guards sorted by ID.
	ids := make([]int, 0, len(asleep))
	for id := range asleep {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	rowOf := make(map[int]int, len(ids))
	for r, id := range ids {
		rowOf[id] = r
	}

	// Max cell count for the brightness scale.
	maxCount := 1
	for _, mins := range asleep {
		for _, n := range mins {
			if n > maxCount {
				maxCount = n
			}
		}
	}

	// Strategy 1: sleepiest guard by total minutes.
	s1Guard, s1Total := 0, -1
	for _, id := range ids {
		total := 0
		for _, n := range asleep[id] {
			total += n
		}
		if total > s1Total {
			s1Guard, s1Total = id, total
		}
	}
	// Strategy 2: single brightest (guard, minute) cell.
	s2Guard, s2Min, s2Count := 0, 0, -1
	for _, id := range ids {
		for m, n := range asleep[id] {
			if n > s2Count {
				s2Guard, s2Min, s2Count = id, m, n
			}
		}
	}

	const (
		cell   = 12
		labelW = 78 // room for "#NNNN" row labels
		padT   = 40
		padB   = 28
		padR   = 20
	)
	cols := 60
	gridW := cols * cell
	W := labelW + gridW + padR
	H := padT + len(ids)*cell + padB

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for i := 0; i < W*H; i++ {
		img.Pix[i*4], img.Pix[i*4+1], img.Pix[i*4+2], img.Pix[i*4+3] = bg.R, bg.G, bg.B, bg.A
	}

	fg := color.RGBA{0xe8, 0xec, 0xf4, 0xff}
	vermilion := color.RGBA{0xD5, 0x5E, 0x00, 0xff}
	yellow := color.RGBA{0xF0, 0xE4, 0x42, 0xff}

	// Sequential dark->light-blue ramp: 0 stays background, higher counts brighten.
	shade := func(n int) color.RGBA {
		if n <= 0 {
			return color.RGBA{0x1a, 0x1f, 0x27, 0xff} // faint empty-cell tint
		}
		t := float64(n) / float64(maxCount)
		r := uint8(0x1c + t*(0xE9-0x1c))
		g := uint8(0x3a + t*(0xEC-0x3a))
		b := uint8(0x5e + t*(0xF4-0x5e))
		return color.RGBA{r, g, b, 0xff}
	}

	fillRect := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if x >= 0 && y >= 0 && x < W && y < H {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}
	strokeRect := func(x0, y0, x1, y1 int, c color.RGBA) {
		for x := x0; x <= x1; x++ {
			img.SetRGBA(x, y0, c)
			img.SetRGBA(x, y1, c)
		}
		for y := y0; y <= y1; y++ {
			img.SetRGBA(x0, y, c)
			img.SetRGBA(x1, y, c)
		}
	}
	label := func(s string, x, y int, c color.RGBA) {
		d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
		d.DrawString(s)
	}

	// Title and minute-axis ticks (every 10 minutes).
	label("Repose Record: guard sleep heatmap (minutes 0-59)", labelW, 18, fg)
	for m := 0; m <= 59; m += 10 {
		x := labelW + m*cell
		label(fmt.Sprintf("%02d", m), x, padT-4, fg)
	}

	// Cells + row labels.
	for _, id := range ids {
		r := rowOf[id]
		y0 := padT + r*cell
		label(fmt.Sprintf("#%d", id), 6, y0+cell-2, fg)
		mins := asleep[id]
		for m := 0; m < cols; m++ {
			x0 := labelW + m*cell
			fillRect(x0, y0, x0+cell-1, y0+cell-1, shade(mins[m]))
		}
	}

	// Strategy 1: outline the sleepiest guard's whole row in yellow, with the
	// label above the row so it never sits on top of data cells.
	{
		r := rowOf[s1Guard]
		y0 := padT + r*cell
		strokeRect(labelW-1, y0-1, labelW+gridW-1, y0+cell-1, yellow)
		label(fmt.Sprintf("S1: #%d (total asleep)", s1Guard), labelW+2, y0-3, yellow)
	}
	// Strategy 2: box the single brightest cell in vermilion. Draw a short
	// connector up to the title band and place the label there, clear of cells.
	{
		r := rowOf[s2Guard]
		x0 := labelW + s2Min*cell
		y0 := padT + r*cell
		strokeRect(x0-2, y0-2, x0+cell, y0+cell, vermilion)
		label(fmt.Sprintf("S2: #%d @ min %02d", s2Guard, s2Min), labelW+gridW/2, 18, vermilion)
	}

	// Legend.
	ly := H - 10
	label("brighter = more days asleep at that minute", labelW, ly, fg)
	strokeRect(labelW+300, ly-9, labelW+312, ly-1, yellow)
	label("S1 row", labelW+318, ly, fg)
	strokeRect(labelW+390, ly-9, labelW+402, ly-1, vermilion)
	label("S2 cell", labelW+408, ly, fg)

	f, err := os.Create(filepath.Join(outdir, "repose-record.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
