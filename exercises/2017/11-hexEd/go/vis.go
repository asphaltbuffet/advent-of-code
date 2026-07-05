package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func hexToPixel(x, z int, size float64) (float64, float64) {
	px := size * (math.Sqrt(3)*float64(x) + math.Sqrt(3)/2*float64(z))
	py := size * (1.5 * float64(z))
	return px, py
}

// Vis draws the full walk across the hex plane as a polyline (SVG). Each segment
// is coloured by distance from the origin (teal → gold → red), and the start,
// furthest-ever, and end points are marked distinctly.
func (e Exercise) Vis(instr, outdir string) error {
	const size = 1.0

	x, y, z := 0, 0, 0
	px, py := hexToPixel(x, z, size)
	verts := []hexVtx{{px, py, 0}}

	furthest, furthestIdx := 0, 0
	for step := range strings.SplitSeq(strings.TrimSpace(instr), ",") {
		d := hexDeltas[strings.TrimSpace(step)]
		x, y, z = x+d[0], y+d[1], z+d[2]
		px, py = hexToPixel(x, z, size)
		dist := hexDist(x, y, z)
		verts = append(verts, hexVtx{px, py, dist})
		if dist > furthest {
			furthest, furthestIdx = dist, len(verts)-1
		}
	}

	minX, maxX := verts[0].px, verts[0].px
	minY, maxY := verts[0].py, verts[0].py
	for _, v := range verts {
		minX = math.Min(minX, v.px)
		maxX = math.Max(maxX, v.px)
		minY = math.Min(minY, v.py)
		maxY = math.Max(maxY, v.py)
	}

	svg := buildHexSVG(verts, furthest, furthestIdx, minX, maxX, minY, maxY)
	return os.WriteFile(filepath.Join(outdir, "hex-ed.svg"), []byte(svg), 0o644)
}

type hexVtx struct {
	px, py float64
	dist   int
}

func buildHexSVG(verts []hexVtx, furthest, furthestIdx int, minX, maxX, minY, maxY float64) string {
	const target = 1400.0
	const pad = 30.0
	const labelPad = 130.0
	span := math.Max(maxX-minX, maxY-minY)
	scale := target / span
	w := (maxX-minX)*scale + 2*pad + labelPad
	h := (maxY-minY)*scale + 2*pad
	tx := func(x float64) float64 { return pad + (x-minX)*scale }
	ty := func(y float64) float64 { return pad + (y-minY)*scale }

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

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f">`+"\n", w, h)
	fmt.Fprintf(&sb, `<rect width="%.0f" height="%.0f" fill="#0d0f18"/>`+"\n", w, h)

	fmt.Fprint(&sb, `<g stroke-width="1.4" stroke-linecap="round" fill="none" stroke-opacity="0.85">`+"\n")
	for i := 1; i < len(verts); i++ {
		a, c := verts[i-1], verts[i]
		t := float64(c.dist) / float64(furthest)
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s"/>`+"\n",
			tx(a.px), ty(a.py), tx(c.px), ty(c.py), heat(t))
	}
	fmt.Fprint(&sb, `</g>`+"\n")

	writeHexMarker(&sb, tx(verts[0].px), ty(verts[0].py), "#ffffff", "#ffffff", "start")
	fv := verts[furthestIdx]
	writeHexMarker(&sb, tx(fv.px), ty(fv.py), "#ff4455", "#ff8090", fmt.Sprintf("furthest %d", furthest))
	ev := verts[len(verts)-1]
	writeHexMarker(&sb, tx(ev.px), ty(ev.py), "#44ff88", "#8effb0", fmt.Sprintf("end %d", ev.dist))

	fmt.Fprint(&sb, `</svg>`+"\n")
	return sb.String()
}

func writeHexMarker(b *strings.Builder, x, y float64, dotColor, textColor, label string) {
	fmt.Fprintf(b, `<circle cx="%.1f" cy="%.1f" r="5" fill="%s"/>`+"\n", x, y, dotColor)
	fmt.Fprintf(b,
		`<text x="%.1f" y="%.1f" font-family="monospace" font-size="13" fill="%s">%s</text>`+"\n",
		x+8, y-6, textColor, label,
	)
}
