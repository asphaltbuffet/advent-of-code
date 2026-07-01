package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 7.
type Exercise struct {
	common.BaseExercise
}

type program struct {
	weight   int
	children []string
}

// parseTower reads the program listing into a name->program map.
func parseTower(instr string) map[string]program {
	tower := map[string]program{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// name (weight) [-> c1, c2, ...]
		name := strings.Fields(line)[0]
		open := strings.IndexByte(line, '(')
		close := strings.IndexByte(line, ')')
		weight, _ := strconv.Atoi(line[open+1 : close])

		var children []string
		if arrow := strings.Index(line, "->"); arrow >= 0 {
			for _, c := range strings.Split(line[arrow+2:], ",") {
				children = append(children, strings.TrimSpace(c))
			}
		}
		tower[name] = program{weight, children}
	}
	return tower
}

// root returns the one program that is never listed as anyone's child.
func root(tower map[string]program) string {
	isChild := map[string]bool{}
	for _, p := range tower {
		for _, c := range p.children {
			isChild[c] = true
		}
	}
	for name := range tower {
		if !isChild[name] {
			return name
		}
	}
	return ""
}

// One returns the name of the bottom program.
func (e Exercise) One(instr string) (any, error) {
	return root(parseTower(instr)), nil
}

// totalWeight returns the combined weight of a program's sub-tower.
func totalWeight(tower map[string]program, name string) int {
	p := tower[name]
	sum := p.weight
	for _, c := range p.children {
		sum += totalWeight(tower, c)
	}
	return sum
}

// balance descends from name looking for the single wrong-weight program. want
// is the total weight this sub-tower is expected to have. It returns the
// corrected own-weight for the offending program and true once found.
func balance(tower map[string]program, name string, want int) (int, bool) {
	p := tower[name]

	// Group children by their total weight.
	totals := make([]int, len(p.children))
	counts := map[int]int{}
	for i, c := range p.children {
		totals[i] = totalWeight(tower, c)
		counts[totals[i]]++
	}

	// If a child's total is the minority, the imbalance lives inside it.
	if len(counts) > 1 {
		var oddTotal, goodTotal int
		for t, n := range counts {
			if n == 1 {
				oddTotal = t
			} else {
				goodTotal = t
			}
		}
		for i, c := range p.children {
			if totals[i] == oddTotal {
				// The odd child should weigh goodTotal in total; recurse in
				// case its own children are further unbalanced.
				return balance(tower, c, goodTotal)
			}
		}
	}

	// Children balance (or there are none): this program is the culprit. Its
	// corrected own-weight makes the sub-tower hit want.
	childSum := 0
	for _, t := range totals {
		childSum += t
	}
	return want - childSum, true
}

// Two returns the corrected weight for the single unbalanced program.
func (e Exercise) Two(instr string) (any, error) {
	tower := parseTower(instr)
	fixed, _ := balance(tower, root(tower), 0)
	return fixed, nil
}

// culpritPath returns the chain of names from the root down to the single
// unbalanced program (inclusive).
func culpritPath(tower map[string]program, name string) []string {
	p := tower[name]
	totals := make([]int, len(p.children))
	counts := map[int]int{}
	for i, c := range p.children {
		totals[i] = totalWeight(tower, c)
		counts[totals[i]]++
	}
	if len(counts) > 1 {
		var oddTotal int
		for t, n := range counts {
			if n == 1 {
				oddTotal = t
			}
		}
		for i, c := range p.children {
			if totals[i] == oddTotal {
				return append([]string{name}, culpritPath(tower, c)...)
			}
		}
	}
	return []string{name} // balanced children: this node is the culprit
}

