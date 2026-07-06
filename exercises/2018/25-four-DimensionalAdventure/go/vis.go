package exercises

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Okabe-Ito colorblind-safe palette for the largest constellations. Everything
// beyond the top few is drawn in neutral gray, so meaning never rests on telling
// two similar hues apart — the big clusters are the ones worth distinguishing.
var visPalette = []string{
	"#e69f00", // orange
	"#56b4e9", // sky blue
	"#009e73", // bluish green
	"#f0e442", // yellow
	"#0072b2", // blue
	"#d55e00", // vermillion
	"#cc79a7", // reddish purple
	"#ffffff", // white
}

const (
	visBG   = "#141418" // near-black background
	visGray = "#5a5a64" // small/singleton constellations
)

// Vis renders the constellation structure as a force-directed graph. Nodes are the
// 4D points; an edge joins any two points within Manhattan distance 3 (the same
// relation that defines a constellation). A Fruchterman-Reingold layout pulls each
// densely-linked constellation into its own blob, so the ~324 clusters read as
// distinct islands on the page — the answer made visible. The largest handful of
// constellations get distinct Okabe-Ito colors; the long tail of small ones is gray.
func (e Exercise) Vis(instr, outdir string) error {
	pts := parse(instr)
	n := len(pts)
	if n == 0 {
		return errors.New("no points to visualize")
	}

	edges, nodeColor := buildConstellationGraph(pts, n)
	xs, ys := layout(n, edges)

	const (
		size    = 1400.0
		margin  = 40.0
		radius  = 1.5
		edgeCol = "#2e2e38"
	)
	sb := renderConstellationSVG(xs, ys, edges, nodeColor, size, margin, radius, edgeCol)
	return os.WriteFile(filepath.Join(outdir, "four-dimensional-adventure.svg"), []byte(sb), 0o600)
}

// layout runs a Fruchterman-Reingold force simulation: every node repels every
// other (spreading separate constellations apart) while edges act as springs
// (pulling each constellation's members together). Positions are returned once the
// per-iteration displacement (temperature) has cooled to near zero.
// visEdge is a distance-3 link between two point indices.
type visEdge struct{ a, b int }

func layout(n int, edges []visEdge) ([]float64, []float64) {
	rng := rand.New(rand.NewPCG(2018, 0)) // fixed seed → reproducible SVG
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := range n {
		a := 2 * math.Pi * rng.Float64()
		r := math.Sqrt(rng.Float64())
		xs[i] = r * math.Cos(a)
		ys[i] = r * math.Sin(a)
	}

	const (
		k       = 0.032
		gravity = 0.02
		iters   = 400
	)
	temp := 0.15
	cool := temp / float64(iters+1)
	dx := make([]float64, n)
	dy := make([]float64, n)

	for range iters {
		frIteration(xs, ys, dx, dy, edges, n, k, gravity, temp)
		temp -= cool
	}
	return xs, ys
}

func frIteration(xs, ys, dx, dy []float64, edges []visEdge, n int, k, gravity, temp float64) {
	for i := range dx {
		dx[i], dy[i] = 0, 0
	}
	for i := range n {
		for j := i + 1; j < n; j++ {
			ddx, ddy := xs[i]-xs[j], ys[i]-ys[j]
			dist := math.Hypot(ddx, ddy)
			if dist < 1e-4 {
				dist = 1e-4
			}
			force := k * k / dist
			ux, uy := ddx/dist, ddy/dist
			dx[i] += ux * force
			dy[i] += uy * force
			dx[j] -= ux * force
			dy[j] -= uy * force
		}
	}
	for _, e := range edges {
		ddx, ddy := xs[e.a]-xs[e.b], ys[e.a]-ys[e.b]
		dist := math.Hypot(ddx, ddy)
		if dist < 1e-4 {
			dist = 1e-4
		}
		force := dist * dist / k
		ux, uy := ddx/dist, ddy/dist
		dx[e.a] -= ux * force
		dy[e.a] -= uy * force
		dx[e.b] += ux * force
		dy[e.b] += uy * force
	}
	for i := range n {
		dx[i] -= xs[i] * gravity
		dy[i] -= ys[i] * gravity
	}
	for i := range n {
		d := math.Hypot(dx[i], dy[i])
		if d < 1e-9 {
			continue
		}
		step := math.Min(d, temp)
		xs[i] += dx[i] / d * step
		ys[i] += dy[i] / d * step
	}
}

func buildConstellationGraph(pts []point, n int) ([]visEdge, []string) {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	var find func(x int) int //nolint:staticcheck // recursive closure must be pre-declared to call itself
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	var edges []visEdge
	for i := range pts {
		for j := i + 1; j < n; j++ {
			if manhattan(pts[i], pts[j]) <= 3 {
				edges = append(edges, visEdge{i, j})
				parent[find(i)] = find(j)
			}
		}
	}

	members := map[int][]int{}
	for i := range pts {
		r := find(i)
		members[r] = append(members[r], i)
	}
	roots := make([]int, 0, len(members))
	for r := range members {
		roots = append(roots, r)
	}
	sort.Slice(roots, func(a, b int) bool {
		if len(members[roots[a]]) != len(members[roots[b]]) {
			return len(members[roots[a]]) > len(members[roots[b]])
		}
		return roots[a] < roots[b]
	})

	nodeColor := make([]string, n)
	for rank, r := range roots {
		c := visGray
		if rank < len(visPalette) {
			c = visPalette[rank]
		}
		for _, i := range members[r] {
			nodeColor[i] = c
		}
	}
	return edges, nodeColor
}

func renderConstellationSVG(
	xs, ys []float64, edges []visEdge, nodeColor []string, size, margin, radius float64, edgeCol string,
) string {
	minX, maxX, minY, maxY := xs[0], xs[0], ys[0], ys[0]
	for i := range xs {
		minX, maxX = math.Min(minX, xs[i]), math.Max(maxX, xs[i])
		minY, maxY = math.Min(minY, ys[i]), math.Max(maxY, ys[i])
	}
	spanX, spanY := maxX-minX, maxY-minY
	span := math.Max(spanX, spanY)
	if span == 0 {
		span = 1
	}
	sc := (size - 2*margin) / span
	offX := margin + (size-2*margin-spanX*sc)/2
	offY := margin + (size-2*margin-spanY*sc)/2
	px := func(i int) float64 { return offX + (xs[i]-minX)*sc }
	py := func(i int) float64 { return offY + (ys[i]-minY)*sc }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f">`, size, size)
	fmt.Fprintf(&sb, `<rect width="%.0f" height="%.0f" fill="%s"/>`, size, size, visBG)
	sb.WriteString(`<g stroke="` + edgeCol + `" stroke-width="0.5">`)
	for _, e := range edges {
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`,
			px(e.a), py(e.a), px(e.b), py(e.b))
	}
	sb.WriteString(`</g>`)
	for i := range xs {
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`,
			px(i), py(i), radius, nodeColor[i])
	}
	sb.WriteString(`</svg>`)
	return sb.String()
}
