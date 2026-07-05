package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

// Vis renders the spiral memory heat-map as a PNG.
func (e Exercise) Vis(instr, outdir string) error {
	target := parseTarget(instr)

	const half = 13 // window spans -half..+half on each axis
	vals, maxV, answer, foundAnswer := buildSpiralVals(target, half)

	const cell = 30
	const pad = 16
	img := renderSpiralImage(vals, maxV, answer, foundAnswer, half, cell, pad)

	f, err := os.Create(filepath.Join(outdir, "spiral-memory.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func buildSpiralVals(target, half int) (map[coord]int, int, coord, bool) {
	limit := (2*half + 1) * (2*half + 1)
	vals := map[coord]int{}
	first := true
	maxV := 1
	var answer coord
	foundAnswer := false
	count := 0
	spiralWalk(func(c coord) bool {
		if c.x < -half || c.x > half || c.y < -half || c.y > half {
			return true
		}
		count++
		if first {
			vals[c] = 1
			first = false
			return count < limit
		}
		sum := neighborSum(vals, c)
		vals[c] = sum
		if sum > maxV {
			maxV = sum
		}
		if !foundAnswer && sum > target {
			answer = c
			foundAnswer = true
		}
		return count < limit
	})
	return vals, maxV, answer, foundAnswer
}

func neighborSum(vals map[coord]int, c coord) int {
	sum := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			sum += vals[coord{c.x + dx, c.y + dy}]
		}
	}
	return sum
}

func renderSpiralImage(vals map[coord]int, maxV int, answer coord, foundAnswer bool, half, cell, pad int) *image.RGBA {
	minX, maxX, minY, maxY := -half, half, -half, half
	cols, rows := maxX-minX+1, maxY-minY+1
	img := image.NewRGBA(image.Rect(0, 0, cols*cell+2*pad, rows*cell+2*pad))

	bg := color.RGBA{0x10, 0x12, 0x1c, 0xff}
	for y := range img.Bounds().Dy() {
		for x := range img.Bounds().Dx() {
			img.SetRGBA(x, y, bg)
		}
	}

	heat := func(t float64) color.RGBA {
		switch {
		case t < 0.5:
			u := t / 0.5
			return color.RGBA{uint8(0x1a + u*0x15), uint8(0x2c + u*0x5e), uint8(0x50 + u*0x36), 0xff}
		default:
			u := (t - 0.5) / 0.5
			return color.RGBA{uint8(0x2f + u*0xd0), uint8(0x8a + u*0x3e), uint8(0x86 - u*0x3c), 0xff}
		}
	}

	fill := func(c coord, col color.RGBA) {
		gx, gy := c.x-minX, maxY-c.y
		x0, y0 := pad+gx*cell, pad+gy*cell
		for yy := y0; yy < y0+cell; yy++ {
			for xx := x0; xx < x0+cell; xx++ {
				img.SetRGBA(xx, yy, col)
			}
		}
	}

	logMax := math.Log(float64(maxV) + 1)
	for c, v := range vals {
		t := math.Log(float64(v)+1) / logMax
		fill(c, heat(t))
	}
	fill(coord{0, 0}, color.RGBA{0xff, 0x44, 0x55, 0xff})
	if foundAnswer {
		fill(answer, color.RGBA{0xff, 0xff, 0xff, 0xff})
	}
	return img
}
