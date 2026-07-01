package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 18.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	steps, _, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	return digArea(steps), nil
}

// Vis renders the Part One trench: the dig plan's boundary polygon with its
// enclosed interior filled, scaled to a comfortable canvas. (Part Two's hex
// plan encloses trillions of cubic metres — far too large to raster — so the
// visualisation uses the Part One dig.)
func (e Exercise) Vis(instr, outdir string) error {
	steps, _, err := parseInput(instr)
	if err != nil {
		return err
	}
	boundary, err := steps.GetBoundaryPoints()
	if err != nil {
		return err
	}

	// Bounds of the trench.
	minX, maxX := boundary[0].X, boundary[0].X
	minY, maxY := boundary[0].Y, boundary[0].Y
	for _, p := range boundary {
		minX, maxX = min(minX, p.X), max(maxX, p.X)
		minY, maxY = min(minY, p.Y), max(maxY, p.Y)
	}
	w := maxX - minX + 1
	h := maxY - minY + 1

	// Integer scale so the longer side is ~900px.
	scale := 900 / max(w, h)
	if scale < 1 {
		scale = 1
	}
	const pad = 12
	img := image.NewRGBA(image.Rect(0, 0, w*scale+2*pad, h*scale+2*pad))

	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	interior := color.RGBA{0x8a, 0x5a, 0x2a, 0xff} // dug lagoon
	trench := color.RGBA{0xff, 0xc8, 0x4a, 0xff}   // boundary
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	plot := func(gx, gy int, col color.RGBA) {
		x0, y0 := pad+(gx-minX)*scale, pad+(gy-minY)*scale
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(x0+dx, y0+dy, col)
			}
		}
	}

	// Fill interior via scanline point-in-polygon (ray cast) per grid row.
	for gy := minY; gy <= maxY; gy++ {
		for gx := minX; gx <= maxX; gx++ {
			if pointInPolygon(gx, gy, boundary) {
				plot(gx, gy, interior)
			}
		}
	}

	// Draw the trench boundary edges on top.
	for i := 0; i < len(boundary); i++ {
		a := boundary[i]
		b := boundary[(i+1)%len(boundary)]
		if a.X == b.X {
			for y := min(a.Y, b.Y); y <= max(a.Y, b.Y); y++ {
				plot(a.X, y, trench)
			}
		} else {
			for x := min(a.X, b.X); x <= max(a.X, b.X); x++ {
				plot(x, a.Y, trench)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "lavaduct-lagoon.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// pointInPolygon reports whether grid cell (x,y) lies inside the boundary
// polygon, using an even-odd ray cast on the rectilinear edges.
func pointInPolygon(x, y int, poly []Point) bool {
	inside := false
	n := len(poly)
	for i := 0; i < n; i++ {
		a, b := poly[i], poly[(i+1)%n]
		if (a.Y > y) != (b.Y > y) {
			// X coordinate of the edge at scanline y.
			cross := a.X + (y-a.Y)*(b.X-a.X)/(b.Y-a.Y)
			if x < cross {
				inside = !inside
			}
		}
	}
	return inside
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	_, steps, err := parseInput(instr)
	if err != nil {
		return nil, err
	}

	return digArea(steps), nil
}
