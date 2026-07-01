package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis draws each navigation line as a nesting-depth profile (SVG): the height at
// each character is how deeply the brackets are nested there. A sample of lines
// is stacked vertically. Corrupted lines (part one) are drawn in vermilion and
// cut off at the offending bracket, marked with an X; incomplete lines (part two)
// are drawn in blue with their completion tail shown faded. Status reads by
// position and marker as well as color, so it survives grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	lines := strings.Split(strings.TrimSpace(instr), "\n")

	// Keep the profiles a manageable count and width.
	const sample = 28
	if len(lines) > sample {
		lines = lines[:sample]
	}
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}

	const (
		cw     = 6  // px per character column
		rowH   = 40 // px per line row
		depthU = 3  // px per nesting level
		mLeft  = 40
		mTop   = 48
		mRight = 30
	)
	W := mLeft + maxLen*cw + mRight
	H := mTop + len(lines)*rowH + 20

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="20" fill="#e8ecf4" font-size="14">Bracket nesting depth per line</text>`, mLeft)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#9aa4b2" font-size="12">corrupted = vermilion, cut off at X · incomplete = blue with faded completion tail</text>`, mLeft)

	const (
		corruptCol = "#D55E00" // vermilion
		okCol      = "#56B4E9" // sky blue
		tailCol    = "#2f4a63" // faded blue for completion
	)

	for li, line := range lines {
		baseY := mTop + li*rowH + rowH - 6
		corrupt, completion := scanLine(line)

		// Walk the line building the depth profile up to the corruption point.
		depth := 0
		var pts []string
		lineCol := okCol
		errAt := -1
		var stack []rune
		for i, c := range line {
			if _, ok := pairs[c]; ok {
				stack = append(stack, pairs[c])
				depth++
			} else {
				if len(stack) == 0 || stack[len(stack)-1] != c {
					errAt = i
					lineCol = corruptCol
					break
				}
				stack = stack[:len(stack)-1]
				depth--
			}
			x := mLeft + i*cw
			pts = append(pts, fmt.Sprintf("%d,%d", x, baseY-depth*depthU))
		}

		// Baseline for the row.
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#242a33"/>`, mLeft, baseY, mLeft+maxLen*cw, baseY)

		if len(pts) > 1 {
			fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="1.5"/>`, strings.Join(pts, " "), lineCol)
		}

		if corrupt != 0 && errAt >= 0 {
			// Mark the corruption point with an X.
			ex := mLeft + errAt*cw
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="14" text-anchor="middle">%c</text>`, ex, baseY-depth*depthU-6, corruptCol, 'x')
			fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="3" fill="%s"/>`, ex, baseY-depth*depthU, corruptCol)
		} else if len(completion) > 0 {
			// Draw the completion tail as a descending faded profile.
			var tail []string
			x := mLeft + len(line)*cw
			tail = append(tail, fmt.Sprintf("%d,%d", x, baseY-depth*depthU))
			for _, c := range completion {
				_ = c
				depth--
				x += cw
				tail = append(tail, fmt.Sprintf("%d,%d", x, baseY-depth*depthU))
			}
			fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="1.5" stroke-dasharray="3 2"/>`, strings.Join(tail, " "), tailCol)
		}
	}

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "syntax-scoring.svg"), []byte(sb.String()), 0o600)
}
