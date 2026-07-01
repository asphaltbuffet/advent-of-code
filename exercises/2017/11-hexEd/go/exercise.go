package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 11.
type Exercise struct {
	common.BaseExercise
}

// hexDeltas maps each hex direction to its cube-coordinate step. Cube
// coordinates keep x+y+z == 0, so distance is a single formula.
var hexDeltas = map[string][3]int{
	"n":  {0, 1, -1},
	"s":  {0, -1, 1},
	"ne": {1, 0, -1},
	"sw": {-1, 0, 1},
	"nw": {-1, 1, 0},
	"se": {1, -1, 0},
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// hexDist returns the number of steps from the origin to a cube coordinate.
func hexDist(x, y, z int) int {
	return (abs(x) + abs(y) + abs(z)) / 2
}

// walk follows the path and returns the final distance from origin and the
// greatest distance reached at any point along the way.
func walk(instr string) (final, furthest int) {
	x, y, z := 0, 0, 0
	for _, step := range strings.Split(strings.TrimSpace(instr), ",") {
		d := hexDeltas[strings.TrimSpace(step)]
		x, y, z = x+d[0], y+d[1], z+d[2]
		if dist := hexDist(x, y, z); dist > furthest {
			furthest = dist
		}
	}
	return hexDist(x, y, z), furthest
}

// One returns the fewest steps back to the origin from the path's end.
func (e Exercise) One(instr string) (any, error) {
	final, _ := walk(instr)
	return final, nil
}

// Two returns the furthest distance from the origin reached during the walk.
func (e Exercise) Two(instr string) (any, error) {
	_, furthest := walk(instr)
	return furthest, nil
}

// hexToPixel converts cube coordinates to pointy-top hex pixel coordinates.
func hexToPixel(x, z int, size float64) (float64, float64) {
	px := size * (math.Sqrt(3)*float64(x) + math.Sqrt(3)/2*float64(z))
	py := size * (1.5 * float64(z))
	return px, py
}

// Vis draws the full walk across the hex plane as a polyline (SVG). Each segment
// is coloured by its distance from the origin — cool near the centre, warm far
// out — so the wandering trajectory and its reach are visible. The origin, the
// furthest point, and the endpoint are marked.
func (e Exercise) Vis(instr, outdir string) error {
	const size = 1.0

	// Walk the path, collecting pixel points and per-vertex distances.
	type vtx struct {
		px, py float64
		dist   int
	}
	x, y, z := 0, 0, 0
	px, py := hexToPixel(x, z, size)
	verts := []vtx{{px, py, 0}}

	furthest, furthestIdx := 0, 0
	for _, step := range strings.Split(strings.TrimSpace(instr), ",") {
		d := hexDeltas[strings.TrimSpace(step)]
		x, y, z = x+d[0], y+d[1], z+d[2]
		px, py := hexToPixel(x, z, size)
		dist := hexDist(x, y, z)
		verts = append(verts, vtx{px, py, dist})
		if dist > furthest {
			furthest, furthestIdx = dist, len(verts)-1
		}
	}

	// Bounding box.
	minX, maxX := verts[0].px, verts[0].px
	minY, maxY := verts[0].py, verts[0].py
	for _, v := range verts {
		minX = math.Min(minX, v.px)
		maxX = math.Max(maxX, v.px)
		minY = math.Min(minY, v.py)
		maxY = math.Max(maxY, v.py)
	}

	// Scale to a comfortable canvas.
	const target = 1400.0
	const pad = 30.0
	const labelPad = 130.0 // extra right/bottom room for marker labels
	span := math.Max(maxX-minX, maxY-minY)
	scale := target / span
	W := (maxX-minX)*scale + 2*pad + labelPad
	H := (maxY-minY)*scale + 2*pad
	tx := func(x float64) float64 { return pad + (x-minX)*scale }
	ty := func(y float64) float64 { return pad + (y-minY)*scale }

	// heat maps a 0..1 fraction onto teal -> gold -> red.
	heat := func(t float64) string {
		var r, g, b int
		if t < 0.5 {
			u := t / 0.5
			r = int(0x2f + u*(0xff-0x2f))
			g = int(0x8a + u*(0xc8-0x8a))
			b = int(0x86 + u*(0x4a-0x86))
		} else {
			u := (t - 0.5) / 0.5
			r = int(0xff)
			g = int(0xc8 - u*(0xc8-0x44))
			b = int(0x4a - u*(0x4a-0x55))
		}
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f">`+"\n", W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0d0f18"/>`+"\n", W, H)

	// Path segments, coloured by distance.
	fmt.Fprint(&b, `<g stroke-width="1.4" stroke-linecap="round" fill="none" stroke-opacity="0.85">`+"\n")
	for i := 1; i < len(verts); i++ {
		a, c := verts[i-1], verts[i]
		t := float64(c.dist) / float64(furthest)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s"/>`+"\n",
			tx(a.px), ty(a.py), tx(c.px), ty(c.py), heat(t))
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Markers: origin (white), furthest (red), endpoint (green).
	o := verts[0]
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#ffffff"/>`+"\n", tx(o.px), ty(o.py))
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="monospace" font-size="13" fill="#ffffff">start</text>`+"\n", tx(o.px)+8, ty(o.py)-6)

	fv := verts[furthestIdx]
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#ff4455"/>`+"\n", tx(fv.px), ty(fv.py))
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="monospace" font-size="13" fill="#ff8090">furthest %d</text>`+"\n", tx(fv.px)+8, ty(fv.py)-6, furthest)

	ev := verts[len(verts)-1]
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#44ff88"/>`+"\n", tx(ev.px), ty(ev.py))
	fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="monospace" font-size="13" fill="#8effb0">end %d</text>`+"\n", tx(ev.px)+8, ty(ev.py)-6, ev.dist)

	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "hex-ed.svg"), []byte(b.String()), 0o644)
}
