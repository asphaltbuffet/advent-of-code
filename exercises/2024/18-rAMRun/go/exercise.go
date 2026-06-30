package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 18.
type Exercise struct {
	common.BaseExercise
}

var numRe = regexp.MustCompile(`\d+`)

// parseBytes returns the falling-byte coordinates as (x, y) pairs.
func parseBytes(instr string) [][2]int {
	nums := numRe.FindAllString(instr, -1)
	var pts [][2]int
	for i := 0; i+1 < len(nums); i += 2 {
		x, _ := strconv.Atoi(nums[i])
		y, _ := strconv.Atoi(nums[i+1])
		pts = append(pts, [2]int{x, y})
	}
	return pts
}

// params picks grid size and the Part-One byte count from the data: the example
// fits in 0..6 (7x7, 12 bytes); the real input is 71x71 with 1024 bytes.
func params(pts [][2]int) (size, count int) {
	for _, p := range pts {
		if p[0] > 6 || p[1] > 6 {
			return 71, 1024
		}
	}
	return 7, 12
}

// bfs returns the shortest path length from (0,0) to (size-1,size-1) avoiding
// corrupted cells, or -1 if unreachable. corrupt is keyed by y*size+x.
func bfs(size int, corrupt map[int]bool) int {
	start, goal := 0, (size-1)*size+(size-1)
	if corrupt[start] || corrupt[goal] {
		return -1
	}
	dist := map[int]int{start: 0}
	queue := []int{start}
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == goal {
			return dist[cur]
		}
		x, y := cur%size, cur/size
		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx < 0 || nx >= size || ny < 0 || ny >= size {
				continue
			}
			nk := ny*size + nx
			if corrupt[nk] {
				continue
			}
			if _, ok := dist[nk]; !ok {
				dist[nk] = dist[cur] + 1
				queue = append(queue, nk)
			}
		}
	}
	return -1
}

// corruptAfter builds the corruption set for the first n fallen bytes.
func corruptAfter(pts [][2]int, size, n int) map[int]bool {
	corrupt := make(map[int]bool, n)
	for i := 0; i < n && i < len(pts); i++ {
		corrupt[pts[i][1]*size+pts[i][0]] = true
	}
	return corrupt
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	pts := parseBytes(instr)
	size, count := params(pts)
	return bfs(size, corruptAfter(pts, size, count)), nil
}

// firstBlocker binary-searches the smallest prefix length that disconnects the
// exit, returning the index of the byte that does it.
func firstBlocker(pts [][2]int, size int) int {
	// reachable(n): is the exit reachable after the first n bytes fall?
	reachable := func(n int) bool {
		return bfs(size, corruptAfter(pts, size, n)) != -1
	}
	lo, hi := 0, len(pts) // lo reachable, hi assumed blocked
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if reachable(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi - 1 // the byte added at index hi-1 caused the block
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	pts := parseBytes(instr)
	size, _ := params(pts)
	idx := firstBlocker(pts, size)
	return fmt.Sprintf("%d,%d", pts[idx][0], pts[idx][1]), nil
}

// --- Visualization ---

// bfsPath returns the cells (keyed by y*size+x) on a shortest path, or nil.
func bfsPath(size int, corrupt map[int]bool) []int {
	start, goal := 0, (size-1)*size+(size-1)
	if corrupt[start] || corrupt[goal] {
		return nil
	}
	prev := map[int]int{start: -1}
	queue := []int{start}
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == goal {
			var path []int
			for c := goal; c != -1; c = prev[c] {
				path = append(path, c)
			}
			return path
		}
		x, y := cur%size, cur/size
		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx < 0 || nx >= size || ny < 0 || ny >= size {
				continue
			}
			nk := ny*size + nx
			if corrupt[nk] {
				continue
			}
			if _, ok := prev[nk]; !ok {
				prev[nk] = cur
				queue = append(queue, nk)
			}
		}
	}
	return nil
}

// Vis renders the grid: the Part-One safe path in green over the early
// corruption, plus all later-falling bytes and the Part-Two blocking byte
// highlighted in red.
func (e Exercise) Vis(instr string, outdir string) error {
	pts := parseBytes(instr)
	size, count := params(pts)
	blockIdx := firstBlocker(pts, size)

	pathCells := map[int]bool{}
	for _, c := range bfsPath(size, corruptAfter(pts, size, count)) {
		pathCells[c] = true
	}

	const scale = 14
	const pad = 10
	img := image.NewRGBA(image.Rect(0, 0, size*scale+2*pad, size*scale+2*pad))

	bg := color.RGBA{0x0c, 0x0c, 0x14, 0xff}
	early := color.RGBA{0x55, 0x55, 0x66, 0xff}  // bytes fallen by Part-One count
	later := color.RGBA{0x2a, 0x2a, 0x36, 0xff}  // bytes that fall afterwards
	pathC := color.RGBA{0x44, 0xe0, 0x80, 0xff}  // Part-One shortest path
	blockC := color.RGBA{0xff, 0x3b, 0x3b, 0xff} // Part-Two blocking byte
	startC := color.RGBA{0x33, 0xff, 0xcc, 0xff}
	goalC := color.RGBA{0xff, 0xd1, 0x66, 0xff}

	fill := func(x, y int, col color.RGBA) {
		x0, y0 := pad+x*scale, pad+y*scale
		for yy := y0; yy < y0+scale; yy++ {
			for xx := x0; xx < x0+scale; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// Bytes that fall after the Part-One count (faint), up to the blocker.
	for i := count; i <= blockIdx && i < len(pts); i++ {
		fill(pts[i][0], pts[i][1], later)
	}
	// Bytes fallen by the Part-One count.
	for i := 0; i < count && i < len(pts); i++ {
		fill(pts[i][0], pts[i][1], early)
	}
	// Part-One shortest path.
	for c := range pathCells {
		fill(c%size, c/size, pathC)
	}
	// Start, goal, and the Part-Two blocking byte.
	fill(0, 0, startC)
	fill(size-1, size-1, goalC)
	fill(pts[blockIdx][0], pts[blockIdx][1], blockC)

	f, err := os.Create(filepath.Join(outdir, "ram-run.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
