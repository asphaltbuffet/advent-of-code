package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders two z=0 cross-sections of the reactor. The left panel zooms into
// the -50..50 initialization region (part one), where the layered on/off steps
// carve visible pockets. The right panel pulls back to the full coordinate extent
// of every step (part two), where the initialization region shrinks to a speck
// against the giant far-field cuboids — the reason a voxel grid is hopeless and a
// signed-cuboid count is needed. On cells are bright sky blue, off (carved-out)
// cells are dark, so brightness carries the meaning without relying on hue alone.
func (e Exercise) Vis(instr, outdir string) error {
	steps, err := parse(instr)
	if err != nil {
		return err
	}

	const (
		res  = 101 // sampling grid per panel
		cell = 6
		pad  = 30
		gap  = 40
	)
	panelPx := res * cell
	W := pad*2 + panelPx*2 + gap
	H := pad*2 + panelPx + 28

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	fillRect(img, image.Rect(0, 0, W, H), color.RGBA{0x11, 0x14, 0x18, 0xff})

	// Okabe-Ito: sky blue for "on", kept bright; off carves back to dark.
	onCol := color.RGBA{0x56, 0xB4, 0xE9, 0xff}
	offCol := color.RGBA{0x22, 0x28, 0x30, 0xff}
	// Reddish-purple frame marks the init region inside the wide view.
	frame := color.RGBA{0xCC, 0x79, 0xA7, 0xff}

	// Left panel: initialization region only (steps clipped to -50..50).
	var initSteps []step
	for _, s := range steps {
		b := s.b
		if b.x1 < -50 || b.x2 > 50 || b.y1 < -50 || b.y2 > 50 || b.z1 < -50 || b.z2 > 50 {
			continue
		}
		initSteps = append(initSteps, s)
	}

	drawSlice(img, pad, pad+28, cell, res, -50, 50, initSteps, onCol, offCol)

	// Right panel: full extent of all steps in x/y.
	minX, maxX, minY, maxY := steps[0].b.x1, steps[0].b.x2, steps[0].b.y1, steps[0].b.y2
	for _, s := range steps {
		minX, maxX = min(minX, s.b.x1), max(maxX, s.b.x2)
		minY, maxY = min(minY, s.b.y1), max(maxY, s.b.y2)
	}
	span := max(maxX-minX, maxY-minY)
	cx, cy := (minX+maxX)/2, (minY+maxY)/2
	wlo, whi := cx-span/2, cx+span/2
	drawSliceXY(img, pad+panelPx+gap, pad+28, cell, res, wlo, whi, cy-span/2, cy+span/2, steps, onCol, offCol)

	// Overlay the ±50 init box on the wide view for scale.
	rcx := pad + panelPx + gap + int(float64(-wlo)/float64(whi-wlo)*float64(panelPx))
	rcy := pad + 28 + int(float64(span/2)/float64(span)*float64(panelPx))
	// The init region is sub-pixel at this zoom; draw a fixed-size marker so it
	// stays visible, with a leader label.
	const m = 5
	strokeRect(img, image.Rect(rcx-m, rcy-m, rcx+m, rcy+m), frame)
	label(img, rcx+m+4, rcy+4, "init region (left)", frame)

	white := color.RGBA{0xe8, 0xec, 0xf4, 0xff}
	label(img, pad, 20, "Initialization region (part one), z=0 slice", white)
	label(img, pad+panelPx+gap, 20, "Full extent (part two), z=0 slice", white)

	f, err := os.Create(filepath.Join(outdir, "reactor-reboot.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// drawSlice paints the z=0 on/off state over a square [lo,hi]^2 window at unit
// (integer-coordinate) resolution.
func drawSlice(img *image.RGBA, x0, y0, cell, res, lo, hi int, steps []step, onCol, offCol color.RGBA) {
	drawSliceXY(img, x0, y0, cell, res, lo, hi, lo, hi, steps, onCol, offCol)
}

// drawSliceXY paints the z=0 on/off state over [xlo,xhi]x[ylo,yhi], sampling res
// points per axis (used when the window is far larger than the cell grid).
func drawSliceXY(img *image.RGBA, x0, y0, cell, res, xlo, xhi, ylo, yhi int, steps []step, onCol, offCol color.RGBA) {
	for iy := 0; iy < res; iy++ {
		for ix := 0; ix < res; ix++ {
			gx := xlo + (xhi-xlo)*ix/(res-1)
			gy := ylo + (yhi-ylo)*iy/(res-1)
			on := false
			for _, s := range steps {
				b := s.b
				if gx >= b.x1 && gx <= b.x2 && gy >= b.y1 && gy <= b.y2 && b.z1 <= 0 && b.z2 >= 0 {
					on = s.on
				}
			}
			c := offCol
			if on {
				c = onCol
			}
			px := x0 + ix*cell
			// Flip y so +y points up.
			py := y0 + (res-1-iy)*cell
			fillRect(img, image.Rect(px, py, px+cell-1, py+cell-1), c)
		}
	}
}

func strokeRect(img *image.RGBA, r image.Rectangle, c color.RGBA) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		img.SetRGBA(x, r.Min.Y, c)
		img.SetRGBA(x, r.Max.Y, c)
	}
	for y := r.Min.Y; y <= r.Max.Y; y++ {
		img.SetRGBA(r.Min.X, y, c)
		img.SetRGBA(r.Max.X, y, c)
	}
}

func fillRect(img *image.RGBA, r image.Rectangle, c color.RGBA) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

func label(img *image.RGBA, x, y int, s string, c color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}
