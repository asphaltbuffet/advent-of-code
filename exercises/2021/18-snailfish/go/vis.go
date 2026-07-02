package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the final summed snailfish number (part one, before taking the
// magnitude) as a binary tree (SVG). Internal nodes are pairs, leaves are the
// regular numbers; depth is shown top to bottom and encoded by node color on a
// colorblind-safe ramp, so the deep four-level nesting the reduce rules fight to
// contain is visible. Leaf values are labeled. Depth reads by vertical position
// as well as color, so the structure survives grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	nums := parse(instr)
	if len(nums) == 0 {
		return nil
	}
	acc := clone(nums[0])
	for _, n := range nums[1:] {
		acc = add(acc, n)
	}

	root := buildTree(acc)
	maxDepth := treeDepth(root)
	leaves := countLeaves(root)

	const (
		W      = 1100
		mTop   = 44
		mBot   = 20
		mSide  = 24
		levelH = 74
	)
	H := mTop + maxDepth*levelH + mBot
	plotW := float64(W - 2*mSide)

	// Assign each leaf an x slot left-to-right; internal nodes sit above the
	// midpoint of their children.
	var xpos []float64
	slot := 0
	var assign func(n *snode)
	assign = func(n *snode) {
		if n.leaf {
			x := float64(mSide) + (float64(slot)+0.5)/float64(leaves)*plotW
			n.x = x
			xpos = append(xpos, x)
			slot++
			return
		}
		assign(n.left)
		assign(n.right)
		n.x = (n.left.x + n.right.x) / 2
	}
	assign(root)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="26" fill="#e8ecf4" font-size="15" text-anchor="middle">Final summed snailfish number as a tree (nodes colored by depth)</text>`, W/2)

	yOf := func(depth int) float64 { return float64(mTop + depth*levelH) }

	// Edges first.
	var edges func(n *snode, depth int)
	edges = func(n *snode, depth int) {
		if n.leaf {
			return
		}
		for _, c := range []*snode{n.left, n.right} {
			fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#3a4250" stroke-width="1.5"/>`,
				n.x, yOf(depth), c.x, yOf(depth+1))
		}
		edges(n.left, depth+1)
		edges(n.right, depth+1)
	}
	edges(root, 0)

	// Nodes.
	var draw func(n *snode, depth int)
	draw = func(n *snode, depth int) {
		x, y := n.x, yOf(depth)
		col := depthColor(depth)
		if n.leaf {
			fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="12" fill="%s" stroke="#0d0f18" stroke-width="1.5"/>`, x, y, col)
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#0d0f18" font-size="12" font-weight="bold" text-anchor="middle" dominant-baseline="middle">%d</text>`, x, y, n.value)
			return
		}
		fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="6" fill="%s" stroke="#0d0f18" stroke-width="1"/>`, x, y, col)
		draw(n.left, depth+1)
		draw(n.right, depth+1)
	}
	draw(root, 0)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "snailfish.svg"), []byte(sb.String()), 0o600)
}

// snode is a reconstructed snailfish tree node.
type snode struct {
	leaf        bool
	value       int
	left, right *snode
	x           float64
}

// buildTree reconstructs the binary tree from the flat (value, depth) token list.
// Consecutive tokens at the same depth that form a pair are merged bottom-up.
func buildTree(toks []token) *snode {
	type item struct {
		node  *snode
		depth int
	}
	var stack []item
	for _, t := range toks {
		stack = append(stack, item{&snode{leaf: true, value: t.value}, t.depth})
		// Merge equal-depth adjacent pairs into a parent one level shallower.
		for len(stack) >= 2 && stack[len(stack)-1].depth == stack[len(stack)-2].depth {
			r := stack[len(stack)-1]
			l := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			parent := &snode{left: l.node, right: r.node}
			stack = append(stack, item{parent, l.depth - 1})
		}
	}
	return stack[0].node
}

func treeDepth(n *snode) int {
	if n.leaf {
		return 1
	}
	l, r := treeDepth(n.left), treeDepth(n.right)
	if r > l {
		l = r
	}
	return l + 1
}

func countLeaves(n *snode) int {
	if n.leaf {
		return 1
	}
	return countLeaves(n.left) + countLeaves(n.right)
}

// depthColor ramps from cool (shallow) to warm (deep) on a colorblind-safe scale.
func depthColor(depth int) string {
	cols := []string{"#0072B2", "#56B4E9", "#009E73", "#F0E442", "#E69F00", "#D55E00"}
	if depth >= len(cols) {
		depth = len(cols) - 1
	}
	return cols[depth]
}
