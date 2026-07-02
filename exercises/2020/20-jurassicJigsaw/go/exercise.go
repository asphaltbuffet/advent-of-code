package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 20.
type Exercise struct {
	common.BaseExercise
}

// tile is one image tile: its ID and its square grid of rows.
type tile struct {
	id   int
	rows []string
}

func parseTiles(instr string) ([]tile, error) {
	var tiles []tile
	for _, block := range strings.Split(strings.TrimSpace(instr), "\n\n") {
		lines := strings.Split(block, "\n")
		var id int
		if _, err := fmt.Sscanf(lines[0], "Tile %d:", &id); err != nil {
			return nil, fmt.Errorf("parsing tile header %q: %w", lines[0], err)
		}
		tiles = append(tiles, tile{id: id, rows: lines[1:]})
	}
	return tiles, nil
}

// edges returns a tile's four borders (top, bottom, left, right) as read left to
// right / top to bottom.
func (t tile) edges() [4]string {
	n := len(t.rows)
	top := t.rows[0]
	bottom := t.rows[n-1]
	var left, right strings.Builder
	for _, r := range t.rows {
		left.WriteByte(r[0])
		right.WriteByte(r[len(r)-1])
	}
	return [4]string{top, bottom, left.String(), right.String()}
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// canonicalEdge normalizes an edge so a border and its flip compare equal.
func canonicalEdge(e string) string {
	if r := reverse(e); r < e {
		return r
	}
	return e
}

// edgeCounts tallies how many tiles carry each canonical edge.
func edgeCounts(tiles []tile) map[string]int {
	counts := map[string]int{}
	for _, t := range tiles {
		for _, e := range t.edges() {
			counts[canonicalEdge(e)]++
		}
	}
	return counts
}

// corners returns the tiles that have exactly two edges not shared with any
// other tile (the four corners of the assembled image).
func corners(tiles []tile) []tile {
	counts := edgeCounts(tiles)
	var cs []tile
	for _, t := range tiles {
		unmatched := 0
		for _, e := range t.edges() {
			if counts[canonicalEdge(e)] == 1 {
				unmatched++
			}
		}
		if unmatched == 2 {
			cs = append(cs, t)
		}
	}
	return cs
}

// One multiplies the four corner tile IDs.
func (e Exercise) One(instr string) (any, error) {
	tiles, err := parseTiles(instr)
	if err != nil {
		return nil, err
	}

	product := 1
	cs := corners(tiles)
	for _, c := range cs {
		product *= c.id
	}
	if len(cs) != 4 {
		return nil, fmt.Errorf("expected 4 corners, found %d", len(cs))
	}

	return strconv.Itoa(product), nil
}

// rotate turns a grid 90 degrees clockwise.
func rotate(g []string) []string {
	n := len(g)
	out := make([]string, n)
	for i := 0; i < n; i++ {
		var b strings.Builder
		for j := n - 1; j >= 0; j-- {
			b.WriteByte(g[j][i])
		}
		out[i] = b.String()
	}
	return out
}

// flip mirrors a grid horizontally.
func flip(g []string) []string {
	out := make([]string, len(g))
	for i, r := range g {
		out[i] = reverse(r)
	}
	return out
}

// orientations returns all 8 rotations/flips of a grid.
func orientations(g []string) [][]string {
	var os [][]string
	cur := g
	for i := 0; i < 4; i++ {
		os = append(os, cur, flip(cur))
		cur = rotate(cur)
	}
	return os
}

func topEdge(g []string) string { return g[0] }
func leftEdge(g []string) string {
	var b strings.Builder
	for _, r := range g {
		b.WriteByte(r[0])
	}
	return b.String()
}

// assemble places tiles into a size x size grid of oriented tile bodies, matching
// shared edges. Returns nil if no arrangement is found.
func assemble(tiles []tile) [][][]string {
	size := 1
	for size*size < len(tiles) {
		size++
	}

	counts := edgeCounts(tiles)
	// Pick a corner tile to seed the top-left, oriented so its two unmatched
	// edges face up and left.
	var start tile
	for _, t := range tiles {
		unmatched := 0
		for _, ed := range t.edges() {
			if counts[canonicalEdge(ed)] == 1 {
				unmatched++
			}
		}
		if unmatched == 2 {
			start = t
			break
		}
	}

	grid := make([][][]string, size)
	for i := range grid {
		grid[i] = make([][]string, size)
	}
	used := map[int]bool{}

	// Orient the seed corner so its top and left edges are outer borders.
	for _, o := range orientations(start.rows) {
		if counts[canonicalEdge(topEdge(o))] == 1 && counts[canonicalEdge(leftEdge(o))] == 1 {
			grid[0][0] = o
			used[start.id] = true
			break
		}
	}

	byID := map[int]tile{}
	for _, t := range tiles {
		byID[t.id] = t
	}

	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if r == 0 && c == 0 {
				continue
			}
			placed := false
			for _, t := range tiles {
				if used[t.id] {
					continue
				}
				for _, o := range orientations(t.rows) {
					ok := true
					if c > 0 && leftEdge(o) != rightEdgeOf(grid[r][c-1]) {
						ok = false
					}
					if r > 0 && topEdge(o) != bottomEdgeOf(grid[r-1][c]) {
						ok = false
					}
					if ok {
						grid[r][c] = o
						used[t.id] = true
						placed = true
						break
					}
				}
				if placed {
					break
				}
			}
			if !placed {
				return nil
			}
		}
	}
	return grid
}

