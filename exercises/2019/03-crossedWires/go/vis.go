package exercises //nolint:cyclop // orbit-tree visualization has inherently high package complexity

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders both wire paths on a dark grid with their intersection points
// marked in white. Wire 1 is drawn in bright blue (Okabe-Ito #56B4E9) and
// wire 2 in bright orange (#E69F00). Intersections are white with a small
// cross marker. The origin is shown as a light-gray dot. Brightness differs
// enough between the two wires that the image reads in grayscale.
//
//nolint:gocognit,gocyclo,cyclop,funlen // wire-path intersection visualization is inherently complex
func (e Exercise) Vis(instr, outdir string) error {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 wire paths, got %d", len(lines))
	}

	wire1 := traceWire(lines[0])
	wire2 := traceWire(lines[1])

	// Compute bounding box over all visited points.
	minX, maxX := 0, 0
	minY, maxY := 0, 0
	for p := range wire1 {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	for p := range wire2 {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	// Add margin.
	const margin = 20
	spanX := maxX - minX
	spanY := maxY - minY
	if spanX <= 0 {
		spanX = 1
	}
	if spanY <= 0 {
		spanY = 1
	}

	// Scale so the longer axis is ~800px (capped to avoid huge images).
	const targetPx = 800
	scaleX := float64(targetPx) / float64(spanX)
	scaleY := float64(targetPx) / float64(spanY)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}
	if scale < 0.1 {
		scale = 0.1
	}
	if scale > 4 {
		scale = 4
	}

	imgW := int(float64(spanX)*scale) + 2*margin
	imgH := int(float64(spanY)*scale) + 2*margin

	// Map a logical point to image pixel (y-axis flipped so Up is up).
	toPixel := func(p point) (int, int) {
		px := int(float64(p.x-minX)*scale) + margin
		py := imgH - margin - int(float64(p.y-minY)*scale)
		return px, py
	}

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	// Fill background dark.
	bg := color.RGBA{0x11, 0x11, 0x11, 0xff}
	for y := range imgH {
		for x := range imgW {
			img.SetRGBA(x, y, bg)
		}
	}

	// Colors: Okabe-Ito palette, tuned for grayscale separability.
	// Wire 1: bright blue  #56B4E9 — perceived luminance ~58% (light in grayscale)
	// Wire 2: deep orange  #7B3F00 — perceived luminance ~22% (clearly darker in grayscale)
	// White intersections at 100% and dark background at ~7% give four distinct tones.
	wire1Color := color.RGBA{0x56, 0xB4, 0xE9, 0xff}  // bright blue  (~58% lum)
	wire2Color := color.RGBA{0xC0, 0x50, 0x00, 0xff}  // vivid orange (~28% lum)
	interColor := color.RGBA{0xFF, 0xFF, 0xFF, 0xff}  // white intersections
	originColor := color.RGBA{0x88, 0x88, 0x88, 0xff} // light gray origin

	setPixel := func(x, y int, c color.RGBA) {
		if x >= 0 && y >= 0 && x < imgW && y < imgH {
			img.SetRGBA(x, y, c)
		}
	}

	// Draw wire paths (1px each).
	for p := range wire1 {
		px, py := toPixel(p)
		setPixel(px, py, wire1Color)
	}
	for p := range wire2 {
		px, py := toPixel(p)
		// Don't overwrite intersection points yet.
		if _, isIntersection := wire1[p]; !isIntersection {
			setPixel(px, py, wire2Color)
		}
	}

	// Draw intersections as small cross markers.
	for p := range wire1 {
		if _, ok := wire2[p]; ok {
			px, py := toPixel(p)
			// 3x3 cross.
			for d := -2; d <= 2; d++ {
				setPixel(px+d, py, interColor)
				setPixel(px, py+d, interColor)
			}
		}
	}

	// Draw origin.
	ox, oy := toPixel(point{0, 0})
	for d := -3; d <= 3; d++ {
		setPixel(ox+d, oy, originColor)
		setPixel(ox, oy+d, originColor)
	}
	setPixel(ox, oy, originColor)

	f, err := os.Create(filepath.Join(outdir, "vis.png"))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
