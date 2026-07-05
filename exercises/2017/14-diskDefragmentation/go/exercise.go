package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 14.
type Exercise struct {
	common.BaseExercise
}

// knotHash returns the 32-character hex Knot Hash of the input (day 10). It is
// duplicated here because each exercise builds as its own binary.
func knotHash(input string) string {
	lengths := make([]int, 0, len(input)+5)
	for _, b := range []byte(input) {
		lengths = append(lengths, int(b))
	}
	lengths = append(lengths, 17, 31, 73, 47, 23)

	list := make([]int, 256)
	for i := range list {
		list[i] = i
	}
	pos, skip := 0, 0
	for range 64 {
		for _, l := range lengths {
			for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
				a, b := (pos+i)%256, (pos+j)%256
				list[a], list[b] = list[b], list[a]
			}
			pos = (pos + l + skip) % 256
			skip++
		}
	}

	var sb strings.Builder
	for block := range 16 {
		x := 0
		for i := range 16 {
			x ^= list[block*16+i]
		}
		fmt.Fprintf(&sb, "%02x", x)
	}
	return sb.String()
}

// buildGrid returns the 128x128 used/free grid for the given key. grid[r][c] is
// true when the square is used (a 1 bit in row r's Knot Hash).
func buildGrid(key string) [128][128]bool {
	var grid [128][128]bool
	for r := range 128 {
		hash := knotHash(fmt.Sprintf("%s-%d", key, r))
		for i, ch := range hash {
			var nib int
			switch {
			case ch >= '0' && ch <= '9':
				nib = int(ch - '0')
			default:
				nib = int(ch-'a') + 10
			}
			for bit := range 4 {
				if nib&(1<<(3-bit)) != 0 {
					grid[r][i*4+bit] = true
				}
			}
		}
	}
	return grid
}

// One returns the total number of used squares.
func (e Exercise) One(instr string) (any, error) {
	grid := buildGrid(strings.TrimSpace(instr))
	used := 0
	for r := range 128 {
		for c := range 128 {
			if grid[r][c] {
				used++
			}
		}
	}
	return used, nil
}

// Two returns the number of connected regions of used squares.
func (e Exercise) Two(instr string) (any, error) {
	grid := buildGrid(strings.TrimSpace(instr))
	var seen [128][128]bool

	var flood func(r, c int)
	flood = func(r, c int) {
		if r < 0 || r >= 128 || c < 0 || c >= 128 || seen[r][c] || !grid[r][c] {
			return
		}
		seen[r][c] = true
		flood(r-1, c)
		flood(r+1, c)
		flood(r, c-1)
		flood(r, c+1)
	}

	regions := 0
	for r := range 128 {
		for c := range 128 {
			if grid[r][c] && !seen[r][c] {
				regions++
				flood(r, c)
			}
		}
	}
	return regions, nil
}

// labelRegions flood-fills the grid, tagging each used square with its region
// index (1-based; 0 means free). Returns the labels and the region count.
func labelRegions(grid [128][128]bool) ([128][128]int, int) {
	var label [128][128]int

	var flood func(r, c, id int)
	flood = func(r, c, id int) {
		if r < 0 || r >= 128 || c < 0 || c >= 128 || !grid[r][c] || label[r][c] != 0 {
			return
		}
		label[r][c] = id
		flood(r-1, c, id)
		flood(r+1, c, id)
		flood(r, c-1, id)
		flood(r, c+1, id)
	}

	regions := 0
	for r := range 128 {
		for c := range 128 {
			if grid[r][c] && label[r][c] == 0 {
				regions++
				flood(r, c, regions)
			}
		}
	}
	return label, regions
}

// hsv converts an HSV triple (h in degrees, s,v in 0..1) to an RGBA colour.
