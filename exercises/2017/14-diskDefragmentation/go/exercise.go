package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
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
	for r := 0; r < 64; r++ {
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
	for block := 0; block < 16; block++ {
		x := 0
		for i := 0; i < 16; i++ {
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
	for r := 0; r < 128; r++ {
		hash := knotHash(fmt.Sprintf("%s-%d", key, r))
		for i, ch := range hash {
			var nib int
			switch {
			case ch >= '0' && ch <= '9':
				nib = int(ch - '0')
			default:
				nib = int(ch-'a') + 10
			}
			for bit := 0; bit < 4; bit++ {
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
	for r := 0; r < 128; r++ {
		for c := 0; c < 128; c++ {
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
	for r := 0; r < 128; r++ {
		for c := 0; c < 128; c++ {
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
	for r := 0; r < 128; r++ {
		for c := 0; c < 128; c++ {
			if grid[r][c] && label[r][c] == 0 {
				regions++
				flood(r, c, regions)
			}
		}
	}
	return label, regions
}

// hsv converts an HSV triple (h in degrees, s,v in 0..1) to an RGBA colour.
func hsv(h, s, v float64) color.RGBA {
	c := v * s
	hp := h / 60
	x := c * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = c, x, 0
	case hp < 2:
		r, g, b = x, c, 0
	case hp < 3:
		r, g, b = 0, c, x
	case hp < 4:
		r, g, b = 0, x, c
	case hp < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	m := v - c
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}

// Vis renders the 128x128 disk grid, colouring each connected region a distinct
// hue over a dark background of free squares.
func (e Exercise) Vis(instr, outdir string) error {
	grid := buildGrid(strings.TrimSpace(instr))
	label, regions := labelRegions(grid)

	// Precompute a colour per region; a golden-angle hue spread keeps adjacent
	// region ids visually distinct.
	regionColor := make([]color.RGBA, regions+1)
	for id := 1; id <= regions; id++ {
		hue := math.Mod(float64(id)*137.508, 360)
		regionColor[id] = hsv(hue, 0.62, 0.95)
	}

	const cell = 6
	const pad = 8
	size := 128*cell + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	free := color.RGBA{0x18, 0x1c, 0x2a, 0xff}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	for r := 0; r < 128; r++ {
		for c := 0; c < 128; c++ {
			col := free
			if id := label[r][c]; id != 0 {
				col = regionColor[id]
			}
			x0, y0 := pad+c*cell, pad+r*cell
			for yy := y0; yy < y0+cell; yy++ {
				for xx := x0; xx < x0+cell; xx++ {
					img.SetRGBA(xx, yy, col)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "disk-defragmentation.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
