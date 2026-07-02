package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// executed returns, for a program, the set of instruction indices the run visits
// (until it loops or terminates) in visit order.
func executed(prog []instruction) map[int]bool {
	visited := map[int]bool{}
	pc := 0
	for pc >= 0 && pc < len(prog) {
		if visited[pc] {
			break
		}
		visited[pc] = true
		switch prog[pc].op {
		case "acc":
			pc++
		case "jmp":
			pc += prog[pc].arg
		default:
			pc++
		}
	}
	return visited
}

// findFix returns the index of the single jmp<->nop flip that makes the program
// terminate, or -1.
func findFix(prog []instruction) int {
	for i := range prog {
		orig := prog[i].op
		switch orig {
		case "jmp":
			prog[i].op = "nop"
		case "nop":
			prog[i].op = "jmp"
		default:
			continue
		}
		_, ok := run(prog)
		prog[i].op = orig
		if ok {
			return i
		}
	}
	return -1
}

// Vis draws the boot program as two side-by-side columns of instruction rows: the
// original run that ends in an infinite loop (left) and the repaired run after
// the single jmp<->nop flip that lets it terminate (right). Each row is colored by
// whether that run executes it; the flipped instruction is marked, and the tail
// the fix newly reaches is visible on the right. Executed rows are the brightest,
// so the two paths read in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	prog, err := parse(instr)
	if err != nil {
		return err
	}
	n := len(prog)

	origRun := executed(prog)
	fixIdx := findFix(prog)
	patched := make([]instruction, n)
	copy(patched, prog)
	if fixIdx >= 0 {
		switch patched[fixIdx].op {
		case "jmp":
			patched[fixIdx].op = "nop"
		case "nop":
			patched[fixIdx].op = "jmp"
		}
	}
	fixedRun := executed(patched)

	const (
		rowH  = 3
		colW  = 220
		gap   = 40
		mL    = 20
		mT    = 50
		mB    = 20
	)
	W := mL*2 + colW*2 + gap
	H := mT + n*rowH + mB

	img := image.NewRGBA(image.Rect(0, 0, W, H))
	bg := color.RGBA{0x11, 0x14, 0x18, 0xff}
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	notRun := color.RGBA{0x24, 0x2b, 0x33, 0xff}  // dark: skipped
	ran := color.RGBA{0x56, 0xB4, 0xE9, 0xff}     // sky blue: executed
	terminated := color.RGBA{0x00, 0x9E, 0x73, 0xff} // green: executed on the terminating run
	fixCol := color.RGBA{0xF0, 0xE4, 0x42, 0xff}  // yellow: the flipped instruction

	drawCol := func(x0 int, runSet map[int]bool, runColor color.RGBA) {
		for i := 0; i < n; i++ {
			y := mT + i*rowH
			c := notRun
			if runSet[i] {
				c = runColor
			}
			if i == fixIdx {
				c = fixCol
			}
			for dy := 0; dy < rowH-1; dy++ {
				for dx := 0; dx < colW; dx++ {
					img.SetRGBA(x0+dx, y+dy, c)
				}
			}
		}
	}

	drawCol(mL, origRun, ran)
	drawCol(mL+colW+gap, fixedRun, terminated)

	label := func(x, y int, s string, c color.RGBA) {
		(&font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(c),
			Face: basicfont.Face7x13,
			Dot:  fixed.P(x, y),
		}).DrawString(s)
	}
	label(mL, 20, "original run (loops)", ran)
	label(mL+colW+gap, 20, "after 1 flip (terminates)", terminated)
	label(mL, 38, "yellow = flipped instruction; dark = not executed", color.RGBA{0x9a, 0xa4, 0xb2, 0xff})

	f, err := os.Create(filepath.Join(outdir, "handheld-halting.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
