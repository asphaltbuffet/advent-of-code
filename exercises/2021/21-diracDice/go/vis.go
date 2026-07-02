package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis shows what drives the quantum game (SVG), in two panels. The left panel is
// the roll-sum distribution: three 3-sided dice give 27 universes that collapse
// to seven sums (3..9) with multiplicities 1,3,6,7,6,3,1 — the branching weights
// the whole search. The right panel compares how many universes each player wins
// (part two), with the winner highlighted. Bars are labeled with their values, so
// the chart reads without relying on color.
func (e Exercise) Vis(instr, outdir string) error {
	p1, p2, err := parse(instr)
	if err != nil {
		return err
	}
	wins := countWins(dstate{p1, p2, 0, 0}, map[dstate][2]int64{})

	const (
		W      = 960
		H      = 460
		gap    = 70
		mLeft  = 60
		mTop   = 56
		mBot   = 60
	)
	panelW := (W - mLeft - gap - 40) / 2
	plotH := H - mTop - mBot

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)

	// ---- Left: roll-sum multiplicities. ----
	fmt.Fprintf(&sb, `<text x="%d" y="30" fill="#e8ecf4" font-size="14">Three-roll sums (27 → 7 universes)</text>`, mLeft)
	maxMult := 7
	barW := panelW / len(diracRolls)
	for i, r := range diracRolls {
		sum, mult := r[0], r[1]
		bh := plotH * mult / maxMult
		x := mLeft + i*barW
		y := mTop + plotH - bh
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="#56B4E9"/>`, x+4, y, barW-8, bh)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="12" font-weight="bold" text-anchor="middle">%d</text>`, x+barW/2, y+14, mult)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">%d</text>`, x+barW/2, mTop+plotH+18, sum)
	}
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">roll sum</text>`, mLeft+panelW/2, mTop+plotH+40)

	// ---- Right: universe wins per player. ----
	ox := mLeft + panelW + gap
	fmt.Fprintf(&sb, `<text x="%d" y="30" fill="#e8ecf4" font-size="14">Universes won (part two)</text>`, ox)
	maxWin := wins[0]
	if wins[1] > maxWin {
		maxWin = wins[1]
	}
	labels := []string{"Player 1", "Player 2"}
	cols := []string{"#E69F00", "#009E73"}
	winner := 0
	if wins[1] > wins[0] {
		winner = 1
	}
	bw := panelW / 3
	for i := 0; i < 2; i++ {
		bh := int(int64(plotH) * wins[i] / maxWin)
		x := ox + bw/2 + i*bw
		y := mTop + plotH - bh
		fill := cols[i]
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x, y, bw-20, bh, fill)
		if i == winner {
			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#e8ecf4" stroke-width="2.5"/>`, x, y, bw-20, bh)
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="12" text-anchor="middle">winner</text>`, x+(bw-20)/2, y-8)
		}
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">%s</text>`, x+(bw-20)/2, mTop+plotH+18, labels[i])
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="11" text-anchor="middle">%s</text>`, x+(bw-20)/2, mTop+plotH+36, commafy(wins[i]))
	}

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "dirac-dice.svg"), []byte(sb.String()), 0o600)
}

// commafy formats n with thousands separators.
func commafy(n int64) string {
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, d := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, d)
	}
	return string(out)
}
