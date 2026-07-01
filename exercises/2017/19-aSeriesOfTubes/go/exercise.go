package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 19.
type Exercise struct {
	common.BaseExercise
}

type cell struct {
	r, c     int
	isLetter bool
	ch       byte
}

// walk follows the routing diagram from the top-row entry and returns the
// letters collected in order, the total number of steps taken, and the sequence
// of cells visited.
func walk(instr string) (string, int, []cell) {
	// Preserve leading spaces — they position the path. Only trim the trailing
	// newline so the final row isn't empty.
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")

	at := func(r, c int) byte {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[r]) {
			return ' '
		}
		return grid[r][c]
	}

	// Entry: the only '|' in the top row.
	r, c := 0, strings.IndexByte(grid[0], '|')
	dr, dc := 1, 0 // moving down

	var letters strings.Builder
	var path []cell
	steps := 0
	for at(r, c) != ' ' {
		ch := at(r, c)
		isLetter := ch >= 'A' && ch <= 'Z'
		if isLetter {
			letters.WriteByte(ch)
		}
		path = append(path, cell{r, c, isLetter, ch})
		if ch == '+' {
			// Turn to whichever perpendicular direction continues.
			if dr != 0 { // moving vertically -> turn horizontally
				if at(r, c-1) != ' ' {
					dr, dc = 0, -1
				} else {
					dr, dc = 0, 1
				}
			} else { // moving horizontally -> turn vertically
				if at(r-1, c) != ' ' {
					dr, dc = -1, 0
				} else {
					dr, dc = 1, 0
				}
			}
		}
		r, c = r+dr, c+dc
		steps++
	}
	return letters.String(), steps, path
}

// One returns the letters collected along the path.
func (e Exercise) One(instr string) (any, error) {
	letters, _, _ := walk(instr)
	return letters, nil
}

// Two returns the total number of steps taken along the path.
func (e Exercise) Two(instr string) (any, error) {
	_, steps, _ := walk(instr)
	return steps, nil
}

// Vis renders the routing diagram (faint) with the traversed path overlaid as a
// crisp progress-coloured polyline (SVG); each collected letter is drawn as a
// labelled glyph so the order GPALMJSOY is legible.
func (e Exercise) Vis(instr, outdir string) error {
	grid := strings.Split(strings.TrimRight(instr, "\n"), "\n")
	_, _, path := walk(instr)

	rows := len(grid)
	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	const step = 6.0
	const pad = 14.0
	W := float64(cols)*step + 2*pad
	H := float64(rows)*step + 2*pad
	px := func(c int) float64 { return pad + float64(c)*step + step/2 }
	py := func(r int) float64 { return pad + float64(r)*step + step/2 }

	// heat maps progress 0..1 onto teal -> gold -> magenta.
	heat := func(t float64) string {
		var r, g, b int
		if t < 0.5 {
			u := t / 0.5
			r = int(0x2f + u*(0xff-0x2f))
			g = int(0x8a + u*(0xc8-0x8a))
			b = int(0x86 + u*(0x4a-0x86))
		} else {
			u := (t - 0.5) / 0.5
			r = 0xff
			g = int(0xc8 - u*(0xc8-0x44))
			b = int(0x4a + u*(0xd0-0x4a))
		}
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" font-family="monospace">`+"\n", W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0d0f18"/>`+"\n", W, H)

	// Faint diagram underneath.
	fmt.Fprint(&b, `<g fill="#232a3e">`+"\n")
	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] != ' ' {
				fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f"/>`+"\n",
					pad+float64(c)*step, pad+float64(r)*step, step, step)
			}
		}
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Path as a progress-coloured polyline, drawn as short segments so the hue
	// shifts along its length.
	fmt.Fprint(&b, `<g stroke-width="2.4" stroke-linecap="round" fill="none">`+"\n")
	n := len(path)
	for i := 1; i < n; i++ {
		t := float64(i) / math.Max(1, float64(n-1))
		a, c := path[i-1], path[i]
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s"/>`+"\n",
			px(a.c), py(a.r), px(c.c), py(c.r), heat(t))
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Letters: a dark disc with the glyph on top, so each waypoint is readable.
	for _, p := range path {
		if !p.isLetter {
			continue
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="7" fill="#0d0f18" stroke="#ffffff" stroke-width="1.2"/>`+"\n", px(p.c), py(p.r))
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="9" font-weight="bold" fill="#ffffff" text-anchor="middle">%c</text>`+"\n",
			px(p.c), py(p.r)+3.2, p.ch)
	}

	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "a-series-of-tubes.svg"), []byte(b.String()), 0o644)
}
