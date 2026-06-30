package exercises

import (
	"container/heap"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 16.
type Exercise struct {
	common.BaseExercise
}

// Directions: 0=East, 1=South, 2=West, 3=North. (d+1)%4 turns clockwise.
var dirVec = [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

type state struct{ r, c, d int }

type pqItem struct {
	st   state
	cost int
}

type priorityQueue []pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].cost < pq[j].cost }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(pqItem)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type maze struct {
	grid       []string
	start, end [2]int
}

func parseMaze(instr string) maze {
	m := maze{grid: strings.Fields(instr)}
	for r := range m.grid {
		for c := 0; c < len(m.grid[r]); c++ {
			switch m.grid[r][c] {
			case 'S':
				m.start = [2]int{r, c}
			case 'E':
				m.end = [2]int{r, c}
			}
		}
	}
	return m
}

func (m maze) wall(r, c int) bool {
	return r < 0 || r >= len(m.grid) || c < 0 || c >= len(m.grid[r]) || m.grid[r][c] == '#'
}

// dijkstra returns the minimum cost to each (r,c,dir) state from S facing East,
// plus the set of predecessor states for each state (for best-path recovery).
func (m maze) dijkstra() (map[state]int, map[state][]state) {
	dist := map[state]int{}
	preds := map[state][]state{}
	start := state{m.start[0], m.start[1], 0}
	dist[start] = 0

	pq := &priorityQueue{{start, 0}}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(pqItem)
		if d, ok := dist[cur.st]; ok && cur.cost > d {
			continue
		}
		// Neighbours: forward (cost 1) and two turns (cost 1000).
		type edge struct {
			next state
			cost int
		}
		var edges []edge
		fr, fc := cur.st.r+dirVec[cur.st.d][0], cur.st.c+dirVec[cur.st.d][1]
		if !m.wall(fr, fc) {
			edges = append(edges, edge{state{fr, fc, cur.st.d}, 1})
		}
		edges = append(edges,
			edge{state{cur.st.r, cur.st.c, (cur.st.d + 1) % 4}, 1000},
			edge{state{cur.st.r, cur.st.c, (cur.st.d + 3) % 4}, 1000},
		)

		for _, e := range edges {
			nc := cur.cost + e.cost
			best, ok := dist[e.next]
			if !ok || nc < best {
				dist[e.next] = nc
				preds[e.next] = []state{cur.st}
				heap.Push(pq, pqItem{e.next, nc})
			} else if nc == best {
				preds[e.next] = append(preds[e.next], cur.st)
			}
		}
	}
	return dist, preds
}

// bestScore returns the minimum cost to reach E in any orientation.
func (m maze) bestScore(dist map[state]int) (int, []state) {
	best := -1
	var ends []state
	for d := 0; d < 4; d++ {
		st := state{m.end[0], m.end[1], d}
		if v, ok := dist[st]; ok {
			if best == -1 || v < best {
				best = v
				ends = []state{st}
			} else if v == best {
				ends = append(ends, st)
			}
		}
	}
	return best, ends
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	m := parseMaze(instr)
	dist, _ := m.dijkstra()
	best, _ := m.bestScore(dist)
	return best, nil
}

// goodTiles returns the set of (r,c) tiles on at least one best path.
func (m maze) goodTiles() map[[2]int]bool {
	dist, preds := m.dijkstra()
	_, ends := m.bestScore(dist)

	tiles := map[[2]int]bool{}
	seen := map[state]bool{}
	stack := append([]state{}, ends...)
	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[st] {
			continue
		}
		seen[st] = true
		tiles[[2]int{st.r, st.c}] = true
		stack = append(stack, preds[st]...)
	}
	return tiles
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	m := parseMaze(instr)
	return len(m.goodTiles()), nil
}

// --- Visualization ---

const (
	visCell = 8
	visPad  = 12
)

// rampColor returns a cool->warm colour for t in [0,1] (path-distance ramp).
func rampColor(t float64) color.RGBA {
	// teal -> yellow -> magenta
	r := uint8(60 + 195*t)
	g := uint8(230 - 150*t)
	b := uint8(180 - 60*t + 80*t*t)
	return color.RGBA{r, g, b, 0xff}
}

// Vis renders the maze with the best-path tiles (Part Two "good seats")
// highlighted on a distance ramp, with S and E marked.
func (e Exercise) Vis(instr string, outdir string) error {
	m := parseMaze(instr)
	dist, preds := m.dijkstra()
	best, ends := m.bestScore(dist)

	// Collect good tiles and the best score at which each is reached.
	tileScore := map[[2]int]int{}
	seen := map[state]bool{}
	stack := append([]state{}, ends...)
	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[st] {
			continue
		}
		seen[st] = true
		key := [2]int{st.r, st.c}
		if cur, ok := tileScore[key]; !ok || dist[st] < cur {
			tileScore[key] = dist[st]
		}
		stack = append(stack, preds[st]...)
	}

	rows := len(m.grid)
	cols := 0
	for _, row := range m.grid {
		if len(row) > cols {
			cols = len(row)
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, cols*visCell+2*visPad, rows*visCell+2*visPad))

	wallC := color.RGBA{0x20, 0x20, 0x2c, 0xff}
	floorC := color.RGBA{0x0c, 0x0c, 0x12, 0xff}
	startC := color.RGBA{0x33, 0xff, 0x66, 0xff}
	endC := color.RGBA{0xff, 0x44, 0x44, 0xff}

	fill := func(r, c int, col color.RGBA) {
		x0, y0 := visPad+c*visCell, visPad+r*visCell
		for y := y0; y < y0+visCell; y++ {
			for x := x0; x < x0+visCell; x++ {
				img.SetRGBA(x, y, col)
			}
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < len(m.grid[r]); c++ {
			ch := m.grid[r][c]
			switch {
			case ch == '#':
				fill(r, c, wallC)
			case ch == 'S':
				fill(r, c, startC)
			case ch == 'E':
				fill(r, c, endC)
			default:
				if s, ok := tileScore[[2]int{r, c}]; ok {
					t := 0.0
					if best > 0 {
						t = float64(s) / float64(best)
					}
					fill(r, c, rampColor(t))
				} else {
					fill(r, c, floorC)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "reindeer-maze.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
