package exercises

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 12.
type Exercise struct {
	common.BaseExercise
}

// parseGraph reads the pipe listing into an adjacency map.
func parseGraph(instr string) map[int][]int {
	graph := map[int][]int{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "<->", 2)
		id, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		for _, p := range strings.Split(parts[1], ",") {
			peer, _ := strconv.Atoi(strings.TrimSpace(p))
			graph[id] = append(graph[id], peer)
		}
	}
	return graph
}

// component returns every program reachable from start, marking them in seen.
func component(graph map[int][]int, start int, seen map[int]bool) int {
	size := 0
	stack := []int{start}
	seen[start] = true
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		size++
		for _, next := range graph[cur] {
			if !seen[next] {
				seen[next] = true
				stack = append(stack, next)
			}
		}
	}
	return size
}

// One returns the size of the group containing program 0.
func (e Exercise) One(instr string) (any, error) {
	graph := parseGraph(instr)
	return component(graph, 0, map[int]bool{}), nil
}

// Two returns the number of distinct groups.
func (e Exercise) Two(instr string) (any, error) {
	graph := parseGraph(instr)
	seen := map[int]bool{}
	groups := 0
	for id := range graph {
		if !seen[id] {
			component(graph, id, seen)
			groups++
		}
	}
	return groups, nil
}

// componentNodes returns the sorted list of nodes reachable from start.
func componentNodes(graph map[int][]int, start int, seen map[int]bool) []int {
	var nodes []int
	stack := []int{start}
	seen[start] = true
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		nodes = append(nodes, cur)
		for _, next := range graph[cur] {
			if !seen[next] {
				seen[next] = true
				stack = append(stack, next)
			}
		}
	}
	sort.Ints(nodes)
	return nodes
}

type vec struct{ x, y float64 }

// layoutComponent runs a small Fruchterman-Reingold force simulation on the
// given nodes and returns their positions (centred on the origin) plus the
// component's bounding radius.
func layoutComponent(nodes []int, adj map[int]map[int]bool, rng *rand.Rand) (map[int]vec, float64) {
	n := len(nodes)
	pos := make(map[int]vec, n)
	if n == 1 {
		pos[nodes[0]] = vec{0, 0}
		return pos, 1
	}

	// Initial positions on a small circle.
	for i, id := range nodes {
		a := 2 * math.Pi * float64(i) / float64(n)
		r := 4 + rng.Float64()*2
		pos[id] = vec{r * math.Cos(a), r * math.Sin(a)}
	}

	area := float64(n) * 100
	k := math.Sqrt(area / float64(n)) // ideal edge length
	temp := k * 2
	iters := 220
	for it := 0; it < iters; it++ {
		disp := make(map[int]vec, n)
		// Repulsion between every pair in the component.
		for i := 0; i < n; i++ {
			a := nodes[i]
			for j := i + 1; j < n; j++ {
				b := nodes[j]
				dx, dy := pos[a].x-pos[b].x, pos[a].y-pos[b].y
				d := math.Hypot(dx, dy) + 0.01
				f := k * k / d
				ux, uy := dx/d*f, dy/d*f
				disp[a] = vec{disp[a].x + ux, disp[a].y + uy}
				disp[b] = vec{disp[b].x - ux, disp[b].y - uy}
			}
		}
		// Attraction along edges.
		for i := 0; i < n; i++ {
			a := nodes[i]
			for b := range adj[a] {
				if b <= a {
					continue // each undirected edge once
				}
				if _, ok := pos[b]; !ok {
					continue
				}
				dx, dy := pos[a].x-pos[b].x, pos[a].y-pos[b].y
				d := math.Hypot(dx, dy) + 0.01
				f := d * d / k
				ux, uy := dx/d*f, dy/d*f
				disp[a] = vec{disp[a].x - ux, disp[a].y - uy}
				disp[b] = vec{disp[b].x + ux, disp[b].y + uy}
			}
		}
		// Apply displacement, capped by temperature; cool down.
		for _, id := range nodes {
			dp := disp[id]
			d := math.Hypot(dp.x, dp.y) + 0.01
			step := math.Min(d, temp)
			pos[id] = vec{pos[id].x + dp.x/d*step, pos[id].y + dp.y/d*step}
		}
		temp *= 0.97
	}

	// Centre and measure radius.
	var cx, cy float64
	for _, id := range nodes {
		cx += pos[id].x
		cy += pos[id].y
	}
	cx, cy = cx/float64(n), cy/float64(n)
	radius := 1.0
	for _, id := range nodes {
		p := vec{pos[id].x - cx, pos[id].y - cy}
		pos[id] = p
		radius = math.Max(radius, math.Hypot(p.x, p.y))
	}
	return pos, radius
}

