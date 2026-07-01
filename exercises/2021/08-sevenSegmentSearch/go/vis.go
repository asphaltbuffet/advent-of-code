package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the puzzle's actual subject — seven-segment displays — as an SVG.
// The top row shows the ten digit glyphs 0-9 for reference; below it, several
// entries from the input are decoded and their four-digit outputs drawn as lit
// displays, with the numeric value beside each. Lit segments use a warm accent,
// unlit segments a faint outline, so the encoding reads by shape and brightness
// (fine in grayscale) rather than color alone.
func (e Exercise) Vis(instr, outdir string) error {
	entries, err := parse(instr)
	if err != nil {
		return err
	}

	// Canonical segment sets for 0-9 (segments a..g) for the reference row.
	canonical := []string{
		"abcefg",  // 0
		"cf",      // 1
		"acdeg",   // 2
		"acdfg",   // 3
		"bcdf",    // 4
		"abdfg",   // 5
		"abdefg",  // 6
		"acf",     // 7
		"abcdefg", // 8
		"abcdfg",  // 9
	}

	const (
		dw    = 46 // digit cell width
		dh    = 80 // digit cell height
		dgap  = 12
		rowH  = dh + 46
		mLeft = 24
		mTop  = 40
	)

	sampleN := 6
	if len(entries) < sampleN {
		sampleN = len(entries)
	}

	W := mLeft*2 + 10*(dw+dgap)
	H := mTop + rowH + 30 + sampleN*rowH + 20

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)

	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="16">Seven-segment digits 0-9</text>`, mLeft)
	for d := 0; d < 10; d++ {
		x := mLeft + d*(dw+dgap)
		drawDigit(&sb, x, mTop, dw, dh, segMask(canonical[d]))
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="13" text-anchor="middle">%d</text>`, x+dw/2, mTop+dh+18, d)
	}

	yStart := mTop + rowH + 30
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="16">Decoded outputs from the input</text>`, mLeft, yStart-12)

	for i := 0; i < sampleN; i++ {
		ent := entries[i]
		y := yStart + i*rowH
		for j, o := range ent.output {
			x := mLeft + j*(dw+dgap)
			drawDigit(&sb, x, y, dw, dh, o)
		}
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#56B4E9" font-size="26" dominant-baseline="middle">= %d</text>`,
			mLeft+4*(dw+dgap)+16, y+dh/2, decode(ent))
	}

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "seven-segment.svg"), []byte(sb.String()), 0o600)
}

// drawDigit draws a seven-segment glyph for the given segment mask at (x,y).
// Lit segments are drawn in the accent color; unlit segments as faint outlines.
func drawDigit(sb *strings.Builder, x, y, w, h int, mask uint8) {
	const lit = "#E69F00"    // amber, warm accent
	const unlit = "#242a33"  // faint
	t := 6                   // segment thickness
	fx, fy := float64(x), float64(y)
	fw, fh := float64(w), float64(h)
	ft := float64(t)
	midY := fy + fh/2

	// Horizontal segments: a (top), d (middle), g (bottom).
	horiz := func(cy float64, on bool) {
		col := unlit
		if on {
			col = lit
		}
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="2" fill="%s"/>`,
			fx+ft, cy-ft/2, fw-2*ft, ft, col)
	}
	// Vertical segments given top y and the column x.
	vert := func(cx, top float64, on bool) {
		col := unlit
		if on {
			col = lit
		}
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="2" fill="%s"/>`,
			cx-ft/2, top+ft, ft, fh/2-1.5*ft, col)
	}

	bit := func(seg byte) bool { return mask&(1<<(seg-'a')) != 0 }

	// Standard layout: a top, b top-right, c bottom-right, d bottom,
	// e bottom-left, f top-left, g middle.
	horiz(fy, bit('a'))         // a top
	vert(fx, fy, bit('f'))      // f top-left
	vert(fx+fw, fy, bit('b'))   // b top-right
	horiz(midY, bit('g'))       // g middle
	vert(fx, midY, bit('e'))    // e bottom-left
	vert(fx+fw, midY, bit('c')) // c bottom-right
	horiz(fy+fh, bit('d'))      // d bottom
}
