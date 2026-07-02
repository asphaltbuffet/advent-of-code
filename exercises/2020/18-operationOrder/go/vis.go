package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// node is a parse-tree node: a leaf value, or an operator with two children.
type node struct {
	op          string // "" for a leaf
	val         int    // leaf value
	left, right *node
	result      int // evaluated subtree value
}

// treeParser builds an actual parse tree (rather than folding into a value) so it
// can be drawn.
type treeParser struct {
	toks []string
	pos  int
	prec map[string]int
}

func (p *treeParser) peek() string {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return ""
}

func (p *treeParser) expr(minPrec int) *node {
	left := p.atom()
	for {
		op := p.peek()
		if op != "+" && op != "*" || p.prec[op] < minPrec {
			break
		}
		p.pos++
		right := p.expr(p.prec[op] + 1)
		n := &node{op: op, left: left, right: right}
		if op == "+" {
			n.result = left.result + right.result
		} else {
			n.result = left.result * right.result
		}
		left = n
	}
	return left
}

func (p *treeParser) atom() *node {
	t := p.toks[p.pos]
	p.pos++
	if t == "(" {
		v := p.expr(1)
		p.pos++ // ')'
		return v
	}
	n := 0
	for _, ch := range t {
		n = n*10 + int(ch-'0')
	}
	return &node{val: n, result: n}
}

// layout assigns x positions to leaves left-to-right and depth to every node.
type placed struct {
	n          *node
	x, y       int
	left, right *placed
}

// Vis draws the two parse trees for one representative expression: under Part
// One's equal-precedence rule and under Part Two's addition-first rule. The same
// expression yields different tree shapes and therefore different results, which
// is the whole puzzle. Operator nodes and leaves use distinct colorblind-safe
// colors, every node is labeled with its value, and the trees are labeled with
// their totals, so the diagram reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	// A compact expression whose two interpretations differ.
	const expr = "3 + 5 * 4 + 2 * 6"
	toks := tokenize(expr)

	build := func(prec map[string]int) *node {
		return (&treeParser{toks: append([]string(nil), toks...), prec: prec}).expr(1)
	}
	t1 := build(map[string]int{"+": 1, "*": 1})
	t2 := build(map[string]int{"+": 2, "*": 1})

	const (
		W       = 940
		H       = 460
		panelW  = 440
		gap     = 40
		mL      = 20
		mT      = 80
		levelH  = 74
		leafGap = 54
	)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="30" fill="#e8ecf4" font-size="16">Operation Order: %s under two precedence rules</text>`, mL, expr)

	opCol := "#F0E442"   // bright yellow: operator node
	leafCol := "#0072B2" // dark blue: number leaf

	drawTree := func(ox int, root *node, title string) {
		// Assign leaf x-order via in-order traversal.
		leafX := 0
		var place func(n *node, depth int) *placed
		place = func(n *node, depth int) *placed {
			if n.op == "" {
				p := &placed{n: n, x: ox + 20 + leafX*leafGap, y: mT + depth*levelH}
				leafX++
				return p
			}
			l := place(n.left, depth+1)
			r := place(n.right, depth+1)
			return &placed{n: n, x: (l.x + r.x) / 2, y: mT + depth*levelH, left: l, right: r}
		}
		root2 := place(root, 0)

		var draw func(p *placed)
		draw = func(p *placed) {
			if p.left != nil {
				fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, p.x, p.y, p.left.x, p.left.y)
				fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#3a424e"/>`, p.x, p.y, p.right.x, p.right.y)
				draw(p.left)
				draw(p.right)
			}
			if p.n.op == "" {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="14" fill="%s"/>`, p.x, p.y, leafCol)
				fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="12" font-weight="bold" text-anchor="middle">%d</text>`, p.x, p.y+4, p.n.val)
			} else {
				fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="16" fill="%s"/>`, p.x, p.y, opCol)
				fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="15" font-weight="bold" text-anchor="middle">%s</text>`, p.x, p.y+5, p.n.op)
				// subtree result under the operator
				fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="10" text-anchor="middle">=%d</text>`, p.x+24, p.y+3, p.n.result)
			}
		}
		draw(root2)

		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14">%s = %d</text>`, ox+20, mT-24, title, root.result)
	}

	drawTree(mL, t1, "Part 1 (equal precedence)")
	drawTree(mL+panelW+gap, t2, "Part 2 (+ before *)")

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "operation-order.svg"), []byte(sb.String()), 0o600)
}
