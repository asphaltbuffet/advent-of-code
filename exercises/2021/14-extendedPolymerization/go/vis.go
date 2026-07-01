package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis shows how the polymer's element composition evolves over the 40 insertion
// steps as an SVG. The left panel plots the total length on a log scale — a
// near-straight line, confirming the exponential growth that makes tracking the
// string directly impossible. The right panel is a normalized stacked band chart
// of each element's share, showing the composition settling into a steady mix
// even as the total explodes. Each band is labeled with its element letter, so it
// reads without relying on color alone.
func (e Exercise) Vis(instr, outdir string) error {
	tmpl, rules, err := parse(instr)
	if err != nil {
		return err
	}

	const steps = 40

	// Track element counts at every step using the pair-count method.
	pairs := map[string]int{}
	for i := 0; i+1 < len(tmpl); i++ {
		pairs[tmpl[i:i+2]]++
	}

	elements := sortedElements(tmpl, rules)
	totals := make([]float64, steps+1)
	// composition[step][elementIndex] = count
	composition := make([][]float64, steps+1)

	record := func(step int) {
		counts := map[byte]int{tmpl[len(tmpl)-1]: 1}
		for pair, n := range pairs {
			counts[pair[0]] += n
		}
		row := make([]float64, len(elements))
		total := 0.0
		for i, el := range elements {
			row[i] = float64(counts[el])
			total += row[i]
		}
		composition[step] = row
		totals[step] = total
	}
	record(0)
	for s := 1; s <= steps; s++ {
		next := make(map[string]int, len(pairs))
		for pair, n := range pairs {
			if ins, ok := rules[pair]; ok {
				next[string([]byte{pair[0], ins})] += n
				next[string([]byte{ins, pair[1]})] += n
			} else {
				next[pair] += n
			}
		}
		pairs = next
		record(s)
	}

	const (
		W      = 1000
		H      = 460
		gap    = 60
		mLeft  = 70
		mTop   = 40
		mBot   = 50
		mRight = 30
	)
	panelW := (W - mLeft - mRight - gap) / 2
	plotH := H - mTop - mBot

	xOf := func(s, ox int) float64 { return float64(ox) + float64(s)/float64(steps)*float64(panelW) }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)

	// ---- Left panel: total length, log scale. ----
	minLog := math.Log10(totals[0])
	maxLog := math.Log10(totals[steps])
	yLog := func(v float64) float64 {
		f := (math.Log10(v) - minLog) / (maxLog - minLog)
		return float64(mTop) + (1-f)*float64(plotH)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="14">Total length (log scale)</text>`, mLeft)
	for p := int(math.Ceil(minLog)); float64(p) <= maxLog; p++ {
		gy := yLog(math.Pow(10, float64(p)))
		fmt.Fprintf(&sb, `<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="#2a2f3a"/>`, mLeft, gy, mLeft+panelW, gy)
		fmt.Fprintf(&sb, `<text x="%d" y="%.1f" fill="#9aa4b2" font-size="11" text-anchor="end">1e%d</text>`, mLeft-6, gy+4, p)
	}
	var lp strings.Builder
	for s := 0; s <= steps; s++ {
		cmd := "L"
		if s == 0 {
			cmd = "M"
		}
		fmt.Fprintf(&lp, "%s%.1f %.1f ", cmd, xOf(s, mLeft), yLog(totals[s]))
	}
	fmt.Fprintf(&sb, `<path d="%s" fill="none" stroke="#56B4E9" stroke-width="2"/>`, strings.TrimSpace(lp.String()))
	axisLabels(&sb, mLeft, panelW, mTop, plotH, steps)

	// ---- Right panel: normalized stacked composition. ----
	ox := mLeft + panelW + gap
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="14">Element composition (share)</text>`, ox)

	pal := bandPalette(len(elements))

	// cumFrac returns the fraction of the total at or below band k at step s.
	cumFrac := func(s, k int) float64 {
		var sum float64
		for i := 0; i <= k; i++ {
			sum += composition[s][i]
		}
		return sum / totals[s]
	}
	yFrac := func(f float64) float64 { return float64(mTop) + (1-f)*float64(plotH) }

	// Build one stacked band polygon per element: top boundary left-to-right,
	// bottom boundary right-to-left.
	for ei := range elements {
		pts := make([]string, 0, 2*(steps+1))
		for s := 0; s <= steps; s++ {
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", xOf(s, ox), yFrac(cumFrac(s, ei))))
		}
		for s := steps; s >= 0; s-- {
			lower := 0.0
			if ei > 0 {
				lower = cumFrac(s, ei-1)
			}
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", xOf(s, ox), yFrac(lower)))
		}
		fmt.Fprintf(&sb, `<polygon points="%s" fill="%s" stroke="#111418" stroke-width="0.4"/>`, strings.Join(pts, " "), pal[ei])

		// Label the band, centered in its final share at the right edge.
		lower := 0.0
		if ei > 0 {
			lower = cumFrac(steps, ei-1)
		}
		midF := (lower + cumFrac(steps, ei)) / 2
		fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#0d0f18" font-size="12" font-weight="bold" text-anchor="middle" dominant-baseline="middle">%c</text>`,
			xOf(steps, ox)-8, yFrac(midF), elements[ei])
	}
	axisLabels(&sb, ox, panelW, mTop, plotH, steps)

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="13" text-anchor="middle">insertion step</text>`, W/2, H-14)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "extended-polymerization.svg"), []byte(sb.String()), 0o600)
}

func axisLabels(sb *strings.Builder, ox, panelW, mTop, plotH, steps int) {
	fmt.Fprintf(sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#e8ecf4"/>`, ox, mTop+plotH, ox+panelW, mTop+plotH)
	for _, s := range []int{0, 10, 20, 30, 40} {
		x := float64(ox) + float64(s)/float64(steps)*float64(panelW)
		fmt.Fprintf(sb, `<text x="%.1f" y="%d" fill="#9aa4b2" font-size="11" text-anchor="middle">%d</text>`, x, mTop+plotH+16, s)
	}
}

// sortedElements returns the distinct elements sorted, so band order is stable.
func sortedElements(tmpl string, rules map[string]byte) []byte {
	set := map[byte]bool{}
	for i := 0; i < len(tmpl); i++ {
		set[tmpl[i]] = true
	}
	for pair, ins := range rules {
		set[pair[0]] = true
		set[pair[1]] = true
		set[ins] = true
	}
	var els []byte
	for b := range set {
		els = append(els, b)
	}
	sort.Slice(els, func(i, j int) bool { return els[i] < els[j] })
	return els
}

// bandPalette returns n visually separable fills; adjacent bands alternate light
// and dark tones so neighbors stay distinct even without color perception.
func bandPalette(n int) []string {
	base := []string{
		"#0072B2", "#E69F00", "#009E73", "#CC79A7", "#56B4E9",
		"#D55E00", "#F0E442", "#8c8c8c", "#7a5195", "#bc5090",
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = base[i%len(base)]
	}
	return out
}
