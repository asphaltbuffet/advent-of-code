package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// path returns the sequence of ship positions for a navigation mode. When
// waypoint is false the ship steers by its own heading (part one); when true it
// steers by a waypoint (part two).
func path(moves []move, waypoint bool) [][2]int {
	x, y := 0, 0
	dx, dy := 1, 0 // heading (part one) or waypoint offset (part two)
	if waypoint {
		dx, dy = 10, 1
	}
	pts := [][2]int{{0, 0}}
	for _, m := range moves {
		switch m.action {
		case 'N':
			if waypoint {
				dy += m.value
			} else {
				y += m.value
			}
		case 'S':
			if waypoint {
				dy -= m.value
			} else {
				y -= m.value
			}
		case 'E':
			if waypoint {
				dx += m.value
			} else {
				x += m.value
			}
		case 'W':
			if waypoint {
				dx -= m.value
			} else {
				x -= m.value
			}
		case 'F':
			x += dx * m.value
			y += dy * m.value
		case 'L':
			dx, dy = rotate(dx, dy, m.value)
		case 'R':
			dx, dy = rotate(dx, dy, 360-m.value)
		}
		pts = append(pts, [2]int{x, y})
	}
	return pts
}

// Vis plots the two ship trajectories in separate, independently scaled panels
// (the waypoint route travels ~40x farther than the direct route, so a shared
// scale would hide the direct one). Each panel shows the route from the origin
// with a start dot and an end marker, titled with its Manhattan distance answer.
// The panels use colorblind-safe colors and distinct end-marker shapes, so they
// read in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	moves, err := parse(instr)
	if err != nil {
		return err
	}
	p1 := path(moves, false)
	p2 := path(moves, true)

	const (
		W       = 940
		H       = 480
		panelW  = 430
		panelH  = 380
		pad     = 40
		gap     = 40
		panelTY = 70
	)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#e8ecf4" font-size="16">Rain Risk: ship routes from the origin</text>`, pad)

	drawPanel := func(px int, pts [][2]int, col, title string, square bool) {
		// Bounds.
		minX, maxX, minY, maxY := 0, 0, 0, 0
		for _, p := range pts {
			minX, maxX = min(minX, p[0]), max(maxX, p[0])
			minY, maxY = min(minY, p[1]), max(maxY, p[1])
		}
		spanX := max(maxX-minX, 1)
		spanY := max(maxY-minY, 1)
		inX, inY := panelW-30, panelH-30
		sc := min(inX*1000/spanX, inY*1000/spanY)
		// Center the drawing in the panel.
		drawnW := spanX * sc / 1000
		drawnH := spanY * sc / 1000
		ox := px + 15 + (inX-drawnW)/2 - minX*sc/1000
		oy := panelTY + 15 + (inY-drawnH)/2 + maxY*sc/1000
		tx := func(x int) int { return ox + x*sc/1000 }
		ty := func(y int) int { return oy - y*sc/1000 }

		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#2a3038"/>`, px, panelTY, panelW, panelH)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="14">%s</text>`, px, panelTY-10, col, title)

		var s []string
		for _, p := range pts {
			s = append(s, fmt.Sprintf("%d,%d", tx(p[0]), ty(p[1])))
		}
		fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="1.5"/>`, strings.Join(s, " "), col)

		fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="4" fill="#e8ecf4"/>`, tx(0), ty(0))
		end := pts[len(pts)-1]
		if square {
			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="11" height="11" fill="%s"/>`, tx(end[0])-5, ty(end[1])-5, col)
		} else {
			fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="6" fill="%s"/>`, tx(end[0]), ty(end[1]), col)
		}
	}

	e1 := p1[len(p1)-1]
	e2 := p2[len(p2)-1]
	drawPanel(pad, p1, "#E69F00", fmt.Sprintf("direct steering (part 1): %d", abs(e1[0])+abs(e1[1])), false)
	drawPanel(pad+panelW+gap, p2, "#0072B2", fmt.Sprintf("waypoint steering (part 2): %d", abs(e2[0])+abs(e2[1])), true)

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11">white dot = start; panels are scaled independently</text>`, pad, H-16)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "rain-risk.svg"), []byte(sb.String()), 0o600)
}
