package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the black tiles on a real hexagonal grid. The tiles flipped by the
// initial paths (Part One) are drawn in one color, and the tiles that are black
// only after 100 days of the flipping rules (Part Two) in another, over the shared
// hex layout — so the small seed pattern and the large grown pattern are both
// visible. The two states differ in brightness as well as color, so the growth
// reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	seed := blackTiles(instr)
	final := blackTiles(instr)
	for day := 0; day < 100; day++ {
		final = step(final)
	}

	// Bounds over all tiles to be drawn (union of seed and final).
	all := map[hex]bool{}
	for t := range seed {
		all[t] = true
	}
	for t := range final {
		all[t] = true
	}
	minQ, maxQ, minR, maxR := 1<<30, -(1 << 30), 1<<30, -(1 << 30)
	for t := range all {
		if t.q < minQ {
			minQ = t.q
		}
		if t.q > maxQ {
			maxQ = t.q
		}
		if t.r < minR {
			minR = t.r
		}
		if t.r > maxR {
			maxR = t.r
		}
	}

	// Pointy-top hex pixel layout from axial coords.
	const size = 7.0
	hexW := math.Sqrt(3) * size
	hexH := 2 * size * 0.75
	px := func(t hex) (float64, float64) {
		x := hexW * (float64(t.q) + float64(t.r)/2)
		y := hexH * float64(t.r)
		return x, y
	}

	// Compute drawing extents.
	minX, maxX, minY, maxY := math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	for t := range all {
		x, y := px(t)
		minX, maxX = math.Min(minX, x), math.Max(maxX, x)
		minY, maxY = math.Min(minY, y), math.Max(maxY, y)
	}
	const pad = 20
	W := int(maxX-minX) + 2*pad + 40
	H := int(maxY-minY) + 2*pad + 60
	ox := pad - minX
	oy := pad + 40 - minY

	seedCol := "#F0E442"  // bright yellow: black at day 0 (part one)
	grownCol := "#0072B2" // blue: black only after day 100 (part two)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="15">Lobby Layout: black tiles after 100 days (%d), seed pattern (%d) in yellow</text>`, pad, len(final), len(seed))

	hexPath := func(cx, cy float64) string {
		var pts []string
		for i := 0; i < 6; i++ {
			ang := math.Pi/180*(60*float64(i)-90)
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", cx+size*math.Cos(ang), cy+size*math.Sin(ang)))
		}
		return strings.Join(pts, " ")
	}

	// Draw grown-only tiles first, then seed tiles on top.
	for t := range final {
		if seed[t] {
			continue
		}
		x, y := px(t)
		fmt.Fprintf(&sb, `<polygon points="%s" fill="%s"/>`, hexPath(x+ox, y+oy), grownCol)
	}
	for t := range seed {
		x, y := px(t)
		fmt.Fprintf(&sb, `<polygon points="%s" fill="%s"/>`, hexPath(x+ox, y+oy), seedCol)
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "lobby-layout.svg"), []byte(sb.String()), 0o600)
}
