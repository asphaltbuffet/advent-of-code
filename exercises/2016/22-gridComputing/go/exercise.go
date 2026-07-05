package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 22.
type Exercise struct {
	common.BaseExercise
}

type node struct {
	x, y, size, used, avail int
}

// parseNodes reads the df listing into a slice of nodes, skipping the two
// header lines.
func parseNodes(instr string) []node {
	var nodes []node
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "/dev/grid/") {
			continue
		}
		f := strings.Fields(line)
		name := f[0]
		// name like /dev/grid/node-x7-y17
		name = name[strings.Index(name, "node-")+5:]
		xy := strings.SplitN(name, "-", 2)
		x, _ := strconv.Atoi(strings.TrimPrefix(xy[0], "x"))
		y, _ := strconv.Atoi(strings.TrimPrefix(xy[1], "y"))
		size := atoiT(f[1])
		used := atoiT(f[2])
		avail := atoiT(f[3])
		nodes = append(nodes, node{x, y, size, used, avail})
	}
	return nodes
}

// atoiT parses a value like "92T" into its integer terabytes.
func atoiT(s string) int {
	n, _ := strconv.Atoi(strings.TrimSuffix(s, "T"))
	return n
}

// One counts viable pairs: ordered (A, B) where A is non-empty, A != B, and
// A's used data fits in B's available space.
func (e Exercise) One(instr string) (any, error) {
	nodes := parseNodes(instr)
	count := 0
	for i := range nodes {
		if nodes[i].used == 0 {
			continue
		}
		for j := range nodes {
			if i == j {
				continue
			}
			if nodes[i].used <= nodes[j].avail {
				count++
			}
		}
	}
	return count, nil
}

// Two returns the fewest steps to move the goal data from the top-right node to
// node (0,0). The grid is a sliding puzzle: one empty node shuffles data, while
// a band of over-full "wall" nodes is impassable.
//
// Cost = (steps to walk the empty node beside the goal) + 1 (slide goal one
// left) + 5 * (remaining leftward goal moves), since each subsequent leftward
// move cycles the empty node around the goal in 5 steps.
func (e Exercise) Two(instr string) (any, error) {
	nodes := parseNodes(instr)

	maxX := 0
	for _, n := range nodes {
		if n.x > maxX {
			maxX = n.x
		}
	}

	// Classify nodes: the empty node, and impassable walls (used larger than
	// any node's capacity, so their data can never move).
	var empty node
	wall := map[[2]int]bool{}
	nodeCap := map[[2]int]int{}
	minSize := 1 << 30
	for _, n := range nodes {
		nodeCap[[2]int{n.x, n.y}] = n.size
		if n.size < minSize {
			minSize = n.size
		}
		if n.used == 0 {
			empty = n
		}
	}
	for _, n := range nodes {
		if n.used > minSize {
			wall[[2]int{n.x, n.y}] = true
		}
	}

	goal := [2]int{maxX, 0}       // current goal-data location
	target := [2]int{maxX - 1, 0} // move empty here, just left of the goal

	walk := bfs([2]int{empty.x, empty.y}, target, wall, nodeCap)
	// 1 step to slide goal into the empty's spot, then 5 per remaining column.
	return walk + 1 + 5*(goal[0]-1), nil
}

// bfs returns the fewest moves for the empty node to travel from start to dst,
// stepping between adjacent in-grid nodes that are not walls.
func bfs(start, dst [2]int, wall map[[2]int]bool, nodeCap map[[2]int]int) int {
	type item struct {
		p [2]int
		d int
	}
	seen := map[[2]int]bool{start: true}
	queue := []item{{start, 0}}
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.p == dst {
			return cur.d
		}
		for _, d := range dirs {
			np := [2]int{cur.p[0] + d[0], cur.p[1] + d[1]}
			if _, ok := nodeCap[np]; !ok || wall[np] || seen[np] {
				continue
			}
			seen[np] = true
			queue = append(queue, item{np, cur.d + 1})
		}
	}
	return -1
}

