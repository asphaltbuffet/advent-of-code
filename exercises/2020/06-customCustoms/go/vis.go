package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis contrasts the two ways of counting a group's yes answers as a function of
// group size. For each group size it plots the mean "anyone" count (union, Part
// One) and the mean "everyone" count (intersection, Part Two): union rises with
// size while intersection falls, because more people means more distinct yeses
// but fewer unanimous ones. Bars behind the curves show how many groups have each
// size. The two series use colorblind-safe colors that also differ in brightness
// and are drawn with distinct markers, so they read in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	groups := parse(instr)

	maxSize := 0
	for _, g := range groups {
		if len(g) > maxSize {
			maxSize = len(g)
		}
	}

	// Per group size: count of groups, sum of union sizes, sum of intersection.
	freq := make([]int, maxSize+1)
	sumUnion := make([]int, maxSize+1)
	sumInter := make([]int, maxSize+1)
	for _, g := range groups {
		n := len(g)
		count := map[rune]int{}
		for _, p := range g {
			for _, q := range p {
				count[q]++
			}
		}
		inter := 0
		for _, c := range count {
			if c == n {
				inter++
			}
		}
		freq[n]++
		sumUnion[n] += len(count)
		sumInter[n] += inter
	}

	const (
		W     = 860
		H     = 460
		mL    = 60
		mR    = 40
		mT    = 60
		mB    = 60
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	// Y axis 0..26 (max distinct questions).
	maxY := 26.0
	xOf := func(size int) int { return mL + (size-1)*plotW/max(maxSize-1, 1) }
	yOf := func(v float64) int { return mT + plotH - int(v/maxY*float64(plotH)) }

	unionCol := "#E69F00" // orange, bright
	interCol := "#0072B2" // blue, dark
	freqCol := "#3a424e"  // gray bars

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="36" fill="#e8ecf4" font-size="16">Custom Customs: answers per group by group size</text>`, mL)

	// Axes.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT+plotH, W-mR, mT+plotH)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT, mL, mT+plotH)
	for yv := 0; yv <= 26; yv += 13 {
		yy := yOf(float64(yv))
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11" text-anchor="end">%d</text>`, mL-6, yy+4, yv)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">group size</text>`, mL+plotW/2, mT+plotH+40)

	// Frequency bars (scaled to plot height by the busiest size).
	maxFreq := 1
	for _, f := range freq {
		if f > maxFreq {
			maxFreq = f
		}
	}
	bw := 14
	for size := 1; size <= maxSize; size++ {
		x := xOf(size)
		// x-axis tick label for every size.
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10" text-anchor="middle">%d</text>`, x, mT+plotH+18, size)
		if freq[size] == 0 {
			continue
		}
		bh := freq[size] * plotH / maxFreq
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x-bw/2, mT+plotH-bh, bw, bh, freqCol)
	}

	// Line series for mean union and mean intersection.
	line := func(sum []int, col, marker string) {
		var pts []string
		for size := 1; size <= maxSize; size++ {
			if freq[size] == 0 {
				continue
			}
			mean := float64(sum[size]) / float64(freq[size])
			x, y := xOf(size), yOf(mean)
			pts = append(pts, fmt.Sprintf("%d,%d", x, y))
			if marker == "circle" {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="4" fill="%s"/>`, x, y, col)
			} else {
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="8" height="8" fill="%s"/>`, x-4, y-4, col)
			}
		}
		fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>`, strings.Join(pts, " "), col)
	}
	line(sumUnion, unionCol, "circle")
	line(sumInter, interCol, "square")

	// Legend.
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="4" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">anyone (union, part 1)</text>`,
		W-mR-250, mT+8, unionCol, W-mR-238, mT+12)
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="8" height="8" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">everyone (intersection, part 2)</text>`,
		W-mR-254, mT+24, interCol, W-mR-238, mT+32)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "custom-customs.svg"), []byte(sb.String()), 0o600)
}
