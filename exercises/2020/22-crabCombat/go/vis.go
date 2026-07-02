package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis plots the two players' deck sizes over the standard Combat game (part one),
// a tug-of-war that always sums to the full 50 cards until one deck empties and
// the game ends. The two series use colorblind-safe colors and distinct markers,
// and the winner is annotated, so the momentum swings read in grayscale as well
// as color.
func (e Exercise) Vis(instr, outdir string) error {
	p1, p2, err := parseDecks(instr)
	if err != nil {
		return err
	}

	// Replay, recording each deck size per round.
	var s1, s2 []int
	a := append([]int(nil), p1...)
	b := append([]int(nil), p2...)
	s1 = append(s1, len(a))
	s2 = append(s2, len(b))
	for len(a) > 0 && len(b) > 0 {
		x, y := a[0], b[0]
		a, b = a[1:], b[1:]
		if x > y {
			a = append(a, x, y)
		} else {
			b = append(b, y, x)
		}
		s1 = append(s1, len(a))
		s2 = append(s2, len(b))
	}
	rounds := len(s1) - 1
	total := len(p1) + len(p2)

	const (
		W  = 960
		H  = 420
		mL = 55
		mR = 130
		mT = 50
		mB = 55
	)
	plotW := W - mL - mR
	plotH := H - mT - mB
	xOf := func(r int) int { return mL + r*plotW/rounds }
	yOf := func(v int) int { return mT + plotH - v*plotH/total }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="32" fill="#e8ecf4" font-size="16">Crab Combat: deck sizes over %d rounds (part one)</text>`, mL, rounds)

	// Axes.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT+plotH, W-mR, mT+plotH)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, mL, mT, mL, mT+plotH)
	for _, v := range []int{0, total / 2, total} {
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11" text-anchor="end">%d</text>`, mL-6, yOf(v)+4, v)
	}
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#2a3038" stroke-dasharray="4 4"/>`, mL, yOf(total/2), W-mR, yOf(total/2))
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">round</text>`, mL+plotW/2, mT+plotH+40)

	line := func(sizes []int, col, dash string) {
		var pts []string
		for r, v := range sizes {
			pts = append(pts, fmt.Sprintf("%d,%d", xOf(r), yOf(v)))
		}
		fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2" stroke-dasharray="%s"/>`, strings.Join(pts, " "), col, dash)
	}
	p1Col := "#E69F00" // orange
	p2Col := "#0072B2" // blue
	line(s1, p1Col, "none")
	line(s2, p2Col, "6 4") // dashed so the mirror lines differ in grayscale

	// Legend + winner.
	winner, wcol := "Player 1", p1Col
	if s2[rounds] > s1[rounds] {
		winner, wcol = "Player 2", p2Col
	}
	lx := W - mR + 10
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">Player 1</text>`, lx, mT, p1Col, lx+16, mT+10)
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="12" height="12" fill="%s"/><text x="%d" y="%d" fill="#e8ecf4" font-size="12">Player 2</text>`, lx, mT+20, p2Col, lx+16, mT+30)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="13">%s wins</text>`, lx, mT+60, wcol, winner)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "crab-combat.svg"), []byte(sb.String()), 0o600)
}
