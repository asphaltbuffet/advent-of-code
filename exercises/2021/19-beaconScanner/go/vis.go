package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis renders the reconstructed beacon map (SVG). Every scanner is located
// relative to scanner 0, and the merged point cloud is drawn with a simple
// isometric projection: small dots are beacons, larger ringed markers are the
// scanner positions. Points are depth-cued — those farther back are smaller and
// dimmer — so the 3D layout reads on a flat image. Scanners use a warm accent and
// beacons a cool one, but they also differ in size and ring, so the distinction
// survives grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	scanners := parse(instr)
	if len(scanners) == 0 {
		return nil
	}
	beaconSet, positions := assemble(scanners)

	beacons := make([]vec, 0, len(beaconSet))
	for b := range beaconSet {
		beacons = append(beacons, b)
	}

	// Isometric projection: rotate 30° so all three axes are visible.
	const angle = math.Pi / 6
	cos, sin := math.Cos(angle), math.Sin(angle)
	project := func(v vec) (float64, float64, float64) {
		x, y, z := float64(v.x), float64(v.y), float64(v.z)
		px := (x - z) * cos
		py := y + (x+z)*sin
		depth := x + z // larger = farther back for cueing
		return px, py, depth
	}

	// Compute projected bounds.
	minPx, maxPx := math.Inf(1), math.Inf(-1)
	minPy, maxPy := math.Inf(1), math.Inf(-1)
	minD, maxD := math.Inf(1), math.Inf(-1)
	all := append(append([]vec(nil), beacons...), positions...)
	for _, v := range all {
		px, py, d := project(v)
		minPx, maxPx = math.Min(minPx, px), math.Max(maxPx, px)
		minPy, maxPy = math.Min(minPy, py), math.Max(maxPy, py)
		minD, maxD = math.Min(minD, d), math.Max(maxD, d)
	}

	const (
		W   = 900
		H   = 900
		pad = 40
	)
	sw := (W - 2*pad) / (maxPx - minPx)
	sh := (H - 2*pad) / (maxPy - minPy)
	scale := math.Min(sw, sh)
	sx := func(px float64) float64 { return pad + (px-minPx)*scale }
	sy := func(py float64) float64 { return H - pad - (py-minPy)*scale }
	depthT := func(d float64) float64 {
		if maxD == minD {
			return 0.5
		}
		return (d - minD) / (maxD - minD) // 0 = near, 1 = far
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#0d0f18"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="28" fill="#e8ecf4" font-size="15" text-anchor="middle">Reconstructed map: %d beacons, %d scanners (isometric)</text>`, W/2, len(beacons), len(positions))

	// Depth-sort so nearer points draw on top.
	type dp struct {
		x, y, d float64
	}
	bpts := make([]dp, len(beacons))
	for i, b := range beacons {
		px, py, d := project(b)
		bpts[i] = dp{sx(px), sy(py), depthT(d)}
	}
	sort.Slice(bpts, func(i, j int) bool { return bpts[i].d > bpts[j].d }) // far first

	for _, p := range bpts {
		r := 1.5 + 2.0*(1-p.d)         // nearer = bigger
		op := 0.35 + 0.55*(1-p.d)      // nearer = brighter
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.2f" fill="#56B4E9" fill-opacity="%.2f"/>`, p.x, p.y, r, op)
	}

	// Scanners on top: larger amber ringed markers.
	for _, s := range positions {
		px, py, d := project(s)
		t := depthT(d)
		r := 4.0 + 3.0*(1-t)
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#E69F00" stroke="#0d0f18" stroke-width="1.5"/>`, sx(px), sy(py), r)
	}

	// Legend.
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="3.5" fill="#56B4E9"/><text x="%d" y="%d" fill="#9aa4b2" font-size="12">beacon</text>`, pad, H-20, pad+12, H-16)
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="5.5" fill="#E69F00" stroke="#0d0f18" stroke-width="1.5"/><text x="%d" y="%d" fill="#9aa4b2" font-size="12">scanner</text>`, pad+110, H-20, pad+124, H-16)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "beacon-scanner.svg"), []byte(sb.String()), 0o600)
}
