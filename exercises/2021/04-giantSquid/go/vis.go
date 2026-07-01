package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Vis animates the two boards that decide the puzzle side by side (GIF): the
// first board to win (part one, left) and the last board to win (part two,
// right). Each frame is one drawn number; cells are lit as they are marked, and
// once a board wins its completing row or column flashes gold and the animation
// holds. The left board fills quickly and stops; the right keeps going, marked
// but "losing", until it finally completes — the drama part two is built on.
func (e Exercise) Vis(instr, outdir string) error {
	draws, boards, err := parse(instr)
	if err != nil {
		return err
	}
	if len(boards) < 2 {
		return nil
	}

	firstIdx, firstDraw := winOrderIndex(draws, boards, true)
	lastIdx, lastDraw := winOrderIndex(draws, boards, false)

	// Fresh copies to animate (parse state above is now marked).
	_, freshFirst, _ := parse(instr)
	_, freshLast, _ := parse(instr)
	left := freshFirst[firstIdx]
	right := freshLast[lastIdx]

	frameCount := lastDraw + 1 // draws are 0-indexed; include the last

	const cell = 34
	const gap = 40
	const pad = 16
	const top = 30
	boardPx := boardSize * cell
	W := pad*2 + boardPx*2 + gap
	H := top + boardPx + pad

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 bg
		color.RGBA{0x1c, 0x24, 0x30, 0xff}, // 1 unmarked cell
		color.RGBA{0x3f, 0xd0, 0x9a, 0xff}, // 2 marked cell
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 3 winning line
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // 4 text/grid
		color.RGBA{0x66, 0x70, 0x80, 0xff}, // 5 dim text
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 6 current draw highlight
	}

	anim := &gif.GIF{}

	for d := 0; d < frameCount; d++ {
		n := draws[d]
		if d <= firstDraw {
			left.mark(n)
		}
		right.mark(n)

		img := image.NewPaletted(image.Rect(0, 0, W, H), pal)

		leftWon := d >= firstDraw
		rightWon := d >= lastDraw
		drawBoard(img, left, pad, top, cell, n, leftWon)
		drawBoard(img, right, pad+boardPx+gap, top, cell, n, rightWon)

		label(img, pad, 20, "FIRST TO WIN", 4)
		label(img, pad+boardPx+gap, 20, "LAST TO WIN", 4)
		label(img, W/2-24, H-6, "draw "+itoa(n), 5)

		anim.Image = append(anim.Image, img)
		delay := 22
		if d == 0 {
			delay = 120
		}
		if d == firstDraw || d == lastDraw {
			delay = 160
		}
		anim.Delay = append(anim.Delay, delay)
	}
	if len(anim.Delay) > 0 {
		anim.Delay[len(anim.Delay)-1] = 400
	}

	f, err := os.Create(filepath.Join(outdir, "giant-squid.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}

// winOrderIndex replays the draws and returns the index of the first (or last)
// board to win and the draw index at which it did.
func winOrderIndex(draws []int, boards []*board, first bool) (int, int) {
	// Work on fresh copies so callers keep clean boards.
	local := make([]*board, len(boards))
	for i, b := range boards {
		nb := &board{cells: b.cells}
		local[i] = nb
	}

	winIdx, winDraw := -1, -1
	remaining := len(local)
	for di, n := range draws {
		for bi, b := range local {
			if b.won {
				continue
			}
			if b.mark(n) {
				b.won = true
				remaining--
				if first {
					return bi, di
				}
				winIdx, winDraw = bi, di
				_ = remaining
			}
		}
	}

	return winIdx, winDraw
}

func drawBoard(img *image.Paletted, b *board, ox, oy, cell, current int, won bool) {
	// Determine the winning line to highlight, if this board has won.
	var winRow, winCol = -1, -1
	if won {
		winRow, winCol = winningLine(b)
	}

	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			idx := uint8(1)
			if b.marked[r][c] {
				idx = 2
			}
			if won && (r == winRow || c == winCol) && b.marked[r][c] {
				idx = 3
			}
			if b.cells[r][c] == current && b.marked[r][c] {
				idx = 6
			}
			x0 := ox + c*cell
			y0 := oy + r*cell
			for dy := 1; dy < cell-1; dy++ {
				for dx := 1; dx < cell-1; dx++ {
					img.SetColorIndex(x0+dx, y0+dy, idx)
				}
			}
			// number label centred-ish in the cell
			textIdx := uint8(5)
			if b.marked[r][c] {
				textIdx = 4
			}
			label(img, x0+cell/2-6, y0+cell/2+4, pad2(b.cells[r][c]), textIdx)
		}
	}
}

func winningLine(b *board) (int, int) {
	for r := 0; r < boardSize; r++ {
		full := true
		for c := 0; c < boardSize; c++ {
			if !b.marked[r][c] {
				full = false
				break
			}
		}
		if full {
			return r, -1
		}
	}
	for c := 0; c < boardSize; c++ {
		full := true
		for r := 0; r < boardSize; r++ {
			if !b.marked[r][c] {
				full = false
				break
			}
		}
		if full {
			return -1, c
		}
	}
	return -1, -1
}

func label(img *image.Paletted, x, y int, s string, idx uint8) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(img.Palette[idx]),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}

func pad2(n int) string {
	if n < 10 {
		return " " + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
