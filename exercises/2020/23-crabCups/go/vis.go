package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis illustrates the crab game's mechanics on the small (Part One) ring by
// drawing the first several moves as circular cup arrangements. In each snapshot
// the current cup is outlined, the three cups about to be picked up are one color,
// and the destination cup they will be placed after is another. Roles are marked
// by outline and label as well as color, so the move reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	cups := parseCups(instr)
	n := len(cups)

	// Build the next-array and replay, capturing a snapshot before each of the
	// first `snaps` moves.
	next := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		next[cups[i]] = cups[i+1]
	}
	next[cups[n-1]] = cups[0]
	current := cups[0]

	const snaps = 6
	type frame struct {
		order              []int
		current, dest      int
		a, b, c            int
	}
	var frames []frame

	for m := 0; m < snaps; m++ {
		// Ring order starting from current.
		order := make([]int, 0, n)
		for c, i := current, 0; i < n; c, i = next[c], i+1 {
			order = append(order, c)
		}
		a := next[current]
		b := next[a]
		cc := next[b]
		dest := current - 1
		if dest < 1 {
			dest = n
		}
		for dest == a || dest == b || dest == cc {
			dest--
			if dest < 1 {
				dest = n
			}
		}
		frames = append(frames, frame{order, current, dest, a, b, cc})

		// Advance one move.
		next[current] = next[cc]
		next[cc] = next[dest]
		next[dest] = a
		current = next[current]
	}

	const (
		perRow = 3
		panelW = 300
		panelH = 210
		mL     = 20
		mT     = 60
		radius = 66
		cupR   = 17
	)
	rows := (len(frames) + perRow - 1) / perRow
	W := mL*2 + perRow*panelW
	H := mT + rows*panelH + 10

	pickCol := "#E69F00" // orange: picked up
	destCol := "#F0E442" // bright yellow (dashed ring): destination
	plain := "#3a4450"

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="30" fill="#e8ecf4" font-size="16">Crab Cups: first %d moves (solid ring = current, orange = picked up, dashed ring = destination)</text>`, mL, snaps)

	for fi, f := range frames {
		col := fi % perRow
		row := fi / perRow
		cx := mL + col*panelW + panelW/2
		cy := mT + row*panelH + panelH/2 - 6
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">move %d</text>`, cx, cy-radius-24, fi+1)

		for i, cup := range f.order {
			ang := -math.Pi/2 + 2*math.Pi*float64(i)/float64(n)
			x := cx + int(float64(radius)*math.Cos(ang))
			y := cy + int(float64(radius)*math.Sin(ang))
			fill := plain
			isDest := cup == f.dest
			switch {
			case cup == f.a || cup == f.b || cup == f.c:
				fill = pickCol
			case isDest:
				fill = destCol
			}
			fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="%d" fill="%s"/>`, x, y, cupR, fill)
			// Destination gets a dashed ring so it differs from picked-up cups by
			// shape as well as color (survives grayscale).
			if isDest {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="%d" fill="none" stroke="#e8ecf4" stroke-width="2" stroke-dasharray="3 3"/>`, x, y, cupR+2)
			}
			if cup == f.current {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="%d" fill="none" stroke="#e8ecf4" stroke-width="3"/>`, x, y, cupR+2)
			}
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="14" font-weight="bold" text-anchor="middle">%d</text>`, x, y+5, cup)
		}
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "crab-cups.svg"), []byte(sb.String()), 0o600)
}
