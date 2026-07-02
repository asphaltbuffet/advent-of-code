package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis shows the Part Two constraint matrix: rows are fields, columns are ticket
// positions, and a cell is filled when that field's ranges accept every valid
// ticket's value in that column (a candidate). Constraint propagation reduces this
// dense matrix to exactly one field per column; those final assignments are
// outlined, and the six "departure" fields whose product is the answer are
// highlighted. Candidates, the solved assignment, and departure fields use
// distinct colorblind-safe colors reinforced by outline and position, so the grid
// reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	n, err := parse(instr)
	if err != nil {
		return err
	}
	cols := len(n.mine)
	nf := len(n.fields)

	// Valid tickets.
	valid := [][]int{n.mine}
	for _, t := range n.tickets {
		ok := true
		for _, v := range t {
			if !validForAny(n.fields, v) {
				ok = false
				break
			}
		}
		if ok {
			valid = append(valid, t)
		}
	}

	// candidate[fi][c] = field fi can occupy column c.
	candidate := make([][]bool, nf)
	for fi, f := range n.fields {
		candidate[fi] = make([]bool, cols)
		for c := 0; c < cols; c++ {
			fits := true
			for _, t := range valid {
				if !f.valid(t[c]) {
					fits = false
					break
				}
			}
			candidate[fi][c] = fits
		}
	}

	// Final assignment (column -> field index).
	assignedName := assignFields(n)
	colField := make([]int, cols)
	for c := 0; c < cols; c++ {
		for fi, f := range n.fields {
			if f.name == assignedName[c] {
				colField[c] = fi
			}
		}
	}

	const (
		cell   = 22
		labelW = 170
		mT     = 60
		mL     = 10
	)
	gridW := cols * cell
	W := mL + labelW + gridW + 20
	H := mT + nf*cell + 20

	candCol := "#56B4E9"  // sky blue: a candidate
	solveCol := "#F0E442" // yellow: the solved assignment
	depCol := "#D55E00"   // vermilion outline: a departure field row

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="26" fill="#e8ecf4" font-size="15">Ticket Translation: field candidates per column (part 2)</text>`, mL)
	fmt.Fprintf(&sb, `<text x="%d" y="44" fill="#9aa4b2" font-size="12">filled = possible; outlined = solved; vermilion label = departure field (answer)</text>`, mL)

	gx := mL + labelW
	// Column index headers.
	for c := 0; c < cols; c++ {
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10" text-anchor="middle">%d</text>`, gx+c*cell+cell/2, mT-4, c)
	}

	for fi, f := range n.fields {
		y := mT + fi*cell
		isDep := strings.HasPrefix(f.name, "departure")
		labelCol := "#c8d0dc"
		if isDep {
			labelCol = depCol
		}
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="11" text-anchor="end">%s</text>`, mL+labelW-6, y+cell/2+4, labelCol, f.name)
		for c := 0; c < cols; c++ {
			cx, cy := gx+c*cell, y
			if candidate[fi][c] {
				col := candCol
				if colField[c] == fi {
					col = solveCol
				}
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, cx+1, cy+1, cell-2, cell-2, col)
			} else {
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="#1b1f26"/>`, cx+1, cy+1, cell-2, cell-2)
			}
			// Outline the solved cell.
			if colField[c] == fi {
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#e8ecf4" stroke-width="1.5"/>`, cx+1, cy+1, cell-2, cell-2)
			}
		}
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "ticket-translation.svg"), []byte(sb.String()), 0o600)
}
