package exercises

import (
	"container/heap"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis produces two animated GIFs:
//   - vis-part1.gif: single robot collecting all keys in Part One
//   - vis-part2.gif: four robots collecting all keys in Part Two
func (e Exercise) Vis(instr, outdir string) error {
	if err := visPart1(instr, outdir); err != nil {
		return err
	}

	return visPart2(instr, outdir)
}

// ─── color palette (Okabe-Ito colorblind-safe) ───────────────────────────────

var palette = color.Palette{
	color.RGBA{R: 80, G: 80, B: 80, A: 255},    // 0 wall
	color.RGBA{R: 20, G: 20, B: 20, A: 255},    // 1 open / passage
	color.RGBA{R: 240, G: 200, B: 0, A: 255},   // 2 key (uncollected) — yellow
	color.RGBA{R: 100, G: 100, B: 100, A: 255}, // 3 key (collected) / door (unlocked) — dim gray
	color.RGBA{R: 213, G: 94, B: 0, A: 255},    // 4 door (locked) — vermillion
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 5 robot 0 — white
	color.RGBA{R: 86, G: 180, B: 233, A: 255},  // 6 robot 1 — sky blue
	color.RGBA{R: 0, G: 158, B: 115, A: 255},   // 7 robot 2 — bluish green
	color.RGBA{R: 204, G: 121, B: 167, A: 255}, // 8 robot 3 — reddish purple
}

const (
	colWall     = 0
	colOpen     = 1
	colKeyFree  = 2
	colKeyDone  = 3
	colDoorLock = 4
	colRobot0   = 5
	colRobot1   = 6
	colRobot2   = 7
	colRobot3   = 8
)

// ─── grid drawing helpers ─────────────────────────────────────────────────────

func gridCellSize(gridW, gridH int) int {
	size := 500 / max(gridW, gridH)
	return max(size, 4)
}

func drawGrid(img *image.Paletted, mg mazeGrid, keyMask uint32, robotPositions []int, robotCols []uint8, cell int) {
	width := mg.width
	height := len(mg.cells) / width

	fill := func(px, py int, idx uint8) {
		for dy := range cell {
			for dx := range cell {
				img.SetColorIndex(px*cell+dx, py*cell+dy, idx)
			}
		}
	}

	robotAt := make(map[int]uint8, len(robotPositions))
	for ri, pos := range robotPositions {
		c := uint8(colRobot0 + ri)
		if ri < len(robotCols) {
			c = robotCols[ri]
		}

		robotAt[pos] = c
	}

	for y := range height {
		for x := range width {
			pos := y*width + x
			c := mg.cells[pos]

			var idx uint8

			if rc, ok := robotAt[pos]; ok {
				idx = rc
			} else {
				idx = cellColorIdx(c, keyMask)
			}

			fill(x, y, idx)
		}
	}
}

func cellColorIdx(c byte, keyMask uint32) uint8 {
	switch {
	case c == '#':
		return colWall
	case c >= 'a' && c <= 'z':
		if keyMask&(1<<uint(c-'a')) != 0 {
			return colKeyDone
		}

		return colKeyFree
	case c >= 'A' && c <= 'Z':
		if keyMask&(1<<uint(c-'A')) != 0 {
			return colKeyDone
		}

		return colDoorLock
	default:
		return colOpen
	}
}

// ─── BFS path through maze ────────────────────────────────────────────────────

// mazePath returns the shortest sequence of cell positions from src to dst.
//
//nolint:gocognit // BFS path reconstruction; inherent complexity
func mazePath(mg mazeGrid, src, dst int) []int {
	if src == dst {
		return []int{src}
	}

	prev := make([]int, len(mg.cells))
	visited := make([]bool, len(mg.cells))

	for i := range prev {
		prev[i] = -1
	}

	visited[src] = true
	queue := []int{src}
	found := false
	dirs := [4]int{-mg.width, mg.width, -1, 1}

	for len(queue) > 0 && !found {
		cur := queue[0]
		queue = queue[1:]

		if cur == dst {
			found = true
			break
		}

		x := cur % mg.width

		for di, off := range dirs {
			if (di == 2 && x == 0) || (di == 3 && x == mg.width-1) {
				continue
			}

			np := cur + off
			if np < 0 || np >= len(mg.cells) || visited[np] || mg.cells[np] == '#' {
				continue
			}

			visited[np] = true
			prev[np] = cur
			queue = append(queue, np)
		}
	}

	if !found {
		return []int{src}
	}

	path := []int{}
	for p := dst; p != -1; p = prev[p] {
		path = append(path, p)
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// ─── Part One path reconstruction ────────────────────────────────────────────

type p1Step struct {
	node    int
	keyMask uint32
}

//nolint:gocognit // Dijkstra with predecessor tracking; inherent complexity
func dijkstraKeysWithPath(graph [][]edge, activeKeys []int, allKeys uint32, startNode int) []p1Step {
	type visitKey struct {
		node    int
		keyMask uint32
	}

	best := make(map[visitKey]int)
	pred := make(map[visitKey]visitKey)

	h := &pq{items: []dijkState{{cost: 0, node: startNode, keyMask: 0}}}
	heap.Init(h)

	var goalVK visitKey
	found := false

	for h.Len() > 0 {
		cur := heap.Pop(h).(dijkState) //nolint:errcheck // heap.Interface contract

		if cur.keyMask == allKeys {
			goalVK = visitKey{node: cur.node, keyMask: cur.keyMask}
			found = true

			break
		}

		vk := visitKey{node: cur.node, keyMask: cur.keyMask}

		if prev, ok := best[vk]; ok && prev <= cur.cost {
			continue
		}

		best[vk] = cur.cost

		for _, e := range graph[cur.node] {
			if cur.keyMask&e.reqKeys != e.reqKeys {
				continue
			}

			ki := activeKeys[e.to]
			if cur.keyMask&(1<<uint(ki)) != 0 {
				continue
			}

			newMask := cur.keyMask | (1 << uint(ki))
			newCost := cur.cost + e.dist
			nvk := visitKey{node: e.to, keyMask: newMask}

			if prev, ok := best[nvk]; ok && prev <= newCost {
				continue
			}

			pred[nvk] = vk
			heap.Push(h, dijkState{cost: newCost, node: e.to, keyMask: newMask})
		}
	}

	if !found {
		return nil
	}

	var path []visitKey

	for vk := goalVK; ; {
		path = append(path, vk)
		pv, ok := pred[vk]

		if !ok {
			break
		}

		vk = pv
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	result := make([]p1Step, len(path))
	for i, vk := range path {
		result[i] = p1Step(vk)
	}

	return result
}

// ─── Part Two path reconstruction ────────────────────────────────────────────

type p2Step struct {
	robots  [4]int8
	keyMask uint32
}

func dijkstraFourRobotsWithPath(graph [][]edge, activeKeys []int, allKeys uint32, startNodes [4]int) []p2Step {
	numKeys := len(activeKeys)
	initRobots := [4]int8{int8(startNodes[0]), int8(startNodes[1]), int8(startNodes[2]), int8(startNodes[3])}

	best := make(map[robotState]int)
	pred := make(map[robotState]robotState)
	h := &pq4{items: []dijkState4{{cost: 0, robots: initRobots, keyMask: 0}}}
	heap.Init(h)

	goalRS, found := dijk4Search(h, best, pred, graph, activeKeys, allKeys, numKeys)
	if !found {
		return nil
	}

	return reconstructP2Path(pred, goalRS)
}

func dijk4Search(
	h *pq4,
	best map[robotState]int,
	pred map[robotState]robotState,
	graph [][]edge,
	activeKeys []int,
	allKeys uint32,
	numKeys int,
) (robotState, bool) {
	for h.Len() > 0 {
		cur := heap.Pop(h).(dijkState4) //nolint:errcheck // heap.Interface contract

		if cur.keyMask == allKeys {
			return robotState{robots: cur.robots, keyMask: cur.keyMask}, true
		}

		rs := robotState{robots: cur.robots, keyMask: cur.keyMask}

		if prev, ok := best[rs]; ok && prev <= cur.cost {
			continue
		}

		best[rs] = cur.cost

		for r := range cur.robots {
			dijk4ExpandRobot(h, best, pred, rs, cur, r, graph, activeKeys, numKeys)
		}
	}

	return robotState{}, false
}

func dijk4ExpandRobot(
	h *pq4,
	best map[robotState]int,
	pred map[robotState]robotState,
	rs robotState,
	cur dijkState4,
	r int,
	graph [][]edge,
	activeKeys []int,
	numKeys int,
) {
	node := int(cur.robots[r])

	for _, e := range graph[node] {
		if cur.keyMask&e.reqKeys != e.reqKeys || e.to >= numKeys+4 {
			continue
		}

		ki := activeKeys[e.to]
		if cur.keyMask&(1<<uint(ki)) != 0 {
			continue
		}

		newMask := cur.keyMask | (1 << uint(ki))
		newCost := cur.cost + e.dist
		newRobots := cur.robots
		newRobots[r] = int8(e.to)
		nrs := robotState{robots: newRobots, keyMask: newMask}

		if prev, ok := best[nrs]; ok && prev <= newCost {
			continue
		}

		pred[nrs] = rs
		heap.Push(h, dijkState4{cost: newCost, robots: newRobots, keyMask: newMask})
	}
}

func reconstructP2Path(pred map[robotState]robotState, goalRS robotState) []p2Step {
	var path []robotState

	for rs := goalRS; ; {
		path = append(path, rs)
		pv, ok := pred[rs]

		if !ok {
			break
		}

		rs = pv
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	result := make([]p2Step, len(path))

	for i, rs := range path {
		result[i] = p2Step(rs)
	}

	return result
}

// ─── node→grid-position maps ──────────────────────────────────────────────────

func buildNodeToPos(mg mazeGrid, activeKeys []int, startPos int) []int {
	nodeToPos := make([]int, len(activeKeys)+1)

	for i, ki := range activeKeys {
		nodeToPos[i] = mg.keyPos[ki]
	}

	nodeToPos[len(activeKeys)] = startPos

	return nodeToPos
}

func buildNodeToPos4(mg mazeGrid, activeKeys []int, startPositions [4]int) []int {
	n := len(activeKeys)
	nodeToPos := make([]int, n+4)

	for i, ki := range activeKeys {
		nodeToPos[i] = mg.keyPos[ki]
	}

	for r := range 4 {
		nodeToPos[n+r] = startPositions[r]
	}

	return nodeToPos
}

// ─── GIF building ─────────────────────────────────────────────────────────────

func buildGIF(frames [][]int, keyMasks []uint32, mg mazeGrid, robotColsPerFrame [][]uint8, cell int) *gif.GIF {
	width := mg.width
	height := len(mg.cells) / width
	imgW := width * cell
	imgH := height * cell

	anim := &gif.GIF{}

	for fi, robotPositions := range frames {
		img := image.NewPaletted(image.Rect(0, 0, imgW, imgH), palette)
		cols := robotColsPerFrame[fi]
		drawGrid(img, mg, keyMasks[fi], robotPositions, cols, cell)
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 4) // 4 centiseconds per frame
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)
	}

	return anim
}

// ─── Part One visualization ───────────────────────────────────────────────────

func visPart1(instr, outdir string) error {
	mg, err := parseMaze(instr)
	if err != nil {
		return err
	}

	activeKeys, compactOf := buildKeyIndex(mg.keyPos, mg.numKeys)
	graph := buildGraph(mg, activeKeys, compactOf)
	allKeys := uint32((1 << mg.numKeys) - 1)

	steps := dijkstraKeysWithPath(graph, activeKeys, allKeys, mg.numKeys)
	if len(steps) == 0 {
		return nil
	}

	nodeToPos := buildNodeToPos(mg, activeKeys, mg.startPos)

	type cellStep struct {
		pos     int
		keyMask uint32
	}

	var cellSteps []cellStep

	for i, s := range steps {
		pos := nodeToPos[s.node]

		if i == 0 {
			cellSteps = append(cellSteps, cellStep{pos: pos, keyMask: s.keyMask})
			continue
		}

		prev := nodeToPos[steps[i-1].node]
		path := mazePath(mg, prev, pos)

		for j, p := range path {
			km := steps[i-1].keyMask
			if j == len(path)-1 {
				km = s.keyMask
			}

			cellSteps = append(cellSteps, cellStep{pos: p, keyMask: km})
		}
	}

	const maxFrames = 500

	stride := 1
	if len(cellSteps) > maxFrames {
		stride = (len(cellSteps) + maxFrames - 1) / maxFrames
	}

	var frames [][]int
	var keyMasks []uint32
	var robotColsPerFrame [][]uint8

	prevMask := uint32(0)

	for i, cs := range cellSteps {
		isKeyFrame := cs.keyMask != prevMask
		isFinal := i == len(cellSteps)-1

		if i%stride == 0 || isKeyFrame || isFinal {
			frames = append(frames, []int{cs.pos})
			keyMasks = append(keyMasks, cs.keyMask)
			robotColsPerFrame = append(robotColsPerFrame, []uint8{colRobot0})
		}

		prevMask = cs.keyMask
	}

	cell := gridCellSize(mg.width, len(mg.cells)/mg.width)
	anim := buildGIF(frames, keyMasks, mg, robotColsPerFrame, cell)

	f, err := os.Create(filepath.Join(outdir, "vis-part1.gif"))
	if err != nil {
		return err
	}

	defer f.Close()

	return gif.EncodeAll(f, anim)
}

// ─── Part Two visualization ───────────────────────────────────────────────────

func visPart2(instr, outdir string) error {
	mg, startPositions, err := parseMazeFour(instr)
	if err != nil {
		return err
	}

	activeKeys, compactOf := buildKeyIndex(mg.keyPos, mg.numKeys)
	graph := buildGraphFour(mg, activeKeys, compactOf, startPositions)
	allKeys := uint32((1 << mg.numKeys) - 1)

	numKeys := len(activeKeys)
	startNodes := [4]int{numKeys, numKeys + 1, numKeys + 2, numKeys + 3}

	steps := dijkstraFourRobotsWithPath(graph, activeKeys, allKeys, startNodes)
	if len(steps) == 0 {
		return nil
	}

	nodeToPos := buildNodeToPos4(mg, activeKeys, startPositions)
	allFrames := expandFrames4(steps, nodeToPos, mg, numKeys)
	frames, keyMasks, robotColsPerFrame := sampleFrames4(allFrames)

	cell := gridCellSize(mg.width, len(mg.cells)/mg.width)
	anim := buildGIF(frames, keyMasks, mg, robotColsPerFrame, cell)

	f2, err := os.Create(filepath.Join(outdir, "vis-part2.gif"))
	if err != nil {
		return err
	}

	defer f2.Close()

	return gif.EncodeAll(f2, anim)
}

type frame4 struct {
	positions [4]int
	keyMask   uint32
}

func expandFrames4(steps []p2Step, nodeToPos []int, mg mazeGrid, numKeys int) []frame4 {
	initPositions := [4]int{
		nodeToPos[numKeys],
		nodeToPos[numKeys+1],
		nodeToPos[numKeys+2],
		nodeToPos[numKeys+3],
	}

	allFrames := []frame4{{positions: initPositions, keyMask: 0}}
	prevPositions := initPositions

	for i := 1; i < len(steps); i++ {
		cur := steps[i]
		prev := steps[i-1]
		movedRobot := -1

		for r := range cur.robots {
			if cur.robots[r] != prev.robots[r] {
				movedRobot = r
				break
			}
		}

		if movedRobot < 0 {
			continue
		}

		srcPos := nodeToPos[prev.robots[movedRobot]]
		dstPos := nodeToPos[cur.robots[movedRobot]]
		path := mazePath(mg, srcPos, dstPos)

		for j, p := range path {
			positions := prevPositions
			positions[movedRobot] = p
			km := prev.keyMask

			if j == len(path)-1 {
				km = cur.keyMask
				prevPositions = positions
			}

			allFrames = append(allFrames, frame4{positions: positions, keyMask: km})
		}
	}

	return allFrames
}

func sampleFrames4(allFrames []frame4) ([][]int, []uint32, [][]uint8) {
	const maxFrames = 500

	stride := 1
	if len(allFrames) > maxFrames {
		stride = (len(allFrames) + maxFrames - 1) / maxFrames
	}

	var frames [][]int
	var keyMasks []uint32
	var robotColsPerFrame [][]uint8
	prevMask := uint32(0)

	for i, fr := range allFrames {
		isKeyFrame := fr.keyMask != prevMask
		isFinal := i == len(allFrames)-1

		if i%stride == 0 || isKeyFrame || isFinal {
			frames = append(frames, []int{fr.positions[0], fr.positions[1], fr.positions[2], fr.positions[3]})
			keyMasks = append(keyMasks, fr.keyMask)
			robotColsPerFrame = append(robotColsPerFrame, []uint8{colRobot0, colRobot1, colRobot2, colRobot3})
		}

		prevMask = fr.keyMask
	}

	return frames, keyMasks, robotColsPerFrame
}
