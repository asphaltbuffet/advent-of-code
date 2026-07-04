package exercises //nolint:cyclop // asteroid monitoring visualization has inherently high package complexity

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
)

// Vis renders the asteroid field as a PNG:
//   - Dark background (#111111)
//   - All asteroids as small blue dots (#56B4E9)
//   - The best monitoring station as a bright white crosshair
//   - Asteroids colored by vaporization order (yellow → red gradient)
//   - The 200th asteroid highlighted in orange (#E69F00) with a ring
//
//nolint:gocognit,gocyclo,cyclop,funlen // asteroid monitoring animated visualization is inherently complex
func (e Exercise) Vis(instr, outdir string) error {
	asteroids := parseAsteroids(instr)

	station, _ := bestVisible(asteroids)

	// Compute grid dimensions.
	maxX, maxY := 0, 0
	for _, a := range asteroids {
		if a[0] > maxX {
			maxX = a[0]
		}
		if a[1] > maxY {
			maxY = a[1]
		}
	}

	const margin = 10
	cellSize := 500 / (maxX + 1)
	if g := 500 / (maxY + 1); g < cellSize {
		cellSize = g
	}
	if cellSize < 8 {
		cellSize = 8
	}
	w := (maxX+1)*cellSize + 2*margin
	h := (maxY+1)*cellSize + 2*margin

	bg := color.RGBA{0x11, 0x11, 0x11, 0xff}
	asteroidCol := color.RGBA{0x56, 0xB4, 0xE9, 0xff}
	stationCol := color.RGBA{0xff, 0xff, 0xff, 0xff}
	target200Col := color.RGBA{0xE6, 0x9F, 0x00, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Fill background.
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, bg)
		}
	}

	// Helper: fill a square cell centered on grid pos (gx,gy) with given radius.
	fillRect := func(gx, gy, r int, col color.RGBA) {
		cx := margin + gx*cellSize + cellSize/2
		cy := margin + gy*cellSize + cellSize/2
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				px, py := cx+dx, cy+dy
				if px >= 0 && px < w && py >= 0 && py < h {
					img.SetRGBA(px, py, col)
				}
			}
		}
	}

	// Draw ring (hollow circle) around a grid cell.
	drawRing := func(gx, gy, radius int, col color.RGBA) {
		cx := margin + gx*cellSize + cellSize/2
		cy := margin + gy*cellSize + cellSize/2
		for angle := 0.0; angle < 2*math.Pi; angle += 0.05 {
			px := cx + int(math.Round(float64(radius)*math.Cos(angle)))
			py := cy + int(math.Round(float64(radius)*math.Sin(angle)))
			if px >= 0 && px < w && py >= 0 && py < h {
				img.SetRGBA(px, py, col)
			}
		}
	}

	// Compute vaporization order from the station.
	type asteroid struct{ x, y int }
	dirMap := make(map[[2]int][]asteroid)
	for _, a := range asteroids {
		if [2]int{a[0], a[1]} == station {
			continue
		}
		dx := a[0] - station[0]
		dy := a[1] - station[1]
		norm := normalize(dx, dy)
		dirMap[norm] = append(dirMap[norm], asteroid{a[0], a[1]})
	}

	abs := func(v int) int {
		if v < 0 {
			return -v
		}
		return v
	}
	for k := range dirMap {
		sort.Slice(dirMap[k], func(i, j int) bool {
			di := abs(dirMap[k][i].x-station[0]) + abs(dirMap[k][i].y-station[1])
			dj := abs(dirMap[k][j].x-station[0]) + abs(dirMap[k][j].y-station[1])
			return di < dj
		})
	}

	clockwiseAngle := func(d [2]int) float64 {
		a := math.Atan2(float64(d[0]), -float64(d[1]))
		if a < 0 {
			a += 2 * math.Pi
		}
		return a
	}
	dirs := make([][2]int, 0, len(dirMap))
	for k := range dirMap {
		dirs = append(dirs, k)
	}
	sort.Slice(dirs, func(i, j int) bool {
		return clockwiseAngle(dirs[i]) < clockwiseAngle(dirs[j])
	})

	// Simulate vaporization and record order for each asteroid.
	type vap struct {
		x, y, order int
	}
	var vaporized []vap
	count := 0
	dirMapCopy := make(map[[2]int][]asteroid)
	for k, v := range dirMap {
		cp := make([]asteroid, len(v))
		copy(cp, v)
		dirMapCopy[k] = cp
	}
	for {
		anyLeft := false
		for _, d := range dirs {
			if len(dirMapCopy[d]) == 0 {
				continue
			}
			a := dirMapCopy[d][0]
			dirMapCopy[d] = dirMapCopy[d][1:]
			count++
			vaporized = append(vaporized, vap{a.x, a.y, count})
			anyLeft = true
		}
		if !anyLeft {
			break
		}
	}

	totalVaporized := len(vaporized)

	// Build a map from (x,y) -> vaporization order.
	vapOrder := make(map[[2]int]int, totalVaporized)
	for _, v := range vaporized {
		vapOrder[[2]int{v.x, v.y}] = v.order
	}

	// Draw all asteroids.
	for _, a := range asteroids {
		ax, ay := a[0], a[1]
		pos := [2]int{ax, ay}
		if pos == station {
			continue
		}
		order, isVaporized := vapOrder[pos]
		if isVaporized {
			// Color by vaporization order: early=yellow, late=red.
			f := float64(order-1) / float64(totalVaporized-1)
			r := uint8(0xff)
			g := uint8(math.Round(float64(0xD5) * (1.0 - f)))
			b := uint8(0x00)
			col := color.RGBA{r, g, b, 0xff}
			fillRect(ax, ay, 2, col)
		} else {
			fillRect(ax, ay, 2, asteroidCol)
		}
	}

	// Highlight the 200th vaporized asteroid distinctly.
	if len(vaporized) >= 200 {
		v200 := vaporized[199]
		fillRect(v200.x, v200.y, 3, target200Col)
		drawRing(v200.x, v200.y, 7, target200Col)
	}

	// Draw station crosshair.
	sx := margin + station[0]*cellSize + cellSize/2
	sy := margin + station[1]*cellSize + cellSize/2
	for d := -5; d <= 5; d++ {
		if sx+d >= 0 && sx+d < w {
			img.SetRGBA(sx+d, sy, stationCol)
		}
		if sy+d >= 0 && sy+d < h {
			img.SetRGBA(sx, sy+d, stationCol)
		}
	}
	fillRect(station[0], station[1], 2, stationCol)

	f, err := os.Create(filepath.Join(outdir, "vis.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
