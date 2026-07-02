package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis plots the number stream and marks the two answers: the first number that
// breaks the pair-sum rule (Part One) and the contiguous range that sums to it
// (Part Two), whose smallest and largest members give the weakness. The range is
// shaded, its min and max are marked, and the invalid value is flagged. Values
// span a wide magnitude, so the y axis is log-scaled. Highlights use
// colorblind-safe colors with distinct shapes and labels, so they read in
// grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	nums, err := parse(instr)
	if err != nil {
		return err
	}
	w := windowFor(nums)
	target, invIdx := firstInvalid(nums, w)

	// Find the contiguous range summing to target.
	rangeLo, rangeHi := -1, -1
	lo, sum := 0, 0
	for hi := 0; hi < len(nums); hi++ {
		sum += nums[hi]
		for sum > target && lo < hi {
			sum -= nums[lo]
			lo++
		}
		if sum == target && hi > lo {
			rangeLo, rangeHi = lo, hi
			break
		}
	}
	// Min and max within the range.
	minIdx, maxIdx := rangeLo, rangeLo
	for i := rangeLo; i <= rangeHi; i++ {
		if nums[i] < nums[minIdx] {
			minIdx = i
		}
		if nums[i] > nums[maxIdx] {
			maxIdx = i
		}
	}

	const (
		W  = 960
		H  = 420
		mL = 60
		mR = 30
		mT = 60
		mB = 50
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	n := len(nums)

	// Log scale over positive values.
	logv := func(v int) float64 {
		if v < 1 {
			v = 1
		}
		return math.Log(float64(v))
	}
	maxL := 0.0
	for _, v := range nums {
		if l := logv(v); l > maxL {
			maxL = l
		}
	}
	xOf := func(i int) int { return mL + i*plotW/(n-1) }
	yOf := func(v int) int { return mT + plotH - int(logv(v)/maxL*float64(plotH)) }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#e8ecf4" font-size="16">Encoding Error: %d values (log scale)</text>`, mL, n)

	// Shade the summing range with a vertical marker line (it is narrow because
	// the values climb, so a thin band plus arrow reads better than a caption).
	if rangeLo >= 0 {
		x0 := xOf(rangeLo)
		x1 := xOf(rangeHi)
		bw := x1 - x0
		if bw < 3 {
			bw = 3
		}
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="#009E73" fill-opacity="0.30"/>`, x0, mT, bw, plotH)
	}

	// Value polyline.
	var pts []string
	for i, v := range nums {
		pts = append(pts, fmt.Sprintf("%d,%d", xOf(i), yOf(v)))
	}
	fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="#56B4E9" stroke-width="1"/>`, strings.Join(pts, " "))

	// Invalid value marker (orange diamond) with a leader up into open space.
	if invIdx >= 0 {
		x, y := xOf(invIdx), yOf(nums[invIdx])
		fmt.Fprintf(&sb, `<polygon points="%d,%d %d,%d %d,%d %d,%d" fill="#E69F00"/>`,
			x, y-7, x+7, y, x, y+7, x-7, y)
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#E69F00" stroke-width="1"/>`, x, y-7, x, mT+40)
	}

	// Min/max within the range (yellow / vermilion squares, no inline labels —
	// values go in the legend to avoid crowding this dense region).
	if rangeLo >= 0 {
		for _, m := range []struct {
			i   int
			col string
		}{{minIdx, "#F0E442"}, {maxIdx, "#D55E00"}} {
			x, y := xOf(m.i), yOf(nums[m.i])
			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="9" height="9" fill="%s"/>`, x-4, y-4, m.col)
		}
	}

	// Legend box in the open upper-left area.
	legend := []struct {
		col, text string
	}{
		{"#E69F00", fmt.Sprintf("invalid value (part 1): %d", target)},
		{"#009E73", fmt.Sprintf("contiguous range: idx %d..%d", rangeLo, rangeHi)},
		{"#F0E442", fmt.Sprintf("range min: %d", nums[minIdx])},
		{"#D55E00", fmt.Sprintf("range max: %d", nums[maxIdx])},
	}
	lx, ly := mL+10, mT+16
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="300" height="%d" fill="#0d0f18" fill-opacity="0.7"/>`, lx-8, ly-14, len(legend)*22+16)
	for i, l := range legend {
		yy := ly + i*22
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/>`, lx, yy-10, l.col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="12">%s</text>`, lx+20, yy, l.text)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14">min + max = %d (part 2)</text>`,
		mL, H-16, nums[minIdx]+nums[maxIdx])

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "encoding-error.svg"), []byte(sb.String()), 0o600)
}