func rightEdgeOf(g []string) string {
	var b strings.Builder
	for _, r := range g {
		b.WriteByte(r[len(r)-1])
	}
	return b.String()
}
func bottomEdgeOf(g []string) string { return g[len(g)-1] }

// stitch trims each tile's border and joins the bodies into one image.
func stitch(grid [][][]string) []string {
	size := len(grid)
	tileN := len(grid[0][0])
	body := tileN - 2

	var image []string
	for r := 0; r < size; r++ {
		for br := 0; br < body; br++ {
			var b strings.Builder
			for c := 0; c < size; c++ {
				b.WriteString(grid[r][c][br+1][1 : tileN-1])
			}
			image = append(image, b.String())
		}
	}
	return image
}

// seaMonster is the pattern to search for; '#' cells must be set.
var seaMonster = []string{
	"                  # ",
	"#    ##    ##    ###",
	" #  #  #  #  #  #   ",
}

// countMonsters returns how many non-overlapping sea monsters appear in image.
func countMonsters(image []string) int {
	var offsets [][2]int
	for r, row := range seaMonster {
		for c, ch := range row {
			if ch == '#' {
				offsets = append(offsets, [2]int{r, c})
			}
		}
	}
	mh, mw := len(seaMonster), len(seaMonster[0])

	found := 0
	for r := 0; r+mh <= len(image); r++ {
		for c := 0; c+mw <= len(image[0]); c++ {
			all := true
			for _, o := range offsets {
				if image[r+o[0]][c+o[1]] != '#' {
					all = false
					break
				}
			}
			if all {
				found++
			}
		}
	}
	return found
}

// Two assembles the image and returns the water roughness: the number of '#'
// cells not part of any sea monster.
func (e Exercise) Two(instr string) (any, error) {
	tiles, err := parseTiles(instr)
	if err != nil {
		return nil, err
	}

	grid := assemble(tiles)
	if grid == nil {
		return nil, fmt.Errorf("could not assemble image")
	}
	image := stitch(grid)

	total := 0
	for _, row := range image {
		total += strings.Count(row, "#")
	}

	monsterCells := 0
	for _, row := range seaMonster {
		monsterCells += strings.Count(row, "#")
	}

	for _, o := range orientations(image) {
		if n := countMonsters(o); n > 0 {
			return strconv.Itoa(total - n*monsterCells), nil
		}
	}

	return strconv.Itoa(total), nil
}