// Vis lays out the pipe graph as component "islands" (SVG). Each connected group
// gets its own force-directed layout and a distinct colour; the group holding
// program 0 is highlighted, and the many singleton groups pepper the field.
func (e Exercise) Vis(instr, outdir string) error {
	graph := parseGraph(instr)

	// Undirected adjacency (pipes are bidirectional).
	adj := map[int]map[int]bool{}
	addEdge := func(a, b int) {
		if adj[a] == nil {
			adj[a] = map[int]bool{}
		}
		adj[a][b] = true
	}
	for a, peers := range graph {
		for _, b := range peers {
			addEdge(a, b)
			addEdge(b, a)
		}
	}

	// Collect components, largest first.
	type comp struct {
		nodes  []int
		pos    map[int]vec
		radius float64
		hasZero bool
	}
	adjSlice := mapToSlice(adj)
	seen := map[int]bool{}
	var comps []comp
	for id := range adj {
		if seen[id] {
			continue
		}
		nodes := componentNodes(adjSlice, id, seen)
		comps = append(comps, comp{nodes: nodes})
	}
	sort.Slice(comps, func(i, j int) bool { return len(comps[i].nodes) > len(comps[j].nodes) })

	rng := rand.New(rand.NewSource(12))
	for i := range comps {
		p, r := layoutComponent(comps[i].nodes, adj, rng)
		comps[i].pos = p
		comps[i].radius = r
		for _, id := range comps[i].nodes {
			if id == 0 {
				comps[i].hasZero = true
			}
		}
	}

	// Shelf-pack variable-size tiles: each component's tile matches its own
	// radius, so the giant group and the many singletons all pack tightly.
	const gap = 16.0
	// Wide enough for the largest island, but roomy enough that several medium
	// components share each shelf instead of standing alone.
	targetW := math.Max(comps[0].radius*2+gap, 2200)

	type placed struct {
		idx    int
		cx, cy float64
	}
	var slots []placed
	x, y, shelfH := gap, gap, 0.0
	for idx, c := range comps {
		w := c.radius*2 + gap
		if x+w > targetW && x > gap {
			x = gap
			y += shelfH
			shelfH = 0
		}
		cx := x + c.radius + gap/2
		cy := y + c.radius + gap/2
		slots = append(slots, placed{idx, cx, cy})
		x += w
		if h := c.radius*2 + gap; h > shelfH {
			shelfH = h
		}
	}
	W := targetW + gap
	H := y + shelfH + gap

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f">`+"\n", W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0d0f18"/>`+"\n", W, H)

	for _, s := range slots {
		idx := s.idx
		c := comps[idx]
		cx, cy := s.cx, s.cy

		hue := float64((idx*47)%360)
		stroke := hslHex(hue, 0.5, 0.5)
		fill := hslHex(hue, 0.6, 0.62)
		if c.hasZero {
			stroke, fill = "#ffd24a", "#ffe27a" // highlight program 0's group
		}

		// Edges.
		fmt.Fprintf(&b, `<g stroke="%s" stroke-width="0.5" stroke-opacity="0.5">`+"\n", stroke)
		for _, a := range c.nodes {
			for peer := range adj[a] {
				if peer <= a {
					continue
				}
				pa, pb := c.pos[a], c.pos[peer]
				fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`+"\n",
					cx+pa.x, cy+pa.y, cx+pb.x, cy+pb.y)
			}
		}
		fmt.Fprint(&b, `</g>`+"\n")

		// Nodes.
		r := 1.3
		if len(c.nodes) == 1 {
			r = 1.8
		}
		for _, id := range c.nodes {
			p := c.pos[id]
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%g" fill="%s"/>`+"\n", cx+p.x, cy+p.y, r, fill)
		}
		if c.hasZero {
			p := c.pos[0]
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="#ff4455"/>`+"\n", cx+p.x, cy+p.y)
		}
	}
	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "digital-plumber.svg"), []byte(b.String()), 0o644)
}

// mapToSlice converts an undirected adjacency-set map into the []int adjacency
// form componentNodes expects.
func mapToSlice(adj map[int]map[int]bool) map[int][]int {
	out := make(map[int][]int, len(adj))
	for a, peers := range adj {
		for b := range peers {
			out[a] = append(out[a], b)
		}
	}
	return out
}

// hslHex renders an HSL colour as a hex string.
func hslHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	hp := h / 60
	x := c * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, bb float64
	switch {
	case hp < 1:
		r, g, bb = c, x, 0
	case hp < 2:
		r, g, bb = x, c, 0
	case hp < 3:
		r, g, bb = 0, c, x
	case hp < 4:
		r, g, bb = 0, x, c
	case hp < 5:
		r, g, bb = x, 0, c
	default:
		r, g, bb = c, 0, x
	}
	m := l - c/2
	return fmt.Sprintf("#%02x%02x%02x", int((r+m)*255), int((g+m)*255), int((bb+m)*255))
}
