package exercises

import (
	"container/heap"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Vis renders the full part-two 5x5-expanded risk grid as a PNG and traces the
// lowest-risk path across it. Risk is drawn as a dark-to-bright relief (low risk
// dark, high risk bright), and the optimal path from the top-left to the
// bottom-right is overlaid in vermilion. The path threads through the low-risk
// (dark) valleys while avoiding the bright high-risk cells. Risk is encoded by
// brightness, so the terrain reads in grayscale; the path stays visible as a
// continuous line.
func (e Exercise) Vis(instr, outdir string) error {
	base, err := parse(instr)
	if err != nil {
		return err
	}
	grid := expand(base, 5)
	rows, cols := len(grid), len(grid[0])

	path := lowestRiskPath(grid)
	onPath := make([][]bool, rows)
	for r := range onPath {
		onPath[r] = make([]bool, cols)
	}
	// halo marks a dark outline one cell around the path so the bright path
	// separates from the risk field regardless of the local brightness.
	halo := make([][]bool, rows)
	for r := range halo {
		halo[r] = make([]bool, cols)
	}
	// Thicken the path to a 2x2 block per cell so it reads as a continuous line.
	for _, p := range path {
		for _, d := range [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
			nr, nc := p[0]+d[0], p[1]+d[1]
			if nr < rows && nc < cols {
				onPath[nr][nc] = true
			}
		}
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if !onPath[r][c] {
				continue
			}
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					nr, nc := r+dr, c+dc
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !onPath[nr][nc] {
						halo[nr][nc] = true
					}
				}
			}
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, cols, rows))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			switch {
			case onPath[r][c]:
				// Path drawn brightest (white) so it stays distinct in grayscale;
				// a warm tint keeps it identifiable in color too.
				img.SetRGBA(c, r, color.RGBA{0xff, 0xf0, 0xd8, 0xff})
			case halo[r][c]:
				img.SetRGBA(c, r, color.RGBA{0x08, 0x06, 0x02, 0xff}) // dark outline
			default:
				// Risk 1..9 mapped to a dark-to-bright blue-gray relief.
				f := float64(grid[r][c]-1) / 8
				v := uint8(24 + 200*f)
				img.SetRGBA(c, r, color.RGBA{uint8(float64(v) * 0.8), uint8(float64(v) * 0.85), v, 0xff})
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "chiton.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// lowestRiskPath runs Dijkstra like lowestRisk but keeps a best-known distance
// and predecessor per cell so the optimal path can be reconstructed.
func lowestRiskPath(grid [][]int) [][2]int {
	rows, cols := len(grid), len(grid[0])
	const inf = 1 << 30
	dist := make([][]int, rows)
	prev := make([][][2]int, rows)
	for r := range dist {
		dist[r] = make([]int, cols)
		prev[r] = make([][2]int, cols)
		for c := range dist[r] {
			dist[r][c] = inf
			prev[r][c] = [2]int{-1, -1}
		}
	}
	dist[0][0] = 0

	pq := &pqueue{}
	heap.Push(pq, node{0, 0, 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(node)
		if cur.risk > dist[cur.r][cur.c] {
			continue // stale queue entry
		}
		if cur.r == rows-1 && cur.c == cols-1 {
			break
		}
		for _, m := range moves {
			nr, nc := cur.r+m[0], cur.c+m[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				continue
			}
			nd := cur.risk + grid[nr][nc]
			if nd < dist[nr][nc] {
				dist[nr][nc] = nd
				prev[nr][nc] = [2]int{cur.r, cur.c}
				heap.Push(pq, node{nd, nr, nc})
			}
		}
	}

	// Walk back from the goal.
	var path [][2]int
	cur := [2]int{rows - 1, cols - 1}
	for cur != [2]int{-1, -1} {
		path = append(path, cur)
		cur = prev[cur[0]][cur[1]]
	}
	return path
}
