package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 21.
type Exercise struct {
	common.BaseExercise
}

// Vis animates the reachable frontier spreading across the garden (GIF). Each
// frame is the set of plots reachable in exactly that many steps — the wave the
// step counter tracks — over the rocks (slate) with the start plot marked. The
// alternating parity produces the characteristic growing diamond.
func (e Exercise) Vis(input, outdir string) error {
	tiles, start, dim, err := parseInput(input)
	if err != nil {
		return err
	}

	const steps = 64
	q := []Point{start}
	frontiers := [][]Point{append([]Point(nil), q...)}
	for s := 0; s < steps; s++ {
		q = walk(q, tiles, dim)
		frontiers = append(frontiers, append([]Point(nil), q...))
	}

	const cell = 4
	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 garden
		color.RGBA{0x2a, 0x30, 0x44, 0xff}, // 1 rock
		color.RGBA{0x2f, 0xd0, 0x9a, 0xff}, // 2 reachable frontier
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 3 start
	}

	anim := &gif.GIF{}
	for fi, front := range frontiers {
		img := image.NewPaletted(image.Rect(0, 0, dim*cell, dim*cell), pal)
		// Base layer: rocks.
		for y := 0; y < dim; y++ {
			for x := 0; x < dim; x++ {
				if tiles[Point{x, y}] == Rock {
					fillCellIdx(img, x, y, cell, 1)
				}
			}
		}
		// Frontier on top.
		for _, p := range front {
			fillCellIdx(img, p.X, p.Y, cell, 2)
		}
		fillCellIdx(img, start.X, start.Y, cell, 3)

		anim.Image = append(anim.Image, img)
		delay := 8
		if fi == 0 || fi == len(frontiers)-1 {
			delay = 150
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "step-counter.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

func fillCellIdx(img *image.Paletted, gx, gy, cell int, idx uint8) {
	for dy := 0; dy < cell; dy++ {
		for dx := 0; dx < cell; dx++ {
			img.SetColorIndex(gx*cell+dx, gy*cell+dy, idx)
		}
	}
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	tiles, start, dim, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	q := []Point{start}

	for steps := 0; steps < 64; steps++ {
		q = walk(q, tiles, dim)
	}

	return len(q), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(input string) (any, error) {
	tiles, start, dim, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	q := []Point{start}

	const stepGoal int = 26_501_365

	poly := make([]int, 0, 3)
	stepFactor := stepGoal % dim

	for steps := 0; ; {
		q = walk(q, tiles, dim)

		steps++

		// build up our polynomial data until we have enough to solve for the answer
		if steps%dim == stepFactor {
			// fmt.Printf("step %3d: polynomial[%d] = %d\n", steps, len(poly), vCount)

			poly = append(poly, len(q)) // this should be appending values for steps 65, 196, and 327

			if len(poly) == 3 {
				// more efficient to check if we have enough data to solve for the answer here
				// rather than checking every iteration
				break
			}
		}
	}

	a := (poly[2] + poly[0] - 2*poly[1]) / 2
	b := poly[1] - poly[0] - a
	c := poly[0]
	n := stepGoal / dim

	// fmt.Printf("%d*%d*%d + %d*%d + %d\n", a, n, n, b, n, c)

	return a*n*n + b*n + c, nil
}

func walk(q []Point, tiles map[Point]TileType, dim int) []Point {
	// fq := list.New() // build up a queue for the next iteration
	visited := make(map[Point]bool, len(q)*4)
	newQueue := make([]Point, 0, len(q)*4)

	for _, pt := range q {
		for _, d := range Moves {
			np := Point{pt.X + d.X, pt.Y + d.Y}

			tmp := Point{((np.Y % dim) + dim) % dim, ((np.X % dim) + dim) % dim}

			if tiles[tmp] != Rock {
				if _, ok := visited[np]; !ok {
					visited[np] = true
					newQueue = append(newQueue, np)
				}
			}
		}
	}

	return newQueue
}

//nolint:unused // this is kept around for debugging purposes
func generateExample(input string) {
	lines := strings.Split(input, "\n")

	// generate a minimized map of the garden for testing
	for i := 0; i < 2; i++ {
		width := len(lines[0])
		height := len(lines)
		mm := []string{}

		for row := 0; row < height; row++ {
			sb := strings.Builder{}
			if row%2 == 0 && row != 0 && row != height-1 { // skip half but keep first, last, and middle
				continue
			}

			for col := 0; col < width; col++ {
				if col%2 == 0 && col != 0 && col != width-1 { // skip half but keep first, last, and middle
					continue
				}

				c := rune(lines[row][col])

				sb.WriteRune(c)
			}

			mm = append(mm, sb.String())
			fmt.Println(sb.String())
		}

		lines = mm
	}
}
