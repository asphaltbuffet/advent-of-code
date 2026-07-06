package exercises

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis draws the reacted polymer length after removing each unit type A-Z as a
// bar chart. The dashed baseline is Part One (remove nothing); each bar shows the
// length once that one letter is stripped and the rest reacted. The shortest bar
// — the best letter to remove, the Part Two answer — is highlighted and labeled.
// Lengths read as bar height and the answer by outline and label, so the chart
// reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	polymer := strings.TrimSpace(instr)
	if polymer == "" {
		return errors.New("empty polymer")
	}

	base := react(polymer, 0)

	lengths := make([]int, 26)
	best, bestLetter := base+1, byte('a')
	for i := range 26 {
		c := byte('a' + i)
		lengths[i] = react(polymer, c)
		if lengths[i] < best {
			best, bestLetter = lengths[i], c
		}
	}

	maxLen := base
	for _, l := range lengths {
		if l > maxLen {
			maxLen = l
		}
	}

	const (
		imgW = 900
		imgH = 420
		padL = 56
		padR = 24
		padT = 56
		padB = 56
	)
	pw := imgW - padL - padR
	ph := imgH - padT - padB
	baseY := float64(padT + ph)

	fg := "#e8ecf4"
	dim := "#c8d0dc"
	grid := "#2a303a"
	barCol := "#56B4E9"  // sky blue: per-letter reacted length
	bestCol := "#D55E00" // vermilion: the best letter to remove (Part Two)
	baseCol := "#F0E442" // yellow: Part One baseline (remove nothing)

	scaleTop := maxLen * 11 / 10 // headroom
	yOf := func(l int) float64 { return float64(padT) + (1-float64(l)/float64(scaleTop))*float64(ph) }

	var sb strings.Builder
	drawPolymerHeader(&sb, imgW, imgH, padL, pw, base, maxLen, baseY, yOf, fg, dim, grid, baseCol)

	drawPolymerBars(&sb, lengths, bestLetter, best, padL, pw, baseY, yOf, barCol, bestCol, dim)
	drawPolymerLegend(&sb, padL, imgH, base, best, fg, barCol, bestCol, baseCol)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "alchemical-reduction.svg"), []byte(sb.String()), 0o600)
}

func drawPolymerHeader(
	sb *strings.Builder, imgW, imgH, padL, pw, base, maxLen int, baseY float64,
	yOf func(int) float64, fg, dim, grid, baseCol string,
) {
	fmt.Fprintf(sb,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`,
		imgW, imgH, imgW, imgH)
	fmt.Fprintf(sb, `<rect width="%d" height="%d" fill="#111418"/>`, imgW, imgH)
	const polyTitleFmt = `<text x="%d" y="30" fill="%s" font-size="16">Alchemical Reduction: reacted length after removing each unit type</text>`
	fmt.Fprintf(sb, polyTitleFmt, padL, fg)
	fmt.Fprintf(sb,
		`<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="%s" stroke-width="1.5"/>`,
		padL, baseY, padL+pw, baseY, grid)
	fmt.Fprintf(sb, `<text x="%d" y="%.1f" fill="%s" font-size="10" text-anchor="end">0</text>`,
		padL-6, baseY+3, dim)
	fmt.Fprintf(sb, `<text x="%d" y="%.1f" fill="%s" font-size="10" text-anchor="end">%d</text>`,
		padL-6, yOf(maxLen)+3, dim, maxLen)
	by := yOf(base)
	const dashedLine = `<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="%s" stroke-width="1.5" stroke-dasharray="7 4"/>`
	fmt.Fprintf(sb, dashedLine, padL, by, padL+pw, by, baseCol)
	fmt.Fprintf(sb,
		`<text x="%d" y="%.1f" fill="%s" font-size="11">Part 1 baseline = %d (remove nothing)</text>`,
		padL+6, by-5, baseCol, base)
}

func drawPolymerBars(
	sb *strings.Builder, lengths []int, bestLetter byte, best, padL, pw int,
	baseY float64, yOf func(int) float64, barCol, bestCol, dim string,
) {
	slot := float64(pw) / 26
	bw := slot * 0.7
	for i := range 26 {
		x := float64(padL) + float64(i)*slot + (slot-bw)/2
		y := yOf(lengths[i])
		col := barCol
		if byte('a'+i) == bestLetter {
			col = bestCol
		}
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`, x, y, bw, baseY-y, col)
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" fill="%s" font-size="10" text-anchor="middle">%c</text>`,
			x+bw/2, baseY+16, dim, 'A'+i)
	}
	i := int(bestLetter - 'a')
	x := float64(padL) + float64(i)*slot + (slot-bw)/2
	y := yOf(best)
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="none" stroke="%s" stroke-width="2"/>`,
		x-1, y-1, bw+2, baseY-y+1, bestCol)
	const winnerLblFmt = `<text x="%.1f" y="%.1f" fill="%s" font-size="12" font-weight="bold" text-anchor="middle">%c: %d</text>`
	fmt.Fprintf(sb, winnerLblFmt, x+bw/2, y-6, bestCol, 'A'+bestLetter-'a', best)
}

func drawPolymerLegend(sb *strings.Builder, padL, imgH, base, best int, fg, barCol, bestCol, baseCol string) {
	ly := imgH - 14
	const reactLegFmt = `<rect x="%d" y="%d" width="12" height="9" fill="%s"/><text x="%d" y="%d" fill="%s" font-size="10">reacted length</text>`
	fmt.Fprintf(sb, reactLegFmt, padL, ly-8, barCol, padL+18, ly, fg)
	const bestLegFmt = `<rect x="%d" y="%d" width="12" height="9" fill="%s"/><text x="%d" y="%d" fill="%s" font-size="10">best to remove (Part 2 = %d)</text>`
	fmt.Fprintf(sb, bestLegFmt, padL+150, ly-8, bestCol, padL+168, ly, fg, best)
	const p1BaseLegFmt = `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1.5" stroke-dasharray="7 4"/><text x="%d" y="%d" fill="%s" font-size="10">Part 1 = %d</text>`
	fmt.Fprintf(sb, p1BaseLegFmt, padL+420, ly-4, padL+450, ly-4, baseCol, padL+456, ly, fg, base)
}
