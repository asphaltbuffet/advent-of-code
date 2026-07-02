package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis renders the final 3D state (part one, after 6 cycles) as a row of z-slices:
// each panel is the 2D x/y grid of active cubes at one z level, laid out from the
// most negative z to the most positive. The seed lives on z=0, so the structure
// is symmetric about the center panel. Active cubes are bright yellow on a dark
// field — a single active/inactive distinction that is inherently grayscale-safe —
// and each panel is labeled with its z.
func (e Exercise) Vis(instr, outdir string) error {
	active := simulateTo(parseActive(instr), 3, 6)

	// Bounds over x, y, z.
	minX, maxX := 1<<30, -(1 << 30)
	minY, maxY := 1<<30, -(1 << 30)
	minZ, maxZ := 1<<30, -(1 << 30)
	for c := range active {
		minX, maxX = min(minX, c[0]), max(maxX, c[0])
		minY, maxY = min(minY, c[1]), max(maxY, c[1])
		minZ, maxZ = min(minZ, c[2]), max(maxZ, c[2])
	}
	gw := maxX - minX + 1
	gh := maxY - minY + 1
	nZ := maxZ - minZ + 1

	const (
		cell = 8
		gap  = 16
		mL   = 12
		mT   = 46
	)
	panelW := gw * cell
	panelH := gh * cell
	W := mL*2 + nZ*panelW + (nZ-1)*gap
	H := mT + panelH + 24

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	cellBG := color.RGBA{0x1b, 0x1f, 0x26, 0xff}
	on := color.RGBA{0xF0, 0xE4, 0x42, 0xff}

	label := func(x, y int, s string, c color.RGBA) {
		(&font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}).DrawString(s)
	}
	white := color.RGBA{0xe8, 0xec, 0xf4, 0xff}
	gray := color.RGBA{0x9a, 0xa4, 0xb2, 0xff}
	label(mL, 24, "Conway Cubes: active cubes by z-slice (part one, after 6 cycles)", white)

	for zi := 0; zi < nZ; zi++ {
		z := minZ + zi
		px0 := mL + zi*(panelW+gap)
		for gy := 0; gy < gh; gy++ {
			for gx := 0; gx < gw; gx++ {
				c := cellBG
				if active[coord{minX + gx, minY + gy, z, 0}] {
					c = on
				}
				for dy := 0; dy < cell-1; dy++ {
					for dx := 0; dx < cell-1; dx++ {
						img.SetRGBA(px0+gx*cell+dx, mT+gy*cell+dy, c)
					}
				}
			}
		}
		zl := "z=0"
		if z != 0 {
			zl = fmt.Sprintf("z=%d", z)
		}
		label(px0+panelW/2-14, mT+panelH+16, zl, gray)
	}

	f, err := os.Create(filepath.Join(outdir, "conway-cubes.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
