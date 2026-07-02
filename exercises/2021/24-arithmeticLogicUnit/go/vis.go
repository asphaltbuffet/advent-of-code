package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the MONAD block structure as a pairing diagram (SVG). Each of the
// 14 blocks is a column; push blocks (div z 1) and pop blocks (div z 26) are
// marked distinctly. An arc joins every pop to the push it cancels, labeled with
// the constraint w_pop = w_push + delta that links the two digits. The nested,
// non-crossing arcs are the stack structure that collapses 14 free digits down to
// 7 constraints, so both extreme model numbers are determined without any search.
// Below each column the largest (part one) and smallest (part two) digits are
// shown. Push vs pop reads by shape and label, not color alone.
func (e Exercise) Vis(instr, outdir string) error {
	blocks, err := parseBlocks(instr)
	if err != nil {
		return err
	}
	pairs := constraints(blocks)
	maxModel := bestModel(pairs, true)
	minModel := bestModel(pairs, false)

	const (
		n     = 14
		colW  = 78
		leftX = 60
		baseY = 420
		topY  = 250
		boxW  = 46
		boxH  = 34
	)
	W := leftX*2 + n*colW
	H := 520

	// Okabe-Ito: sky blue for push, vermilion for pop.
	pushCol := "#56B4E9"
	popCol := "#D55E00"

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="16">MONAD block pairing: 14 digits, 7 constraints</text>`, leftX)

	cx := func(i int) int { return leftX + i*colW + colW/2 }

	// Arcs (draw first so boxes sit on top). Nesting depth sets arc height.
	// Compute a nesting level per pair for pleasant separation.
	for _, p := range pairs {
		x1, x2 := cx(p.i), cx(p.j)
		span := p.j - p.i
		lift := 40 + span*14
		midX := (x1 + x2) / 2
		ctrlY := topY - lift
		fmt.Fprintf(&sb, `<path d="M %d %d Q %d %d %d %d" fill="none" stroke="#8a94a2" stroke-width="1.5"/>`,
			x1, topY, midX, ctrlY, x2, topY)
		sign := "+"
		d := p.delta
		if d < 0 {
			sign, d = "-", -d
		}
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="11" text-anchor="middle">w%d = w%d %s %d</text>`,
			midX, ctrlY-4, p.j+1, p.i+1, sign, d)
	}

	// Column boxes and labels.
	for i, b := range blocks {
		x := leftX + i*colW + (colW-boxW)/2
		col := pushCol
		kind := "push"
		if b.div == 26 {
			col = popCol
			kind = "pop"
		}
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" rx="4" fill="%s"/>`, x, topY, boxW, boxH, col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="12" font-weight="bold" text-anchor="middle">%d</text>`, cx(i), topY+15, i+1)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="10" text-anchor="middle">%s</text>`, cx(i), topY+28, kind)

		// digit rows
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="18" font-weight="bold" text-anchor="middle">%c</text>`,
			cx(i), baseY, pushCol, maxModel[i])
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="18" font-weight="bold" text-anchor="middle">%c</text>`,
			cx(i), baseY+34, popCol, minModel[i])
	}

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="12" text-anchor="start">largest</text>`, 6, baseY, pushCol)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="12" text-anchor="start">smallest</text>`, 6, baseY+34, popCol)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "monad-pairing.svg"), []byte(sb.String()), 0o600)
}
