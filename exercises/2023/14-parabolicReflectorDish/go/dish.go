package exercises

import (
	"crypto/md5" //nolint:gosec // not for security
	"encoding/hex"
	"fmt"
	"strings"
)

type Dish struct {
	Rocks [][]byte
}

func parseInput(input string) (*Dish, error) {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no lines in input")
	}

	r := make([][]byte, len(lines))

	for y, line := range lines {
		r[y] = []byte(line)
	}

	return &Dish{Rocks: r}, nil
}

func (d *Dish) tiltNorth() {
	d.transpose()
	d.settle()

	d.transpose()
}

// rotate CCW 4x sorting each time
func (d *Dish) spin() {
	for i := 0; i < 4; i++ {
		d.settle()
		d.RotateCCW()
	}
}

func (d *Dish) calcLoad() int {
	var load int
	for y, row := range d.Rocks {
		// fmt.Printf("%s\n", row)
		for _, rock := range row {
			if rock == 'O' {
				load += len(row) - y
			}
		}
	}

	return load
}

func (d *Dish) transpose() {
	t := make([][]byte, 0, len(d.Rocks[0]))

	for x := 0; x < len(d.Rocks[0]); x++ {
		tt := make([]byte, len(d.Rocks))

		for y := 0; y < len(d.Rocks); y++ {
			tt[y] = d.Rocks[y][x]
		}

		t = append(t, tt)
	}

	d.Rocks = t
}

func (d *Dish) RotateCCW() {
	r, c := len(d.Rocks), len(d.Rocks[0])
	rot := make([][]byte, c)

	for i := 0; i < c; i++ {
		rot[i] = make([]byte, r)

		for j := 0; j < r; j++ {
			rot[i][j] = d.Rocks[j][c-i-1]
		}
	}

	d.Rocks = rot
}

func (d *Dish) settle() {
	for _, row := range d.Rocks {
		// sortRocks(row, cmpRock)
		countSort(row)
	}
}

func cmpRock(a, b byte) int {
	if a == '#' || b == '#' {
		return 0
	}

	if a == 'O' && b == '.' {
		return -1
	}

	if a == '.' && b == 'O' {
		return 1
	}

	return 0
}

func sortRocks[E any](data []E, cmp func(a, b E) int) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && (cmp(data[j], data[j-1]) < 0); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func countSort(r []byte) {
	start := 0
	// Process each segment between '#' characters
	for start < len(r) {
		end := start
		// Find the end of the current segment
		for end < len(r) && r[end] != '#' {
			end++
		}

		// Count 'O's in the current segment
		countO := 0
		for i := start; i < end; i++ {
			if r[i] == 'O' {
				countO++
			}
		}

		// Place 'O's and then '.'s in the current segment
		for i := start; i < end; i++ {
			if countO > 0 {
				r[i] = 'O'
				countO--
			} else {
				r[i] = '.'
			}
		}

		// Move to the next segment
		start = end + 1
	}
}

func (d *Dish) hash() string {
	h := md5.New() //nolint:gosec // not for security
	for _, row := range d.Rocks {
		h.Write(row)
	}

	return hex.EncodeToString(h.Sum(nil))
}
