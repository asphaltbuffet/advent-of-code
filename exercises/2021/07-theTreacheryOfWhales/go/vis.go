package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis plots the total fuel cost against every candidate alignment position as an
// SVG, one panel per part stacked over a shared x-axis. Part one's linear cost is
// a piecewise-linear "V" minimized at the median; part two's triangular cost is a
// smooth parabola minimized near the mean. Each panel marks its minimum — the
// answer — and a faint tick row at the bottom shows where the crabs actually sit.
// Colors come from a colorblind-safe palette and the minima are marked by
// position too, so the chart reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	positions, err := parse(instr)
	if err != nil {
		return err
	}

	lo, hi := positions[0], positions[0]
	for _, p := range positions {
		if p < lo {
			lo = p
		}
		if p > hi {
			hi = p
		}
	}
	span := hi - lo

	// Cost curves across the range for both parts.
	linear := make([]int, span+1)
	triangular := make([]int, span+1)
	for i := range linear {
		target := lo + i
		var l, t int
		for _, p := range positions {
			d := abs(p - target)
			l += d
			t += d * (d + 1) / 2
		}
		linear[i] = l
		triangular[i] = t
	}
	linMin, linAt := argmin(linear, lo)
	triMin, triAt := argmin(triangular, lo)

	const (
		W       = 960
		panelH  = 240
		gap     = 40
		mLeft   = 90
		mRight  = 30
		mTop    = 34
		tickRow = 22
	)
	H := mTop + panelH*2 + gap + tickRow + 44
	plotW := W - mLeft - mRight

	xOf := func(pos int) float64 {
		if span == 0 {
			return float64(mLeft)
		}
		return float64(mLeft) + float64(pos-lo)/float64(span)*float64(plotW)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)

	// Two panels.
	drawPanel(&sb, "Part One - linear fuel", linear, lo, xOf, mTop, panelH, plotW, mLeft, "#56B4E9", linMin, linAt)
	top2 := mTop + panelH + gap
	drawPanel(&sb, "Part Two - triangular fuel", triangular, lo, xOf, top2, panelH, plotW, mLeft, "#E69F00", triMin, triAt)

	// Crab position ticks along the bottom.
	tickY := top2 + panelH + 26
	seen := map[int]bool{}
	fmt.Fprintf(&sb, `<text x="%d" y="%.1f" fill="#9aa4b2" font-size="12" text-anchor="end">crabs</text>`, mLeft-8, float64(tickY)+10)
	for _, p := range positions {
		if seen[p] {
			continue
		}
		seen[p] = true
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="#009E73" stroke-opacity="0.5"/>`,
			xOf(p), tickY, xOf(p), tickY+14)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14" text-anchor="middle">horizontal position</text>`, mLeft+plotW/2, H-8)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "crab-alignment.svg"), []byte(sb.String()), 0o600)
}

func drawPanel(sb *strings.Builder, title string, cost []int, lo int, xOf func(int) float64, top, h, plotW, mLeft int, col string, minVal, minAt int) {
	maxVal := 0
	for _, v := range cost {
		if v > maxVal {
			maxVal = v
		}
	}
	yOf := func(v int) float64 {
		return float64(top) + (1-float64(v)/float64(maxVal))*float64(h)
	}

	// Frame + baseline.
	fmt.Fprintf(sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#e8ecf4" stroke-width="1.5"/>`, mLeft, top+h, mLeft+plotW, top+h)
	fmt.Fprintf(sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#e8ecf4" stroke-width="1.5"/>`, mLeft, top, mLeft, top+h)
	fmt.Fprintf(sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14">%s</text>`, mLeft, top-8, title)

	// Cost curve.
	var path strings.Builder
	for i, v := range cost {
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&path, "%s%.1f %.1f ", cmd, xOf(lo+i), yOf(v))
	}
	fmt.Fprintf(sb, `<path d="%s" fill="none" stroke="%s" stroke-width="2"/>`, strings.TrimSpace(path.String()), col)

	// Minimum marker.
	mx := xOf(minAt)
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="#e8ecf4" stroke-width="1" stroke-dasharray="4 3"/>`, mx, top, mx, top+h)
	fmt.Fprintf(sb, `<circle cx="%.1f" cy="%.1f" r="4" fill="#e8ecf4"/>`, mx, yOf(minVal))
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" fill="#e8ecf4" font-size="13" text-anchor="middle">min at %d: %s fuel</text>`,
		mx, yOf(minVal)-10, minAt, commafy(int64(minVal)))
}

func argmin(cost []int, lo int) (val, at int) {
	val, at = cost[0], lo
	for i, v := range cost {
		if v < val {
			val, at = v, lo+i
		}
	}
	return val, at
}

// commafy formats n with thousands separators for readable labels.
func commafy(n int64) string {
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, digit := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, digit)
	}
	return string(out)
}
