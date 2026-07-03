package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// job records one step's placement on the worker schedule.
type job struct {
	step         byte
	worker       int
	start, end   int
}

// schedule replays the Part Two simulation, recording which worker ran each step
// and when, so the timeline can be drawn. It mirrors Two exactly.
func schedule(instr string) ([]job, int, int, error) {
	edges, err := parse(instr)
	if err != nil {
		return nil, 0, 0, err
	}
	d := deps(edges)

	workers, base := 5, 60
	if len(d) <= 6 {
		workers, base = 2, 0
	}

	done := map[byte]bool{}
	inProgress := map[byte]int{} // step -> finish second
	assignedWorker := map[byte]int{}
	free := make([]bool, workers)
	for i := range free {
		free[i] = true
	}

	var jobs []job

	for t := 0; ; t++ {
		for step, finish := range inProgress {
			if finish <= t {
				done[step] = true
				free[assignedWorker[step]] = true
				delete(inProgress, step)
			}
		}

		if len(done) == len(d) {
			return jobs, t, workers, nil
		}

		var ready []byte
		for step, prereqs := range d {
			if done[step] || inProgress[step] != 0 {
				continue
			}
			blocked := false
			for p := range prereqs {
				if !done[p] {
					blocked = true
					break
				}
			}
			if !blocked {
				ready = append(ready, step)
			}
		}
		sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })

		for _, step := range ready {
			w := -1
			for i := 0; i < workers; i++ {
				if free[i] {
					w = i
					break
				}
			}
			if w == -1 {
				break
			}
			end := t + base + int(step-'A'+1)
			inProgress[step] = end
			assignedWorker[step] = w
			free[w] = false
			jobs = append(jobs, job{step: step, worker: w, start: t, end: end})
		}
	}
}

// Vis renders the Part Two worker schedule as a Gantt chart: one row per worker,
// time running left to right, each step drawn as a labeled bar from its start to
// its finish. Bars alternate between two colorblind-safe tints (Okabe-Ito blue
// and orange) that also differ in brightness, and idle gaps stay dark, so the
// utilization and dependency stalls that make the total 914 seconds read clearly
// even in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	jobs, total, workers, err := schedule(instr)
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return fmt.Errorf("no steps to visualize")
	}

	const (
		leftPad  = 40 // room for worker labels
		topPad   = 24 // room for title
		botPad   = 20 // room for the time axis
		rowH     = 26
		rowGap   = 6
		plotW    = 900
	)
	plotH := workers*rowH + (workers-1)*rowGap
	W := leftPad + plotW + 12
	H := topPad + plotH + botPad

	img := image.NewRGBA(image.Rect(0, 0, W, H))

	bg := color.RGBA{0x0c, 0x0f, 0x14, 0xff}
	lane := color.RGBA{0x18, 0x1d, 0x25, 0xff} // idle lane background
	axis := color.RGBA{0x55, 0x5c, 0x66, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	// Okabe-Ito blue and orange: distinct in hue AND brightness for grayscale.
	barA := color.RGBA{0x56, 0xB4, 0xE9, 0xff} // sky blue
	barB := color.RGBA{0xE6, 0x9F, 0x00, 0xff} // orange
	inkA := color.RGBA{0x06, 0x1a, 0x26, 0xff} // dark label on blue
	inkB := color.RGBA{0x1a, 0x12, 0x00, 0xff} // dark label on orange

	fill := func(x0, y0, x1, y1 int, c color.RGBA) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if x >= 0 && y >= 0 && x < W && y < H {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}

	fill(0, 0, W, H, bg)

	xOf := func(t int) int { return leftPad + t*plotW/total }
	rowY := func(w int) int { return topPad + w*(rowH+rowGap) }

	// Idle lane backgrounds and worker labels.
	for w := 0; w < workers; w++ {
		y := rowY(w)
		fill(leftPad, y, leftPad+plotW, y+rowH, lane)
		drawText7(img, fmt.Sprintf("W%d", w+1), 8, y+rowH/2+4, white)
	}

	// Step bars.
	for _, j := range jobs {
		y := rowY(j.worker)
		x0, x1 := xOf(j.start), xOf(j.end)
		if x1 <= x0 {
			x1 = x0 + 1
		}
		bar, ink := barA, inkA
		if (j.step-'A')%2 == 1 {
			bar, ink = barB, inkB
		}
		fill(x0, y+2, x1, y+rowH-2, bar)
		// Label the step letter if the bar is wide enough to hold it.
		if x1-x0 >= 8 {
			drawText7(img, string(j.step), x0+3, y+rowH/2+4, ink)
		}
	}

	// Time axis with a few ticks.
	ay := topPad + plotH + 6
	fill(leftPad, ay, leftPad+plotW, ay+1, axis)
	for k := 0; k <= 4; k++ {
		t := total * k / 4
		x := xOf(t)
		fill(x, ay, x+1, ay+5, axis)
		drawText7(img, fmt.Sprintf("%d", t), x+2, ay+14, axis)
	}

	drawText7(img, fmt.Sprintf("Part Two schedule: %d workers, %ds total", workers, total), leftPad, 16, white)

	f, err := os.Create(filepath.Join(outdir, "the-sum-of-its-parts.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawText7(img *image.RGBA, s string, x, y int, c color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}
