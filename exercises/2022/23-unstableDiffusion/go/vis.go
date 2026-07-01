package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the elves spreading out (GIF). Starting from the packed initial
// clump, each frame is one round of the diffusion rules; elves are colored by
// how far they sit from the center of mass so the outward "boiling" motion is
// easy to follow. The animation runs until no elf moves — the round that is the
// part-two answer — which is also when the bounding rectangle stops growing.
func (e Exercise) Vis(instr, outdir string) error {
	elfLocations, err := parse(instr)
	if err != nil {
		return err
	}

	// Record every round's positions until the field settles.
	frames := [][]point{snapshot(elfLocations)}
	startDirection := 0
	const maxFrames = 160
	for len(frames) < maxFrames {
		plannedMoves, targetCounts := planElfMoves(elfLocations, startDirection)
		next, moved := updateElfLocations(plannedMoves, targetCounts)
		startDirection++
		elfLocations = next
		frames = append(frames, snapshot(elfLocations))
		if !moved {
			break
		}
	}

	// Fixed bounds from the final (widest) frame.
	last := frames[len(frames)-1]
	minX, maxX := last[0].x, last[0].x
	minY, maxY := last[0].y, last[0].y
	cx, cy := 0, 0
	for _, p := range last {
		minX = min(minX, p.x)
		maxX = max(maxX, p.x)
		minY = min(minY, p.y)
		maxY = max(maxY, p.y)
		cx += p.x
		cy += p.y
	}
	cx /= len(last)
	cy /= len(last)

	const pad = 2
	gw := maxX - minX + 1 + 2*pad
	gh := maxY - minY + 1 + 2*pad
	maxR := 1
	for _, p := range last {
		if r := abs23(p.x-cx) + abs23(p.y-cy); r > maxR {
			maxR = r
		}
	}

	const cell = 4
	pal := color.Palette{color.RGBA{0x0d, 0x0f, 0x18, 0xff}} // 0 background
	for i := 1; i <= 32; i++ {
		f := float64(i-1) / 31
		pal = append(pal, hsv23(f*280, 0.8, 0.95))
	}

	anim := &gif.GIF{}
	for _, fr := range frames {
		img := image.NewPaletted(image.Rect(0, 0, gw*cell, gh*cell), pal)
		for _, p := range fr {
			gx := p.x - minX + pad
			gy := p.y - minY + pad
			r := abs23(p.x-cx) + abs23(p.y-cy)
			idx := uint8(1 + 31*r/maxR)
			if idx > 32 {
				idx = 32
			}
			for dy := 0; dy < cell; dy++ {
				for dx := 0; dx < cell; dx++ {
					img.SetColorIndex(gx*cell+dx, gy*cell+dy, idx)
				}
			}
		}
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 10)
	}
	// Linger on the final settled frame.
	if len(anim.Delay) > 0 {
		anim.Delay[len(anim.Delay)-1] = 200
	}

	f, err := os.Create(filepath.Join(outdir, "unstable-diffusion.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func snapshot(m map[point]string) []point {
	pts := make([]point, 0, len(m))
	for p := range m {
		pts = append(pts, p)
	}
	return pts
}

func abs23(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func hsv23(h, s, v float64) color.RGBA {
	c := v * s
	hp := h / 60
	x := c * (1 - absf23(mod23(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = c, x, 0
	case hp < 2:
		r, g, b = x, c, 0
	case hp < 3:
		r, g, b = 0, c, x
	case hp < 4:
		r, g, b = 0, x, c
	case hp < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	m := v - c
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}

func absf23(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func mod23(a, b float64) float64 {
	for a >= b {
		a -= b
	}
	return a
}
