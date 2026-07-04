package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"maps"
	"math"
	"os"
	"path/filepath"
)

// Vis renders an animated GIF of the robot painting the registration identifier
// (Part Two: start on white). Each frame shows the current hull state + robot position.
// Sampled to ~200 frames so the GIF is manageable.
//
//nolint:gocognit,funlen,nestif // robot painting animated visualization requires first-pass bounding box branch
func (e Exercise) Vis(instr, outdir string) error {
	prog := parseProgram(instr)
	steps := runRobotHistory(prog, 1)

	// Compute bounding box across all steps (use final state).
	final := steps[len(steps)-1].grid
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for p := range final {
		if first {
			minX, maxX, minY, maxY = p[0], p[0], p[1], p[1]
			first = false
		} else {
			if p[0] < minX {
				minX = p[0]
			}
			if p[0] > maxX {
				maxX = p[0]
			}
			if p[1] < minY {
				minY = p[1]
			}
			if p[1] > maxY {
				maxY = p[1]
			}
		}
	}
	// Expand bounding box a little for the robot marker.
	minX--
	minY--
	maxX++
	maxY++

	gridW := maxX - minX + 1
	gridH := maxY - minY + 1

	// Scale so the image is at least 500px wide.
	cellSize := max((500+gridW-1)/gridW, 6)
	const margin = 8
	w := gridW*cellSize + 2*margin
	h := gridH*cellSize + 2*margin

	// Palette: background, white panel, robot.
	palette := color.Palette{
		color.RGBA{0x11, 0x11, 0x11, 0xff}, // 0: background / black panel
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // 1: white panel
		color.RGBA{0xE6, 0x9F, 0x00, 0xff}, // 2: robot (Okabe-Ito orange)
	}

	// Sample frames — aim for ~150 frames spread evenly, always include last.
	const maxFrames = 150
	total := len(steps)
	stride := 1
	if total > maxFrames {
		stride = int(math.Ceil(float64(total) / float64(maxFrames)))
	}

	var frames []*image.Paletted
	var delays []int

	renderFrame := func(s robotStep, delay int) {
		img := image.NewPaletted(image.Rect(0, 0, w, h), palette)
		// Fill background.
		for y := range h {
			for x := range w {
				img.SetColorIndex(x, y, 0)
			}
		}
		// Draw panels.
		for p, col := range s.grid {
			if col != 1 {
				continue
			}
			gx := p[0] - minX
			gy := maxY - p[1]
			px0 := margin + gx*cellSize
			py0 := margin + gy*cellSize
			for dy := range cellSize {
				for dx := range cellSize {
					img.SetColorIndex(px0+dx, py0+dy, 1)
				}
			}
		}
		// Draw robot as a filled square in orange.
		gx := s.robotPos[0] - minX
		gy := maxY - s.robotPos[1]
		px0 := margin + gx*cellSize
		py0 := margin + gy*cellSize
		for dy := range cellSize {
			for dx := range cellSize {
				img.SetColorIndex(px0+dx, py0+dy, 2)
			}
		}
		frames = append(frames, img)
		delays = append(delays, delay)
	}

	for i := 0; i < total; i += stride {
		renderFrame(steps[i], 3) // 3/100 s per frame = ~33 fps
	}
	// Always include final frame with a longer hold.
	renderFrame(steps[total-1], 200)

	f, err := os.Create(filepath.Join(outdir, "vis.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, &gif.GIF{
		Image: frames,
		Delay: delays,
	})
}

// robotStep records one moment in the robot's execution.
type robotStep struct {
	grid     map[pos]int
	robotPos pos
}

// runRobotHistory runs the Part Two robot and returns one step per paint action,
// recording the grid state and robot position at each step.
func runRobotHistory(program mem, startColor int) []robotStep {
	in := make(chan int, 1)
	out := make(chan int)
	go intcode(program, in, out)

	grid := make(map[pos]int)
	grid[pos{0, 0}] = startColor

	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	dir := 0
	cur := pos{0, 0}

	var steps []robotStep

	// Record initial state.
	snap := make(map[pos]int, len(grid))
	maps.Copy(snap, grid)
	steps = append(steps, robotStep{grid: snap, robotPos: cur})

	for {
		in <- grid[cur]

		paintColor, ok := <-out
		if !ok {
			break
		}
		turn, ok := <-out
		if !ok {
			break
		}

		grid[cur] = paintColor

		if turn == 0 {
			dir = (dir + 3) % 4
		} else {
			dir = (dir + 1) % 4
		}
		cur[0] += dirs[dir][0]
		cur[1] += dirs[dir][1]

		// Snapshot grid (shallow copy is fine — map values are ints).
		snapCopy := make(map[pos]int, len(grid))
		maps.Copy(snapCopy, grid)
		steps = append(steps, robotStep{grid: snapCopy, robotPos: cur})
	}

	return steps
}
