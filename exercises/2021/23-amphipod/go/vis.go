package exercises

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// solvePath runs Dijkstra like solve but records predecessors, returning the
// sequence of states along one least-energy sort (start first, goal last) and
// the cumulative energy after each state.
func solvePath(start state) ([]state, []int) {
	depth := len(start.rooms[0])
	dist := map[state]int{start: 0}
	prev := map[state]state{}
	q := &pq{{start, 0}}
	var goal state
	found := false

	for q.Len() > 0 {
		cur := heap.Pop(q).(item)
		if cur.cost > dist[cur.s] {
			continue
		}
		if done(cur.s) {
			goal = cur.s
			found = true
			break
		}
		for _, m := range moves(cur.s, depth) {
			nc := cur.cost + m.cost
			if d, ok := dist[m.st]; !ok || nc < d {
				dist[m.st] = nc
				prev[m.st] = cur.s
				heap.Push(q, item{m.st, nc})
			}
		}
	}
	if !found {
		return nil, nil
	}

	var chain []state
	for s := goal; ; s = prev[s] {
		chain = append(chain, s)
		if s == start {
			break
		}
	}
	// Reverse to start-first order.
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	costs := make([]int, len(chain))
	for i, s := range chain {
		costs[i] = dist[s]
	}
	return chain, costs
}

// Okabe-Ito colors per amphipod type; brightness increases with cost so the
// expensive movers read darker/warmer and the chart survives grayscale.
var amphColor = map[byte]string{
	'A': "#56B4E9", // sky blue   (cost 1)
	'B': "#009E73", // bluish green(cost 10)
	'C': "#E69F00", // orange     (cost 100)
	'D': "#D55E00", // vermilion  (cost 1000)
}

// Vis renders the least-energy sort of the folded (Part One) burrow as a strip
// of successive states, each burrow drawn as a hallway over four rooms, with the
// cumulative energy under each. Amphipods are labeled A..D and colored on a
// colorblind-safe palette, so type reads by letter as well as hue.
func (e Exercise) Vis(instr, outdir string) error {
	start := parseBurrow(instr)
	chain, costs := solvePath(start)
	if chain == nil {
		return fmt.Errorf("no solution to visualize")
	}

	depth := len(start.rooms[0])
	const (
		cellW  = 22
		cellH  = 22
		bMargX = 10
		bMargY = 26
		gapX   = 26
		topPad = 40
	)
	burrowW := 11*cellW + 2*bMargX
	burrowH := (1+depth)*cellH + 2*bMargY
	perRow := 6
	rows := (len(chain) + perRow - 1) / perRow
	W := perRow*(burrowW+gapX) - gapX + 40
	H := topPad + rows*(burrowH+gapX) + 20

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="20" y="26" fill="#e8ecf4" font-size="16">Least-energy sort (part one): %d moves, total %d energy</text>`, len(chain)-1, costs[len(costs)-1])

	for idx, s := range chain {
		col := idx % perRow
		row := idx / perRow
		ox := 20 + col*(burrowW+gapX)
		oy := topPad + row*(burrowH+gapX)
		drawBurrow(&sb, ox, oy, cellW, cellH, bMargX, bMargY, depth, s)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="12" text-anchor="middle">%d energy</text>`,
			ox+burrowW/2, oy+burrowH+2, costs[idx])
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "amphipod.svg"), []byte(sb.String()), 0o600)
}

// drawBurrow draws one burrow state: an 11-cell hallway over four rooms.
func drawBurrow(sb *strings.Builder, ox, oy, cw, ch, mx, my, depth int, s state) {
	wall := "#2a3038"
	// Hallway row.
	for c := 0; c < 11; c++ {
		x := ox + mx + c*cw
		y := oy + my
		fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="#181c22" stroke="%s"/>`, x, y, cw, ch, wall)
		if v := s.hall[c]; v != empty {
			drawAmph(sb, x, y, cw, ch, v)
		}
	}
	// Rooms.
	for j := 0; j < 4; j++ {
		col := roomHall[j]
		for i := 0; i < depth; i++ {
			x := ox + mx + col*cw
			y := oy + my + (i+1)*ch
			fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="#181c22" stroke="%s"/>`, x, y, cw, ch, wall)
			if v := s.rooms[j][i]; v != empty {
				drawAmph(sb, x, y, cw, ch, v)
			}
		}
	}
}

func drawAmph(sb *strings.Builder, x, y, cw, ch int, v byte) {
	fmt.Fprintf(sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x+2, y+2, cw-4, ch-4, amphColor[v])
	fmt.Fprintf(sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="13" font-weight="bold" text-anchor="middle">%c</text>`,
		x+cw/2, y+ch/2+5, v)
}
