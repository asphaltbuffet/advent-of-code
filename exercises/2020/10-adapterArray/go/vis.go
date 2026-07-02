package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis shows the sorted adapter chain and where its arrangement count comes from.
// The top panel plots each adapter as a step, colored by the jolt gap to the
// previous one (1 or 3); the 3-gaps are forced links, while runs of 1-gaps are
// where choices — and the branching of Part Two — happen. The bottom panel plots
// the running number of arrangements (log scale), which jumps inside each 1-gap
// run and plateaus across every forced 3-gap. Gap sizes use distinct
// colorblind-safe colors and are labeled, so the plot reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	nums, err := chain(instr)
	if err != nil {
		return err
	}
	n := len(nums)

	// Running arrangement counts (same DP as part two).
	ways := make([]int, n)
	ways[0] = 1
	for i := 1; i < n; i++ {
		for j := i - 1; j >= 0 && nums[i]-nums[j] <= 3; j-- {
			ways[i] += ways[j]
		}
	}

	const (
		W    = 960
		H    = 500
		mL   = 60
		mR   = 30
		mT   = 50
		midG = 40
	)
	topH := 190
	botH := 150
	plotW := W - mL - mR
	xOf := func(i int) int { return mL + i*plotW/(n-1) }

	maxJolt := nums[n-1]
	yTop := func(v int) int { return mT + topH - v*topH/maxJolt }

	gap1 := "#56B4E9" // sky blue: 1-jolt step
	gap3 := "#D55E00" // vermilion: 3-jolt step (forced)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="30" fill="#e8ecf4" font-size="16">Adapter Array: joltage chain and arrangement count</text>`, mL)

	// Top: joltage steps colored AND weighted by gap, so 3-jolt links stay
	// distinct in grayscale (thicker line plus a dot marker).
	for i := 1; i < n; i++ {
		x1, x2 := xOf(i-1), xOf(i)
		y1, y2 := yTop(nums[i-1]), yTop(nums[i])
		col, sw := gap1, 2
		three := nums[i]-nums[i-1] == 3
		if three {
			col, sw = gap3, 4
		}
		// Horizontal tread then vertical riser.
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="%d"/>`, x1, y1, x2, y1, col, sw)
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="%d"/>`, x2, y1, x2, y2, col, sw)
		if three {
			fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="3.5" fill="%s"/>`, x2, (y1+y2)/2, col)
		}
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12">joltage (0 → %d)</text>`, mL, mT+topH+16, maxJolt)

	// Bottom: running arrangement count on a log scale.
	botTop := mT + topH + midG + 20
	maxL := math.Log(float64(ways[n-1]))
	yBot := func(w int) int {
		l := 0.0
		if w > 0 {
			l = math.Log(float64(w))
		}
		return botTop + botH - int(l/maxL*float64(botH))
	}
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, botTop+botH, W-mR, botTop+botH)
	var pts []string
	for i := 0; i < n; i++ {
		pts = append(pts, fmt.Sprintf("%d,%d", xOf(i), yBot(ways[i])))
	}
	fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="#009E73" stroke-width="2"/>`, strings.Join(pts, " "))
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#009E73" font-size="12">arrangements so far (log): final = %d</text>`, mL, botTop-6, ways[n-1])

	// Legend.
	lx := W - mR - 250
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">1-jolt gap (choices here)</text>`, lx, mT-2, gap1, lx+18, mT+9)
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="5" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">3-jolt gap (forced link)</text>`, lx+6, mT+22, gap3, lx+18, mT+27)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "adapter-array.svg"), []byte(sb.String()), 0o600)
}
