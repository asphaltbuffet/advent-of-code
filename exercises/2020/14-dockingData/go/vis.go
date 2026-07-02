package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis explains why Part Two fans out. Each mask has some number of floating `X`
// bits, and a write under it touches 2^(that many) addresses. The chart shows the
// distribution of X-counts across the program's masks (bars) and, on a second
// axis, how many addresses each X-count expands to (2^X), making clear that a few
// high-X masks dominate the total writes. Bars and the exponential curve use
// distinct colorblind-safe colors and the curve has point markers, so the chart
// reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	ops, err := parse(instr)
	if err != nil {
		return err
	}

	// Count masks by their number of X (floating) bits.
	xCount := map[int]int{}
	for _, o := range ops {
		if o.isMemory {
			continue
		}
		x := strings.Count(o.mask, "X")
		xCount[x]++
	}
	xs := make([]int, 0, len(xCount))
	for x := range xCount {
		xs = append(xs, x)
	}
	sort.Ints(xs)
	minX, maxX := xs[0], xs[len(xs)-1]

	const (
		W  = 900
		H  = 460
		mL = 60
		mR = 60
		mT = 60
		mB = 70
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	span := maxX - minX
	if span == 0 {
		span = 1
	}
	xOf := func(x int) int { return mL + (x-minX)*plotW/span }

	maxCount := 1
	for _, c := range xCount {
		if c > maxCount {
			maxCount = c
		}
	}
	// Headroom so near-equal bars don't fill the whole frame.
	scaleCount := maxCount * 13 / 10

	barCol := "#0072B2"   // blue: number of masks
	curveCol := "#E69F00" // orange: addresses per write (2^X)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#e8ecf4" font-size="16">Docking Data: floating bits per mask (part 2 fan-out)</text>`, mL)

	// Axes.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT+plotH, W-mR, mT+plotH)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT, mL, mT+plotH)

	// Bars: mask count per X value.
	bw := plotW/(span+1) - 8
	if bw < 6 {
		bw = 6
	}
	for x := minX; x <= maxX; x++ {
		c := xCount[x]
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10" text-anchor="middle">%d</text>`, xOf(x), mT+plotH+16, x)
		if c == 0 {
			continue
		}
		bh := c * plotH / scaleCount
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, xOf(x)-bw/2, mT+plotH-bh, bw, bh, barCol)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="10" text-anchor="middle">%d</text>`, xOf(x), mT+plotH-bh-4, c)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">number of X (floating) bits in the mask</text>`, mL+plotW/2, mT+plotH+42)

	// Curve: 2^X addresses per write, on its own log-ish mapping (use bit height
	// = X directly, since 2^X on a log2 axis is linear — label with the value).
	var pts []string
	for x := minX; x <= maxX; x++ {
		// Map 2^X onto plot height by X/maxX (log2 scale).
		y := mT + plotH - x*plotH/maxX
		pts = append(pts, fmt.Sprintf("%d,%d", xOf(x), y))
		fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="3.5" fill="%s"/>`, xOf(x), y, curveCol)
	}
	fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>`, strings.Join(pts, " "), curveCol)
	// Annotate the largest fan-out.
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="11" text-anchor="end">%d addresses / write</text>`,
		W-mR-16, mT+30, curveCol, 1<<uint(maxX))

	// Legend along the bottom.
	ly := H - 16
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">masks with this many X bits</text>`, mL, ly-10, barCol, mL+18, ly)
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="5" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">addresses per write = 2^X (log2 axis)</text>`, mL+296, ly-4, curveCol, mL+308, ly)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "docking-data.svg"), []byte(sb.String()), 0o600)
}
