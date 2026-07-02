package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis draws the containment DAG rooted at shiny gold — the bags that must go
// inside it (Part Two). Each distinct color is one node, placed on the row for
// its deepest position below shiny gold; edges are labeled with how many of the
// child bag each parent holds. Each node also shows the total number of bags
// contained within one bag of that color, so the multiplication cascading up to
// 20189 is visible. Node color encodes depth on a colorblind-safe ramp reinforced
// by row position, so depth reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	rules := parse(instr)

	// Collect the descendant colors reachable from the target and each color's
	// deepest level (root = 0).
	level := map[string]int{}
	// Compute deepest level via DFS, keeping the max.
	var setLevel func(color string, depth int)
	setLevel = func(color string, depth int) {
		if d, ok := level[color]; !ok || depth > d {
			level[color] = depth
		}
		for _, c := range rules[color] {
			setLevel(c.color, depth+1)
		}
	}
	setLevel(target, 0)

	// Group colors by level.
	maxLevel := 0
	for _, d := range level {
		if d > maxLevel {
			maxLevel = d
		}
	}
	rows := make([][]string, maxLevel+1)
	for color, d := range level {
		rows[d] = append(rows[d], color)
	}
	for _, r := range rows {
		sort.Strings(r)
	}

	memo := map[string]int{}
	inside := func(c string) int { return countInside(c, rules, memo) }

	const (
		rowH   = 92
		nodeW  = 118
		nodeH  = 40
		mL     = 40
		mT     = 70
		colGap = 20
	)
	// Width from the widest row.
	widest := 0
	for _, r := range rows {
		if len(r) > widest {
			widest = len(r)
		}
	}
	W := mL*2 + widest*(nodeW+colGap)
	if W < 720 {
		W = 720
	}
	H := mT + (maxLevel+1)*rowH + 20

	// Position of each node.
	pos := map[string][2]int{}
	for d, r := range rows {
		total := len(r)
		rowWidth := total * (nodeW + colGap)
		startX := (W - rowWidth) / 2
		for i, color := range r {
			x := startX + i*(nodeW+colGap) + nodeW/2
			y := mT + d*rowH
			pos[color] = [2]int{x, y}
		}
	}

	// Okabe-Ito depth ramp.
	ramp := []string{"#F0E442", "#E69F00", "#56B4E9", "#009E73", "#0072B2", "#D55E00", "#CC79A7", "#8a94a2", "#7a869a"}
	depthCol := func(d int) string { return ramp[d%len(ramp)] }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="36" fill="#e8ecf4" font-size="16">Handy Haversacks: bags inside one shiny gold bag (part 2 = %d)</text>`, mL, inside(target))

	// Edges first (label with counts).
	for color, contents := range rules {
		pp, ok := pos[color]
		if !ok {
			continue
		}
		for _, c := range contents {
			cp, ok := pos[c.color]
			if !ok {
				continue
			}
			x1, y1 := pp[0], pp[1]+nodeH/2
			x2, y2 := cp[0], cp[1]-nodeH/2
			fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e" stroke-width="1"/>`, x1, y1, x2, y2)
			mx, my := (x1+x2)/2, (y1+y2)/2
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10">x%d</text>`, mx+2, my, c.count)
		}
	}

	// Nodes.
	for color, p := range pos {
		x, y := p[0], p[1]
		col := depthCol(level[color])
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" rx="5" fill="%s"/>`, x-nodeW/2, y-nodeH/2, nodeW, nodeH, col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="10" font-weight="bold" text-anchor="middle">%s</text>`, x, y-3, color)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="9" text-anchor="middle">holds %d</text>`, x, y+11, inside(color))
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "handy-haversacks.svg"), []byte(sb.String()), 0o600)
}
