package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis shows how the Part Two incremental sieve converges. Each bus locked in
// raises the running timestamp t and multiplies the step by that bus's id, so t
// climbs toward 825305207525452 while the stride grows exponentially — the whole
// search finishing in as many iterations as there are buses. Two log-scale series
// (t and step) are plotted against the bus lock-in order, with distinct
// colorblind-safe colors and markers so they read in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	_, buses, err := parse(instr)
	if err != nil {
		return err
	}

	// Replay the sieve, recording t and step after each bus.
	type stage struct {
		id, offset int
		t, step    int
	}
	var stages []stage
	t, step := 0, 1
	for _, b := range buses {
		for (t+b.offset)%b.id != 0 {
			t += step
		}
		step *= b.id
		stages = append(stages, stage{b.id, b.offset, t, step})
	}
	n := len(stages)

	const (
		W  = 900
		H  = 480
		mL = 70
		mR = 40
		mT = 60
		mB = 70
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	xOf := func(i int) int { return mL + i*plotW/max(n-1, 1) }

	maxL := 0.0
	for _, s := range stages {
		if l := math.Log(float64(s.step)); l > maxL {
			maxL = l
		}
	}
	yOf := func(v int) int {
		l := 0.0
		if v > 0 {
			l = math.Log(float64(v))
		}
		return mT + plotH - int(l/maxL*float64(plotH))
	}

	tCol := "#E69F00"    // orange: timestamp t
	stepCol := "#0072B2" // blue: step size

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="32" fill="#e8ecf4" font-size="16">Shuttle Search: sieve converging to t = %d (log scale)</text>`, mL, stages[n-1].t)

	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT+plotH, W-mR, mT+plotH)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT, mL, mT+plotH)

	series := func(get func(stage) int, col, marker string) {
		var pts []string
		for i, s := range stages {
			x, y := xOf(i), yOf(get(s))
			pts = append(pts, fmt.Sprintf("%d,%d", x, y))
			if marker == "circle" {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="4" fill="%s"/>`, x, y, col)
			} else {
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="8" height="8" fill="%s"/>`, x-4, y-4, col)
			}
		}
		fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>`, strings.Join(pts, " "), col)
	}
	series(func(s stage) int { return s.step }, stepCol, "square")
	series(func(s stage) int { return s.t }, tCol, "circle")

	// X axis: bus id at each lock-in.
	for i, s := range stages {
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10" text-anchor="middle">%d</text>`, xOf(i), mT+plotH+16, s.id)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">bus id locked in (in schedule order)</text>`, mL+plotW/2, mT+plotH+40)

	// Legend.
	lx := W - mR - 210
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="4" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">timestamp t</text>`, lx, mT, tCol, lx+14, mT+4)
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="8" height="8" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">step (product of ids)</text>`, lx-4, mT+16, stepCol, lx+14, mT+24)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "shuttle-search.svg"), []byte(sb.String()), 0o600)
}
