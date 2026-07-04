package exercises //nolint:cyclop // orbit-tree visualization has inherently high package complexity

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// Vis renders the orbit tree as a top-down layered diagram. Each depth level
// is a horizontal band; nodes are small dots connected to their parent by a
// line. The YOU→SAN path is highlighted in Okabe-Ito orange (#E69F00); all
// other nodes are drawn in blue (#56B4E9) against a dark background (#111111).
// YOU and SAN nodes are drawn with a larger radius.
//
//nolint:gocognit,gocyclo,cyclop,funlen // orbit-tree layout and rendering is inherently complex
func (e Exercise) Vis(instr, outdir string) error {
	parent := make(map[string]string)
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ")", 2)
		if len(parts) == 2 {
			parent[parts[1]] = parts[0]
		}
	}

	// Build children map.
	children := make(map[string][]string)
	root := ""
	allNodes := make(map[string]bool)
	for child, par := range parent {
		children[par] = append(children[par], child)
		allNodes[child] = true
		allNodes[par] = true
	}
	// Root is the node with no parent (COM).
	for node := range allNodes {
		if _, hasParent := parent[node]; !hasParent {
			root = node
			break
		}
	}
	// Sort children for deterministic layout.
	for node := range children {
		sort.Strings(children[node])
	}

	// Assign depth to each node via BFS.
	depth := make(map[string]int)
	queue := []string{root}
	depth[root] = 0
	maxDepth := 0
	order := []string{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		d := depth[cur]
		if d > maxDepth {
			maxDepth = d
		}
		for _, ch := range children[cur] {
			depth[ch] = d + 1
			queue = append(queue, ch)
		}
	}

	// Assign horizontal positions using a leaf-count-based layout.
	// Each node's width is proportional to the number of leaf descendants.
	leafCount := make(map[string]int)
	// Compute leaf counts in reverse BFS order.
	for _, v := range slices.Backward(order) {
		node := v
		if len(children[node]) == 0 {
			leafCount[node] = 1
		} else {
			total := 0
			for _, ch := range children[node] {
				total += leafCount[ch]
			}
			leafCount[node] = total
		}
	}

	// Assign x-offset (in leaf units) to each node recursively.
	xOffset := make(map[string]float64) // left edge of node's subtree in leaf units
	xOffset[root] = 0
	for _, node := range order {
		cur := xOffset[node]
		for _, ch := range children[node] {
			xOffset[ch] = cur
			cur += float64(leafCount[ch])
		}
	}
	// Center of each node in leaf units = xOffset[node] + leafCount[node]/2.
	totalLeaves := float64(leafCount[root])

	// Compute YOU→SAN path nodes (set of nodes on the path).
	pathSet := make(map[string]bool)
	// Walk YOU's ancestors.
	youAnc := make(map[string]int)
	dist := 0
	cur := parent["YOU"]
	for cur != "" {
		youAnc[cur] = dist
		dist++
		cur = parent[cur]
	}
	// Walk SAN's ancestors until we find common ancestor.
	lca := ""
	dist = 0
	cur = parent["SAN"]
	for cur != "" {
		if d, ok := youAnc[cur]; ok {
			lcaDist := d + dist
			lca = cur
			_ = lcaDist
			break
		}
		dist++
		cur = parent[cur]
	}
	// Mark path: YOU's chain from parent(YOU) down to lca.
	if lca != "" {
		cur = parent["YOU"]
		for cur != lca {
			pathSet[cur] = true
			cur = parent[cur]
		}
		pathSet[lca] = true
		// SAN's chain from parent(SAN) down to lca.
		cur = parent["SAN"]
		for cur != lca {
			pathSet[cur] = true
			cur = parent[cur]
		}
		pathSet["YOU"] = true
		pathSet["SAN"] = true
	}

	// Image dimensions.
	const (
		marginX = 30
		marginY = 30
		rowH    = 6 // pixels per depth level (tree is ~350 deep)
		minW    = 1200
	)
	numRows := maxDepth + 1
	w := minW
	h := marginY*2 + numRows*rowH

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Fill background.
	bg := color.RGBA{0x11, 0x11, 0x11, 0xff}
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, bg)
		}
	}

	// Node center pixel coordinates.
	nodeX := func(node string) int {
		cx := xOffset[node] + float64(leafCount[node])/2.0
		return marginX + int(cx/totalLeaves*float64(w-2*marginX))
	}
	nodeY := func(node string) int {
		return marginY + depth[node]*rowH
	}

	// Colors (Okabe-Ito).
	colEdge := color.RGBA{0x33, 0x44, 0x55, 0xff}
	colEdgePath := color.RGBA{0xE6, 0x9F, 0x00, 0xff} // orange path
	colNode := color.RGBA{0x56, 0xB4, 0xE9, 0xff}     // blue
	colPath := color.RGBA{0xE6, 0x9F, 0x00, 0xff}     // orange
	colSpecial := color.RGBA{0xF0, 0xE4, 0x42, 0xff}  // yellow for YOU/SAN

	drawDot := func(cx, cy, r int, c color.RGBA) {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				if dx*dx+dy*dy <= r*r {
					px, py := cx+dx, cy+dy
					if px >= 0 && px < w && py >= 0 && py < h {
						img.SetRGBA(px, py, c)
					}
				}
			}
		}
	}

	drawLine := func(x1, y1, x2, y2 int, c color.RGBA) {
		dx := x2 - x1
		dy := y2 - y1
		steps := dx
		if steps < 0 {
			steps = -steps
		}
		if dy < 0 {
			dy = -dy
		}
		if dy > steps {
			steps = dy
		}
		if steps == 0 {
			return
		}
		for i := 0; i <= steps; i++ {
			px := x1 + (x2-x1)*i/steps
			py := y1 + (y2-y1)*i/steps
			if px >= 0 && px < w && py >= 0 && py < h {
				img.SetRGBA(px, py, c)
			}
		}
	}

	// Draw edges first.
	for _, node := range order {
		px, py := nodeX(node), nodeY(node)
		for _, ch := range children[node] {
			cx, cy := nodeX(ch), nodeY(ch)
			ec := colEdge
			if pathSet[node] && pathSet[ch] {
				ec = colEdgePath
			}
			drawLine(px, py, cx, cy, ec)
		}
	}

	// Draw nodes.
	for _, node := range order {
		nx, ny := nodeX(node), nodeY(node)
		r := 2
		col := colNode
		if pathSet[node] {
			col = colPath
			r = 3
		}
		if node == "YOU" || node == "SAN" {
			col = colSpecial
			r = 4
		}
		drawDot(nx, ny, r, col)
	}

	f, err := os.Create(filepath.Join(outdir, "vis.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
