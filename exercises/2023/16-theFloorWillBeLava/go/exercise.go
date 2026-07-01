package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 16.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(input string) (any, error) {
	m, err := parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	m.Power(Point{0, 0}, Right)

	// m.DebugPrintContraption()
	// m.DebugPrintEnergized()

	return m.CountEnergized(), nil
}

// beam is a starting position and direction to test.
type beam struct {
	at  Point
	dir Direction
}

// Two returns the answer to the second part of the exercise. Every edge start
// is an independent simulation, so they are fanned out across worker goroutines.
func (e Exercise) Two(input string) (any, error) {
	m, err := parseInput(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Collect every edge entry point.
	var starts []beam
	for x := 0; x < m.Width; x++ {
		starts = append(starts, beam{Point{x, 0}, Down})
		starts = append(starts, beam{Point{x, m.Height - 1}, Up})
	}
	for y := 0; y < m.Height; y++ {
		starts = append(starts, beam{Point{0, y}, Right})
		starts = append(starts, beam{Point{m.Width - 1, y}, Left})
	}

	// Fan out across GOMAXPROCS workers; each clones the map and simulates.
	jobs := make(chan beam)
	results := make(chan int)
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for b := range jobs {
				c := m.Clone()
				c.Power(b.at, b.dir)
				results <- c.CountEnergized()
			}
		}()
	}
	go func() {
		for _, b := range starts {
			jobs <- b
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	best := 0
	for e := range results {
		if e > best {
			best = e
		}
	}
	return best, nil
}

// Vis renders the energized contraption for the best Part Two start beam: the
// beam floods a subset of tiles, glowing gold, while mirrors and splitters are
// picked out over the dark, un-energized floor.
func (e Exercise) Vis(input, outdir string) error {
	m, err := parseInput(input)
	if err != nil {
		return err
	}

	// Find the strongest start beam.
	var starts []beam
	for x := 0; x < m.Width; x++ {
		starts = append(starts, beam{Point{x, 0}, Down}, beam{Point{x, m.Height - 1}, Up})
	}
	for y := 0; y < m.Height; y++ {
		starts = append(starts, beam{Point{0, y}, Right}, beam{Point{m.Width - 1, y}, Left})
	}
	best, bestN := starts[0], -1
	for _, b := range starts {
		c := m.Clone()
		c.Power(b.at, b.dir)
		if n := c.CountEnergized(); n > bestN {
			best, bestN = b, n
		}
	}

	// Re-run the winning beam to capture its energized set.
	c := m.Clone()
	c.Power(best.at, best.dir)

	const cell = 6
	img := image.NewRGBA(image.Rect(0, 0, m.Width*cell, m.Height*cell))

	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	mirror := color.RGBA{0x4a, 0x5a, 0x7a, 0xff}
	glow := color.RGBA{0xff, 0xc8, 0x4a, 0xff}
	glowTile := color.RGBA{0xff, 0xf0, 0xc0, 0xff}

	fill := func(x, y int, col color.RGBA) {
		for dy := 0; dy < cell; dy++ {
			for dx := 0; dx < cell; dx++ {
				img.SetRGBA(x*cell+dx, y*cell+dy, col)
			}
		}
	}

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			p := Point{x, y}
			energized := len(c.Energized[p]) > 0
			tile := m.Points[p]
			switch {
			case tile != Empty && energized:
				fill(x, y, glowTile) // an active mirror/splitter
			case tile != Empty:
				fill(x, y, mirror)
			case energized:
				fill(x, y, glow)
			default:
				fill(x, y, bg)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "floor-will-be-lava.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
