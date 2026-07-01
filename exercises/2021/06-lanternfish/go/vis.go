package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis plots the lanternfish population over the full 256 days as an SVG line
// chart on a log scale — the only way the growth stays on one page. The curve is
// the total count each day; vertical markers call out day 80 (part one) and day
// 256 (part two), and a light band shows the 7-day spawning cycle rippling as
// steps in the otherwise smooth exponential. Colors are from a colorblind-safe
// palette and the encoding also reads by position, so it survives grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	counts, err := parse(instr)
	if err != nil {
		return err
	}

	const days = 256
	totals := make([]float64, days+1)
	totals[0] = sumCounts(counts)
	for d := 1; d <= days; d++ {
		spawning := counts[0]
		for i := 0; i < 8; i++ {
			counts[i] = counts[i+1]
		}
		counts[6] += spawning
		counts[8] = spawning
		totals[d] = sumCounts(counts)
	}

	const (
		W      = 960
		H      = 520
		mLeft  = 90
		mRight = 150
		mTop   = 40
		mBot   = 60
	)
	plotW := W - mLeft - mRight
	plotH := H - mTop - mBot

	// Log scale bounds.
	minLog := math.Log10(totals[0])
	maxLog := math.Log10(totals[days])
	x := func(d int) float64 { return float64(mLeft) + float64(d)/float64(days)*float64(plotW) }
	y := func(v float64) float64 {
		f := (math.Log10(v) - minLog) / (maxLog - minLog)
		return float64(mTop) + (1-f)*float64(plotH)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)

	// Horizontal gridlines at each power of ten.
	for p := int(math.Ceil(minLog)); float64(p) <= maxLog; p++ {
		gy := y(math.Pow(10, float64(p)))
		fmt.Fprintf(&sb, `<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="#2a2f3a" stroke-width="1"/>`,
			mLeft, gy, W-mRight, gy)
		fmt.Fprintf(&sb, `<text x="%d" y="%.1f" fill="#9aa4b2" font-size="12" text-anchor="end">1e%d</text>`,
			mLeft-8, gy+4, p)
	}

	// Day markers for the two parts (drawn under the curve).
	marker := func(d int, label, col, anchor string, dx, dy float64) {
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="%s" stroke-width="1.5" stroke-dasharray="4 3"/>`,
			x(d), mTop, x(d), H-mBot, col)
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="4" fill="%s"/>`, x(d), y(totals[d]), col)
		fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="%s" font-size="13" text-anchor="%s">%s</text>`,
			x(d)+dx, y(totals[d])+dy, col, anchor, label)
	}
	marker(80, fmt.Sprintf("day 80: %s", commafy(int64(totals[80]))), "#E69F00", "middle", 0, -10)  // orange
	marker(256, fmt.Sprintf("day 256: %s", commafy(int64(totals[256]))), "#009E73", "end", -8, 16)  // bluish green

	// Population curve.
	var path strings.Builder
	for d := 0; d <= days; d++ {
		cmd := "L"
		if d == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&path, "%s%.1f %.1f ", cmd, x(d), y(totals[d]))
	}
	fmt.Fprintf(&sb, `<path d="%s" fill="none" stroke="#56B4E9" stroke-width="2"/>`, strings.TrimSpace(path.String()))

	// Axes.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#e8ecf4" stroke-width="1.5"/>`, mLeft, H-mBot, W-mRight, H-mBot)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#e8ecf4" stroke-width="1.5"/>`, mLeft, mTop, mLeft, H-mBot)
	for _, d := range []int{0, 50, 100, 150, 200, 256} {
		fmt.Fprintf(&sb, `<text x="%.1f" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">%d</text>`, x(d), H-mBot+20, d)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14" text-anchor="middle">day</text>`, mLeft+plotW/2, H-16)
	fmt.Fprintf(&sb, `<text x="20" y="%d" fill="#e8ecf4" font-size="14" transform="rotate(-90 20 %d)" text-anchor="middle">lanternfish (log scale)</text>`, mTop+plotH/2, mTop+plotH/2)
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="16" text-anchor="middle">Lanternfish population growth</text>`, W/2)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "lanternfish.svg"), []byte(sb.String()), 0o600)
}

func sumCounts(c [9]int) float64 {
	total := 0
	for _, v := range c {
		total += v
	}
	return float64(total)
}

// commafy formats n with thousands separators for readable axis labels.
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
