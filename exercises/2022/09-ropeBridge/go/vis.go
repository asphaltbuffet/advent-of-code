package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Vis replays the 10-knot rope (part two) and renders the path swept by the
// tail as a PNG. Each cell the tail visits is colored by the order it was first
// reached — a hue gradient from the start (red) through the spectrum — so the
// looping, self-crossing route the rope drags out across the grid is visible.
// The start position is marked white.
func (c Exercise) Vis(instr, outdir string) error {
	data := strings.Split(instr, "\n")

	order := map[Point]int{}
	knots := make([]Point, 10)
	step := 0
	order[knots[9]] = step

	for _, line := range data {
		direction, right, _ := strings.Cut(line, " ")
		distance, err := strconv.Atoi(right)
		if err != nil {
			return fmt.Errorf("invalid distance %q: %w", right, err)
		}

		mv := getMovement(direction)
		for i := 0; i < distance; i++ {
			knots[0].X += mv.X
			knots[0].Y += mv.Y
			for j := 1; j < 10; j++ {
				t := CalculateMovement(knots[j-1], knots[j])
				knots[j].X += t.X
				knots[j].Y += t.Y
			}
			step++
			if _, seen := order[knots[9]]; !seen {
				order[knots[9]] = step
			}
		}
	}

	// Bounds over every visited tail cell.
	loX, hiX, loY, hiY := 0, 0, 0, 0
	for p := range order {
		loX = min(loX, p.X)
		hiX = max(hiX, p.X)
		loY = min(loY, p.Y)
		hiY = max(hiY, p.Y)
	}

	const scale = 3
	const pad = 6
	gw := hiX - loX + 1
	gh := hiY - loY + 1
	W := gw*scale + 2*pad
	H := gh*scale + 2*pad

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// y grows downward on screen but upward in grid coords; flip for a natural view.
	put := func(gx, gy int, col color.RGBA) {
		sx := pad + (gx-loX)*scale
		sy := pad + (hiY-gy)*scale
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(sx+dx, sy+dy, col)
			}
		}
	}

	for p, s := range order {
		f := float64(s) / float64(step)
		put(p.X, p.Y, hsv(f*300, 0.85, 0.95))
	}

	// Mark the origin.
	put(0, 0, color.RGBA{0xff, 0xff, 0xff, 0xff})

	f, err := os.Create(filepath.Join(outdir, "rope-bridge.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// hsv converts HSV (h in degrees, s,v in 0..1) to RGBA.
func hsv(h, s, v float64) color.RGBA {
	cc := v * s
	hp := math.Mod(h, 360) / 60
	x := cc * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = cc, x, 0
	case hp < 2:
		r, g, b = x, cc, 0
	case hp < 3:
		r, g, b = 0, cc, x
	case hp < 4:
		r, g, b = 0, x, cc
	case hp < 5:
		r, g, b = x, 0, cc
	default:
		r, g, b = cc, 0, x
	}
	m := v - cc
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}
