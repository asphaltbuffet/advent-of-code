package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis scatter-plots the spoken value against the turn number for the opening
// stretch of the game. The sequence has a distinctive texture: frequent 0s along
// the bottom (every time a brand-new number is spoken) punctuated by rising
// diagonal streaks (a value's gap grows the longer ago it was last said). The
// 2020th turn — the Part One answer — is marked. Points brighten with value and
// the answer marker is the brightest, so the structure reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	start, err := parse(instr)
	if err != nil {
		return err
	}

	const turns = 2020
	seq := make([]int, turns)
	copy(seq, start)
	lastSeen := map[int]int{}
	for i := 0; i < len(start)-1; i++ {
		lastSeen[start[i]] = i + 1
	}
	last := start[len(start)-1]
	maxVal := 0
	for _, v := range start {
		if v > maxVal {
			maxVal = v
		}
	}
	for turn := len(start); turn < turns; turn++ {
		var next int
		if prev, ok := lastSeen[last]; ok {
			next = turn - prev
		}
		lastSeen[last] = turn
		seq[turn] = next
		last = next
		if next > maxVal {
			maxVal = next
		}
	}
	answer := seq[turns-1]

	const (
		W   = 960
		H   = 520
		mL  = 60
		mR  = 30
		mT  = 50
		mB  = 50
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	xOf := func(t int) int { return mL + t*plotW/turns }
	yOf := func(v int) int { return mT + plotH - v*plotH/maxVal }

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// Value-brightness ramp from dim blue (low) to bright yellow (high).
	plot := func(t, v int, big bool, override *color.RGBA) {
		x, y := xOf(t), yOf(v)
		var c color.RGBA
		if override != nil {
			c = *override
		} else {
			frac := float64(v) / float64(maxVal)
			c = color.RGBA{
				uint8(0x30 + frac*0xC0),
				uint8(0x60 + frac*0x84),
				uint8(0xB2 - frac*0x70),
				0xff,
			}
		}
		r := 1
		if big {
			r = 5
		}
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				px, py := x+dx, y+dy
				if px >= 0 && px < W && py >= 0 && py < H {
					img.SetRGBA(px, py, c)
				}
			}
		}
	}

	for t := 0; t < turns; t++ {
		plot(t, seq[t], false, nil)
	}
	answerCol := color.RGBA{0xF0, 0xE4, 0x42, 0xff}
	plot(turns-1, answer, true, &answerCol)

	label := func(x, y int, s string, c color.RGBA) {
		(&font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}).DrawString(s)
	}
	white := color.RGBA{0xe8, 0xec, 0xf4, 0xff}
	label(mL, 30, "Rambunctious Recitation: spoken value per turn (first 2020)", white)
	label(mL, H-16, "turn number  →", color.RGBA{0x9a, 0xa4, 0xb2, 0xff})
	// Answer callout near the marker.
	ax, ay := xOf(turns-1), yOf(answer)
	label(ax-160, ay-6, fmt.Sprintf("2020th = %d", answer), answerCol)

	f, err := os.Create(filepath.Join(outdir, "rambunctious-recitation.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
