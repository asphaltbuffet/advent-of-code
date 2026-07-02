package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis plots every expense entry on a shared value axis and marks the two sets
// that sum to 2020: the Part One pair and the Part Two triple. Each set sits in
// its own lane with connector arcs to the entries and a label giving the sum and
// product. The background entries, the pair, and the triple use three
// colorblind-safe colors that also differ in brightness, so the highlights read
// in grayscale; the marked entries are additionally drawn larger.
func (e Exercise) Vis(instr, outdir string) error {
	entries, err := parse(instr)
	if err != nil {
		return err
	}
	pair := findPair(entries)
	triple := findTriple(entries)

	sorted := append([]int(nil), entries...)
	sort.Ints(sorted)
	lo, hi := sorted[0], sorted[len(sorted)-1]

	const (
		W        = 960
		H        = 420
		mLeft    = 50
		mRight   = 50
		axisY    = 300
		pairY    = 200
		tripleY  = 110
	)
	plotW := W - mLeft - mRight
	xOf := func(v int) int {
		return mLeft + int(float64(v-lo)/float64(hi-lo)*float64(plotW))
	}

	// Okabe-Ito: gray-blue background, orange pair, bluish-green triple.
	bgCol := "#7a869a"
	pairCol := "#E69F00"
	tripCol := "#009E73"

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#e8ecf4" font-size="16">Report Repair: entries summing to 2020</text>`, mLeft)

	// Axis line and end labels.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e" stroke-width="1.5"/>`, mLeft, axisY, W-mRight, axisY)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11" text-anchor="middle">%d</text>`, mLeft, axisY+20, lo)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11" text-anchor="middle">%d</text>`, W-mRight, axisY+20, hi)

	// Background ticks for all entries.
	for _, v := range sorted {
		x := xOf(v)
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1"/>`, x, axisY-6, x, axisY+6, bgCol)
	}

	drawSet := func(vals []int, y int, col, label string) {
		for _, v := range vals {
			x := xOf(v)
			// Connector from lane point down to the axis.
			fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1.5" stroke-dasharray="3 3"/>`, x, y, x, axisY, col)
			fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="6" fill="%s"/>`, x, y, col)
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="11" text-anchor="middle">%d</text>`, x, y-12, v)
		}
		// Bracket line across the set.
		x0, x1 := xOf(minOf(vals)), xOf(maxOf(vals))
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2"/>`, x0, y, x1, y, col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="12" text-anchor="start">%s</text>`, W-mRight+2, y+4, col, label)
	}

	prodPair := 1
	for _, v := range pair {
		prodPair *= v
	}
	prodTrip := 1
	for _, v := range triple {
		prodTrip *= v
	}

	drawSet(triple, tripleY, tripCol, "triple")
	drawSet(pair, pairY, pairCol, "pair")

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="13">Part 1: %d + %d = 2020 → %d</text>`,
		mLeft, 372, pairCol, pair[0], pair[1], prodPair)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="13">Part 2: %d + %d + %d = 2020 → %d</text>`,
		mLeft, 394, tripCol, triple[0], triple[1], triple[2], prodTrip)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "report-repair.svg"), []byte(sb.String()), 0o600)
}

func findPair(entries []int) []int {
	seen := map[int]bool{}
	for _, a := range entries {
		if seen[targetSum-a] {
			return []int{targetSum - a, a}
		}
		seen[a] = true
	}
	return nil
}

func findTriple(entries []int) []int {
	for i, a := range entries {
		need := targetSum - a
		seen := map[int]bool{}
		for _, b := range entries[i+1:] {
			if b < need && seen[need-b] {
				return []int{a, need - b, b}
			}
			seen[b] = true
		}
	}
	return nil
}

func minOf(vs []int) int {
	m := vs[0]
	for _, v := range vs {
		if v < m {
			m = v
		}
	}
	return m
}

func maxOf(vs []int) int {
	m := vs[0]
	for _, v := range vs {
		if v > m {
			m = v
		}
	}
	return m
}
