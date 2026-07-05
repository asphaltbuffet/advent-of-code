package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis draws the whole program tower as a layered dendrogram (SVG). The balanced
// subtrees are rendered in a clearly-visible slate, while the culprit path —
// the chain of unbalanced nodes leading to the offending disc — is highlighted
// in bright gold. The root is marked green and the culprit node is marked red.
func (e Exercise) Vis(instr, outdir string) error {
	tower := parseTower(instr)
	r := root(tower)

	path := culpritPath(tower, r)
	onPath := map[string]bool{}
	for _, n := range path {
		onPath[n] = true
	}
	culprit := path[len(path)-1]

	pos, maxLevel, nextLeaf := layoutTree(tower, r)

	const (
		levelGap = 300.0
		leafGap  = 15.0
		marginX  = 40.0
		marginY  = 40.0
	)
	w := float64(maxLevel)*levelGap + 2*marginX + 120
	h := nextLeaf*leafGap + 2*marginY
	tx := func(level float64) float64 { return marginX + level*levelGap }
	ty := func(spread float64) float64 { return marginY + spread*leafGap }

	names := make([]string, 0, len(tower))
	for n := range tower {
		names = append(names, n)
	}
	sort.Strings(names)

	svg := buildCircusSVG(tower, pos, names, onPath, culprit, r, w, h, tx, ty)
	return os.WriteFile(filepath.Join(outdir, "recursive-circus.svg"), []byte(svg), 0o644)
}

type treePos struct{ level, spread float64 }

func layoutTree(tower map[string]program, r string) (map[string]treePos, int, float64) {
	pos := map[string]treePos{}
	maxLevel := 0
	nextLeaf := 0.0
	var place func(name string, d int) float64
	place = func(name string, d int) float64 {
		if d > maxLevel {
			maxLevel = d
		}
		p := tower[name]
		if len(p.children) == 0 {
			s := nextLeaf
			nextLeaf++
			pos[name] = treePos{float64(d), s}
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
		pos[name] = treePos{float64(d), s}
		return s
	}
	place(r, 0)
	return pos, maxLevel, nextLeaf
}

func nodeStyle(n, culprit, rootName string, onPath map[string]bool) (string, string, string, float64, float64) {
	switch {
	case n == culprit:
		return "#ff4455", "#ff9aa5", "bold", 6, 13
	case n == rootName:
		return "#44ff88", "#8effb0", "bold", 6, 13
	case onPath[n]:
		return "#ffc84a", "#ffe0a0", "normal", 4, 11
	default:
		return "#7f92b5", "#9aa8c4", "normal", 2, 8
	}
}

func buildCircusSVG(
	tower map[string]program,
	pos map[string]treePos,
	names []string,
	onPath map[string]bool,
	culprit, rootName string,
	w, h float64,
	tx, ty func(float64) float64,
) string {
	var b strings.Builder
	fmt.Fprintf(
		&b,
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.0f %.0f" font-family="monospace">`+"\n",
		w, h,
	)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#161a26"/>`+"\n", w, h)

	fmt.Fprint(&b, `<g stroke-linecap="round">`+"\n")
	for _, n := range names {
		pp := pos[n]
		for _, c := range tower[n].children {
			if onPath[n] && onPath[c] {
				continue
			}
			cp := pos[c]
			fmt.Fprintf(
				&b,
				`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#5a6a8a" stroke-width="1" stroke-opacity="0.75"/>`+"\n",
				tx(pp.level),
				ty(pp.spread),
				tx(cp.level),
				ty(cp.spread),
			)
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
			fmt.Fprintf(
				&b,
				`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffc84a" stroke-width="2.5" stroke-opacity="0.98"/>`+"\n",
				tx(pp.level),
				ty(pp.spread),
				tx(cp.level),
				ty(cp.spread),
			)
		}
	}
	fmt.Fprint(&b, `</g>`+"\n")

	for _, n := range names {
		pp := pos[n]
		x, y := tx(pp.level), ty(pp.spread)
		dot, label, weight, r0, fs := nodeStyle(n, culprit, rootName, onPath)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%g" fill="%s"/>`+"\n", x, y, r0, dot)
		fmt.Fprintf(&b,
			`<text x="%.1f" y="%.1f" font-size="%g" font-weight="%s" fill="%s" text-anchor="start">%s</text>`+"\n",
			x+r0+3, y+fs*0.35, fs, weight, label, n,
		)
	}
	fmt.Fprint(&b, `</svg>`+"\n")
	return b.String()
}
