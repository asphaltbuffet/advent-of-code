package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
)

// Vis animates the octopus grid step by step (GIF). Each cell is shaded by its
// energy level on a colorblind-safe dark-to-bright ramp; an octopus that just
// flashed is drawn white so the ripples of a cascade stand out. The animation
// runs up to the first step where every octopus flashes at once — the part-two
// answer — and holds on that fully-white synchronized frame. Energy is encoded
// by brightness, so the buildup reads in grayscale too.
func (e Exercise) Vis(instr, outdir string) error {
	g, err := parse(instr)
	if err != nil {
		return err
	}

	all := g.rows * g.cols

	// Palette: index 0..9 energy ramp (dark to warm), index 10 = just flashed.
	pal := color.Palette{}
	for lvl := 0; lvl <= 9; lvl++ {
		f := float64(lvl) / 9
		pal = append(pal, color.RGBA{
			uint8(0x14 + 0xb0*f),
			uint8(0x18 + 0x80*f),
			uint8(0x28 + 0x20*f),
			0xff,
		})
	}
	pal = append(pal, color.RGBA{0xff, 0xff, 0xff, 0xff}) // flashed
	pal = append(pal, color.RGBA{0x0d, 0x0f, 0x18, 0xff}) // background
	const bgIdx = 11

	const cell = 22
	const pad = 8
	W := g.cols*cell + 2*pad
	H := g.rows*cell + 2*pad

	anim := &gif.GIF{}

	// Render the current grid state; flashedNow marks cells that flashed this step.
	render := func(flashedNow [][]bool) *image.Paletted {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		for i := range img.Pix {
			img.Pix[i] = bgIdx
		}
		for r := 0; r < g.rows; r++ {
			for c := 0; c < g.cols; c++ {
				idx := uint8(g.e[r][c])
				if flashedNow != nil && flashedNow[r][c] {
					idx = 10
				}
				for dy := 1; dy < cell-1; dy++ {
					for dx := 1; dx < cell-1; dx++ {
						img.SetColorIndex(pad+c*cell+dx, pad+r*cell+dy, idx)
					}
				}
			}
		}
		return img
	}

	// Initial frame.
	anim.Image = append(anim.Image, render(nil))
	anim.Delay = append(anim.Delay, 60)

	const maxSteps = 400
	for s := 1; s <= maxSteps; s++ {
		flashedNow := g.stepTracked()
		anim.Image = append(anim.Image, render(flashedNow))

		delay := 12
		if countTrue(flashedNow) == all {
			delay = 300 // linger on the synchronized flash
			anim.Delay = append(anim.Delay, delay)
			break
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "dumbo-octopus.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

// stepTracked runs one step like step() but also returns which cells flashed, so
// the animation can highlight them.
func (g *grid) stepTracked() [][]bool {
	flashed := make([][]bool, g.rows)
	for r := range flashed {
		flashed[r] = make([]bool, g.cols)
	}

	var queue [][2]int
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			g.e[r][c]++
			if g.e[r][c] > 9 {
				queue = append(queue, [2]int{r, c})
				flashed[r][c] = true
			}
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range neighbors8 {
			nr, nc := cur[0]+d[0], cur[1]+d[1]
			if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
				continue
			}
			g.e[nr][nc]++
			if g.e[nr][nc] > 9 && !flashed[nr][nc] {
				flashed[nr][nc] = true
				queue = append(queue, [2]int{nr, nc})
			}
		}
	}

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			if flashed[r][c] {
				g.e[r][c] = 0
			}
		}
	}

	return flashed
}

func countTrue(m [][]bool) int {
	n := 0
	for _, row := range m {
		for _, v := range row {
			if v {
				n++
			}
		}
	}
	return n
}
