package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 22.
type Exercise struct {
	common.BaseExercise
}

// Vis renders the settled brick pile as two side elevations (SVG-free PNG): the
// view looking down the y axis (x vs z) and down the x axis (y vs z), the two
// projections the puzzle itself uses. Each brick gets a stable colour so the
// same brick can be picked out in both views.
func (e Exercise) Vis(instr, outdir string) error {
	bricks := parseInput(instr)
	getBrickOrder(bricks) // settle the pile in place

	maxX, maxY, maxZ := 0, 0, 0
	for _, b := range bricks {
		maxX = max(maxX, b.max.X)
		maxY = max(maxY, b.max.Y)
		maxZ = max(maxZ, b.max.Z)
	}

	const scale = 8
	const gap = 40
	const pad = 20
	viewW := (maxX + 1) * scale // x/z view width
	view2W := (maxY + 1) * scale
	H := (maxZ+1)*scale + 2*pad
	W := viewW + gap + view2W + 2*pad

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	ground := color.RGBA{0x2a, 0x2e, 0x3e, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// z grows upward, so flip: screen y = H - pad - z*scale.
	rect := func(x0, x1, z0, z1, offX int, col color.RGBA) {
		px0 := pad + offX + x0*scale
		px1 := pad + offX + (x1+1)*scale
		py1 := H - pad - z0*scale
		py0 := H - pad - (z1+1)*scale
		for py := py0; py < py1; py++ {
			for px := px0; px < px1; px++ {
				if px >= 0 && px < W && py >= 0 && py < H {
					img.SetRGBA(px, py, col)
				}
			}
		}
	}

	// Ground line under each view.
	for x := pad; x < pad+viewW; x++ {
		img.SetRGBA(x, H-pad, ground)
	}
	for x := pad + viewW + gap; x < pad+viewW+gap+view2W; x++ {
		img.SetRGBA(x, H-pad, ground)
	}

	for i, b := range bricks {
		hue := math.Mod(float64(i)*137.508, 360)
		col := hsvColor(hue, 0.55, 0.92)
		// Left view: x horizontal, z vertical (looking along y).
		rect(b.min.X, b.max.X, b.min.Z, b.max.Z, 0, col)
		// Right view: y horizontal, z vertical (looking along x).
		rect(b.min.Y, b.max.Y, b.min.Z, b.max.Z, viewW+gap, col)
	}

	f, err := os.Create(filepath.Join(outdir, "sand-slabs.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// hsvColor converts HSV (h degrees, s,v in 0..1) to an RGBA colour.
func hsvColor(h, s, v float64) color.RGBA {
	c := v * s
	hp := h / 60
	x := c * (1 - math.Abs(math.Mod(hp, 2)-1))
	var r, g, b float64
	switch {
	case hp < 1:
		r, g, b = c, x, 0
	case hp < 2:
		r, g, b = x, c, 0
	case hp < 3:
		r, g, b = 0, c, x
	case hp < 4:
		r, g, b = 0, x, c
	case hp < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	m := v - c
	return color.RGBA{uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255), 0xff}
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	bricks := parseInput(instr)

	bricksBelow, bricksAbove := getBrickOrder(bricks)

	canDisintegrate := getNonSupportingBricks(bricks, bricksBelow, bricksAbove)

	return len(canDisintegrate), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	bricks := parseInput(instr)

	bricksBelow, bricksAbove := getBrickOrder(bricks)

	canDisintegrate := getNonSupportingBricks(bricks, bricksBelow, bricksAbove)

	total := countDisintegratable(bricks, bricksBelow, bricksAbove, canDisintegrate)

	return total, nil
}
