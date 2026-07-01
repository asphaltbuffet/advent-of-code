package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis draws the cave network as an SVG graph laid out with a small force-directed
// simulation. Big caves (uppercase, revisitable) are drawn as rounded squares;
// small caves (lowercase, visit-limited) as circles — so the type reads by shape,
// not just color. start and end are highlighted. This is the graph the two path
// counts explore; the shape distinction is exactly the rule that separates the
// parts.
func (e Exercise) Vis(instr, outdir string) error {
	adj, err := parse(instr)
	if err != nil {
		return err
	}

	// Stable node ordering.
	nodes := make([]string, 0, len(adj))
	for n := range adj {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	idx := map[string]int{}
	for i, n := range nodes {
		idx[n] = i
	}

	// Undirected edge set.
	type edge struct{ a, b int }
	edgeSet := map[edge]bool{}
	for a, nbrs := range adj {
		for _, b := range nbrs {
			i, j := idx[a], idx[b]
			if i > j {
				i, j = j, i
			}
			edgeSet[edge{i, j}] = true
		}
	}
	var edges []edge
	for e := range edgeSet {
		edges = append(edges, e)
	}

	// Initial positions on a circle.
	n := len(nodes)
	px := make([]float64, n)
	py := make([]float64, n)
	for i := range nodes {
		ang := 2 * math.Pi * float64(i) / float64(n)
		px[i] = 300 + 220*math.Cos(ang)
		py[i] = 300 + 220*math.Sin(ang)
	}

	// Force-directed relaxation: repulsion between all nodes, springs on edges.
	const iters = 400
	const k = 150.0 // ideal edge length
	for it := 0; it < iters; it++ {
		fx := make([]float64, n)
		fy := make([]float64, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				dx := px[i] - px[j]
				dy := py[i] - py[j]
				d := math.Hypot(dx, dy) + 0.01
				rep := k * k / (d * d)
				fx[i] += dx / d * rep
				fy[i] += dy / d * rep
			}
		}
		for _, e := range edges {
			dx := px[e.b] - px[e.a]
			dy := py[e.b] - py[e.a]
			d := math.Hypot(dx, dy) + 0.01
			att := (d - k) / d * 0.1
			fx[e.a] += dx * att
			fy[e.a] += dy * att
			fx[e.b] -= dx * att
			fy[e.b] -= dy * att
		}
		damp := 0.85
		for i := 0; i < n; i++ {
			px[i] += clampF(fx[i], -20, 20) * damp
			py[i] += clampF(fy[i], -20, 20) * damp
		}
	}

	// Normalize into the canvas.
	minX, minY := px[0], py[0]
	maxX, maxY := px[0], py[0]
	for i := 0; i < n; i++ {
		minX = math.Min(minX, px[i])
		maxX = math.Max(maxX, px[i])
		minY = math.Min(minY, py[i])
		maxY = math.Max(maxY, py[i])
	}
	const pad = 60
	const size = 620
	sx := (size - 2*pad) / (maxX - minX + 1)
	sy := (size - 2*pad) / (maxY - minY + 1)
	scale := math.Min(sx, sy)
	tx := func(x float64) float64 { return pad + (x-minX)*scale }
	ty := func(y float64) float64 { return pad + (y-minY)*scale }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, size, size+40, size, size+40)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, size, size+40)
	fmt.Fprintf(&sb, `<text x="%d" y="24" fill="#e8ecf4" font-size="13" text-anchor="middle">Cave network: squares = big (revisitable), circles = small (visit-limited)</text>`, size/2)

	// Edges.
	for _, e := range edges {
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#3a4250" stroke-width="1.5"/>`,
			tx(px[e.a]), ty(py[e.a]), tx(px[e.b]), ty(py[e.b]))
	}

	// Nodes.
	for i, name := range nodes {
		x, y := tx(px[i]), ty(py[i])
		fill, stroke := "#56B4E9", "#0072B2" // small cave: sky blue
		switch {
		case name == "start":
			fill, stroke = "#009E73", "#006b4f" // green
		case name == "end":
			fill, stroke = "#D55E00", "#933f00" // vermilion
		case !isSmall(name):
			fill, stroke = "#E69F00", "#9c6c00" // big cave: amber
		}

		if isSmall(name) && name != "start" && name != "end" {
			fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="20" fill="%s" stroke="%s" stroke-width="2"/>`, x, y, fill, stroke)
		} else {
			// start/end and big caves as rounded squares.
			fmt.Fprintf(&sb, `<rect x="%.1f" y="%.1f" width="40" height="40" rx="8" fill="%s" stroke="%s" stroke-width="2"/>`, x-20, y-20, fill, stroke)
		}
		fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#0d0f18" font-size="13" font-weight="bold" text-anchor="middle" dominant-baseline="middle">%s</text>`, x, y, name)
	}

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="12" text-anchor="middle">start (green) · end (vermilion) · big caves amber squares · small caves blue circles</text>`, size/2, size+24)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "passage-pathing.svg"), []byte(sb.String()), 0o600)
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
