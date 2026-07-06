package exercises

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Vis draws the space of visited frequencies as a coverage histogram along a
// number line. Starting from 0, the device walks the signed changes over and
// over; Part One is the frequency after one pass, and Part Two is the first
// frequency reached twice. Because each pass drifts the frequency by a fixed
// amount (+416) while the intra-pass path swings across a wide band, the walk
// gradually fills a range of frequencies — and the puzzle ends the moment the
// path lands on a frequency it already occupied. The histogram shows how densely
// each frequency band was visited (why a collision is inevitable), and the one
// value hit twice is marked with a tall pointer and its two step labels. Roles
// are carried by marker shape, brightness, and labels — not hue alone — so the
// chart reads in grayscale.
func (c Exercise) Vis(instr, outdir string) error {
	changes, err := parse(instr)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		return errors.New("no changes to visualize")
	}

	part1, firstStep, minF, maxF, repeatFreq, priorStep, repeatStep := replayFrequencies(changes)

	const (
		w    = 980
		h    = 340
		padL = 60
		padR = 30
		padT = 78
		padB = 70
		bins = 180
	)
	pw := w - padL - padR
	ph := h - padT - padB

	// Bucket visited frequencies into bins across [minF, maxF]; each bin's
	// height is the number of distinct visited frequencies falling in it.
	span := maxF - minF
	if span == 0 {
		span = 1
	}
	binOf := func(f int) int {
		b := (f - minF) * bins / span
		if b >= bins {
			b = bins - 1
		}
		if b < 0 {
			b = 0
		}
		return b
	}
	counts := make([]int, bins)
	for f := range firstStep {
		counts[binOf(f)]++
	}
	maxCount := 1
	for _, ct := range counts {
		if ct > maxCount {
			maxCount = ct
		}
	}
	scaleTop := maxCount * 12 / 10 // headroom

	const (
		fg        = "#e8ecf4"
		dim       = "#c8d0dc"
		gridC     = "#2a303a"
		barCol    = "#56B4E9"
		repeatCol = "#D55E00"
		p1Col     = "#F0E442"
	)

	xOf := func(f int) float64 { return float64(padL) + float64(f-minF)/float64(span)*float64(pw) }
	baseY := padT + ph

	sb := buildFreqHistogramSVG(w, h, padL, padT, pw, bins, counts, minF, maxF, scaleTop, baseY,
		part1, repeatFreq, priorStep, repeatStep, len(firstStep), xOf, fg, dim, gridC, barCol, p1Col, repeatCol)

	return os.WriteFile(filepath.Join(outdir, "chronal-calibration.svg"), []byte(sb), 0o600)
}

func replayFrequencies(changes []int) (int, map[int]int, int, int, int, int, int) {
	part1 := 0
	for _, n := range changes {
		part1 += n
	}
	firstStep := map[int]int{0: 0}
	sum, step, minF, maxF := 0, 0, 0, 0
	var repeatFreq, priorStep, repeatStep int
outer:
	for {
		for _, n := range changes {
			sum += n
			step++
			if sum < minF {
				minF = sum
			}
			if sum > maxF {
				maxF = sum
			}
			if s, ok := firstStep[sum]; ok {
				repeatFreq, repeatStep, priorStep = sum, step, s
				break outer
			}
			firstStep[sum] = step
		}
	}
	return part1, firstStep, minF, maxF, repeatFreq, priorStep, repeatStep
}

