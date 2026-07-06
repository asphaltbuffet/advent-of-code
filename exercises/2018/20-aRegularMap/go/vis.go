package exercises

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Vis renders the facility as a maze — rooms and the doors between them — with
// each room shaded by its shortest distance from the origin. Distance maps
// monotonically to brightness (dark near the start, bright far away), so the
// furthest reaches and the deep ≥1000-door regions read directly, and because the
// encoding is brightness alone it is inherently grayscale-safe.
func (e Exercise) Vis(instr, outdir string) error {
	dist, doors := walkWithDoors(instr)

	// Bounds of the room grid.
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	maxDist := 0
	for p, d := range dist {
		minX, maxX = min(minX, p.x), max(maxX, p.x)
		minY, maxY = min(minY, p.y), max(maxY, p.y)
		maxDist = max(maxDist, d)
	}

	// Each room is one pixel with a one-pixel gap for a possible door between it
	// and its neighbor: room (x,y) lives at (2*(x-minX)+1, 2*(y-minY)+1).
	w := 2*(maxX-minX) + 3
	h := 2*(maxY-minY) + 3
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	wall := color.RGBA{12, 12, 16, 255}
	for y := range h {
		for x := range w {
			img.Set(x, y, wall)
		}
	}

	shade := func(d int) color.RGBA {
		// Dark navy (near) → bright yellow (far); brightness rises with distance,
		// so the ramp is monotonic in luma and reads in grayscale.
		t := 0.0
		if maxDist > 0 {
			t = float64(d) / float64(maxDist)
		}
		r := uint8(20 + t*220)
		g := uint8(30 + t*198)
		b := uint8(70 - t*40)
		return color.RGBA{r, g, b, 255}
	}

	rx := func(p point) int { return 2*(p.x-minX) + 1 }
	ry := func(p point) int { return 2*(p.y-minY) + 1 }

	// Rooms.
	for p, d := range dist {
		img.Set(rx(p), ry(p), shade(d))
	}
	// Doors: the pixel between two connected rooms, shaded by the nearer room.
	for d := range doors {
		ax, ay := rx(d.a), ry(d.a)
		bx, by := rx(d.b), ry(d.b)
		mx, my := (ax+bx)/2, (ay+by)/2
		nd := min(dist[d.a], dist[d.b])
		img.Set(mx, my, shade(nd))
	}

	f, err := os.Create(filepath.Join(outdir, "a-regular-map.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

type door struct{ a, b point }

// walkWithDoors is walk() extended to also record every door (the connection
// between two adjacent rooms) so the maze walls can be drawn.
func walkWithDoors(instr string) (map[point]int, map[door]struct{}) {
	route := strings.TrimSpace(instr)
	dist := map[point]int{{0, 0}: 0}
	doors := map[door]struct{}{}

	var stack []point
	pos := point{0, 0}

	for i := range len(route) {
		switch route[i] {
		case '(':
			stack = append(stack, pos)
		case '|':
			pos = stack[len(stack)-1]
		case ')':
			pos = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case 'N', 'S', 'E', 'W':
			prev := pos
			switch route[i] {
			case 'N':
				pos.y--
			case 'S':
				pos.y++
			case 'E':
				pos.x++
			case 'W':
				pos.x--
			}
			doors[canonicalDoor(prev, pos)] = struct{}{}
			nd := dist[prev] + 1
			if d, ok := dist[pos]; !ok || nd < d {
				dist[pos] = nd
			}
		}
	}
	return dist, doors
}

// canonicalDoor orders the two rooms so each door has one key regardless of
// traversal direction.
func canonicalDoor(a, b point) door {
	if b.y < a.y || (b.y == a.y && b.x < a.x) {
		a, b = b, a
	}
	return door{a, b}
}
