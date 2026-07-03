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

// Vis renders the Part One game as a bar chart of every player's final score, in
// player order. The winning player's bar is drawn bright and labeled; the rest
// are a dim uniform tint. Because only the winner is highlighted (by brightness,
// not hue) and the bars share one color, the chart reads in grayscale. The jagged
// profile reflects how the every-23rd scoring marble lands on players in a fixed
// cycle, so scores cluster by a player's position modulo the scoring period.
func (e Exercise) Vis(instr, outdir string) error {
	players, last, err := parse(instr)
	if err != nil {
		return err
	}

	scores := playGame(players, last)

	winner, best := 0, 0
	for i, s := range scores {
		if s > best {
			best, winner = s, i
		}
	}
	if best == 0 {
		return fmt.Errorf("no scoring players to visualize")
	}

	const (
		leftPad = 8
		topPad  = 24
		botPad  = 18
		plotH   = 260
		barGap  = 1
	)
	barW := 2
	plotW := players * (barW + barGap)
	if plotW < 600 {
		barW = 600/players + 1
		plotW = players * (barW + barGap)
	}
	W := leftPad*2 + plotW
	H := topPad + plotH + botPad

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	bg := color.RGBA{0x0c, 0x0f, 0x14, 0xff}
	dim := color.RGBA{0x3a, 0x52, 0x66, 0xff}  // ordinary player
	win := color.RGBA{0x86, 0xc8, 0xf0, 0xff}  // winning player (bright)
	axis := color.RGBA{0x55, 0x5c, 0x66, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}

	fill := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if x >= 0 && y >= 0 && x < W && y < H {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}

	fill(0, 0, W, H, bg)

	baseY := topPad + plotH
	for i, s := range scores {
		x0 := leftPad + i*(barW+barGap)
		h := s * plotH / best
		c := dim
		if i == winner {
			c = win
		}
		fill(x0, baseY-h, x0+barW, baseY, c)
	}

	// Baseline.
	fill(leftPad, baseY, leftPad+plotW, baseY+1, axis)

	drawText9(img, fmt.Sprintf("Part One scores by player (%d players)", players), leftPad, 16, white)

	// Point at the winner.
	wx := leftPad + winner*(barW+barGap)
	lbl := fmt.Sprintf("winner: player %d = %d", winner, best)
	lx := wx + 6
	if w := 7 * len(lbl); lx+w > W-4 {
		lx = W - 4 - w
	}
	drawText9(img, lbl, lx, topPad+14, win)

	f, err := os.Create(filepath.Join(outdir, "marble-mania.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawText9(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
