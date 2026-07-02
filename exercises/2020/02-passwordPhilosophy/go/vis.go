package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis breaks the password list down by how the two policies judge each entry: the
// count policy (Part One) and the position policy (Part Two). The four groups —
// valid under both, Part One only, Part Two only, and neither — are shown as a
// stacked bar plus a small table, revealing that the two rules mostly disagree
// about which passwords pass even though their totals (666 vs 670) are close. Each
// group uses a colorblind-safe color that also differs in brightness, and every
// segment is labeled with its count, so the chart reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	entries, err := parse(instr)
	if err != nil {
		return err
	}

	var both, oneOnly, twoOnly, neither int
	for _, en := range entries {
		n := strings.Count(en.pass, string(en.ch))
		p1 := n >= en.lo && n <= en.hi
		a := en.lo-1 < len(en.pass) && en.pass[en.lo-1] == en.ch
		b := en.hi-1 < len(en.pass) && en.pass[en.hi-1] == en.ch
		p2 := a != b
		switch {
		case p1 && p2:
			both++
		case p1:
			oneOnly++
		case p2:
			twoOnly++
		default:
			neither++
		}
	}
	total := len(entries)
	p1total := both + oneOnly
	p2total := both + twoOnly

	const (
		W     = 820
		H     = 380
		barX  = 60
		barY  = 90
		barW  = 700
		barH  = 70
	)

	type seg struct {
		label string
		count int
		col   string
	}
	// Okabe-Ito, ordered by brightness so the stack also reads in grayscale.
	segs := []seg{
		{"valid under both", both, "#F0E442"},    // yellow (bright)
		{"Part 1 only (count)", oneOnly, "#E69F00"}, // orange
		{"Part 2 only (position)", twoOnly, "#009E73"}, // green
		{"neither", neither, "#0072B2"},           // blue (dark)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="40" fill="#e8ecf4" font-size="16">Password Philosophy: how the two policies judge %d passwords</text>`, barX, total)

	// Stacked bar.
	x := barX
	for _, s := range segs {
		w := barW * s.count / total
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x, barY, w, barH, s.col)
		if w > 34 {
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="13" font-weight="bold" text-anchor="middle">%d</text>`,
				x+w/2, barY+barH/2+5, s.count)
		}
		x += w
	}

	// Legend + table.
	ly := barY + barH + 46
	for i, s := range segs {
		yy := ly + i*30
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="16" height="16" fill="%s"/>`, barX, yy-13, s.col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="13">%s</text>`, barX+26, yy, s.label)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="13" text-anchor="end">%d</text>`, barX+360, yy, s.count)
	}

	sy := ly + 4*30 + 20
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#E69F00" font-size="14">Part 1 total (count policy):    %d</text>`, barX, sy, p1total)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#009E73" font-size="14">Part 2 total (position policy): %d</text>`, barX, sy+24, p2total)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "password-philosophy.svg"), []byte(sb.String()), 0o600)
}
