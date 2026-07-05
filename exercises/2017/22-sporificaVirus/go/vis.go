package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"maps"
	"os"
	"path/filepath"
)

// Vis animates the four-state virus spreading (GIF). It runs the burst sequence
// and snapshots the grid at regular intervals, then renders each snapshot as a
// paletted frame: clean=dark, weakened=blue, infected=teal, flagged=gold,
// carrier=red.
func (e Exercise) Vis(instr, outdir string) error {
	const bursts = 400_000 // enough to grow a rich bloom without a huge canvas
	const wantFrames = 300

	snaps, minX, maxX, minY, maxY := runVirusSim(instr, bursts, wantFrames)

	w := maxX - minX + 1
	h := maxY - minY + 1
	cell := max(720/max(w, h), 1)
	imgW, imgH := w*cell, h*cell

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 clean
		color.RGBA{0x3a, 0x6a, 0x9a, 0xff}, // 1 weakened
		color.RGBA{0x2f, 0xd0, 0x9a, 0xff}, // 2 infected
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 3 flagged
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 4 carrier
	}

	anim := &gif.GIF{}
	for si, s := range snaps {
		img := renderVirusFrame(s.state, s.pos, minX, minY, cell, imgW, imgH, pal)
		anim.Image = append(anim.Image, img)
		delay := 4
		if si == len(snaps)-1 {
			delay = 200
		}
		anim.Delay = append(anim.Delay, delay)
	}

	f, err := os.Create(filepath.Join(outdir, "sporifica-virus.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

type virusSnap struct {
	state map[pt]int
	pos   pt
}

func runVirusSim(instr string, bursts, wantFrames int) ([]virusSnap, int, int, int, int) {
	grid := parseGrid(instr)
	pos := pt{0, 0}
	dir := 0
	interval := bursts / wantFrames

	var snaps []virusSnap
	minX, maxX, minY, maxY := 0, 0, 0, 0
	for i := range bursts {
		switch grid[pos] {
		case clean:
			dir = (dir + 3) & 3
			grid[pos] = weakened
		case weakened:
			grid[pos] = infected
		case infected:
			dir = (dir + 1) & 3
			grid[pos] = flagged
		case flagged:
			dir = (dir + 2) & 3
			delete(grid, pos)
		}
		pos = pt{pos.x + dirs[dir].x, pos.y + dirs[dir].y}

		if pos.x < minX {
			minX = pos.x
		}
		if pos.x > maxX {
			maxX = pos.x
		}
		if pos.y < minY {
			minY = pos.y
		}
		if pos.y > maxY {
			maxY = pos.y
		}

		if (i+1)%interval == 0 {
			cp := make(map[pt]int, len(grid))
			maps.Copy(cp, grid)
			snaps = append(snaps, virusSnap{cp, pos})
		}
	}
	return snaps, minX, maxX, minY, maxY
}

func renderVirusFrame(state map[pt]int, pos pt, minX, minY, cell, imgW, imgH int, pal color.Palette) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, imgW, imgH), pal)
	fillCell := func(gx, gy, idx int) {
		x0, y0 := (gx-minX)*cell, (gy-minY)*cell
		for dy := range cell {
			for dx := range cell {
				img.SetColorIndex(x0+dx, y0+dy, uint8(idx))
			}
		}
	}
	for p, st := range state {
		fillCell(p.x, p.y, st)
	}
	fillCell(pos.x, pos.y, 4) // carrier on top
	return img
}
