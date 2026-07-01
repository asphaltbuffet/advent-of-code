package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the fully settled sand pile from part two (the version with a
// floor) as a PNG: rock is slate, settled sand is gold shaded by depth so the
// internal packing is visible, and the source at 500,0 is marked red. The
// characteristic pyramid fanning out from the source down to the floor is the
// shape the two counts are measuring.
func (c Exercise) Vis(instr, outdir string) error {
	data := strings.Split(instr, "\n")
	d := Day14{Tiles: make(map[Point]Tile)}

	rocks, err := InputToPoints(data)
	if err != nil {
		return err
	}
	if err = d.AddRocks(rocks); err != nil {
		return err
	}

	d.MaxY += 2 // floor
	if err = d.BuildGraphWithFloor(root.Coord); err != nil && err != ErrVoidPath {
		return err
	}

	// Bounds from the settled tiles (sand fans wider than the rock scan).
	minX, maxX := root.Coord.X, root.Coord.X
	maxY := d.MaxY
	for p := range d.Tiles {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
	}

	const scale = 3
	const pad = 4
	gw := maxX - minX + 1
	gh := maxY + 1
	W := gw*scale + 2*pad
	H := gh*scale + 2*pad

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	rockCol := color.RGBA{0x2a, 0x30, 0x44, 0xff}
	srcCol := color.RGBA{0xff, 0x44, 0x55, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	put := func(gx, gy int, col color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(pad+(gx-minX)*scale+dx, pad+gy*scale+dy, col)
			}
		}
	}

	// Floor line spanning the whole width.
	for x := minX; x <= maxX; x++ {
		put(x, maxY, rockCol)
	}

	for p, t := range d.Tiles {
		switch t.Type {
		case Rock:
			put(p.X, p.Y, rockCol)
		case Sand:
			// Shade sand by depth: lighter near the source, warmer deeper.
			f := float64(p.Y) / float64(maxY)
			r := uint8(0xff)
			g := uint8(0xd0 - 0x50*f)
			b := uint8(0x60 - 0x40*f)
			put(p.X, p.Y, color.RGBA{r, g, b, 0xff})
		case Source:
			put(p.X, p.Y, srcCol)
		}
	}
	put(root.Coord.X, root.Coord.Y, srcCol)

	f, err := os.Create(filepath.Join(outdir, "regolith-reservoir.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
