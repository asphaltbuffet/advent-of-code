package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the decoded packet tree as an icicle diagram (SVG). Depth runs top
// to bottom: the root operator spans the full width, and each child occupies a
// slice of its parent proportional to the number of leaves beneath it. Operators
// and literals are colored by type from a colorblind-safe palette and labeled
// with an operator glyph (+, ×, min, max, >, <, =) or "lit", so the type reads
// from the symbol as well as the color. This is the whole nested transmission —
// version-sum (part one) and evaluation (part two) are both walks of this tree.
func (e Exercise) Vis(instr, outdir string) error {
	root, err := decode(instr)
	if err != nil {
		return err
	}

	// Depth of the tree bounds the number of horizontal bands.
	maxDepth := treeDepth(root)

	const (
		W       = 1100
		rowH    = 26
		mTop    = 44
		mSide   = 10
		mBottom = 16
	)
	H := mTop + maxDepth*rowH + mBottom

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="26" fill="#e8ecf4" font-size="15" text-anchor="middle">Decoded packet tree (icicle): operators fan out over their sub-packets down to literals</text>`, W/2)

	plotW := float64(W - 2*mSide)
	var draw func(p packet, x, width float64, depth int)
	draw = func(p packet, x, width float64, depth int) {
		y := mTop + depth*rowH
		fill, glyph := typeStyle(p)
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%d" width="%.2f" height="%d" fill="%s" stroke="#111418" stroke-width="1"/>`,
			x, y, width, rowH-2, fill)
		// Only label when there's room for text.
		if width > 22 {
			label := glyph
			if p.typeID == 4 {
				label = "lit"
			}
			fmt.Fprintf(&sb, `<text x="%.2f" y="%d" fill="#0d0f18" font-size="12" font-weight="bold" text-anchor="middle" dominant-baseline="middle">%s</text>`,
				x+width/2, y+(rowH-2)/2, label)
		}

		if len(p.children) == 0 {
			return
		}
		total := 0
		for _, c := range p.children {
			total += leafCount(c)
		}
		cx := x
		for _, c := range p.children {
			cw := width * float64(leafCount(c)) / float64(total)
			draw(c, cx, cw, depth+1)
			cx += cw
		}
	}
	draw(root, mSide, plotW, 0)

	// Legend.
	legend := []struct{ sym, col, name string }{
		{"+", "#0072B2", "sum"}, {"x", "#E69F00", "product"}, {"min", "#009E73", "min"},
		{"max", "#CC79A7", "max"}, {"lit", "#56B4E9", "literal"},
		{"&gt;", "#D55E00", "gt"}, {"&lt;", "#F0E442", "lt"}, {"=", "#8c8c8c", "eq"},
	}
	lx := mSide
	ly := H - 4
	for _, l := range legend {
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/>`, lx, ly-11, l.col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11">%s %s</text>`, lx+16, ly-1, l.sym, l.name)
		lx += 130
	}

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "packet-decoder.svg"), []byte(sb.String()), 0o600)
}

// typeStyle maps a packet's type id to its fill color and operator glyph.
func typeStyle(p packet) (fill, glyph string) {
	// Glyphs are XML-escaped: & < > must not appear literally in SVG text.
	switch p.typeID {
	case 0:
		return "#0072B2", "+"
	case 1:
		return "#E69F00", "x"
	case 2:
		return "#009E73", "min"
	case 3:
		return "#CC79A7", "max"
	case 4:
		return "#56B4E9", ""
	case 5:
		return "#D55E00", "&gt;"
	case 6:
		return "#F0E442", "&lt;"
	case 7:
		return "#8c8c8c", "="
	}
	return "#666666", "?"
}

func treeDepth(p packet) int {
	d := 0
	for _, c := range p.children {
		if cd := treeDepth(c); cd > d {
			d = cd
		}
	}
	return d + 1
}

// leafCount counts literal packets under p (at least 1) to size icicle slices.
func leafCount(p packet) int {
	if len(p.children) == 0 {
		return 1
	}
	n := 0
	for _, c := range p.children {
		n += leafCount(c)
	}
	return n
}
