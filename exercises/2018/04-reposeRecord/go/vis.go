package exercises

import (
	"errors"
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
		return errors.New("no sleep records to visualize")
	}

	ids, rowOf, maxCount := guardOrder(asleep)
	s1Guard, s2Guard, s2Min := findSleepStrategies(ids, asleep)

	const (
		cell   = 12
		labelW = 78 // room for "#NNNN" row labels
		padT   = 40
		padB   = 28
		padR   = 20
	)
	cols := 60
	gridW := cols * cell
	imgW := labelW + gridW + padR
	imgH := padT + len(ids)*cell + padB

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for i := range imgW * imgH {
		img.Pix[i*4], img.Pix[i*4+1], img.Pix[i*4+2], img.Pix[i*4+3] = bg.R, bg.G, bg.B, bg.A
	}

	fg := color.RGBA{0xe8, 0xec, 0xf4, 0xff}
	vermilion := color.RGBA{0xD5, 0x5E, 0x00, 0xff}
	yellow := color.RGBA{0xF0, 0xE4, 0x42, 0xff}

	drawSleepChart(
		img, ids, asleep, rowOf, maxCount, s1Guard, s2Guard, s2Min,
		cell, labelW, gridW, padT, cols, imgW, imgH, fg, yellow, vermilion,
	)

	f, err := os.Create(filepath.Join(outdir, "repose-record.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func guardOrder(asleep map[int][60]int) ([]int, map[int]int, int) {
	ids := make([]int, 0, len(asleep))
	for id := range asleep {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	rowOf := make(map[int]int, len(ids))
	for r, id := range ids {
		rowOf[id] = r
	}
	maxCount := 1
	for _, mins := range asleep {
		for _, n := range mins {
			if n > maxCount {
				maxCount = n
			}
		}
	}
	return ids, rowOf, maxCount
}

func drawSleepChart(
	img *image.RGBA, ids []int, asleep map[int][60]int, rowOf map[int]int,
	maxCount, s1Guard, s2Guard, s2Min, cell, labelW, gridW, padT, cols, imgW, imgH int,
	fg, yellow, vermilion color.RGBA,
) {
	drawLabel4(img, "Repose Record: guard sleep heatmap (minutes 0-59)", labelW, 18, fg)
	for m := 0; m <= 59; m += 10 {
		drawLabel4(img, fmt.Sprintf("%02d", m), labelW+m*cell, padT-4, fg)
	}
	drawSleepHeatmap(img, ids, asleep, rowOf, maxCount, labelW, padT, cell, cols, imgW, imgH, fg)

	r1 := rowOf[s1Guard]
	y1 := padT + r1*cell
	strokeRect4(img, labelW-1, y1-1, labelW+gridW-1, y1+cell-1, yellow)
	drawLabel4(img, fmt.Sprintf("S1: #%d (total asleep)", s1Guard), labelW+2, y1-3, yellow)

	r2 := rowOf[s2Guard]
	x2 := labelW + s2Min*cell
	y2 := padT + r2*cell
	strokeRect4(img, x2-2, y2-2, x2+cell, y2+cell, vermilion)
	drawLabel4(img, fmt.Sprintf("S2: #%d @ min %02d", s2Guard, s2Min), labelW+gridW/2, 18, vermilion)

	ly := imgH - 10
	drawLabel4(img, "brighter = more days asleep at that minute", labelW, ly, fg)
	strokeRect4(img, labelW+300, ly-9, labelW+312, ly-1, yellow)
	drawLabel4(img, "S1 row", labelW+318, ly, fg)
	strokeRect4(img, labelW+390, ly-9, labelW+402, ly-1, vermilion)
	drawLabel4(img, "S2 cell", labelW+408, ly, fg)
}

func findSleepStrategies(ids []int, asleep map[int][60]int) (int, int, int) {
	s1Guard, s1Total := 0, -1
	s2Guard, s2Min, s2Count := 0, 0, -1
	for _, id := range ids {
		total := 0
		for m, n := range asleep[id] {
			total += n
			if n > s2Count {
				s2Guard, s2Min, s2Count = id, m, n
			}
		}
		if total > s1Total {
			s1Guard, s1Total = id, total
		}
	}
	return s1Guard, s2Guard, s2Min
}

func sleepShade4(n, maxCount int) color.RGBA {
	if n <= 0 {
		return color.RGBA{0x1a, 0x1f, 0x27, 0xff}
	}
	t := float64(n) / float64(maxCount)
	r := uint8(0x1c + t*(0xE9-0x1c))
	g := uint8(0x3a + t*(0xEC-0x3a))
	b := uint8(0x5e + t*(0xF4-0x5e))
	return color.RGBA{r, g, b, 0xff}
}

func fillRect4(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if x >= 0 && y >= 0 && x < imgW && y < imgH {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func strokeRect4(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	for x := x0; x <= x1; x++ {
		img.SetRGBA(x, y0, c)
		img.SetRGBA(x, y1, c)
	}
	for y := y0; y <= y1; y++ {
		img.SetRGBA(x0, y, c)
		img.SetRGBA(x1, y, c)
	}
}

func drawLabel4(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}

func drawSleepHeatmap(
	img *image.RGBA, ids []int, asleep map[int][60]int, rowOf map[int]int,
	maxCount, labelW, padT, cell, cols, imgW, imgH int, fg color.RGBA,
) {
	for _, id := range ids {
		r := rowOf[id]
		y0 := padT + r*cell
		drawLabel4(img, fmt.Sprintf("#%d", id), 6, y0+cell-2, fg)
		mins := asleep[id]
		for m := range cols {
			x0 := labelW + m*cell
			fillRect4(img, x0, y0, x0+cell-1, y0+cell-1, imgW, imgH, sleepShade4(mins[m], maxCount))
		}
	}
}
