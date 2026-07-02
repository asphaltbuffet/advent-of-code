package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// Vis renders the fully assembled image (part two) in the orientation where sea
// monsters appear, highlighting them. Calm sea is dark, rough water ('#') that is
// not part of a monster is mid blue, and cells belonging to a sea monster are
// bright yellow. The three states have a wide brightness gap, so the monsters
// stand out in grayscale as well as color; the water roughness answer is the
// count of the mid-blue cells.
func (e Exercise) Vis(instr, outdir string) error {
	tiles, err := parseTiles(instr)
	if err != nil {
		return err
	}
	grid := assemble(tiles)
	if grid == nil {
		return nil
	}

	// Find the orientation containing monsters.
	base := stitch(grid)
	monsterOffsets := monsterCellOffsets()
	var img []string
	for _, o := range orientations(base) {
		if countMonsters(o) > 0 {
			img = o
			break
		}
	}
	if img == nil {
		img = base
	}

	h := len(img)
	w := len(img[0])

	// Mark which cells belong to a monster.
	monster := make([][]bool, h)
	for i := range monster {
		monster[i] = make([]bool, w)
	}
	mh, mw := len(seaMonster), len(seaMonster[0])
	for r := 0; r+mh <= h; r++ {
		for c := 0; c+mw <= w; c++ {
			all := true
			for _, o := range monsterOffsets {
				if img[r+o[0]][c+o[1]] != '#' {
					all = false
					break
				}
			}
			if all {
				for _, o := range monsterOffsets {
					monster[r+o[0]][c+o[1]] = true
				}
			}
		}
	}

	const scale = 6
	W, H := w*scale, h*scale
	out := image.NewRGBA(image.Rect(0, 0, W, H))

	calm := color.RGBA{0x0d, 0x12, 0x1a, 0xff}    // dark: calm sea
	rough := color.RGBA{0x00, 0x72, 0xB2, 0xff}   // blue: rough water
	monCol := color.RGBA{0xF0, 0xE4, 0x42, 0xff}  // bright yellow: sea monster

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := calm
			if monster[y][x] {
				c = monCol
			} else if img[y][x] == '#' {
				c = rough
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					out.SetRGBA(x*scale+dx, y*scale+dy, c)
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "jurassic-jigsaw.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, out)
}

func monsterCellOffsets() [][2]int {
	var offsets [][2]int
	for r, row := range seaMonster {
		for c, ch := range row {
			if ch == '#' {
				offsets = append(offsets, [2]int{r, c})
			}
		}
	}
	return offsets
}
