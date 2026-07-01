package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"
)

// Vis renders the heightmap as a PNG. Height is drawn as a grayscale relief (low
// ground dark, ridges bright) so the terrain reads on its own; the height-9
// ridge lines that wall off the basins stay brightest. The three largest basins
// (the ones part two multiplies) are flooded with distinct colorblind-safe
// accent colors, and every low point (part one) is marked. Because height is
// encoded by brightness, the relief still reads without color.
func (e Exercise) Vis(instr, outdir string) error {
	g, err := parse(instr)
	if err != nil {
		return err
	}

	// Label every non-9 cell with its basin id and collect sizes.
	label := make([][]int, g.rows)
	for i := range label {
		label[i] = make([]int, g.cols)
		for j := range label[i] {
			label[i][j] = -1
		}
	}

	type basin struct {
		id, size int
	}
	var basins []basin
	nextID := 0
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			if label[r][c] != -1 || g.h[r][c] == 9 {
				continue
			}
			id := nextID
			nextID++
			size := 0
			queue := [][2]int{{r, c}}
			label[r][c] = id
			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]
				size++
				for _, d := range dirs {
					nr, nc := cur[0]+d[0], cur[1]+d[1]
					if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
						continue
					}
					if label[nr][nc] != -1 || g.h[nr][nc] == 9 {
						continue
					}
					label[nr][nc] = id
					queue = append(queue, [2]int{nr, nc})
				}
			}
			basins = append(basins, basin{id, size})
		}
	}

	// Rank basins; flag the top three.
	ranked := append([]basin(nil), basins...)
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].size > ranked[j].size })
	topColor := map[int]color.RGBA{}
	accents := []color.RGBA{
		{0xE6, 0x9F, 0x00, 0xff}, // orange
		{0x56, 0xB4, 0xE9, 0xff}, // sky blue
		{0x00, 0x9E, 0x73, 0xff}, // bluish green
	}
	for i := 0; i < 3 && i < len(ranked); i++ {
		topColor[ranked[i].id] = accents[i]
	}

	lowSet := map[[2]int]bool{}
	for _, p := range g.lowPoints() {
		lowSet[p] = true
	}

	const scale = 8
	const pad = 4
	W := g.cols*scale + 2*pad
	H := g.rows*scale + 2*pad
	img := image.NewRGBA(image.Rect(0, 0, W, H))

	put := func(r, c int, col color.RGBA) {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(pad+c*scale+dx, pad+r*scale+dy, col)
			}
		}
	}

	wall := color.RGBA{0x20, 0x24, 0x2e, 0xff} // height-9 ridge lines
	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			h := g.h[r][c]
			var col color.RGBA
			if h == 9 {
				// Walls: dark slate so the basin network reads as light on dark.
				col = wall
			} else {
				// Basin floor: brighter for deeper (lower) ground so pools stand out.
				base := uint8(210 - h*16)
				col = color.RGBA{base, base, uint8(int(base) * 4 / 5), 0xff}
				if acc, ok := topColor[label[r][c]]; ok {
					// Blend the accent over the relief so depth still shows through.
					col = color.RGBA{
						uint8((int(acc.R)*2 + int(base)) / 3),
						uint8((int(acc.G)*2 + int(base)) / 3),
						uint8((int(acc.B)*2 + int(base)) / 3),
						0xff,
					}
				}
			}
			put(r, c, col)

			if lowSet[[2]int{r, c}] {
				// Mark each low point (part one) with a high-contrast vermilion dot
				// against the light basin floor.
				cx := pad + c*scale + scale/2
				cy := pad + r*scale + scale/2
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						img.SetRGBA(cx+dx, cy+dy, color.RGBA{0xD5, 0x2E, 0x00, 0xff})
					}
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "smoke-basin.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
