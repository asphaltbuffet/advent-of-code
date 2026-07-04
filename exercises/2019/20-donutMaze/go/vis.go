package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis produces two animated GIFs:
//   - vis-part1.gif: BFS shortest path through the flat (non-recursive) maze
//   - vis-part2.gif: recursive BFS path, position colored by recursion depth
func (e Exercise) Vis(instr, outdir string) error {
	if err := visMazePart1(instr, outdir); err != nil {
		return err
	}
	return visMazePart2(instr, outdir)
}

// ─── color palette (Okabe-Ito colorblind-safe) ───────────────────────────────

var mazePalette = color.Palette{
	color.RGBA{R: 80, G: 80, B: 80, A: 255},    // 0 wall — dark gray
	color.RGBA{R: 20, G: 20, B: 20, A: 255},    // 1 open passage — near-black
	color.RGBA{R: 0, G: 158, B: 115, A: 255},   // 2 portal tile — teal
	color.RGBA{R: 86, G: 180, B: 233, A: 255},  // 3 AA start — sky-blue
	color.RGBA{R: 213, G: 94, B: 0, A: 255},    // 4 ZZ end — vermillion
	color.RGBA{R: 180, G: 160, B: 0, A: 255},   // 5 breadcrumb visited — dim yellow
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 6 current position (level 0) — white
	color.RGBA{R: 86, G: 180, B: 233, A: 255},  // 7 current position level 1 — sky-blue
	color.RGBA{R: 0, G: 158, B: 115, A: 255},   // 8 current position level 2 — bluish-green
	color.RGBA{R: 230, G: 159, B: 0, A: 255},   // 9 current position level 3 — orange
	color.RGBA{R: 213, G: 94, B: 0, A: 255},    // 10 current position level 4 — vermillion
	color.RGBA{R: 150, G: 100, B: 200, A: 255}, // 11 current position level 5+ — purple
}

const (
	colMazeWall       = 0
	colMazeOpen       = 1
	colMazePortal     = 2
	colMazeStart      = 3
	colMazeEnd        = 4
	colMazeBreadcrumb = 5
	colMazeCurLevel0  = 6
	colMazeCurLevel1  = 7
	colMazeCurLevel2  = 8
	colMazeCurLevel3  = 9
	colMazeCurLevel4  = 10
	colMazeCurLevel5  = 11
)

func levelColor(level int) uint8 {
	switch {
	case level <= 0:
		return colMazeCurLevel0
	case level == 1:
		return colMazeCurLevel1
	case level == 2:
		return colMazeCurLevel2
	case level == 3:
		return colMazeCurLevel3
	case level == 4:
		return colMazeCurLevel4
	default:
		return colMazeCurLevel5
	}
}

func mazeGridCellSize(gridW, gridH int) int {
	size := 500 / max(gridW, gridH)
	return max(size, 4)
}

// ─── path reconstruction ──────────────────────────────────────────────────────