// Vis draws the whole program tower as a layered dendrogram (SVG). The balanced
// structure is drawn faintly; the path from the root down to the single
// unbalanced program is highlighted, with the culprit marked in red.
func (e Exercise) Vis(instr, outdir string) error {
	tower := parseTower(instr)
	r := root(tower)

	path := culpritPath(tower, r)
	onPath := map[string]bool{}
	for _, n := range path {
		onPath[n] = true
	}
	culprit := path[len(path)-1]

	// Layout, in abstract coordinates:
	//   level  = depth from the root (0 at the root)
	//   spread = position along the fan; leaves get increasing spread in DFS
	//            order, internal nodes sit at the midpoint of their children.
	// These map to the drawing as level->x (left to right) and spread->y (top
	// to bottom), giving a left-to-right tree.
	type pt struct{ level, spread float64 }
	pos := map[string]pt{}
	maxLevel := 0
	var nextLeaf float64

	var place func(name string, d int) float64
	place = func(name string, d int) float64 {
		if d > maxLevel {
			maxLevel = d
		}
		p := tower[name]
		if len(p.children) == 0 {
			s := nextLeaf
			nextLeaf++
			pos[name] = pt{float64(d), s}
			return s
		}
		lo, hi := 0.0, 0.0
		for i, c := range p.children {
			cs := place(c, d+1)
			if i == 0 || cs < lo {
				lo = cs
			}
			if i == 0 || cs > hi {
				hi = cs
			}
		}
		s := (lo + hi) / 2
		pos[name] = pt{float64(d), s}
		return s
	}
	place(r, 0)

	const (
		levelGap = 300.0 // horizontal room between depths (fits labels)
		leafGap  = 15.0  // vertical spacing between stacked leaves
		marginX  = 40.0
		marginY  = 40.0
	)
	W := float64(maxLevel)*levelGap + 2*marginX + 120 // +room for rightmost labels
	H := nextLeaf*leafGap + 2*marginY
	tx := func(level float64) float64 { return marginX + level*levelGap }
	ty := func(spread float64) float64 { return marginY + spread*leafGap }

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" font-family="monospace">`+"\n", W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#161a26"/>`+"\n", W, H)

	names := make([]string, 0, len(tower))
	for n := range tower {
		names = append(names, n)
	}
	sort.Strings(names) // deterministic output

	// Edges: the balanced tower in a clearly-visible slate, the culprit path in
	// bright gold on top.
	fmt.Fprint(&b, `<g stroke-linecap="round">`+"\n")
	for _, n := range names {
		pp := pos[n]
		for _, c := range tower[n].children {
			if onPath[n] && onPath[c] {
				continue // drawn later, on top
			}
			cp := pos[c]
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#5a6a8a" stroke-width="1" stroke-opacity="0.75"/>`+"\n",
				tx(pp.level), ty(pp.spread), tx(cp.level), ty(cp.spread))
		}
	}
	for _, n := range names {
		if !onPath[n] {
			continue
		}
		pp := pos[n]
		for _, c := range tower[n].children {
			if !onPath[c] {
				continue
			}
			cp := pos[c]
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffc84a" stroke-width="2.5" stroke-opacity="0.98"/>`+"\n",
				tx(pp.level), ty(pp.spread), tx(cp.level), ty(cp.spread))
		}
	}
	fmt.Fprint(&b, `</g>`+"\n")

	// Nodes + labels. Every program is labelled just to the right of its dot.
	for _, n := range names {
		pp := pos[n]
		x, y := tx(pp.level), ty(pp.spread)
		var dot, label string
		var r0, fs float64
		var weight string
		switch {
		case n == culprit:
			dot, label, r0, fs = "#ff4455", "#ff9aa5", 6, 13
			weight = "bold"
		case n == r:
			dot, label, r0, fs = "#44ff88", "#8effb0", 6, 13
			weight = "bold"
		case onPath[n]:
			dot, label, r0, fs = "#ffc84a", "#ffe0a0", 4, 11
			weight = "normal"
		default:
			dot, label, r0, fs = "#7f92b5", "#9aa8c4", 2, 8
			weight = "normal"
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%g" fill="%s"/>`+"\n", x, y, r0, dot)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%g" font-weight="%s" fill="%s" text-anchor="start">%s</text>`+"\n",
			x+r0+3, y+fs*0.35, fs, weight, label, n)
	}
	fmt.Fprint(&b, `</svg>`+"\n")

	return os.WriteFile(filepath.Join(outdir, "recursive-circus.svg"), []byte(b.String()), 0o644)
}
