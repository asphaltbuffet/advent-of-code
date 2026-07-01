package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the vent field as an SVG with two panels side by side: part one
// (horizontal and vertical lines only, left) and part two (diagonals added,
// right). Each vent is drawn as a translucent line so overlaps read as brighter
// stacked strokes, and every danger cell — where two or more vents cross, the
// cells each part counts — is marked with a gold dot. Splitting the parts makes
// the extra crossings the diagonals introduce obvious; SVG keeps the thousands
// of thin lines crisp at any zoom.
func (e Exercise) Vis(instr, outdir string) error {
	segs, err := parse(instr)
	if err != nil {
		return err
	}

	maxX, maxY := 0, 0
	for _, s := range segs {
		maxX = max(maxX, max(s.a.x, s.b.x))
		maxY = max(maxY, max(s.a.y, s.b.y))
	}
	gw := maxX + 1
	gh := maxY + 1

	const pad = 10
	const gap = 40
	const top = 30
	W := gw*2 + gap + 2*pad
	H := gh + top + pad

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#080a10"/>`, W, H)

	panel(&sb, segs, false, pad, top, "PART ONE - H/V only")
	panel(&sb, segs, true, pad+gw+gap, top, "PART TWO - with diagonals")

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "hydrothermal-venture.svg"), []byte(sb.String()), 0o600)
}

// panel writes one heatmap panel: the vent lines as translucent strokes plus a
// gold dot on every cell covered by two or more lines.
func panel(sb *strings.Builder, segs []segment, diagonals bool, ox, oy int, title string) {
	fmt.Fprintf(sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="16">%s</text>`, ox, oy-10, title)
	fmt.Fprintf(sb, `<g transform="translate(%d,%d)">`, ox, oy)

	// Lines: color by orientation, translucent so stacking brightens overlaps.
	fmt.Fprint(sb, `<g stroke-width="1" fill="none">`)
	for _, s := range segs {
		dx := s.b.x - s.a.x
		dy := s.b.y - s.a.y
		if !diagonals && dx != 0 && dy != 0 {
			continue
		}

		// Okabe-Ito colorblind-safe palette; H vs V differ in brightness too so
		// the distinction survives grayscale, and diagonals use vermilion which is
		// well separated from both blues under common deficiencies.
		var stroke string
		switch {
		case dy == 0:
			stroke = "#0072B2" // horizontal: blue
		case dx == 0:
			stroke = "#56B4E9" // vertical: sky blue
		default:
			stroke = "#D55E00" // diagonal: vermilion
		}
		fmt.Fprintf(sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-opacity="0.55"/>`,
			s.a.x, s.a.y, s.b.x, s.b.y, stroke)
	}
	fmt.Fprint(sb, `</g>`)

	// Danger cells: overlap of two or more lines. Yellow (#F0E442) is the
	// brightest Okabe-Ito color, so these solid dots stay distinct from the
	// thinner lines beneath them even in grayscale.
	grid := buildGrid(segs, diagonals)
	fmt.Fprint(sb, `<g fill="#F0E442">`)
	for p, n := range grid {
		if n >= 2 {
			fmt.Fprintf(sb, `<rect x="%d" y="%d" width="1.6" height="1.6"/>`, p.x, p.y)
		}
	}
	fmt.Fprint(sb, `</g>`)

	fmt.Fprint(sb, `</g>`)
}

// buildGrid rasterizes the segments into a coverage grid. When diagonals is
// false, 45° lines are skipped (part one).
func buildGrid(segs []segment, diagonals bool) map[point]int {
	grid := map[point]int{}
	for _, s := range segs {
		dx := s.b.x - s.a.x
		dy := s.b.y - s.a.y
		if !diagonals && dx != 0 && dy != 0 {
			continue
		}
		steps := abs(dx)
		if abs(dy) > steps {
			steps = abs(dy)
		}
		sx, sy := sign(dx), sign(dy)
		p := s.a
		for i := 0; i <= steps; i++ {
			grid[p]++
			p.x += sx
			p.y += sy
		}
	}
	return grid
}