// reconstructPath rebuilds a path from start to end using the visited-predecessor map.
func reconstructPath(visited map[pos]pos, start, end pos) []pos {
	path := []pos{}
	for p := end; p != start; p = visited[p] {
		path = append(path, p)
	}
	path = append(path, start)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// expandP1Neighbors appends unvisited BFS neighbors (cardinal + portal) of cur.
func expandP1Neighbors(
	grid []string,
	portals map[pos]pos,
	cur pos,
	visited map[pos]pos,
	next []pos,
) []pos {
	for _, d := range dirs {
		nb := pos{cur.x + d.x, cur.y + d.y}
		if gridCell(grid, nb) != '.' {
			continue
		}
		if _, seen := visited[nb]; seen {
			continue
		}
		visited[nb] = cur
		next = append(next, nb)
	}
	partner, ok := portals[cur]
	if !ok {
		return next
	}
	if _, seen := visited[partner]; !seen {
		visited[partner] = cur
		next = append(next, partner)
	}
	return next
}

// bfsMazeWithPath runs BFS from start to end and returns the actual path as a
// slice of positions (in order from start to end).
func bfsMazeWithPath(grid []string, portals map[pos]pos, start, end pos) []pos {
	visited := map[pos]pos{start: start}
	queue := []pos{start}

	for len(queue) > 0 {
		next := make([]pos, 0, len(queue))
		for _, cur := range queue {
			if cur == end {
				return reconstructPath(visited, start, end)
			}
			next = expandP1Neighbors(grid, portals, cur, visited, next)
		}
		queue = next
	}
	return nil
}

// p2State is the BFS state for Part Two (position + recursion level).
type p2State struct {
	p     pos
	level int
}

// reconstructP2Path rebuilds a Part Two path from start to end.
func reconstructP2Path(visited map[p2State]p2State, startState, endState p2State) []p2State {
	path := []p2State{}
	for s := endState; s != startState; s = visited[s] {
		path = append(path, s)
	}
	path = append(path, startState)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// expandP2VisNeighbors appends unvisited BFS neighbors for Part Two path reconstruction.
func expandP2VisNeighbors(
	grid []string,
	portals map[pos]pos,
	cur p2State,
	isOuter func(pos) bool,
	visited map[p2State]p2State,
	next []p2State,
) []p2State {
	for _, d := range dirs {
		nb := pos{cur.p.x + d.x, cur.p.y + d.y}
		if gridCell(grid, nb) != '.' {
			continue
		}
		ns := p2State{nb, cur.level}
		if _, seen := visited[ns]; !seen {
			visited[ns] = cur
			next = append(next, ns)
		}
	}
	ns := portalJumpTarget(portals, p2BFSState(cur), isOuter)
	if ns == (p2BFSState{}) {
		return next
	}
	target := p2State(ns)
	if _, seen := visited[target]; !seen {
		visited[target] = cur
		next = append(next, target)
	}
	return next
}

// bfsMazePart2WithPath runs the recursive BFS and returns the path as
// a slice of (pos, level) steps.
func bfsMazePart2WithPath(grid []string, portals map[pos]pos, start, end pos) []p2State {
	minX, minY, maxX, maxY := mazeBounds(grid)
	isOuter := outerChecker(minX, minY, maxX, maxY)

	startState := p2State{start, 0}
	endState := p2State{end, 0}

	visited := map[p2State]p2State{startState: startState}
	queue := []p2State{startState}

	for len(queue) > 0 {
		next := make([]p2State, 0, len(queue))
		for _, cur := range queue {
			if cur == endState {
				return reconstructP2Path(visited, startState, endState)
			}
			next = expandP2VisNeighbors(grid, portals, cur, isOuter, visited, next)
		}
		queue = next
	}
	return nil
}

// ─── GIF rendering ────────────────────────────────────────────────────────────

type mazeFrame struct {
	cur     pos
	visited map[pos]int // pos → level at which it was visited (for coloring)
	level   int         // current recursion level (Part Two only)
}

// cellBaseColor returns the palette index for a maze cell (wall, portal, start, end, open).
func cellBaseColor(grid []string, portals map[pos]pos, start, end pos, p pos) uint8 {
	c := gridCell(grid, p)
	switch c {
	case '#':
		return colMazeWall
	case '.':
		_, isPortal := portals[p]
		switch {
		case p == start:
			return colMazeStart
		case p == end:
			return colMazeEnd
		case isPortal:
			return colMazePortal
		default:
			return colMazeOpen
		}
	default:
		return colMazeOpen
	}
}

func drawMazeFrame(
	grid []string,
	portals map[pos]pos,
	start, end pos,
	cell int,
	frame mazeFrame,
) *image.Paletted {
	gridH := len(grid)
	gridW := 0
	for _, row := range grid {
		if len(row) > gridW {
			gridW = len(row)
		}
	}

	imgW := gridW * cell
	imgH := gridH * cell

	img := image.NewPaletted(image.Rect(0, 0, imgW, imgH), mazePalette)

	fill := func(p pos, idx uint8) {
		for dy := range cell {
			for dx := range cell {
				img.SetColorIndex(p.x*cell+dx, p.y*cell+dy, idx)
			}
		}
	}

	// Draw base maze.
	for y, row := range grid {
		for x := range row {
			p := pos{x, y}
			fill(p, cellBaseColor(grid, portals, start, end, p))
		}
	}

	// Draw breadcrumbs.
	for p := range frame.visited {
		if p != frame.cur {
			fill(p, colMazeBreadcrumb)
		}
	}

	// Draw current position.
	fill(frame.cur, levelColor(frame.level))

	return img
}

// ─── Part One ─────────────────────────────────────────────────────────────────

// isPortalJump reports whether two consecutive path positions represent a portal teleport.
func isPortalJump(prev, cur pos) bool {
	dx := prev.x - cur.x
	dy := prev.y - cur.y
	return dx < -1 || dx > 1 || dy < -1 || dy > 1
}

// buildP1Keyframes returns the set of frame indices that must always be rendered.
func buildP1Keyframes(path []pos, portalSet map[pos]bool) map[int]bool {
	keyframes := map[int]bool{0: true, len(path) - 1: true}
	for i, p := range path {
		if !portalSet[p] || i == 0 {
			continue
		}
		if isPortalJump(path[i-1], p) {
			keyframes[i] = true
			keyframes[i-1] = true
		}
	}
	return keyframes
}

func visMazePart1(instr, outdir string) error {
	grid, portals, start, end := parseMaze(instr)

	path := bfsMazeWithPath(grid, portals, start, end)
	if len(path) == 0 {
		return nil
	}

	gridH := len(grid)
	gridW := 0
	for _, row := range grid {
		if len(row) > gridW {
			gridW = len(row)
		}
	}
	cell := mazeGridCellSize(gridW, gridH)

	// Build portal position set for keyframe detection.
	portalSet := map[pos]bool{}
	for p := range portals {
		portalSet[p] = true
	}

	keyframes := buildP1Keyframes(path, portalSet)

	const maxFrames = 400
	stride := 1
	if len(path) > maxFrames {
		stride = (len(path) + maxFrames - 1) / maxFrames
	}

	anim := &gif.GIF{}
	visited := map[pos]int{}

	for i, p := range path {
		if i%stride != 0 && !keyframes[i] {
			visited[p] = 0
			continue
		}
		frame := mazeFrame{
			cur:     p,
			visited: visited,
			level:   0,
		}
		img := drawMazeFrame(grid, portals, start, end, cell, frame)
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 4)
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)
		visited[p] = 0
	}

	f, err := os.Create(filepath.Join(outdir, "vis-part1.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

// ─── Part Two ─────────────────────────────────────────────────────────────────

func visMazePart2(instr, outdir string) error {
	grid, portals, start, end := parseMaze(instr)

	path := bfsMazePart2WithPath(grid, portals, start, end)
	if len(path) == 0 {
		return nil
	}

	gridH := len(grid)
	gridW := 0
	for _, row := range grid {
		if len(row) > gridW {
			gridW = len(row)
		}
	}
	cell := mazeGridCellSize(gridW, gridH)

	// Keyframes: start, end, level changes.
	keyframes := map[int]bool{0: true, len(path) - 1: true}
	prevLevel := 0
	for i, s := range path {
		if s.level != prevLevel {
			keyframes[i] = true
			if i > 0 {
				keyframes[i-1] = true
			}
			prevLevel = s.level
		}
	}

	const maxFrames = 400
	stride := 1
	if len(path) > maxFrames {
		stride = (len(path) + maxFrames - 1) / maxFrames
	}

	anim := &gif.GIF{}
	visited := map[pos]int{}

	for i, s := range path {
		if i%stride != 0 && !keyframes[i] {
			visited[s.p] = s.level
			continue
		}
		frame := mazeFrame{
			cur:     s.p,
			visited: visited,
			level:   s.level,
		}
		img := drawMazeFrame(grid, portals, start, end, cell, frame)
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 4)
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)
		visited[s.p] = s.level
	}

	f, err := os.Create(filepath.Join(outdir, "vis-part2.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}
