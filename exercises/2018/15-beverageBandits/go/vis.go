package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"sort"
)

// Vis animates the part-two winning battle — the lowest elf attack power at which
// no elf dies — round by round (GIF). Walls are dim; elves are drawn as bright
// filled blocks and goblins as hollow rings, so the two sides differ by shape as
// well as color and stay distinct in grayscale. A unit dims as its hit points
// fall and disappears when it dies. One frame per round, holding on the final
// victory.
func (e Exercise) Vis(instr, outdir string) error {
	ap := 4
	for combat(instr, ap, true) == -1 {
		ap++
	}

	grid, units := parseCave(instr)
	rows := len(grid)
	cols := 0
	for _, r := range grid {
		if len(r) > cols {
			cols = len(r)
		}
	}

	const (
		cell = 12
		pad  = 6
	)
	imgW := cols*cell + 2*pad
	imgH := rows*cell + 2*pad
	pal := buildBattlePalette()

	anim := &gif.GIF{}
	anim.Image = append(anim.Image, renderBattleFrame(grid, units, pal, rows, imgW, imgH, cell, pad))
	anim.Delay = append(anim.Delay, 80)

	for {
		sort.Slice(units, func(i, j int) bool {
			if units[i].y != units[j].y {
				return units[i].y < units[j].y
			}
			return units[i].x < units[j].x
		})

		combatOver := runBattleTick(grid, units, ap)
		anim.Image = append(anim.Image, renderBattleFrame(grid, units, pal, rows, imgW, imgH, cell, pad))
		if combatOver {
			anim.Delay = append(anim.Delay, 500)
			break
		}
		anim.Delay = append(anim.Delay, 18)
	}

	f, err := os.Create(filepath.Join(outdir, "beverage-bandits.gif"))
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, anim)
}

func buildBattlePalette() color.Palette {
	pal := color.Palette{
		color.RGBA{0x0a, 0x0c, 0x12, 0xff},
		color.RGBA{0x2a, 0x30, 0x3c, 0xff},
	}
	for i := range 5 {
		f := 0.4 + 0.6*float64(i)/4
		pal = append(pal, color.RGBA{uint8(0x50 * f), uint8(0xb0 * f), uint8(0xff * f), 0xff})
	}
	for i := range 5 {
		f := 0.4 + 0.6*float64(i)/4
		pal = append(pal, color.RGBA{uint8(0xff * f), uint8(0x90 * f), uint8(0x10 * f), 0xff})
	}
	return pal
}

func unitCellIdx(u *unit) uint8 {
	const (
		elfBase = 2
		gobBase = 7
	)
	shade := min(max((u.hp-1)*5/200, 0), 4)
	if u.kind == 'E' {
		return uint8(elfBase + shade)
	}
	return uint8(gobBase + shade)
}

func setCellPaletted(img *image.Paletted, x, y, cell, pad int, idx uint8) {
	for dy := range cell {
		for dx := range cell {
			img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
		}
	}
}

func drawRingPaletted(img *image.Paletted, x, y, cell, pad int, idx uint8) {
	for dy := 1; dy < cell-1; dy++ {
		for dx := 1; dx < cell-1; dx++ {
			if dx <= 1 || dy <= 1 || dx >= cell-2 || dy >= cell-2 {
				img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
			}
		}
	}
}

func drawSolidPaletted(img *image.Paletted, x, y, cell, pad int, idx uint8) {
	for dy := 1; dy < cell-1; dy++ {
		for dx := 1; dx < cell-1; dx++ {
			img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
		}
	}
}

func renderBattleFrame(
	grid [][]byte, units []*unit, pal color.Palette, rows, imgW, imgH, cell, pad int,
) *image.Paletted {
	const bgIdx, wallIdx = 0, 1
	img := image.NewPaletted(image.Rect(0, 0, imgW, imgH), pal)
	for i := range img.Pix {
		img.Pix[i] = bgIdx
	}
	for y := range rows {
		for x := range len(grid[y]) {
			if grid[y][x] == '#' {
				setCellPaletted(img, x, y, cell, pad, wallIdx)
			}
		}
	}
	for _, u := range units {
		if !u.alive {
			continue
		}
		idx := unitCellIdx(u)
		if u.kind == 'E' {
			drawSolidPaletted(img, u.x, u.y, cell, pad, idx)
		} else {
			drawRingPaletted(img, u.x, u.y, cell, pad, idx)
		}
	}
	return img
}

func visAttack(units []*unit, u *unit, ap int) {
	target := adjacentEnemy(units, u.x, u.y, u.kind)
	if target == nil {
		return
	}
	power := 3
	if u.kind == 'E' {
		power = ap
	}
	target.hp -= power
	if target.hp <= 0 {
		target.alive = false
	}
}

func runBattleTick(grid [][]byte, units []*unit, ap int) bool {
	for _, u := range units {
		if !u.alive {
			continue
		}
		if !anyEnemyAlive(units, u.kind) {
			return true
		}
		occ := occupied(units)
		if adjacentEnemy(units, u.x, u.y, u.kind) == nil {
			if fx, fy, ok := stepToward(grid, occ, units, u); ok {
				u.x, u.y = fx, fy
			}
		}
		visAttack(units, u, ap)
	}
	return false
}
