package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

// Vis renders the lava droplet as a stack of z-layer slices — a CT-scan style
// contact sheet. Each panel is one z plane: filled lava cubes are shown, tinted
// by depth (z) so the shape reads across panels, and the interior air pockets
// that part two subtracts are highlighted in red. Reading the panels left to
// right, top to bottom walks through the droplet one layer at a time.
func (c Exercise) Vis(instr, outdir string) error {
	cubes, err := parse(instr)
	if err != nil {
		return err
	}

	set := make(map[cube]bool, len(cubes))
	for _, cb := range cubes {
		set[cb] = true
	}

	// Flood-fill the exterior from a corner one cell outside the bounding box.
	// Everything reachable is "outside"; interior air cells never touched by the
	// flood are the trapped pockets that part two subtracts.
	lo := cube{minX - 1, minY - 1, minZ - 1}
	hi := cube{maxX + 1, maxY + 1, maxZ + 1}
	outside := map[cube]bool{lo: true}
	queue := []cube{lo}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, a := range adjacent {
			n := addCubes(cur, a)
			if n.x < lo.x || n.x > hi.x || n.y < lo.y || n.y > hi.y || n.z < lo.z || n.z > hi.z {
				continue
			}
			if set[n] || outside[n] {
				continue
			}
			outside[n] = true
			queue = append(queue, n)
		}
	}

	trapped := map[cube]bool{}
	for z := minZ; z <= maxZ; z++ {
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				cb := cube{x, y, z}
				if !set[cb] && !outside[cb] {
					trapped[cb] = true
				}
			}
		}
	}

	sx := maxX - minX + 1
	sy := maxY - minY + 1
	sz := maxZ - minZ + 1

	const scale = 4
	const gap = 6
	const pad = 8
	// Arrange z-slices in a roughly square contact sheet.
	perRow := int(math.Ceil(math.Sqrt(float64(sz))))
	nRows := int(math.Ceil(float64(sz) / float64(perRow)))

	panelW := sx*scale + gap
	panelH := sy*scale + gap
	W := perRow*panelW + 2*pad
	H := nRows*panelH + 2*pad

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	panelBg := color.RGBA{0x14, 0x18, 0x24, 0xff}
	trap := color.RGBA{0xff, 0x44, 0x55, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	for z := minZ; z <= maxZ; z++ {
		panel := z - minZ
		px0 := pad + (panel%perRow)*panelW
		py0 := pad + (panel/perRow)*panelH

		// Panel background.
		for yy := 0; yy < sy*scale; yy++ {
			for xx := 0; xx < sx*scale; xx++ {
				img.SetRGBA(px0+xx, py0+yy, panelBg)
			}
		}

		f := float64(z-minZ) / float64(sz) // depth 0..1
		lava := color.RGBA{
			uint8(0x40 + 0xb0*f),
			uint8(0xa0 - 0x40*f),
			uint8(0xf0 - 0x80*f),
			0xff,
		}

		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				cb := cube{x, y, z}
				var col color.RGBA
				switch {
				case set[cb]:
					col = lava
				case trapped[cb]:
					col = trap
				default:
					continue
				}
				bx := px0 + (x-minX)*scale
				by := py0 + (y-minY)*scale
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.SetRGBA(bx+dx, by+dy, col)
					}
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "boiling-boulders.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
