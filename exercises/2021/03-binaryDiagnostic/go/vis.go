package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"
)

// Vis shows the part-two oxygen filter as a sequence of panels, one per filter
// step, laid left to right. The first panel is the full diagnostic report; each
// following panel is the shrinking set of candidates that still match the
// most-common bit in the current column, so the funnel narrowing 1000 numbers
// down to one reads across the image. Within every panel a number's bits are
// drawn (1 bright, 0 dark), and the column being filtered on is tinted gold.
func (e Exercise) Vis(instr, outdir string) error {
	nums := parse(instr)
	if len(nums) == 0 {
		return nil
	}
	width := len(nums[0])

	// Capture the surviving candidate set at the start of each filter step.
	steps := [][]string{sortNums(nums)}
	cur := append([]string(nil), nums...)
	for c := 0; c < width && len(cur) > 1; c++ {
		want := byte('0')
		if onesMajority(cur, c) {
			want = '1'
		}
		var kept []string
		for _, n := range cur {
			if n[c] == want {
				kept = append(kept, n)
			}
		}
		cur = kept
		steps = append(steps, sortNums(cur))
	}

	const cellW = 6
	const cellH = 2
	const panelGap = 18
	const pad = 10

	maxRows := len(steps[0])
	panelW := width * cellW
	W := len(steps)*panelW + (len(steps)-1)*panelGap + 2*pad
	H := maxRows*cellH + 2*pad

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x0d, 0x0f, 0x18, 0xff}
	one := color.RGBA{0x3f, 0xd0, 0x9a, 0xff}
	zero := color.RGBA{0x1c, 0x24, 0x30, 0xff}
	filterOne := color.RGBA{0xff, 0xc8, 0x4a, 0xff}
	filterZero := color.RGBA{0x5a, 0x48, 0x1c, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	for si, set := range steps {
		px0 := pad + si*(panelW+panelGap)
		filterCol := si // panel si is filtered on column si (last panel: none)

		for ry, n := range set {
			y0 := pad + ry*cellH
			for c := 0; c < width; c++ {
				var col color.RGBA
				switch {
				case c == filterCol && n[c] == '1':
					col = filterOne
				case c == filterCol:
					col = filterZero
				case n[c] == '1':
					col = one
				default:
					col = zero
				}
				x0 := px0 + c*cellW
				for dy := 0; dy < cellH; dy++ {
					for dx := 0; dx < cellW; dx++ {
						img.SetRGBA(x0+dx, y0+dy, col)
					}
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "binary-diagnostic.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func sortNums(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