func svgAxisTick(sb *strings.Builder, x float64, baseY int, label, col string) {
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="%s" stroke-width="1"/>`,
		x, baseY, x, baseY+5, col)
	fmt.Fprintf(sb, `<text x="%.1f" y="%d" fill="%s" font-size="10" text-anchor="middle">%s</text>`,
		x, baseY+18, col, label)
}

func buildFreqHistogramSVG(
	w, h, padL, padT, pw, bins int,
	counts []int,
	minF, maxF, scaleTop, baseY int,
	part1, repeatFreq, priorStep, repeatStep, nVisited int,
	xOf func(int) float64,
	fg, dim, gridC, barCol, p1Col, repeatCol string,
) string {
	yOf := func(ct int) float64 { return float64(padT) + (1-float64(ct)/float64(scaleTop))*float64(baseY-padT) }

	var sb strings.Builder
	const svgOpen = `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`
	fmt.Fprintf(&sb, svgOpen, w, h, w, h)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, w, h)
	const titleFmt = `<text x="%d" y="30" fill="%s" font-size="16">Chronal Calibration: visited frequencies (first repeat %d)</text>`
	fmt.Fprintf(&sb, titleFmt, padL, fg, repeatFreq)
	const subtitleFmt = `<text x="%d" y="48" fill="%s" font-size="11">%d distinct frequencies visited before the walk lands on one twice</text>`
	fmt.Fprintf(&sb, subtitleFmt, padL, dim, nVisited)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1.5"/>`,
		padL, baseY, padL+pw, baseY, gridC)

	bw := float64(pw) / float64(bins)
	for i, ct := range counts {
		if ct == 0 {
			continue
		}
		x := float64(padL) + float64(i)*bw
		y := yOf(ct)
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.1f" width="%.2f" height="%.1f" fill="%s"/>`,
			x, y, bw*0.9, float64(baseY)-y, barCol)
	}

	svgAxisTick(&sb, xOf(minF), baseY, strconv.Itoa(minF), fg)
	svgAxisTick(&sb, xOf(maxF), baseY, strconv.Itoa(maxF), fg)
	if 0 > minF && 0 < maxF {
		svgAxisTick(&sb, xOf(0), baseY, "0", dim)
	}

	p1x := xOf(part1)
	fmt.Fprintf(&sb, `<path d="M%.1f %d l6 -10 l-12 0 z" fill="%s"/>`, p1x, baseY-2, p1Col)
	fmt.Fprintf(&sb, `<text x="%.1f" y="%d" fill="%s" font-size="10" text-anchor="middle">P1 %d</text>`,
		p1x, baseY-16, p1Col, part1)

	rx := xOf(repeatFreq)
	fmt.Fprintf(&sb, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="%s" stroke-width="2" stroke-dasharray="3 3"/>`,
		rx, padT, rx, baseY, repeatCol)
	fmt.Fprintf(&sb, `<path d="M%.1f %d l7 10 l-7 10 l-7 -10 z" fill="%s"/>`, rx, padT+2, repeatCol)
	fmt.Fprintf(&sb, `<text x="%.1f" y="%d" fill="%s" font-size="12" font-weight="bold" text-anchor="middle">%d</text>`,
		rx, padT-4, repeatCol, repeatFreq)
	const hitTwiceFmt = `<text x="%.1f" y="%d" fill="%s" font-size="10" text-anchor="middle">hit twice: step %d, then step %d</text>`
	fmt.Fprintf(&sb, hitTwiceFmt, rx, baseY+36, repeatCol, priorStep, repeatStep)

	ly := h - 12
	const covLegendFmt = `<rect x="%d" y="%d" width="14" height="9" fill="%s"/><text x="%d" y="%d" fill="%s" font-size="10">visited frequencies (coverage)</text>`
	fmt.Fprintf(&sb, covLegendFmt, padL, ly-8, barCol, padL+20, ly, fg)
	const p1LegFmt = `<path d="M%d %d l5 -9 l-10 0 z" fill="%s"/><text x="%d" y="%d" fill="%s" font-size="10">Part 1 (%d)</text>`
	fmt.Fprintf(&sb, p1LegFmt, padL+280, ly, p1Col, padL+292, ly, fg, part1)
	const p2LegFmt = `<path d="M%d %d l5 6 l-5 6 l-5 -6 z" fill="%s"/><text x="%d" y="%d" fill="%s" font-size="10">Part 2 repeat (%d)</text>`
	fmt.Fprintf(&sb, p2LegFmt, padL+430, ly-6, repeatCol, padL+442, ly, fg, repeatFreq)

	sb.WriteString(`</svg>`)
	return sb.String()
}