// bfsPath returns the sequence of cells the empty node visits travelling from
// start to dst, inclusive of both ends, or nil if unreachable.
func bfsPath(start, dst [2]int, wall map[[2]int]bool, nodeCap map[[2]int]int) [][2]int {
	prev := map[[2]int][2]int{}
	seen := map[[2]int]bool{start: true}
	queue := [][2]int{start}
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == dst {
			var path [][2]int
			for p := dst; ; p = prev[p] {
				path = append([][2]int{p}, path...)
				if p == start {
					break
				}
			}
			return path
		}
		for _, d := range dirs {
			np := [2]int{cur[0] + d[0], cur[1] + d[1]}
			if _, ok := nodeCap[np]; !ok || wall[np] || seen[np] {
				continue
			}
			seen[np] = true
			prev[np] = cur
			queue = append(queue, np)
		}
	}
	return nil
}

// Vis renders the storage grid as a sliding puzzle: the impassable wall band,
// the empty node, the goal data, and the route the empty node walks to reach
// the cell beside the goal.
type visData struct {
	w, h    int
	empty   [2]int
	goal    [2]int
	wall    map[[2]int]bool
	pathSet map[[2]int]bool
}

func buildVisData(nodes []node) visData {
	maxX, maxY := 0, 0
	for _, n := range nodes {
		if n.x > maxX {
			maxX = n.x
		}
		if n.y > maxY {
			maxY = n.y
		}
	}

	var empty [2]int
	wall := map[[2]int]bool{}
	nodeCap := map[[2]int]int{}
	minSize := 1 << 30
	for _, n := range nodes {
		nodeCap[[2]int{n.x, n.y}] = n.size
		if n.size < minSize {
			minSize = n.size
		}
		if n.used == 0 {
			empty = [2]int{n.x, n.y}
		}
	}
	for _, n := range nodes {
		if n.used > minSize {
			wall[[2]int{n.x, n.y}] = true
		}
	}

	goal := [2]int{maxX, 0}
	target := [2]int{maxX - 1, 0}
	path := bfsPath(empty, target, wall, nodeCap)
	pathSet := map[[2]int]bool{}
	for _, p := range path {
		pathSet[p] = true
	}

	return visData{
		w: maxX + 1, h: maxY + 1,
		empty:   empty,
		goal:    goal,
		wall:    wall,
		pathSet: pathSet,
	}
}

// Vis renders the grid with walls, the empty-node path, and key positions.
func (e Exercise) Vis(instr, outdir string) error {
	vd := buildVisData(parseNodes(instr))

	const cell = 18
	const pad = 12
	img := image.NewRGBA(image.Rect(0, 0, vd.w*cell+2*pad, vd.h*cell+2*pad))

	bg := color.RGBA{0x12, 0x14, 0x20, 0xff}
	open := color.RGBA{0x26, 0x2c, 0x3e, 0xff}
	wallC := color.RGBA{0x6a, 0x32, 0x3a, 0xff}
	pathC := color.RGBA{0x2f, 0x5a, 0x86, 0xff}
	emptyC := color.RGBA{0x44, 0xff, 0x88, 0xff}
	goalC := color.RGBA{0xff, 0xc8, 0x4a, 0xff}
	originC := color.RGBA{0xff, 0x44, 0x55, 0xff}

	fill := func(x, y int, col color.RGBA) {
		x0, y0 := pad+x*cell, pad+y*cell
		for yy := y0; yy < y0+cell; yy++ {
			for xx := x0; xx < x0+cell; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	ih, iw := img.Bounds().Dy(), img.Bounds().Dx()
	for y := range ih {
		for x := range iw {
			img.SetRGBA(x, y, bg)
		}
	}

	for y := range vd.h {
		for x := range vd.w {
			p := [2]int{x, y}
			switch {
			case vd.wall[p]:
				fill(x, y, wallC)
			case vd.pathSet[p]:
				fill(x, y, pathC)
			default:
				fill(x, y, open)
			}
		}
	}
	fill(vd.empty[0], vd.empty[1], emptyC)
	fill(vd.goal[0], vd.goal[1], goalC)
	fill(0, 0, originC)

	f, err := os.Create(filepath.Join(outdir, "grid-computing.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
