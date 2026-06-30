package exercises

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 23.
type Exercise struct {
	common.BaseExercise
}

// graph maps each computer to the set of its neighbours.
type graph map[string]map[string]bool

func parseGraph(instr string) graph {
	g := graph{}
	add := func(a, b string) {
		if g[a] == nil {
			g[a] = map[string]bool{}
		}
		g[a][b] = true
	}
	for _, line := range strings.Fields(instr) {
		ab := strings.SplitN(line, "-", 2)
		add(ab[0], ab[1])
		add(ab[1], ab[0])
	}
	return g
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	g := parseGraph(instr)

	triangles := map[string]bool{}
	for a, nbrs := range g {
		neigh := keys(nbrs)
		for i := 0; i < len(neigh); i++ {
			for j := i + 1; j < len(neigh); j++ {
				b, c := neigh[i], neigh[j]
				if g[b][c] {
					tri := []string{a, b, c}
					sort.Strings(tri)
					triangles[strings.Join(tri, ",")] = true
				}
			}
		}
	}

	count := 0
	for tri := range triangles {
		for _, n := range strings.Split(tri, ",") {
			if n[0] == 't' {
				count++
				break
			}
		}
	}
	return count, nil
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// bronKerbosch finds the maximum clique via pivoting, updating *best.
func (g graph) bronKerbosch(r, p, x map[string]bool, best *[]string) {
	if len(p) == 0 && len(x) == 0 {
		if len(r) > len(*best) {
			*best = keys(r)
		}
		return
	}
	// Choose a pivot u from p∪x maximising |neighbours in p|.
	var pivot string
	maxN := -1
	for _, u := range append(keys(p), keys(x)...) {
		n := 0
		for v := range p {
			if g[u][v] {
				n++
			}
		}
		if n > maxN {
			maxN, pivot = n, u
		}
	}

	// Iterate vertices in p not adjacent to the pivot.
	for _, v := range keys(p) {
		if g[pivot][v] {
			continue
		}
		nr := copySet(r)
		nr[v] = true
		np, nx := map[string]bool{}, map[string]bool{}
		for w := range p {
			if g[v][w] {
				np[w] = true
			}
		}
		for w := range x {
			if g[v][w] {
				nx[w] = true
			}
		}
		g.bronKerbosch(nr, np, nx, best)
		delete(p, v)
		x[v] = true
	}
}

func copySet(m map[string]bool) map[string]bool {
	out := make(map[string]bool, len(m))
	for k := range m {
		out[k] = true
	}
	return out
}

// maxClique returns the sorted node names of the largest clique.
func (g graph) maxClique() []string {
	p := map[string]bool{}
	for n := range g {
		p[n] = true
	}
	var best []string
	g.bronKerbosch(map[string]bool{}, p, map[string]bool{}, &best)
	sort.Strings(best)
	return best
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	g := parseGraph(instr)
	return strings.Join(g.maxClique(), ","), nil
}

// --- Visualization ---

// Vis draws the network on a circular layout with the maximum clique (Part Two)
// pulled to the centre and highlighted.
func (e Exercise) Vis(instr string, outdir string) error {
	g := parseGraph(instr)
	clique := g.maxClique()
	inClique := map[string]bool{}
	for _, n := range clique {
		inClique[n] = true
	}

	nodes := keys(graphNodeSet(g))
	sort.Strings(nodes) // stable layout

	const size = 1200.0
	const cx, cy = size / 2, size / 2
	outerR := size/2 - 40
	innerR := size / 3.2 // large inner ring so clique chords separate

	pos := map[string][2]float64{}
	var ring []string
	for _, n := range nodes {
		if !inClique[n] {
			ring = append(ring, n)
		}
	}
	for i, n := range ring {
		a := 2 * math.Pi * float64(i) / float64(len(ring))
		pos[n] = [2]float64{cx + outerR*math.Cos(a), cy + outerR*math.Sin(a)}
	}
	for i, n := range clique {
		a := 2 * math.Pi * float64(i) / float64(len(clique))
		pos[n] = [2]float64{cx + innerR*math.Cos(a), cy + innerR*math.Sin(a)}
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %g %g" font-family="monospace">`+"\n", size, size)
	fmt.Fprintf(&b, `<rect width="%g" height="%g" fill="#080810"/>`+"\n", size, size)

	line := func(a, c [2]float64, stroke string, w, op float64) {
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%g" stroke-opacity="%g"/>`+"\n",
			a[0], a[1], c[0], c[1], stroke, w, op)
	}

	// Faint background edges.
	fmt.Fprint(&b, `<g stroke="#3a4a6a" stroke-width="0.5" stroke-opacity="0.18">`+"\n")
	seen := map[string]bool{}
	for a, nbrs := range g {
		for c := range nbrs {
			if inClique[a] && inClique[c] {
				continue
			}
			ek := a + c
			if a > c {
				ek = c + a
			}
			if seen[ek] {
				continue
			}
			seen[ek] = true
			pa, pc := pos[a], pos[c]
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f"/>`+"\n", pa[0], pa[1], pc[0], pc[1])
		}
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Clique edges: each a crisp vector line, sharp at any zoom.
	for i := 0; i < len(clique); i++ {
		for j := i + 1; j < len(clique); j++ {
			line(pos[clique[i]], pos[clique[j]], "#ffc84a", 1.2, 0.95)
		}
	}

	// Nodes.
	for _, n := range nodes {
		p := pos[n]
		if inClique[n] {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="6" fill="#ffe080"/>`+"\n", p[0], p[1])
			// Offset the label radially outward from centre so it clears the node.
			dx, dy := p[0]-cx, p[1]-cy
			d := math.Hypot(dx, dy)
			lx, ly := p[0]+dx/d*16, p[1]+dy/d*16+4
			fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="13" font-weight="bold" fill="#ff3ad0" text-anchor="middle">%s</text>`+"\n", lx, ly, n)
		} else {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2" fill="#6a7a9a"/>`+"\n", p[0], p[1])
		}
	}
	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "lan-party.svg"), []byte(b.String()), 0o644)
}

func graphNodeSet(g graph) map[string]bool {
	s := map[string]bool{}
	for n := range g {
		s[n] = true
	}
	return s
}
