package exercises

import (
	"container/heap"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Region shades, tiered by brightness so terrain reads in grayscale:
// rocky darkest, wet mid, narrow brightest.
var regionColor = [3]color.RGBA{
	{40, 44, 52, 255},   // rocky
	{72, 110, 150, 255}, // wet
	{150, 176, 210, 255}, // narrow
}

// The rescue path and its halo. The path is a warm bright line; the halo is a dark
// outline so the line stays legible over the brightest (narrow) terrain in
// grayscale.
var (
	pathColor = color.RGBA{240, 200, 40, 255}
	haloColor = color.RGBA{16, 16, 20, 255}
)

// Vis draws the cave terrain (rocky/wet/narrow) with the optimal rescue path traced
// from the mouth to the target. The full grid is a wide, thin strip, so it is
// cropped vertically to the band the path travels through (plus a margin).
func (e Exercise) Vis(instr, outdir string) error {
	c := parse(instr)
	path := rescuePath(c)

	// Vertical band the path occupies, padded a little.
	minY, maxY := c.ty, c.ty
	for _, p := range path {
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	const pad = 20
	minY = max(minY-pad, 0)
	maxY = maxY + pad
	minX, maxX := 0, c.tx+pad

	w := maxX - minX + 1
	h := maxY - minY + 1
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			img.Set(x-minX, y-minY, regionColor[c.erosion(x, y)%3])
		}
	}

	// Halo first, then the path over it.
	onPath := map[[2]int]bool{}
	for _, p := range path {
		onPath[p] = true
	}
	for _, p := range path {
		for _, d := range [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
			hx, hy := p[0]+d[0], p[1]+d[1]
			if hx < minX || hy < minY || hx > maxX || hy > maxY || onPath[[2]int{hx, hy}] {
				continue
			}
			img.Set(hx-minX, hy-minY, haloColor)
		}
	}
	for _, p := range path {
		img.Set(p[0]-minX, p[1]-minY, pathColor)
	}
	// Mark mouth and target with a small bright block.
	markBlock(img, 0-minX, 0-minY, color.RGBA{255, 255, 255, 255})
	markBlock(img, c.tx-minX, c.ty-minY, color.RGBA{255, 80, 80, 255})

	f, err := os.Create(filepath.Join(outdir, "mode-maze.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func markBlock(img *image.RGBA, x, y int, col color.RGBA) {
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

// rescuePath reruns the Dijkstra with predecessor tracking and returns the (x, y)
// coordinates along the fastest route (ignoring tool state, so consecutive
// same-cell tool switches collapse to one point).
func rescuePath(c *cave) [][2]int {
	w, h := c.tx+margin, c.ty+margin
	region := make([][]int, h+1)
	for y := 0; y <= h; y++ {
		region[y] = make([]int, w+1)
		for x := 0; x <= w; x++ {
			region[y][x] = c.erosion(x, y) % 3
		}
	}

	idx := func(x, y, tool int) int { return (y*(w+1)+x)*3 + tool }
	dist := make([]int, (w+1)*(h+1)*3)
	prev := make([]int, (w+1)*(h+1)*3)
	for i := range dist {
		dist[i] = 1 << 30
		prev[i] = -1
	}
	dist[idx(0, 0, torch)] = 0
	goal := idx(c.tx, c.ty, torch)

	pq := &priorityQueue{{x: 0, y: 0, tool: torch, cost: 0}}
	found := false
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		ci := idx(cur.x, cur.y, cur.tool)
		if cur.cost > dist[ci] {
			continue
		}
		if ci == goal {
			found = true
			break
		}
		relaxP := func(nx, ny, tool, nd int) {
			ni := idx(nx, ny, tool)
			if nd < dist[ni] {
				dist[ni] = nd
				prev[ni] = ci
				heap.Push(pq, item{x: nx, y: ny, tool: tool, cost: nd})
			}
		}
		for tool := 0; tool < 3; tool++ {
			if tool != region[cur.y][cur.x] && tool != cur.tool {
				relaxP(cur.x, cur.y, tool, cur.cost+7)
			}
		}
		for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < 0 || ny < 0 || nx > w || ny > h {
				continue
			}
			if region[ny][nx] != cur.tool {
				relaxP(nx, ny, cur.tool, cur.cost+1)
			}
		}
	}
	if !found {
		return nil
	}

	// Walk predecessors back from the goal, collecting distinct cells.
	var path [][2]int
	seen := map[[2]int]bool{}
	for i := goal; i != -1; i = prev[i] {
		cell := i / 3
		x, y := cell%(w+1), cell/(w+1)
		if !seen[[2]int{x, y}] {
			seen[[2]int{x, y}] = true
			path = append(path, [2]int{x, y})
		}
	}
	return path
}
