package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the routing diagram (faint) with the traversed path overlaid as a
// progress-coloured polyline (teal → gold → magenta). Letter waypoints are
// circled and labelled.
func (e Exercise) Vis(instr, outdir string) error {
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	_, _, path := walk(instr)

	rows := len(grid)
	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	const step = 6.0
	const pad = 14.0
	w := float64(cols)*step + 2*pad
	h := float64(rows)*step + 2*pad
	px := func(c int) float64 { return pad + float64(c)*step + step/2 }
	py := func(r int) float64 { return pad + float64(r)*step + step/2 }

	// heat maps progress 0..1 onto teal -> gold -> magenta.
	heat := func(t float64) string {
		var r, g, b int
		if t < 0.5 {
			u := t / 0.5
			r = int(0x2f + u*(0xff-0x2f))
			g = int(0x8a + u*(0xc8-0x8a))
			b = int(0x86 + u*(0x4a-0x86))
		} else {
			u := (t - 0.5) / 0.5
			r = 0xff
			g = int(0xc8 - u*(0xc8-0x44))
			b = int(0x4a + u*(0xd0-0x4a))
		}
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" font-family="monospace">`+"\n",
		w,
		h,
	)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0d0f18"/>`+"\n", w, h)

	// Faint diagram underneath.
	fmt.Fprint(&b, `<g fill="#232a3e">`+"\n")
	for r := range rows {
		for c := range len(grid[r]) {
			if grid[r][c] != ' ' {
				fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f"/>`+"\n",
					pad+float64(c)*step, pad+float64(r)*step, step, step)
			}
		}
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Path as a progress-coloured polyline, drawn as short segments so the hue
	// shifts along its length.
	fmt.Fprint(&b, `<g stroke-width="2.4" stroke-linecap="round" fill="none">`+"\n")
	n := len(path)
	for i := 1; i < n; i++ {
		t := float64(i) / math.Max(1, float64(n-1))
		a, c := path[i-1], path[i]
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s"/>`+"\n",
			px(a.c), py(a.r), px(c.c), py(c.r), heat(t))
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Letters: a dark disc with the glyph on top, so each waypoint is readable.
	for _, p := range path {
		if !p.isLetter {
			continue
		}
		fmt.Fprintf(
			&b,
			`<circle cx="%.1f" cy="%.1f" r="7" fill="#0d0f18" stroke="#ffffff" stroke-width="1.2"/>`+"\n",
			px(p.c),
			py(p.r),
		)
		fmt.Fprintf(
			&b,
			`<text x="%.1f" y="%.1f" font-size="9" font-weight="bold" fill="#ffffff" text-anchor="middle">%c</text>`+"\n",
			px(p.c),
			py(p.r)+3.2,
			p.ch,
		)
	}

	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "a-series-of-tubes.svg"), []byte(b.String()), 0o644)
}
