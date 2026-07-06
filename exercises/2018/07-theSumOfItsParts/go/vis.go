package exercises

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// job records one step's placement on the worker schedule.
type job struct {
	step       byte
	worker     int
	start, end int
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
	inProgress := map[byte]int{}
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

		for _, step := range readySteps(d, done, inProgress) {
			w := freeWorker(free)
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

func freeWorker(free []bool) int {
	for i, f := range free {
		if f {
			return i
		}
	}
	return -1
}

func fill(img *image.RGBA, x0, y0, x1, y1, imgW, imgH int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if x >= 0 && y >= 0 && x < imgW && y < imgH {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func drawSchedule(img *image.RGBA, jobs []job, workers, total, imgW, imgH, leftPad, topPad, plotW, rowH, rowGap int) {
	lane := color.RGBA{0x18, 0x1d, 0x25, 0xff}
	axis := color.RGBA{0x55, 0x5c, 0x66, 0xff}
	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
	barA := color.RGBA{0x56, 0xB4, 0xE9, 0xff}
	barB := color.RGBA{0xE6, 0x9F, 0x00, 0xff}
	inkA := color.RGBA{0x06, 0x1a, 0x26, 0xff}
	inkB := color.RGBA{0x1a, 0x12, 0x00, 0xff}

	xOf := func(t int) int { return leftPad + t*plotW/total }
	rowY := func(w int) int { return topPad + w*(rowH+rowGap) }

	for w := range workers {
		y := rowY(w)
		fill(img, leftPad, y, leftPad+plotW, y+rowH, imgW, imgH, lane)
		drawText7(img, fmt.Sprintf("W%d", w+1), 8, y+rowH/2+4, white)
	}
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
		fill(img, x0, y+2, x1, y+rowH-2, imgW, imgH, bar)
		if x1-x0 >= 8 {
			drawText7(img, string(j.step), x0+3, y+rowH/2+4, ink)
		}
	}
	plotH := workers*rowH + (workers-1)*rowGap
	ay := topPad + plotH + 6
	fill(img, leftPad, ay, leftPad+plotW, ay+1, imgW, imgH, axis)
	for k := range 5 {
		t := total * k / 4
		x := xOf(t)
		fill(img, x, ay, x+1, ay+5, imgW, imgH, axis)
		drawText7(img, strconv.Itoa(t), x+2, ay+14, axis)
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
		return errors.New("no steps to visualize")
	}

	const (
		leftPad = 40
		topPad  = 24
		botPad  = 20
		rowH    = 26
		rowGap  = 6
		plotW   = 900
	)
	plotH := workers*rowH + (workers-1)*rowGap
	imgW := leftPad + plotW + 12
	imgH := topPad + plotH + botPad

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	bg := color.RGBA{0x0c, 0x0f, 0x14, 0xff}
	fill(img, 0, 0, imgW, imgH, imgW, imgH, bg)

	drawSchedule(img, jobs, workers, total, imgW, imgH, leftPad, topPad, plotW, rowH, rowGap)

	white := color.RGBA{0xf0, 0xf4, 0xfa, 0xff}
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
