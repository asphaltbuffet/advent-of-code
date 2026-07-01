package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 14.
type Exercise struct {
	common.BaseExercise
}

// Vis animates the spin cycle (GIF): starting from the parsed dish, each frame
// applies one full N→W→S→E spin, so the round rocks (gold) migrate around the
// fixed cube rocks (slate) and settle into the repeating steady state.
func (e Exercise) Vis(input, outdir string) error {
	d, err := parseInput(input)
	if err != nil {
		return err
	}

	const frames = 48
	snapshot := func() [][]byte {
		cp := make([][]byte, len(d.Rocks))
		for i, row := range d.Rocks {
			cp[i] = append([]byte(nil), row...)
		}
		return cp
	}

	grids := [][][]byte{snapshot()}
	for f := 0; f < frames; f++ {
		d.spin() // one full cycle returns the grid to its original orientation
		grids = append(grids, snapshot())
	}

	h := len(grids[0])
	w := len(grids[0][0])
	const cell = 6
	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 empty
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 1 round rock 'O'
		color.RGBA{0x3a, 0x44, 0x5e, 0xff}, // 2 cube rock '#'
	}

	anim := &gif.GIF{}
	for fi, g := range grids {
		img := image.NewPaletted(image.Rect(0, 0, w*cell, h*cell), pal)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var idx uint8
				switch g[y][x] {
				case 'O':
					idx = 1
				case '#':
					idx = 2
				}
				if idx == 0 {
					continue
				}
				for dy := 0; dy < cell; dy++ {
					for dx := 0; dx < cell; dx++ {
						img.SetColorIndex(x*cell+dx, y*cell+dy, idx)
					}
				}
			}
		}
		anim.Image = append(anim.Image, img)
		delay := 12
		if fi == 0 || fi == len(grids)-1 {
			delay = 150
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "reflector-dish.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

// One returns the answer to the first part of the exercise.
// not: 166856 (too high)
func (e Exercise) One(input string) (any, error) {
	data, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	data.tiltNorth()

	return data.calcLoad(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	const cycles = 1_000_000_000

	memos := make(map[string]int)

	data, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	data.transpose()

	for i := 0; i < cycles; i++ {
		data.spin()

		hash := data.hash()
		if v, ok := memos[hash]; ok {
			// fmt.Printf("hash: %s seen before at i=%d and again at i=%d\n", hash, v, i)
			// fmt.Printf("cycle length: %d\n", i-v)

			// fmt.Printf("skipping %d cycles\n", (cycles - i) / (i - v))
			i += ((cycles - i) / (i - v)) * (i - v)

			// reset memos
			memos = map[string]int{}
		}

		memos[hash] = i
	}

	data.transpose()

	return data.calcLoad(), nil
}
