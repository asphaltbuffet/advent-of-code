package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis shows how making rules 8 and 11 recursive (Part Two) enlarges the set of
// messages that match rule 0. Each message is classified: matched under both rule
// sets, matched only after the Part Two rewrite, or matched by neither. The
// stacked bar and counts make clear that Part Two's matches are a strict superset
// of Part One's — the recursive rules only ever accept more. Groups use
// colorblind-safe colors ordered by brightness and every segment is labeled, so
// the chart reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	g1, messages, err := parse(instr)
	if err != nil {
		return err
	}
	// Part Two grammar: a copy with 8 and 11 rewritten.
	g2, _, _ := parse(instr)
	if _, ok := g2[8]; ok {
		g2[8] = rule{alts: [][]int{{42}, {42, 8}}}
	}
	if _, ok := g2[11]; ok {
		g2[11] = rule{alts: [][]int{{42, 31}, {42, 11, 31}}}
	}

	var both, twoOnly, neither int
	for _, m := range messages {
		m1 := matches(g1, m)
		m2 := matches(g2, m)
		switch {
		case m1 && m2:
			both++
		case m2:
			twoOnly++
		default:
			neither++
		}
	}
	total := len(messages)
	p1total := both
	p2total := both + twoOnly

	const (
		W    = 820
		H    = 360
		barX = 60
		barY = 90
		barW = 700
		barH = 64
	)

	type seg struct {
		label string
		count int
		col   string
	}
	segs := []seg{
		{"matched by both rule sets", both, "#F0E442"},          // yellow (bright)
		{"matched only after part 2 rewrite", twoOnly, "#E69F00"}, // orange
		{"matched by neither", neither, "#0072B2"},               // blue (dark)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="40" fill="#e8ecf4" font-size="16">Monster Messages: %d messages under both rule sets</text>`, barX, total)

	x := barX
	for _, s := range segs {
		w := barW * s.count / total
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x, barY, w, barH, s.col)
		if w > 30 {
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="13" font-weight="bold" text-anchor="middle">%d</text>`, x+w/2, barY+barH/2+5, s.count)
		}
		x += w
	}

	ly := barY + barH + 44
	for i, s := range segs {
		yy := ly + i*30
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="16" height="16" fill="%s"/>`, barX, yy-13, s.col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="13">%s</text>`, barX+26, yy, s.label)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="13" text-anchor="end">%d</text>`, barX+400, yy, s.count)
	}

	sy := ly + 3*30 + 22
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#F0E442" font-size="14">Part 1 matches: %d</text>`, barX, sy, p1total)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#E69F00" font-size="14">Part 2 matches: %d  (superset of part 1)</text>`, barX, sy+24, p2total)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "monster-messages.svg"), []byte(sb.String()), 0o600)
}
