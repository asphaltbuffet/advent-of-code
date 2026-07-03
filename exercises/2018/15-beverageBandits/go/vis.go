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
	// Find the winning elf attack power exactly as part two does.
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
	W := cols*cell + 2*pad
	H := rows*cell + 2*pad

	// Palette: background, wall, then HP ramps for elf (blue) and goblin (orange).
	pal := color.Palette{
		color.RGBA{0x0a, 0x0c, 0x12, 0xff}, // 0 background / open floor
		color.RGBA{0x2a, 0x30, 0x3c, 0xff}, // 1 wall
	}
	const bgIdx, wallIdx = 0, 1
	// elf ramp indices 2..6 (dim->bright blue), goblin ramp 7..11 (dim->bright orange)
	for i := 0; i < 5; i++ {
		f := 0.4 + 0.6*float64(i)/4
		pal = append(pal, color.RGBA{uint8(0x50 * f), uint8(0xb0 * f), uint8(0xff * f), 0xff})
	}
	for i := 0; i < 5; i++ {
		f := 0.4 + 0.6*float64(i)/4
		pal = append(pal, color.RGBA{uint8(0xff * f), uint8(0x90 * f), uint8(0x10 * f), 0xff})
	}
	elfBase, gobBase := 2, 7

	hpShade := func(hp int) int {
		// map hp 1..200 to shade 0..4
		s := (hp - 1) * 5 / 200
		if s < 0 {
			s = 0
		}
		if s > 4 {
			s = 4
		}
		return s
	}

	setCell := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := 0; dy < cell; dy++ {
			for dx := 0; dx < cell; dx++ {
				img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
			}
		}
	}
	// ring draws a hollow square (border only) for goblins.
	ring := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := 1; dy < cell-1; dy++ {
			for dx := 1; dx < cell-1; dx++ {
				if dx <= 1 || dy <= 1 || dx >= cell-2 || dy >= cell-2 {
					img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
				}
			}
		}
	}
	// solid draws a filled block (inset) for elves.
	solid := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := 1; dy < cell-1; dy++ {
			for dx := 1; dx < cell-1; dx++ {
				img.SetColorIndex(pad+x*cell+dx, pad+y*cell+dy, idx)
			}
		}
	}

	render := func(us []*unit) *image.Paletted {
		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)
		for i := range img.Pix {
			img.Pix[i] = bgIdx
		}
		for y := 0; y < rows; y++ {
			for x := 0; x < len(grid[y]); x++ {
				if grid[y][x] == '#' {
					setCell(img, x, y, wallIdx)
				}
			}
		}
		for _, u := range us {
			if !u.alive {
				continue
			}
			if u.kind == 'E' {
				solid(img, u.x, u.y, uint8(elfBase+hpShade(u.hp)))
			} else {
				ring(img, u.x, u.y, uint8(gobBase+hpShade(u.hp)))
			}
		}
		return img
	}

	anim := &gif.GIF{}
	anim.Image = append(anim.Image, render(units))
	anim.Delay = append(anim.Delay, 80)

	// Replay the winning battle, capturing one frame per completed round.
	rounds := 0
	for {
		sort.Slice(units, func(i, j int) bool {
			if units[i].y != units[j].y {
				return units[i].y < units[j].y
			}
			return units[i].x < units[j].x
		})

		combatOver := false
		for _, u := range units {
			if !u.alive {
				continue
			}
			enemiesExist := false
			for _, en := range units {
				if en.alive && en.kind != u.kind {
					enemiesExist = true
					break
				}
			}
			if !enemiesExist {
				combatOver = true
				break
			}

			occ := occupied(units)
			if adjacentEnemy(units, u.x, u.y, u.kind) == nil {
				if fx, fy, ok := stepToward(grid, occ, units, u); ok {
					u.x, u.y = fx, fy
				}
			}
			if target := adjacentEnemy(units, u.x, u.y, u.kind); target != nil {
				power := 3
				if u.kind == 'E' {
					power = ap
				}
				target.hp -= power
				if target.hp <= 0 {
					target.alive = false
				}
			}
		}

		anim.Image = append(anim.Image, render(units))
		if combatOver {
			anim.Delay = append(anim.Delay, 500) // hold on victory
			break
		}
		anim.Delay = append(anim.Delay, 18)
		rounds++
	}

	f, err := os.Create(filepath.Join(outdir, "beverage-bandits.gif"))
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, anim)
}
